package controller

import (
	"context"

	domain "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/organization/datastore"
	user_d "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/user/datastore"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/config/constants"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log/slog"
)

func (c *OrganizationControllerImpl) ListByFilter(ctx context.Context, f *domain.OrganizationListFilter) (*domain.OrganizationListResult, error) {
	// Extract from our session the following data.
	userID := ctx.Value(constants.SessionUserID).(primitive.ObjectID)
	userRole := ctx.Value(constants.SessionUserRole).(int8)

	// Apply filtering based on ownership and role.
	if userRole != user_d.UserRoleRoot {
		f.UserID = userID
		f.UserRole = userRole
	}

	m, err := c.OrganizationStorer.ListByFilter(ctx, f)
	if err != nil {
		c.Logger.Error("database list by filter error", slog.Any("error", err))
		return nil, err
	}
	return m, err
}
