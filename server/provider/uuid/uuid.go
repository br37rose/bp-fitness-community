package uuid

import uuid "github.com/segmentio/ksuid"

// Provider provides interface for abstracting UUID generation.
type Provider interface {
	NewUUID() string
}

type uuidProvider struct {
}

// NewProvider constructor that returns the default UUID generator.
func NewProvider() Provider {
	return uuidProvider{}
}

// NewUUID generates a new UUID.
func (u uuidProvider) NewUUID() string {
	return uuid.New().String()
}
