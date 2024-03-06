package controller

import (
	"context"
	"log/slog"

	domain "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/nutritionplan/datastore"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/config/constants"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (c *NutritionPlanControllerImpl) ListAsSelectOptionByFilter(ctx context.Context, f *domain.NutritionPlanListFilter) ([]*domain.NutritionPlanAsSelectOption, error) {
	// Extract from our session the following data.
	userID, _ := ctx.Value(constants.SessionUserID).(primitive.ObjectID)
	userOrgID, _ := ctx.Value(constants.SessionUserOrganizationID).(primitive.ObjectID)

	// Force tenancy on list.
	c.Logger.Debug("fetching nutritionplans now...", slog.Any("userID", userID))
	f.OrganizationID = userOrgID

	// List.
	m, err := c.NutritionPlanStorer.ListAsSelectOptionByFilter(ctx, f)
	if err != nil {
		c.Logger.Error("database list by filter error", slog.Any("error", err))
		return nil, err
	}
	c.Logger.Debug("fetched nutritionplans", slog.Any("m", m))
	return m, err
}
