package controller

import (
	"context"
	"fmt"
	"log/slog"
	"sort"
	"time"

	"github.com/bartmika/timekit"
	"go.mongodb.org/mongo-driver/bson/primitive"

	gcp_a "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/adapter/cloudprovider/google"
	ag_s "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/aggregatepoint/datastore"
	gfa_ds "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/googlefitapp/datastore"
	rp_s "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/rankpoint/datastore"
)

func (impl *RankPointControllerImpl) processGlobalRanksForGoogleFitAppsV2(ctx context.Context, gfas []*gfa_ds.GoogleFitApp, period int8) error {
	// // Variable will be used to get all the rankings we have for the particular
	// // metric we are processing.
	// rpsAvg := make([]*rp_s.RankPoint, 0)
	// rpsSum := make([]*rp_s.RankPoint, 0)
	// rpsCount := make([]*rp_s.RankPoint, 0)
	rpsAvg := make(map[string][]*rp_s.RankPoint, 0)
	rpsSum := make(map[string][]*rp_s.RankPoint, 0)
	rpsCount := make(map[string][]*rp_s.RankPoint, 0)

	var start time.Time
	var end time.Time
	switch period {
	case rp_s.PeriodDay:
		start = timekit.Midnight(time.Now)
		end = timekit.MidnightTomorrow(time.Now)
	case rp_s.PeriodWeek:
		start = timekit.FirstDayOfThisISOWeek(time.Now)
		end = timekit.FirstDayOfNextISOWeek(time.Now)
	case rp_s.PeriodMonth:
		start = timekit.FirstDayOfThisMonth(time.Now)
		end = timekit.FirstDayOfNextMonth(time.Now)
	case rp_s.PeriodYear:
		start = timekit.FirstDayOfThisYear(time.Now)
		end = timekit.FirstDayOfNextYear(time.Now)
	default:
		err := fmt.Errorf("period does not exist for value: %v", period)
		return err
	}

	metricDataTypeNames := []string{
		gcp_a.DataTypeNameCaloriesBurned,
		gcp_a.DataTypeNameStepCountDelta,
		gcp_a.DataTypeNameDistanceDelta,
		gcp_a.DataTypeNameHeartRateBPM,
	}

	impl.Logger.Debug("processing rankings per data name type",
		slog.Any("period", period),
		slog.Time("start", start),
		slog.Time("end", end),
	)

	for _, metricDataTypeName := range metricDataTypeNames {
		f := &ag_s.AggregatePointPaginationListFilter{
			PageSize:            1_0000_0000,
			SortField:           "created_at",
			SortOrder:           -1, // Descending
			MetricDataTypeNames: []string{metricDataTypeName},
			Period:              period,
			StartGTE:            start,
			EndLTE:              end,
		}
		agag, err := impl.AggregatePointStorer.ListByFilter(ctx, f)
		if err != nil {
			impl.Logger.Error("aggregate point list returned error",
				slog.Any("f", f),
				slog.Any("metric_data_type_name", metricDataTypeName),
				slog.Int("period", int(period)),
				slog.Any("error", err))
			return err
		}
		for _, ag := range agag.Results {
			//
			// User account.
			//

			u, err := impl.UserStorer.GetByID(ctx, ag.UserID)
			if err != nil {
				impl.Logger.Error("failed getting user",
					slog.Any("metric_data_type_name", metricDataTypeName),
					slog.Int("period", int(period)),
					slog.Any("error", err))
				return err
			}
			if u == nil {
				err := fmt.Errorf("user does not exist id: %v", ag.UserID)
				impl.Logger.Error("",
					slog.Any("metric_data_type_name", metricDataTypeName),
					slog.Int("period", int(period)),
					slog.Any("error", err))
				return err
			}

			//
			// Average
			//

			rpAvg, err := impl.RankPointStorer.GetByCompositeKey(ctx, ag.MetricID, rp_s.FunctionAverage, period, start, end)
			if err != nil {
				impl.Logger.Error("rank point returned for composite key",
					slog.Any("metric_data_type_name", metricDataTypeName),
					slog.Any("metric_id", ag.MetricID),
					slog.Any("function", rp_s.FunctionAverage),
					slog.Int("period", int(period)),
					slog.Any("error", err))
				return err
			}

			if rpAvg == nil {
				rpAvg = &rp_s.RankPoint{
					ID:                     primitive.NewObjectID(),
					UserID:                 u.ID,
					UserFirstName:          u.FirstName,
					UserLastName:           u.LastName,
					UserAvatarObjectExpiry: u.AvatarObjectExpiry,
					UserAvatarObjectURL:    u.AvatarObjectURL,
					UserAvatarObjectKey:    u.AvatarObjectKey,
					UserAvatarFileType:     u.AvatarFileType,
					UserAvatarFileName:     u.AvatarFileName,
					Place:                  1_000_000_000, // Last Place.
					MetricID:               ag.MetricID,
					MetricDataTypeName:     ag.MetricDataTypeName,
					Period:                 period,
					Start:                  start,
					End:                    end,
					Function:               rp_s.FunctionAverage,
					OrganizationID:         u.OrganizationID,
					OrganizationName:       u.OrganizationName,
					Value:                  ag.Average,
				}
				if err := impl.RankPointStorer.Create(ctx, rpAvg); err != nil {
					impl.Logger.Error("failed creating rank point",
						slog.Any("error", err),
						slog.Int("period", int(period)),
						slog.Any("start", start),
						slog.Any("end", end))
					return err
				}
			} else {
				rpAvg.Value = ag.Average
				if err := impl.RankPointStorer.UpdateByID(ctx, rpAvg); err != nil {
					impl.Logger.Error("failed updating rank point",
						slog.Any("error", err),
						slog.Int("period", int(period)),
						slog.Any("start", start),
						slog.Any("end", end))
					return err
				}
			}

			rpsAvg[ag.MetricDataTypeName] = append(rpsAvg[ag.MetricDataTypeName], rpAvg)

			//
			// Sum
			//

			rpSum, err := impl.RankPointStorer.GetByCompositeKey(ctx, ag.MetricID, rp_s.FunctionSum, period, start, end)
			if err != nil {
				impl.Logger.Error("rank point returned for composite key",
					slog.Any("metric_data_type_name", metricDataTypeName),
					slog.Any("metric_id", ag.MetricID),
					slog.Any("function", rp_s.FunctionSum),
					slog.Int("period", int(period)),
					slog.Any("error", err))
				return err
			}

			if rpSum == nil {
				rpSum = &rp_s.RankPoint{
					ID:                     primitive.NewObjectID(),
					UserID:                 u.ID,
					UserFirstName:          u.FirstName,
					UserLastName:           u.LastName,
					UserAvatarObjectExpiry: u.AvatarObjectExpiry,
					UserAvatarObjectURL:    u.AvatarObjectURL,
					UserAvatarObjectKey:    u.AvatarObjectKey,
					UserAvatarFileType:     u.AvatarFileType,
					UserAvatarFileName:     u.AvatarFileName,
					Place:                  1_000_000_000, // Last Place.
					MetricID:               ag.MetricID,
					MetricDataTypeName:     ag.MetricDataTypeName,
					Period:                 period,
					Start:                  start,
					End:                    end,
					Function:               rp_s.FunctionSum,
					OrganizationID:         u.OrganizationID,
					OrganizationName:       u.OrganizationName,
					Value:                  ag.Sum,
				}
				if err := impl.RankPointStorer.Create(ctx, rpSum); err != nil {
					impl.Logger.Error("failed creating rank point",
						slog.Any("error", err),
						slog.Int("period", int(period)),
						slog.Any("start", start),
						slog.Any("end", end))
					return err
				}
			} else {
				rpSum.Value = ag.Sum
				if err := impl.RankPointStorer.UpdateByID(ctx, rpSum); err != nil {
					impl.Logger.Error("failed updating rank point",
						slog.Any("error", err),
						slog.Int("period", int(period)),
						slog.Any("start", start),
						slog.Any("end", end))
					return err
				}
			}

			rpsSum[ag.MetricDataTypeName] = append(rpsSum[ag.MetricDataTypeName], rpSum)

			//
			// Count
			//

			rpCount, err := impl.RankPointStorer.GetByCompositeKey(ctx, ag.MetricID, rp_s.FunctionCount, period, start, end)
			if err != nil {
				impl.Logger.Error("rank point returned for composite key",
					slog.Any("metric_data_type_name", metricDataTypeName),
					slog.Any("metric_id", ag.MetricID),
					slog.Any("function", rp_s.FunctionCount),
					slog.Int("period", int(period)),
					slog.Any("error", err))
				return err
			}

			if rpCount == nil {
				rpCount = &rp_s.RankPoint{
					ID:                     primitive.NewObjectID(),
					UserID:                 u.ID,
					UserFirstName:          u.FirstName,
					UserLastName:           u.LastName,
					UserAvatarObjectExpiry: u.AvatarObjectExpiry,
					UserAvatarObjectURL:    u.AvatarObjectURL,
					UserAvatarObjectKey:    u.AvatarObjectKey,
					UserAvatarFileType:     u.AvatarFileType,
					UserAvatarFileName:     u.AvatarFileName,
					Place:                  1_000_000_000, // Last Place.
					MetricID:               ag.MetricID,
					MetricDataTypeName:     ag.MetricDataTypeName,
					Period:                 period,
					Start:                  start,
					End:                    end,
					Function:               rp_s.FunctionCount,
					OrganizationID:         u.OrganizationID,
					OrganizationName:       u.OrganizationName,
					Value:                  ag.Count,
				}
				if err := impl.RankPointStorer.Create(ctx, rpCount); err != nil {
					impl.Logger.Error("failed creating rank point",
						slog.Any("error", err),
						slog.Int("period", int(period)),
						slog.Any("start", start),
						slog.Any("end", end))
					return err
				}
			} else {
				rpCount.Value = ag.Count
				if err := impl.RankPointStorer.UpdateByID(ctx, rpCount); err != nil {
					impl.Logger.Error("failed updating rank point",
						slog.Any("error", err),
						slog.Int("period", int(period)),
						slog.Any("start", start),
						slog.Any("end", end))
					return err
				}
			}

			rpsCount[ag.MetricDataTypeName] = append(rpsCount[ag.MetricDataTypeName], rpCount)

		}
	}

	for metricDataTypeName, rps := range rpsAvg {
		impl.Logger.Debug("ranking start",
			slog.String("metric_data_type_name", metricDataTypeName),
			slog.Int("function", int(rp_s.FunctionAverage)),
			slog.Int("period", int(period)),
			slog.Int("rps_count", len(rps)),
		)

		if err := impl.sortForRankPoints(ctx, metricDataTypeName, rps, rp_s.FunctionAverage, period); err != nil {
			return err
		}

		impl.Logger.Debug("ranking done",
			slog.String("metric_data_type_name", metricDataTypeName),
			slog.Int("period", int(period)),
		)
	}

	for metricDataTypeName, rps := range rpsSum {
		impl.Logger.Debug("ranking start",
			slog.String("metric_data_type_name", metricDataTypeName),
			slog.Int("function", int(rp_s.FunctionSum)),
			slog.Int("period", int(period)),
			slog.Int("rps_count", len(rps)),
		)

		if err := impl.sortForRankPoints(ctx, metricDataTypeName, rps, rp_s.FunctionSum, period); err != nil {
			return err
		}

		impl.Logger.Debug("ranking done",
			slog.String("metric_data_type_name", metricDataTypeName),
			slog.Int("period", int(period)),
		)
	}

	for metricDataTypeName, rps := range rpsCount {
		impl.Logger.Debug("ranking start",
			slog.String("metric_data_type_name", metricDataTypeName),
			slog.Int("function", int(rp_s.FunctionCount)),
			slog.Int("period", int(period)),
			slog.Int("rps_count", len(rps)),
		)
		if err := impl.sortForRankPoints(ctx, metricDataTypeName, rps, rp_s.FunctionCount, period); err != nil {
			return err
		}

		impl.Logger.Debug("ranking done",
			slog.String("metric_data_type_name", metricDataTypeName),
			slog.Int("period", int(period)),
		)
	}

	return nil
}

