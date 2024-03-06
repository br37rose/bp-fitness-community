package controller

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"log/slog"

	domain "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/exercise/datastore"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/utils/httperror"
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

	// Refresh the presigned URL if it has expired or else skip this step.
	if m.VideoType == domain.ExerciseVideoTypeS3 && m.VideoObjectKey != "" {
		if time.Now().After(m.VideoObjectExpiry) {
			c.Kmutex.Lockf("exercise-%v", id.Hex()) // Step 1
			defer func() {
				c.Kmutex.Unlockf("exercise-%v", id.Hex()) // Step 2
			}()

			// Generate a presigned URL for today.
			expiryDur := time.Hour * 12
			videoObjectURL, presignErr := c.S3.GetPresignedURL(ctx, m.VideoObjectKey, expiryDur)
			if presignErr != nil {
				c.Logger.Error("video s3 presign url error", slog.Any("presignErr", presignErr))
				return nil, err
			}

			// Update the exercise.
			m.VideoObjectURL = videoObjectURL
			m.VideoObjectExpiry = time.Now().Add(expiryDur)
			if err := c.ExerciseStorer.UpdateByID(ctx, m); err != nil {
				c.Logger.Error("exercise database update by id error", slog.Any("error", err))
				return nil, err
			}
		}
	}

	// Refresh the presigned URL if it has expired or else skip this step.
	if m.ThumbnailType == domain.ExerciseThumbnailTypeS3 && m.ThumbnailObjectKey != "" {
		if time.Now().After(m.ThumbnailObjectExpiry) {
			c.Kmutex.Lockf("exercise-%v", id.Hex()) // Step 1
			defer func() {
				c.Kmutex.Unlockf("exercise-%v", id.Hex()) // Step 2
			}()

			// Generate a presigned URL for today.
			expiryDur := time.Hour * 12
			thumbnailObjectURL, presignErr := c.S3.GetPresignedURL(ctx, m.ThumbnailObjectKey, expiryDur)
			if presignErr != nil {
				c.Logger.Error("thumbnail s3 presign url error", slog.Any("presignErr", presignErr))
				return nil, err
			}

			// Update the exercise.
			m.ThumbnailObjectURL = thumbnailObjectURL
			m.ThumbnailObjectExpiry = time.Now().Add(expiryDur)
			if err := c.ExerciseStorer.UpdateByID(ctx, m); err != nil {
				c.Logger.Error("exercise database update by id error", slog.Any("error", err))
				return nil, err
			}
		}
	}

	return m, err
}
