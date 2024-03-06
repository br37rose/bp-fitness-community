package controller

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"log/slog"

	vcon_s "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/videocontent/datastore"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/utils/httperror"
)

func (c *VideoContentControllerImpl) GetByID(ctx context.Context, id primitive.ObjectID) (*vcon_s.VideoContent, error) {
	// Retrieve from our database the record for the specific id.
	m, err := c.VideoContentStorer.GetByID(ctx, id)
	if err != nil {
		c.Logger.Error("database get by id error", slog.Any("error", err))
		return nil, err
	}
	if m == nil {
		return nil, httperror.NewForBadRequestWithSingleField("id", "videocontent does not exist")
	}

	// Refresh the presigned URL if it has expired or else skip this step.
	if m.VideoType == vcon_s.VideoContentVideoTypeS3 && m.VideoObjectKey != "" {
		if time.Now().After(m.VideoObjectExpiry) {
			c.Kmutex.Lockf("videocontent-%v", id.Hex()) // Step 1
			defer func() {
				c.Kmutex.Unlockf("videocontent-%v", id.Hex()) // Step 2
			}()

			// Generate a presigned URL for today.
			expiryDur := time.Hour * 12
			videoObjectURL, presignErr := c.S3.GetPresignedURL(ctx, m.VideoObjectKey, expiryDur)
			if presignErr != nil {
				c.Logger.Error("video s3 presign url error", slog.Any("presignErr", presignErr))
				return nil, err
			}

			// Update the videocontent.
			m.VideoObjectURL = videoObjectURL
			m.VideoObjectExpiry = time.Now().Add(expiryDur)
			if err := c.VideoContentStorer.UpdateByID(ctx, m); err != nil {
				c.Logger.Error("videocontent database update by id error", slog.Any("error", err))
				return nil, err
			}
		}
	}

	// Refresh the presigned URL if it has expired or else skip this step.
	if m.ThumbnailType == vcon_s.VideoContentThumbnailTypeS3 && m.ThumbnailObjectKey != "" {
		if time.Now().After(m.ThumbnailObjectExpiry) {
			c.Kmutex.Lockf("videocontent-%v", id.Hex()) // Step 1
			defer func() {
				c.Kmutex.Unlockf("videocontent-%v", id.Hex()) // Step 2
			}()

			// Generate a presigned URL for today.
			expiryDur := time.Hour * 12
			thumbnailObjectURL, presignErr := c.S3.GetPresignedURL(ctx, m.ThumbnailObjectKey, expiryDur)
			if presignErr != nil {
				c.Logger.Error("thumbnail s3 presign url error", slog.Any("presignErr", presignErr))
				return nil, err
			}

			// Update the videocontent.
			m.ThumbnailObjectURL = thumbnailObjectURL
			m.ThumbnailObjectExpiry = time.Now().Add(expiryDur)
			if err := c.VideoContentStorer.UpdateByID(ctx, m); err != nil {
				c.Logger.Error("videocontent database update by id error", slog.Any("error", err))
				return nil, err
			}
		}
	}

	return m, err
}
