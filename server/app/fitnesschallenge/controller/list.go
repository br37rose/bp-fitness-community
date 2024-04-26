package controller

import (
	"context"
	"log/slog"

	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/fitnesschallenge/datastore"
)

func (c *FitnessChallengeControllerImpl) ListByFilter(
	ctx context.Context,
	f *datastore.FitnessChallengeListFilter) (*datastore.FitnessChallengeListResult, error) {

	result, err := c.Storer.ListByFilter(ctx, f)
	if err != nil {
		c.Logger.Error("database list by filter error", slog.Any("error", err))
		return nil, err
	}

	return result, nil
}
