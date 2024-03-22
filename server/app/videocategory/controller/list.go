package controller

import (
	"context"

	user_d "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/user/datastore"
	domain "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/videocategory/datastore"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/config/constants"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log/slog"
)

func (c *VideoCategoryControllerImpl) ListByFilter(ctx context.Context, f *domain.VideoCategoryListFilter) (*domain.VideoCategoryListResult, error) {
	// Extract from our session the following data.
	orgID := ctx.Value(constants.SessionUserOrganizationID).(primitive.ObjectID)
	userID := ctx.Value(constants.SessionUserID).(primitive.ObjectID)
	userRole := ctx.Value(constants.SessionUserRole).(int8)

	// Apply protection based on ownership and role.
	if userRole != user_d.UserRoleRoot {
		f.OrganizationID = orgID // Force organization tenancy restrictions.
	}

	c.Logger.Debug("fetching videocategorys now...", slog.Any("userID", userID))

	aa, err := c.VideoCategoryStorer.ListByFilter(ctx, f)
	if err != nil {
		c.Logger.Error("database list by filter error", slog.Any("error", err))
		return nil, err
	}
	c.Logger.Debug("fetched videocategorys", slog.Any("aa", aa))
	return aa, err
}