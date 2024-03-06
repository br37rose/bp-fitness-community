package controller

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"log/slog"

	vcol_d "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/videocollection/datastore"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/utils/httperror"
)

func (c *VideoCollectionControllerImpl) GetByID(ctx context.Context, id primitive.ObjectID) (*vcol_d.VideoCollection, error) {
	// Retrieve from our database the record for the specific id.
	m, err := c.VideoCollectionStorer.GetByID(ctx, id)
	if err != nil {
		c.Logger.Error("database get by id error", slog.Any("error", err))
		return nil, err
	}
	if m == nil {
		return nil, httperror.NewForBadRequestWithSingleField("id", "videocollection does not exist")
	}

	// Refresh the presigned URL if it has expired or else skip this step.
	if m.ThumbnailType == vcol_d.VideoCollectionThumbnailTypeS3 && m.ThumbnailObjectKey != "" {
		if time.Now().After(m.ThumbnailObjectExpiry) {
			c.Kmutex.Lockf("videocollection-%v", id.Hex()) // Step 1
			defer func() {
				c.Kmutex.Unlockf("videocollection-%v", id.Hex()) // Step 2
			}()

			// Generate a presigned URL for today.
			expiryDur := time.Hour * 12
			thumbnailObjectURL, presignErr := c.S3.GetPresignedURL(ctx, m.ThumbnailObjectKey, expiryDur)
			if presignErr != nil {
				c.Logger.Error("thumbnail s3 presign url error", slog.Any("presignErr", presignErr))
				return nil, err
			}

			// Update the videocollection.
			m.ThumbnailObjectURL = thumbnailObjectURL
			m.ThumbnailObjectExpiry = time.Now().Add(expiryDur)
			if err := c.VideoCollectionStorer.UpdateByID(ctx, m); err != nil {
				c.Logger.Error("videocollection database update by id error", slog.Any("error", err))
				return nil, err
			}
		}
	}

	return m, err
}
