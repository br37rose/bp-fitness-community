package controller

import (
	"context"
	"log/slog"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	ap_s "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/aggregatepoint/datastore"
)

func (impl *AggregatePointControllerImpl) aggregateForMetric(
	ctx context.Context,
	metricID primitive.ObjectID,
	metricDataTypeName string,
	period int8,
	startAt time.Time,
	endAt time.Time,
) error {
	response, err := impl.DataPointStorer.Aggregate(ctx, metricID, startAt, endAt)
	if err != nil {
		impl.Logger.Error("database list by filter error", slog.Any("error", err))
		return err
	}

	ap, err := impl.AggregatePointStorer.GetByCompositeKey(ctx, metricID, period, startAt, endAt)
	if err != nil {
		impl.Logger.Error("failed getting by composite key", slog.Any("error", err))
		return err
	}

	if ap == nil {
		// CASE 1 OF 2: Create
		ap = &ap_s.AggregatePoint{
			ID:         primitive.NewObjectID(),
			MetricID:   metricID,
			MetricDataTypeName: metricDataTypeName,
			Period:     period,
			Start:      startAt,
			End:        endAt,
			Count:      response.Count,
			Average:    response.Average,
			Min:        response.Min,
			Max:        response.Max,
			Sum:        response.Sum,
		}
		if err := impl.AggregatePointStorer.Create(ctx, ap); err != nil {
			impl.Logger.Error("failed creating",
				slog.Any("error", err))
			return err
		}
		// For debugging purposes only.
		impl.Logger.Debug("created aggregate point",
			slog.String("metric_id", metricID.Hex()),
			slog.Int("period", int(period)),
			slog.Time("start", startAt),
			slog.Time("end", endAt),
			slog.Any("count", response.Count),
			slog.Any("avg", response.Average),
			slog.Any("min", response.Min),
			slog.Any("max", response.Max),
			slog.Any("sum", response.Sum))
	} else {
		// CASE 2 OF 2: Update
		ap.MetricID = metricID
		ap.Period = period
		ap.Start = startAt
		ap.End = endAt
		ap.Count = response.Count
		ap.Average = response.Average
		ap.Min = response.Min
		ap.Max = response.Max
		ap.Sum = response.Sum
		if err := impl.AggregatePointStorer.UpdateByID(ctx, ap); err != nil {
			impl.Logger.Error("failed creating",
				slog.Any("error", err))
			return err
		}
		// // For debugging purposes only.
		// impl.Logger.Debug("updated aggregate point",
		// 	slog.String("metric_id", metricID.Hex()),
		// 	slog.Int("period", int(period)),
		// 	slog.Time("start", startAt),
		// 	slog.Time("end", endAt),
		// 	slog.Any("count", response.Count),
		// 	slog.Any("avg", response.Average),
		// 	slog.Any("min", response.Min),
		// 	slog.Any("max", response.Max),
		// 	slog.Any("sum", response.Sum))
	}
	return nil
}
