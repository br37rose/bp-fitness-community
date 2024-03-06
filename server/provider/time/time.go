package time

import (
	"time"
)

// Provider provides an interface for abstracting time.
type Provider interface {
	Now() time.Time
}

type timeProvider struct{}

// NewProvider Provider contructor that returns the default time provider.
func NewProvider() Provider {
	return timeProvider{}
}

// Now returns the current time.
func (t timeProvider) Now() time.Time {
	return time.Now()
}