func (impl *RankPointControllerImpl) sortForRankPoints(ctx context.Context, dtn string, rps []*rp_s.RankPoint, fun int8, period int8) error {
	////
	//// Sorting the array by Place in descending order
	////

	// Sort all the values from greatest value to lowest value.
	sort.Slice(rps, func(i, j int) bool {
		return rps[i].Value > rps[j].Value
	})

	////
	//// Iterate through the sorted array and attach the place in our rank order.
	////

	// Because we are starting from largest `value` to lowest `value` then
	// we can start the rank `place` from the value `1` and increase going
	// through the list.
	for i, rp := range rps {
		rp.Place = uint64(i + 1)
		if err := impl.RankPointStorer.UpdateByID(ctx, rp); err != nil {
			impl.Logger.Error("failed updating rank point",
				slog.Int("period", int(period)),
				slog.Any("error", err))
			return err
		}
		impl.Logger.Debug("ranked",
			slog.String("user_first_name", rp.UserFirstName),
			slog.String("metric_data_type_name", rp.MetricDataTypeName),
			slog.Time("start", rp.Start),
			slog.Time("end", rp.End),
			slog.Int("function", int(rp.Function)),
			slog.Int("period", int(rp.Period)),
			slog.Int("function", int(rp.Function)),
			slog.Int("place", int(rp.Place)),
			slog.Float64("value", rp.Value),
		)
	}

	return nil
}
