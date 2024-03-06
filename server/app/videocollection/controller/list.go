package controller

import (
	"context"
	"time"

	vcol_d "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/videocollection/datastore"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/config/constants"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log/slog"
)

func (c *VideoCollectionControllerImpl) ListByFilter(ctx context.Context, f *vcol_d.VideoCollectionListFilter) (*vcol_d.VideoCollectionListResult, error) {
	// // Extract from our session the following data.
	// userID := ctx.Value(constants.SessionUserID).(primitive.ObjectID)
	// userRole := ctx.Value(constants.SessionUserRole).(int8)
	orgID, _ := ctx.Value(constants.SessionUserOrganizationID).(primitive.ObjectID)

	// Apply filtering organization tenancy.
	f.OrganizationID = orgID

	listRes, err := c.VideoCollectionStorer.ListByFilter(ctx, f)
	if err != nil {
		c.Logger.Error("database list by filter error", slog.Any("error", err))
		return nil, err
	}

	// Iterate through all the videocollections and refresh the presigned URL if it has expired.
	for _, e := range listRes.Results {

		// Refresh the presigned URL if it has expired or else skip this step.
		if e.ThumbnailType == vcol_d.VideoCollectionThumbnailTypeS3 && e.ThumbnailObjectKey != "" {
			if time.Now().After(e.ThumbnailObjectExpiry) {
				c.Kmutex.Lockf("videocollection-%v", e.ID.Hex()) // Step 1
				defer func() {
					c.Kmutex.Unlockf("videocollection-%v", e.ID.Hex()) // Step 2
				}()

				// Generate a presigned URL for today.
				expiryDur := time.Hour * 12
				thumbnailObjectURL, presignErr := c.S3.GetPresignedURL(ctx, e.ThumbnailObjectKey, expiryDur)
				if presignErr != nil {
					c.Logger.Error("thumbnail s3 presign url error", slog.Any("presignErr", presignErr))
					return nil, err
				}

				// Update the videocollection.
				e.ThumbnailObjectURL = thumbnailObjectURL
				e.ThumbnailObjectExpiry = time.Now().Add(expiryDur)
				if err := c.VideoCollectionStorer.UpdateByID(ctx, e); err != nil {
					c.Logger.Error("videocollection database update by id error", slog.Any("error", err))
					return nil, err
				}
			}
		}
	}
	return listRes, err
}
