package controller

import (
	"context"
	"errors"
	"log/slog"
	"time"

	c_s "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/rankpoint/datastore"
)

func (c *RankPointControllerImpl) ListByFilter(ctx context.Context, f *c_s.RankPointPaginationListFilter) (*c_s.RankPointPaginationListResult, error) {
	c.Logger.Debug("listing using filter options:",
		slog.Any("Cursor", f.Cursor),
		slog.Int64("PageSize", f.PageSize),
		slog.String("SortField", f.SortField),
		slog.Int("SortOrder", int(f.SortOrder)),
		slog.Any("MetricIDs", f.MetricIDs),
		slog.Any("MetricDataTypeNames", f.MetricDataTypeNames),
		slog.Any("Function", f.Function),
		slog.Any("Period", f.Period),
		slog.Time("CreatedAtGTE", f.CreatedAtGTE))

	list, err := c.RankPointStorer.ListByFilter(ctx, f)
	if err != nil {
		c.Logger.Error("database list by filter error", slog.Any("error", err))
		return nil, err
	}
	if list == nil {
		err := errors.New("list does not exist error")
		c.Logger.Error("database list error", slog.Any("error", err))
		return nil, err
	}

	// Iterate through all the rankings and if a s3 url expires then refresh.
	for _, rp := range list.Results {
		if rp.UserAvatarObjectKey != "" && time.Now().After(rp.UserAvatarObjectExpiry) {
			c.Kmutex.Lockf("rankpoint_%v", rp.ID.Hex())         // Step 1
			defer c.Kmutex.Unlockf("rankpoint_%v", rp.ID.Hex()) // Step 2

			u, err := c.UserStorer.GetByID(ctx, rp.UserID)
			if err != nil {
				c.Logger.Error("failed getting user", slog.Any("error", err))
				return nil, err
			}

			// Generate a presigned URL for today lasting 12 hours before expiring to becoming invalid.
			expiryDur := time.Hour * 12
			userAvatarObjectURL, presignErr := c.S3.GetPresignedURL(ctx, u.AvatarObjectKey, expiryDur)
			if presignErr != nil {
				c.Logger.Error("video s3 presign url error", slog.Any("presignErr", presignErr))
				return nil, err
			}
			rp.UserAvatarObjectURL = userAvatarObjectURL
			rp.UserAvatarObjectExpiry = time.Now().Add(expiryDur)
			if err := c.RankPointStorer.UpdateByID(ctx, rp); err != nil {
				c.Logger.Error("rank point failed updating", slog.Any("error", err))
				return nil, err
			}

			u.AvatarObjectURL = userAvatarObjectURL
			u.AvatarObjectExpiry = time.Now().Add(expiryDur)
			if err := c.UserStorer.UpdateByID(ctx, u); err != nil {
				c.Logger.Error("failed updating user", slog.Any("error", err))
				return nil, err
			}
		}
	}

	//
	return list, err
}
