package controller

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"log/slog"

	tag_s "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/tag/datastore"
)

func (c *TagControllerImpl) GetByID(ctx context.Context, id primitive.ObjectID) (*tag_s.Tag, error) {
	// Retrieve from our database the record for the specific id.
	m, err := c.TagStorer.GetByID(ctx, id)
	if err != nil {
		c.Logger.Error("database get by id error", slog.Any("error", err))
		return nil, err
	}
	return m, err
}
