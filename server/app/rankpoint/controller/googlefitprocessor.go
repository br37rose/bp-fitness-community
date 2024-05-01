package controller

import (
	"context"
	"fmt"
	"log/slog"
	"sort"
	"time"

	"github.com/bartmika/timekit"
	"go.mongodb.org/mongo-driver/bson/primitive"

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
	default:
		err := fmt.Errorf("period does not exist for value: %v", period)
		return err
	}

	impl.Logger.Debug("processing rankings",
		slog.Any("period_int", period),
		slog.Time("start", start),
		slog.Time("end", end),
	)

	for _, gfa := range gfas {
		// // Lock this google fit gfa for modification and unlock when we are
		// // finished with it.
		impl.Kmutex.Lockf("gfa_%v", gfa.ID.Hex())
		defer impl.Kmutex.Unlockf("gfa_%v", gfa.ID.Hex())

		u, err := impl.UserStorer.GetByID(ctx, gfa.UserID)
		if err != nil {
			impl.Logger.Error("failed getting user",
				slog.Any("google_fit_app_id", gfa.ID),
				slog.Int("period", int(period)),
				slog.Any("error", err))
			return err
		}
		if u == nil {
			err := fmt.Errorf("user does not exist id: %v", gfa.UserID)
			impl.Logger.Error("",
				slog.Any("google_fit_app_id", gfa.ID),
				slog.Int("period", int(period)),
				slog.Any("error", err))
			return err
		}

		// Variable defines all the biometric sensors we want to process for
		// this aggregation function.
		metricIDs := []primitive.ObjectID{
			gfa.CaloriesBurnedMetricID,
			gfa.StepCountDeltaMetricID,
			gfa.DistanceDeltaMetricID,
			gfa.HeartRateBPMMetricID,
			//TODO: Add more health sensors here...
		}

		impl.Logger.Debug("aggregation starting...",
			slog.String("gfa_id", gfa.ID.Hex()))

		for _, metricID := range metricIDs {
			agg, err := impl.AggregatePointStorer.GetByCompositeKey(ctx, metricID, period, start, end)
			if err != nil {
				impl.Logger.Error("aggregate point returned for composite key",
					slog.Any("google_fit_app_id", gfa.ID),
					slog.Any("metric_id", metricID),
					slog.Int("period", int(period)),
					slog.Any("error", err))
				return err
			}
			if agg != nil {

				//
				// Average
				//
				rpAvg, err := impl.RankPointStorer.GetByCompositeKey(ctx, metricID, rp_s.FunctionAverage, period, start, end)
				if err != nil {
					impl.Logger.Error("rank point returned for composite key",
						slog.Any("google_fit_app_id", gfa.ID),
						slog.Any("metric_id", metricID),
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
						MetricID:               agg.MetricID,
						MetricDataTypeName:     agg.MetricDataTypeName,
						Period:                 period,
						Start:                  start,
						End:                    end,
						Function:               rp_s.FunctionAverage,
						OrganizationID:         gfa.OrganizationID,
						OrganizationName:       gfa.OrganizationName,
						Value:                  agg.Average,
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
					rpAvg.Value = agg.Average
					if err := impl.RankPointStorer.UpdateByID(ctx, rpAvg); err != nil {
						impl.Logger.Error("failed updating rank point",
							slog.Any("error", err),
							slog.Int("period", int(period)),
							slog.Any("start", start),
							slog.Any("end", end))
						return err
					}
				}

				rpsAvg[agg.MetricDataTypeName] = append(rpsAvg[agg.MetricDataTypeName], rpAvg)

				//
				// Sum
				//

				rpSum, err := impl.RankPointStorer.GetByCompositeKey(ctx, metricID, rp_s.FunctionSum, period, start, end)
				if err != nil {
					impl.Logger.Error("rank point returned for composite key",
						slog.Any("google_fit_app_id", gfa.ID),
						slog.Any("metric_id", metricID),
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
						MetricID:               agg.MetricID,
						MetricDataTypeName:     agg.MetricDataTypeName,
						Period:                 period,
						Start:                  start,
						End:                    end,
						Function:               rp_s.FunctionSum,
						OrganizationID:         gfa.OrganizationID,
						OrganizationName:       gfa.OrganizationName,
						Value:                  agg.Sum,
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
					rpSum.Value = agg.Sum
					if err := impl.RankPointStorer.UpdateByID(ctx, rpSum); err != nil {
						impl.Logger.Error("failed updating rank point",
							slog.Any("error", err),
							slog.Int("period", int(period)),
							slog.Any("start", start),
							slog.Any("end", end))
						return err
					}
				}

				rpsSum[agg.MetricDataTypeName] = append(rpsSum[agg.MetricDataTypeName], rpSum)

				//
				// Count
				//

				rpCount, err := impl.RankPointStorer.GetByCompositeKey(ctx, metricID, rp_s.FunctionCount, period, start, end)
				if err != nil {
					impl.Logger.Error("rank point returned for composite key",
						slog.Any("google_fit_app_id", gfa.ID),
						slog.Any("metric_id", metricID),
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
						MetricID:               agg.MetricID,
						MetricDataTypeName:     agg.MetricDataTypeName,
						Period:                 period,
						Start:                  start,
						End:                    end,
						Function:               rp_s.FunctionCount,
						OrganizationID:         gfa.OrganizationID,
						OrganizationName:       gfa.OrganizationName,
						Value:                  agg.Count,
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
					rpCount.Value = agg.Count
					if err := impl.RankPointStorer.UpdateByID(ctx, rpCount); err != nil {
						impl.Logger.Error("failed updating rank point",
							slog.Any("error", err),
							slog.Int("period", int(period)),
							slog.Any("start", start),
							slog.Any("end", end))
						return err
					}
				}

				rpsCount[agg.MetricDataTypeName] = append(rpsCount[agg.MetricDataTypeName], rpCount)
			}
		}
	}

	for metricDataTypeName, rps := range rpsAvg {
		if err := impl.sortForRankPoints(ctx, metricDataTypeName, rps, period); err != nil {
			return err
		}
	}
	for metricDataTypeName, rps := range rpsSum {
		if err := impl.sortForRankPoints(ctx, metricDataTypeName, rps, period); err != nil {
			return err
		}
	}
	for metricDataTypeName, rps := range rpsCount {
		if err := impl.sortForRankPoints(ctx, metricDataTypeName, rps, period); err != nil {
			return err
		}
	}

	return nil
}

func (impl *RankPointControllerImpl) sortForRankPoints(ctx context.Context, dtn string, rps []*rp_s.RankPoint, period int8) error {
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
	}

	return nil
}
