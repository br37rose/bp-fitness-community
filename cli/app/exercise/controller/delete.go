package controller

import (
	"context"

	"log/slog"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (impl *ExerciseControllerImpl) DeleteByID(ctx context.Context, id primitive.ObjectID) error {
	exercise, err := impl.GetByID(ctx, id)
	if err != nil {
		impl.Logger.Error("database get by id error", slog.Any("error", err))
		return err
	}
	if exercise == nil {
		impl.Logger.Error("database returns nothing from get by id")
		return err
	}
	if err := impl.ExerciseStorer.DeleteByID(ctx, id); err != nil {
		impl.Logger.Error("database delete by id error", slog.Any("error", err))
		return err
	}
	return nil
}
