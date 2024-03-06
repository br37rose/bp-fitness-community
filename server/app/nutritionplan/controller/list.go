package controller

import (
	"context"
	"log/slog"

	domain "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/nutritionplan/datastore"
	user_d "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/user/datastore"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/config/constants"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (c *NutritionPlanControllerImpl) ListByFilter(ctx context.Context, f *domain.NutritionPlanListFilter) (*domain.NutritionPlanListResult, error) {
	// Extract from our session the following data.
	orgID := ctx.Value(constants.SessionUserOrganizationID).(primitive.ObjectID)
	userID := ctx.Value(constants.SessionUserID).(primitive.ObjectID)
	userRole := ctx.Value(constants.SessionUserRole).(int8)

	// Apply protection based on ownership and role.
	switch userRole {
	case user_d.UserRoleRoot:
		break
	case user_d.UserRoleAdmin:
		f.OrganizationID = orgID // Force organization tenancy restrictions.
		break
	case user_d.UserRoleTrainer:
		f.OrganizationID = orgID // Force organization tenancy restrictions.
		break
	case user_d.UserRoleMember:
		f.OrganizationID = orgID // Force organization tenancy restrictions.
		f.UserID = userID        // Force filtering only for specific logged in user.
		break
	}

	c.Logger.Debug("fetching nutritionplans now...", slog.Any("userID", userID))

	aa, err := c.NutritionPlanStorer.ListByFilter(ctx, f)
	if err != nil {
		c.Logger.Error("database list by filter error", slog.Any("error", err))
		return nil, err
	}
	c.Logger.Debug("fetched nutritionplans", slog.Any("aa", aa))
	return aa, err
}
