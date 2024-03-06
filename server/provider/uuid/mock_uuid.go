package uuid

import (
	uuid "github.com/segmentio/ksuid"
	"github.com/stretchr/testify/mock"
)

// MockProvider mocks uuid provider
type MockProvider struct {
	mock.Mock
}

// NewUUID returns the mocked uuid
func (m MockProvider) NewUUID() uuid.KSUID {
	args := m.Called()
	return args.Get(0).(uuid.KSUID)
}
