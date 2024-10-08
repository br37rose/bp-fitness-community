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
	startDT := time.Now()
	a.Logger.Debug(fmt.Sprintf("locking for key: %v", k))

	// Retry every 250ms, for up-to 20x
	backoff := redislock.LimitRetry(redislock.LinearBackoff(250*time.Millisecond), 20)

	// Obtain lock with retry
	lock, err := a.Locker.Obtain(ctx, k, time.Minute, &redislock.Options{
		RetryStrategy: backoff,
	})
	if err == redislock.ErrNotObtained {
		nowDT := time.Now()
		diff := nowDT.Sub(startDT)
		a.Logger.Error("could not obtain lock",
			slog.String("key", k),
			slog.Time("start_dt", startDT),
			slog.Time("now_dt", nowDT),
			slog.Any("duration_in_minutes", diff.Minutes()))
		return
	} else if err != nil {
		a.Logger.Error("failed obtaining lock because of the following error: %v", err, slog.String("key", k))
		return
	}

	// DEVELOPERS NOTE:
	// The `map` datastructure in Golang is not concurrently safe, therefore we
	// need to use mutex to coordinate access of our `LockInstances` map
	// resource between all the goroutines.
	a.Mutex.Lock()
	defer a.Mutex.Unlock()

	if a.LockInstances != nil { // Defensive code.
		a.LockInstances[k] = lock
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
		a.Logger.Error("could not obtain to unlock", slog.String("key", k))
	}
	return
}

// Unlockf
func (u distributedLockerAdapter) Unlockf(ctx context.Context, format string, a ...any) {
	k := fmt.Sprintf(format, a...) //TODO: https://github.com/bsm/redislock/blob/main/README.md
	u.Unlock(ctx, k)
	return
}
