package kmutex

import (
	uuid "github.com/segmentio/ksuid"
	"github.com/stretchr/testify/mock"
)

// MockProvider mocks uuid provider
type MockProvider struct {
	mock.Mock
}

// NewProvider returns the mocked uuid
func (m MockProvider) NewProvider() uuid.KSUID {
	args := m.Called()
	return args.Get(0).(uuid.KSUID)
}
