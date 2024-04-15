package controller

import (
	"context"
	"log/slog"

	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/trainingprogram/datastore"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (c *TrainingprogramControllerImpl) GetByID(ctx context.Context, id primitive.ObjectID) (*datastore.TrainingProgram, error) {
	tp, err := c.TpStorer.GetByID(ctx, id)
	if err != nil {
		c.Logger.Error("database get by id error", slog.Any("error", err))
		return nil, err
	}
	return tp, err
}
