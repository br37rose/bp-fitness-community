package controller

import (
	"context"
	"log/slog"

	domain "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/nutritionplan/datastore"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (c *NutritionPlanControllerImpl) GetByID(ctx context.Context, id primitive.ObjectID) (*domain.NutritionPlan, error) {
	// Retrieve from our database the record for the specific id.
	m, err := c.NutritionPlanStorer.GetByID(ctx, id)
	if err != nil {
		c.Logger.Error("database get by id error", slog.Any("error", err))
		return nil, err
	}

	return m, err
}
