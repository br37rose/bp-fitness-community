package controller

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"log/slog"

	domain "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/organization/datastore"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/utils/httperror"
)

func (c *OrganizationControllerImpl) GetByID(ctx context.Context, id primitive.ObjectID) (*domain.Organization, error) {
	// Retrieve from our database the record for the specific id.
	m, err := c.OrganizationStorer.GetByID(ctx, id)
	if err != nil {
		c.Logger.Error("database get by id error", slog.Any("error", err))
		return nil, err
	}
	if m == nil {
		return nil, httperror.NewForBadRequestWithSingleField("id", "organization does not exist")
	}
	return m, err
}
