package controller

import (
	"context"
	"time"

	domain "github.com/bci-innovation-labs/bp8fitnesscommunity-cli/app/exercise/datastore"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-cli/utils/httperror"
	"log/slog"
)

func (c *ExerciseControllerImpl) UpdateByID(ctx context.Context, ns *domain.Exercise) (*domain.Exercise, error) {
	// Fetch the original exercise.
	os, err := c.ExerciseStorer.GetByID(ctx, ns.ID)
	if err != nil {
		c.Logger.Error("database get by id error", slog.Any("error", err))
		return nil, err
	}
	if os == nil {
		return nil, httperror.NewForBadRequestWithSingleField("id", "exercise does not exist")
	}

	// Modify our original exercise.
	os.ModifiedAt = time.Now()
	os.Type = ns.Type
	os.Status = ns.Status
	os.Name = ns.Name

	// Save to the database the modified exercise.
	if err := c.ExerciseStorer.UpdateByID(ctx, os); err != nil {
		c.Logger.Error("database update by id error", slog.Any("error", err))
		return nil, err
	}

	return os, nil
}
