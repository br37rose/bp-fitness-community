package controller

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"log/slog"

	domain "github.com/bci-innovation-labs/bp8fitnesscommunity-cli/app/exercise/datastore"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-cli/utils/httperror"
)

func (c *ExerciseControllerImpl) GetByID(ctx context.Context, id primitive.ObjectID) (*domain.Exercise, error) {
	// Retrieve from our database the record for the specific id.
	m, err := c.ExerciseStorer.GetByID(ctx, id)
	if err != nil {
		c.Logger.Error("database get by id error", slog.Any("error", err))
		return nil, err
	}
	if m == nil {
		return nil, httperror.NewForBadRequestWithSingleField("id", "exercise does not exist")
	}
	return m, err
}
