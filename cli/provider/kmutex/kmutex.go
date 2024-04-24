package kmutex

import (
	"fmt"

	"github.com/im7mortal/kmutex"
)

// Provider provides interface for abstracting KMutex generation.
type Provider interface {
	Lock(key string)
	Lockf(format string, a ...any)
	Unlock(key string)
	Unlockf(format string, a ...any)
}

type kmutexProvider struct {
	KMutex *kmutex.Kmutex
}

// NewProvider constructor that returns the default KMutex generator.
func NewProvider() Provider {
	kmux := kmutex.New()
	return kmutexProvider{
		KMutex: kmux,
	}
}

// Lock function blocks the current thread if the lock key is currently locked.
func (u kmutexProvider) Lock(k string) {
	u.KMutex.Lock(k)
	return
}

// Lockf function blocks the current thread if the lock key is currently locked.
func (u kmutexProvider) Lockf(format string, a ...any) {
	k := fmt.Sprintf(format, a...)
	u.KMutex.Lock(k)
	return
}

// Unlock function blocks the current thread if the lock key is currently locked.
func (u kmutexProvider) Unlock(k string) {
	u.KMutex.Unlock(k)
	return
}

// Unlockf
func (u kmutexProvider) Unlockf(format string, a ...any) {
	k := fmt.Sprintf(format, a...)
	u.KMutex.Unlock(k)
	return
}
