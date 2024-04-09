package controller

import (
	"context"
	"log/slog"

	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/workout/datastore"
)

func (c *WorkoutControllerImpl) ListByFilter(ctx context.Context, f *datastore.WorkoutListFilter) (*datastore.WorkoutistResult, error) {
	workpouts, err := c.WorkoutStorer.ListByFilter(ctx, f)
	if err != nil {
		c.Logger.Error("database list by filter error", slog.Any("error", err))
		return nil, err
	}
	return workpouts, err
}
