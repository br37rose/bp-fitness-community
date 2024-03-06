package controller

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/bartmika/timekit"
	rp_s "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/rankpoint/datastore"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/utils/httperror"
)

type LeaderboardRequest struct {
	// Pagination related.
	Cursor   string
	PageSize int64

	// Filter related.
	MetricType int8
	Function   int8
	Period     int8
}

func validateLeaderboardRequest(dirtyData *LeaderboardRequest) error {
	e := make(map[string]string)
	if dirtyData.PageSize == 0 {
		e["page_size"] = "missing value"
	}
	if dirtyData.MetricType == 0 {
		e["metric_type"] = "missing value"
	}
	if dirtyData.Function == 0 {
		e["function"] = "missing value"
	}
	if dirtyData.Period == 0 {
		e["period"] = "missing value"
	}
	if dirtyData.Function == rp_s.FunctionSum && dirtyData.MetricType == rp_s.MetricTypeHeartRate {
		e["function"] = "sum cannot be used on heart rate"
	}
	if len(e) != 0 {
		return httperror.NewForBadRequest(&e)
	}
	return nil
}

func (c *BiometricControllerImpl) Leaderboard(ctx context.Context, req *LeaderboardRequest) (*rp_s.RankPointPaginationListResult, error) {
	// Defensive code: Enforce input restrictions.
	if err := validateLeaderboardRequest(req); err != nil {
		return nil, err
	}

	c.Logger.Debug("leaderboard request",
		slog.Any("cursor", req.Cursor),
		slog.Int64("page_size", req.PageSize),
		slog.Any("metric_type", req.MetricType),
		slog.Any("function", req.Function),
		slog.Any("period", req.Period))

	// Create our custom filter in which we order from top ranked user data
	// to the lowest ranked user data for the specific metric.
	f := &rp_s.RankPointPaginationListFilter{
		Cursor:      req.Cursor,
		PageSize:    req.PageSize,
		SortField:   "place",
		SortOrder:   rp_s.OrderAscending,
		MetricTypes: []int8{req.MetricType},
		Function:    req.Function,
		Period:      req.Period,
	}

	switch req.Period {
	case rp_s.PeriodDay:
		f.StartGTE = timekit.Midnight(time.Now)
		f.EndLTE = timekit.MidnightTomorrow(time.Now)
		break
	case rp_s.PeriodWeek:
		f.StartGTE = timekit.FirstDayOfThisISOWeek(time.Now)
		f.EndLTE = timekit.FirstDayOfNextISOWeek(time.Now)
		break
	case rp_s.PeriodMonth:
		f.StartGTE = timekit.FirstDayOfThisMonth(time.Now)
		f.EndLTE = timekit.FirstDayOfNextMonth(time.Now)
		break
	case rp_s.PeriodYear:
		f.StartGTE = timekit.FirstDayOfThisYear(time.Now)
		f.EndLTE = timekit.FirstDayOfNextYear(time.Now)
		break
	default:
		break
	}

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
		// For debugging purposes only.
		c.Logger.Debug("rankpoint entry",
			slog.Any("id", rp.ID),
			slog.Any("value", rp.Value),
			slog.Any("plance", rp.Place))

		if rp.UserAvatarObjectKey != "" && time.Now().After(rp.UserAvatarObjectExpiry) {
			c.Kmutex.Lockf("rankpoint_%v_avatar_image", rp.ID.Hex())         // Step 1
			defer c.Kmutex.Unlockf("rankpoint_%v_avatar_image", rp.ID.Hex()) // Step 2

			// Generate a presigned URL for today lasting 12 hours before expiring to becoming invalid.
			expiryDur := time.Hour * 12
			userAvatarObjectURL, presignErr := c.S3.GetPresignedURL(ctx, rp.UserAvatarObjectKey, expiryDur)
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

			if err := c.RankPointStorer.UpdateByID(ctx, rp); err != nil {
				c.Logger.Error("failed updating rankpoint", slog.Any("error", err))
				return nil, err
			}
		}
	}

	//
	return list, err
}
