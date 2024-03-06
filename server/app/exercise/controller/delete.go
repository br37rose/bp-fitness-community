package controller

import (
	"context"

	domain "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/exercise/datastore"
	s_d "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/exercise/datastore"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/utils/httperror"
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

	// Defensive code: Do not edit system.
	if exercise.Type == domain.ExerciseTypeSystem {
		impl.Logger.Error("system exercise delete error")
		return httperror.NewForBadRequestWithSingleField("message", "the `system` exercises are read-only and thus access is denied")
	}

	if exercise.VideoType == s_d.ExerciseVideoTypeS3 {
		keys := []string{exercise.VideoObjectKey}
		if err := impl.S3.DeleteByKeys(ctx, keys); err != nil {
			impl.Logger.Error("s3 delete error", slog.Any("error", err))
			// Do not return an error, simply continue this function as there might
			// be a case were the file was removed on the s3 bucket by ourselves
			// or some other reason.
		}
		impl.Logger.Error("s3 deleted video", slog.Any("keys", keys))
	}

	if exercise.ThumbnailType == s_d.ExerciseThumbnailTypeS3 {
		keys := []string{exercise.ThumbnailObjectKey}
		if err := impl.S3.DeleteByKeys(ctx, keys); err != nil {
			impl.Logger.Error("s3 delete error", slog.Any("error", err))
			// Do not return an error, simply continue this function as there might
			// be a case were the file was removed on the s3 bucket by ourselves
			// or some other reason.
		}
		impl.Logger.Error("s3 deleted thumbnail", slog.Any("keys", keys))
	}

	if err := impl.ExerciseStorer.DeleteByID(ctx, id); err != nil {
		impl.Logger.Error("database delete by id error", slog.Any("error", err))
		return err
	}
	return nil
}
