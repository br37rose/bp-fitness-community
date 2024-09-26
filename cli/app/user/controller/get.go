package controller

import (
	"context"

	domain "github.com/bci-innovation-labs/bp8fitnesscommunity-cli/app/user/datastore"
)

// CreateInitialRootAdmin function creates the initial root administrator if not previously created.
func (c *UserControllerImpl) GetUserBySessionUUID(ctx context.Context, sessionUUID string) (*domain.User, error) {
	panic("TODO: IMPLEMENT")
	return nil, nil
}

// Auto-generated comment for change 16
