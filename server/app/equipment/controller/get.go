package controller

import (
	"context"

	domain "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/equipment/datastore"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log/slog"
)

func (c *EquipmentControllerImpl) GetByID(ctx context.Context, id primitive.ObjectID) (*domain.Equipment, error) {
	// Retrieve from our database the record for the specific id.
	m, err := c.EquipmentStorer.GetByID(ctx, id)
	if err != nil {
		c.Logger.Error("database get by id error", slog.Any("error", err))
		return nil, err
	}

	return m, err
}
