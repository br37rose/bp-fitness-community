package controller

import (
	"context"
	"log/slog"

	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/trainingprogram/datastore"
)

func (c *TrainingprogramControllerImpl) ListByFilter(ctx context.Context, f *datastore.TrainingProgramListFilter) (*datastore.TrainingProgramListResult, error) {
	result, err := c.TpStorer.ListByFilter(ctx, f)
	if err != nil {
		c.Logger.Error("database list by filter error", slog.Any("error", err))
		return nil, err
	}

	return result, nil
}
