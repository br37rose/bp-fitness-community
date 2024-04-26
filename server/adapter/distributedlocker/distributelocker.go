package distributedlocker

import (
	"log/slog"

	"github.com/bsm/redislock"
	"github.com/redis/go-redis/v9"
)

// Adapter provides interface for abstracting DistributedLocker generation.
type Adapter interface {
	Lock(key string)
	Lockf(format string, a ...any)
	Unlock(key string)
	Unlockf(format string, a ...any)
}

type distributedLockerAdapter struct {
	Logger *slog.Logger
	Redis  redis.UniversalClient
	Locker *redislock.Client
}

// NewAdapter constructor that returns the default DistributedLocker generator.
func NewAdapter(loggerp *slog.Logger, redisClient redis.UniversalClient) Adapter {

	// Create a new lock client.
	locker := redislock.New(redisClient)

	return distributedLockerAdapter{
		Logger: loggerp,
		Redis:  redisClient,
		Locker: locker,
	}
}

// Lock function blocks the current thread if the lock key is currently locked.
func (a distributedLockerAdapter) Lock(k string) {
	// u.DistributedLocker.Lock(k) //TODO: https://github.com/bsm/redislock/blob/main/README.md
	return
}

// Lockf function blocks the current thread if the lock key is currently locked.
func (u distributedLockerAdapter) Lockf(format string, a ...any) {
	// k := fmt.Sprintf(format, a...) //TODO: https://github.com/bsm/redislock/blob/main/README.md
	// u.DistributedLocker.Lock(k)
	return
}

// Unlock function blocks the current thread if the lock key is currently locked.
func (u distributedLockerAdapter) Unlock(k string) {
	// u.DistributedLocker.Unlock(k) //TODO: https://github.com/bsm/redislock/blob/main/README.md
	return
}

// Unlockf
func (u distributedLockerAdapter) Unlockf(format string, a ...any) {
	// k := fmt.Sprintf(format, a...)//TODO: https://github.com/bsm/redislock/blob/main/README.md
	// u.DistributedLocker.Unlock(k)
	return
}
