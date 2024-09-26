package controller

import (
	"context"

	domain "github.com/bci-innovation-labs/bp8fitnesscommunity-cli/app/videocategory/datastore"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log/slog"
)

func (c *VideoCategoryControllerImpl) GetByID(ctx context.Context, id primitive.ObjectID) (*domain.VideoCategory, error) {
	// Retrieve from our database the record for the specific id.
	m, err := c.VideoCategoryStorer.GetByID(ctx, id)
	if err != nil {
		c.Logger.Error("database get by id error", slog.Any("error", err))
		return nil, err
	}

	return m, err
}

// Auto-generated comment for change 26
