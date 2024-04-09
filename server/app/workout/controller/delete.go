package controller

import (
	"context"
	"log/slog"
	"time"

	w_s "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/workout/datastore"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/utils/httperror"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (impl *WorkoutControllerImpl) DeleteByID(ctx context.Context, id primitive.ObjectID) error {
	workout, err := impl.GetByID(ctx, id)
	if err != nil {
		impl.Logger.Error("database get by id error", slog.Any("error", err))
		return err
	}
	if workout == nil {
		impl.Logger.Error("database returns nothing from get by id")
		return httperror.NewForBadRequestWithSingleField("id", "workout does not exist")
	}
	workout.Status = w_s.WorkoutStatusArchived
	workout.ModifiedAt = time.Now()

	if err := impl.WorkoutStorer.UpdateByID(ctx, workout); err != nil {
		impl.Logger.Error("database update by id error", slog.Any("error", err))
		return err
	}
	return nil
}
