package controller

import (
	"context"
	"log/slog"

	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/fitnesschallenge/datastore"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (c *FitnessChallengeControllerImpl) GetByID(ctx context.Context, id primitive.ObjectID) (*datastore.FitnessChallenge, error) {
	tp, err := c.Storer.GetByID(ctx, id)
	if err != nil {
		c.Logger.Error("database get by id error", slog.Any("error", err))
		return nil, err
	}
	return tp, err
}
