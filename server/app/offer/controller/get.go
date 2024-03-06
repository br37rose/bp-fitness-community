package controller

import (
	"context"

	domain "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/offer/datastore"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/utils/httperror"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log/slog"
)

func (c *OfferControllerImpl) GetByID(ctx context.Context, id primitive.ObjectID) (*domain.Offer, error) {
	// Retrieve from our database the record for the specific id.
	m, err := c.OfferStorer.GetByID(ctx, id)
	if err != nil {
		c.Logger.Error("database get by id error", slog.Any("error", err))
		return nil, err
	}
	if m == nil {
		return nil, httperror.NewForBadRequestWithSingleField("id", "workout program type does not exist")
	}
	return m, err
}
