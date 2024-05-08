package distributedmutex

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/bsm/redislock"
	"github.com/redis/go-redis/v9"
)

// Adapter provides interface for abstracting distributedmutex generation.
type Adapter interface {
	Lock(ctx context.Context, key string)
	Lockf(ctx context.Context, format string, a ...any)
	Unlock(ctx context.Context, key string)
	Unlockf(ctx context.Context, format string, a ...any)
}

type distributedLockerAdapter struct {
	Logger        *slog.Logger
	Redis         redis.UniversalClient
	Locker        *redislock.Client
	LockInstances map[string]*redislock.Lock
	Mutex         *sync.Mutex // Add a mutex for synchronization with goroutines
}

// NewAdapter constructor that returns the default DistributedLocker generator.
func NewAdapter(loggerp *slog.Logger, redisClient redis.UniversalClient) Adapter {
	loggerp.Debug("distributed mutex starting and connecting...")

	// Create a new lock client.
	locker := redislock.New(redisClient)

	loggerp.Debug("distributed mutex initialized")

	return distributedLockerAdapter{
		Logger:        loggerp,
		Redis:         redisClient,
		Locker:        locker,
		LockInstances: make(map[string]*redislock.Lock, 0),
		Mutex:         &sync.Mutex{}, // Initialize the mutex
	}
}

// Lock function blocks the current thread if the lock key is currently locked.
func (a distributedLockerAdapter) Lock(ctx context.Context, k string) {
	a.Logger.Debug(fmt.Sprintf("locking fo key: %v", k))

	// Retry every 250ms, for up-to 10x
	backoff := redislock.LimitRetry(redislock.LinearBackoff(250*time.Millisecond), 10)

	// Obtain lock with retry
	lock, err := a.Locker.Obtain(ctx, k, time.Minute, &redislock.Options{
		RetryStrategy: backoff,
	})
	if err == redislock.ErrNotObtained {
		a.Logger.Error("could not obtain lock!")
	} else if err != nil {
		a.Logger.Error("failed obtaining lock because of the following error: %v", err)
		return
	}

	a.Mutex.Lock()
	defer a.Mutex.Unlock()

	if a.LockInstances != nil { // Defensive code.
		a.LockInstances[k] = lock //TODO: Why am I getting error "fatal error: concurrent map writes" here?
	}
}

// Lockf function blocks the current thread if the lock key is currently locked.
func (u distributedLockerAdapter) Lockf(ctx context.Context, format string, a ...any) {
	k := fmt.Sprintf(format, a...)
	u.Lock(ctx, k)
	return
}

// Unlock function blocks the current thread if the lock key is currently locked.
func (a distributedLockerAdapter) Unlock(ctx context.Context, k string) {
	a.Logger.Debug(fmt.Sprintf("unlocking for key: %v", k))

	lockInstance, ok := a.LockInstances[k]
	if ok {
		defer lockInstance.Release(ctx)
	} else {
		a.Logger.Error("Could not obtain lock to unlock!")
	}
	return
}

// Unlockf
func (u distributedLockerAdapter) Unlockf(ctx context.Context, format string, a ...any) {
	k := fmt.Sprintf(format, a...) //TODO: https://github.com/bsm/redislock/blob/main/README.md
	u.Unlock(ctx, k)
	return
}
