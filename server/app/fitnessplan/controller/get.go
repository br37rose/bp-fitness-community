package controller

import (
	"context"
	"log/slog"

	domain "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/fitnessplan/datastore"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (c *FitnessPlanControllerImpl) GetByID(ctx context.Context, id primitive.ObjectID) (*domain.FitnessPlan, error) {
	// Retrieve from our database the record for the specific id.
	m, err := c.FitnessPlanStorer.GetByID(ctx, id)
	if err != nil {
		c.Logger.Error("database get by id error", slog.Any("error", err))
		return nil, err
	}

	return m, err
}
