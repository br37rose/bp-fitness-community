package controller

import (
	"context"
	"fmt"
	"log/slog"
	"sort"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	gfa_ds "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/googlefitapp/datastore"
	rp_s "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/rankpoint/datastore"
)

func (impl *RankPointControllerImpl) processGlobalRanksForGoogleFitApps(ctx context.Context, gfas []*gfa_ds.GoogleFitApp, metricType int8, function int8, period int8, start time.Time, end time.Time) error {
	// Variable will be used to get all the rankings we have for the particular
	// metric we are processing.
	rankPoints := make([]*rp_s.RankPoint, 0)

	for _, gfa := range gfas {
		// // Lock this google fit gfa for modification and unlock when we are
		// // finished with it.
		impl.Kmutex.Lockf("gfa_%v", gfa.ID.Hex())
		defer impl.Kmutex.Unlockf("gfa_%v", gfa.ID.Hex())

		// Pick the metric ID based on metric type selected.
		var metricID primitive.ObjectID
		switch metricType {
		case rp_s.MetricTypeHeartRate:
			metricID = gfa.HeartRateBPMMetricID
			// case rp_s.MetricTypeActivitySteps:
			// 	metricID = gfa.StepsCountMetricID
		default:
			err := fmt.Errorf("does not exist for metric type: %v", metricType)
			impl.Logger.Error("",
				slog.Any("google_fit_app_id", gfa.ID),
				slog.Any("metric_id", metricID),
				slog.Int("period", int(period)),
				slog.Any("error", err))
			return err
		}

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
			rp, err := impl.RankPointStorer.GetByCompositeKey(ctx, metricID, function, period, start, end)
			if err != nil {
				impl.Logger.Error("rank point returned for composite key",
					slog.Any("google_fit_app_id", gfa.ID),
					slog.Any("metric_id", metricID),
					slog.Any("function", function),
					slog.Int("period", int(period)),
					slog.Any("error", err))
				return err
			}

			if rp == nil {

				////
				//// CASE 1 OF 2: No rank record, therefore create rank record.
				////

				u, err := impl.UserStorer.GetByID(ctx, gfa.UserID)
				if err != nil {
					impl.Logger.Error("failed getting user",
						slog.Any("error", err))
					return err
				}

				rp = &rp_s.RankPoint{
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
					MetricID:               metricID,
					MetricType:             metricType,
					Period:                 period,
					Start:                  start,
					End:                    end,
					Function:               function,
					OrganizationID:         gfa.OrganizationID,
					OrganizationName:       gfa.OrganizationName,
				}

				switch function {
				case rp_s.FunctionAverage:
					rp.Value = agg.Average
				case rp_s.FunctionSum:
					rp.Value = agg.Sum
				}

				if err := impl.RankPointStorer.Create(ctx, rp); err != nil {
					impl.Logger.Error("failed creating rank point",
						slog.Any("error", err),
						slog.Any("function", function),
						slog.Int("period", int(period)),
						slog.Any("start", start),
						slog.Any("end", end))
					return err
				}

				rankPoints = append(rankPoints, rp)

				impl.Logger.Debug("created rank point",
					slog.Any("rp_id", rp.ID),
					slog.Any("function", function),
					slog.Int("period", int(period)),
					slog.Any("start", start),
					slog.Any("end", end))
			} else {

				////
				//// CASE 2 OF 2: Update existing rank record.
				////

				switch function {
				case rp_s.FunctionAverage:
					rp.Value = agg.Average
					rp.MetricType = metricType
				case rp_s.FunctionSum:
					rp.Value = agg.Sum
					rp.MetricType = metricType
				}

				if err := impl.RankPointStorer.UpdateByID(ctx, rp); err != nil {
					impl.Logger.Error("failed updating rank point",
						slog.Any("error", err))
					return err
				}

				rankPoints = append(rankPoints, rp)

				// For debugging purposes only.
				impl.Logger.Debug("updated rank point",
					slog.Any("rp_id", rp.ID),
					slog.Any("function", function),
					slog.Int("period", int(period)),
					slog.Any("start", start),
					slog.Any("end", end))
			}

		}
	}

	////
	//// Sorting the array by Place in descending order
	////

	// for _, rp := range rankPoints {
	// 	// For debugging purposes only.
	// 	impl.Logger.Error("pre-sorting rankpoint entry",
	// 		slog.Any("id", rp.ID),
	// 		slog.Any("value", rp.Value),
	// 		slog.Any("plance", rp.Place))
	// }

	// Sort all the values from greatest value to lowest value.
	sort.Slice(rankPoints, func(i, j int) bool {
		return rankPoints[i].Value > rankPoints[j].Value
	})

	// for _, rp := range rankPoints {
	// 	// For debugging purposes only.
	// 	impl.Logger.Error("post-sorting rankpoint entry",
	// 		slog.Any("id", rp.ID),
	// 		slog.Any("value", rp.Value),
	// 		slog.Any("plance", rp.Place))
	// }

	////
	//// Iterate through the sorted array and attach the place in our rank order.
	////

	// Because we are starting from largest `value` to lowest `value` then
	// we can start the rank `place` from the value `1` and increase going
	// through the list.
	for i, rp := range rankPoints {
		rp.Place = uint64(i + 1)
		if err := impl.RankPointStorer.UpdateByID(ctx, rp); err != nil {
			impl.Logger.Error("failed updating rank point",
				slog.Any("error", err))
			return err
		}
	}

	// for _, rp := range rankPoints {
	// 	// For debugging purposes only.
	// 	impl.Logger.Error("exit-sorting rankpoint entry",
	// 		slog.Any("id", rp.ID),
	// 		slog.Any("value", rp.Value),
	// 		slog.Any("plance", rp.Place))
	// }

	// For debugging purposes only.
	impl.Logger.Debug("ranked all active fitibt apps",
		slog.Any("records_count", len(rankPoints)),
		slog.Any("function", function),
		slog.Int("period", int(period)),
		slog.Any("start", start),
		slog.Any("end", end))

	return nil
}
