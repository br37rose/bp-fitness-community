package datastore

import (
	"context"
	"log/slog"
	"time"

	"github.com/bartmika/timekit"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func (impl RankPointStorerImpl) GetByID(ctx context.Context, id primitive.ObjectID) (*RankPoint, error) {
	filter := bson.M{"_id": id}

	var result RankPoint
	err := impl.Collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// This error means your query did not match any documents.
			return nil, nil
		}
		impl.Logger.Error("database get by id error", slog.Any("error", err))
		return nil, err
	}
	return &result, nil
}

func (impl RankPointStorerImpl) GetByCompositeKey(ctx context.Context, metricID primitive.ObjectID, function int8, period int8, start time.Time, end time.Time) (*RankPoint, error) {
	filter := bson.M{
		"metric_id": metricID,
		"function":  function,
		"period":    period,
		"start":     start,
		"end":       end,
	}

	var result RankPoint
	err := impl.Collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// This error means your query did not match any documents.
			return nil, nil
		}
		impl.Logger.Error("database get by composite key error", slog.Any("error", err))
		return nil, err
	}
	return &result, nil
}

func (impl RankPointStorerImpl) GetByCompositeKeyForToday(ctx context.Context, metricID primitive.ObjectID, function int8) (*RankPoint, error) {
	thisDayStart := timekit.Midnight(time.Now)
	thisDayEnd := timekit.MidnightTomorrow(time.Now)
	return impl.GetByCompositeKey(ctx, metricID, function, PeriodDay, thisDayStart, thisDayEnd)
}

func (impl RankPointStorerImpl) GetByCompositeKeyForYesterday(ctx context.Context, metricID primitive.ObjectID, function int8) (*RankPoint, error) {
	lastDayStart := timekit.MidnightYesterday(time.Now)
	lastDayEnd := timekit.Midnight(time.Now)
	return impl.GetByCompositeKey(ctx, metricID, function, PeriodDay, lastDayStart, lastDayEnd)
}

func (impl RankPointStorerImpl) GetByCompositeKeyForThisISOWeek(ctx context.Context, metricID primitive.ObjectID, function int8) (*RankPoint, error) {
	thisWeekStart := timekit.FirstDayOfThisISOWeek(time.Now)
	thisWeekEnd := timekit.FirstDayOfNextISOWeek(time.Now)
	return impl.GetByCompositeKey(ctx, metricID, function, PeriodWeek, thisWeekStart, thisWeekEnd)
}

func (impl RankPointStorerImpl) GetByCompositeKeyForLastISOWeek(ctx context.Context, metricID primitive.ObjectID, function int8) (*RankPoint, error) {
	lastWeekStart := timekit.FirstDayOfLastISOWeek(time.Now)
	lastWeekEnd := timekit.FirstDayOfThisISOWeek(time.Now)
	return impl.GetByCompositeKey(ctx, metricID, function, PeriodWeek, lastWeekStart, lastWeekEnd)
}

func (impl RankPointStorerImpl) GetByCompositeKeyForThisMonth(ctx context.Context, metricID primitive.ObjectID, function int8) (*RankPoint, error) {
	thisMonthStart := timekit.FirstDayOfThisMonth(time.Now)
	thisMonthEnd := timekit.FirstDayOfNextMonth(time.Now)
	return impl.GetByCompositeKey(ctx, metricID, function, PeriodMonth, thisMonthStart, thisMonthEnd)
}

func (impl RankPointStorerImpl) GetByCompositeKeyForLastMonth(ctx context.Context, metricID primitive.ObjectID, function int8) (*RankPoint, error) {
	lastMonthStart := timekit.FirstDayOfLastMonth(time.Now)
	lastMonthEnd := timekit.FirstDayOfThisMonth(time.Now)
	return impl.GetByCompositeKey(ctx, metricID, function, PeriodMonth, lastMonthStart, lastMonthEnd)
}

func (impl RankPointStorerImpl) GetByCompositeKeyForThisYear(ctx context.Context, metricID primitive.ObjectID, function int8) (*RankPoint, error) {
	thisYearStart := timekit.FirstDayOfThisYear(time.Now)
	thisYearEnd := timekit.FirstDayOfNextYear(time.Now)
	return impl.GetByCompositeKey(ctx, metricID, function, PeriodYear, thisYearStart, thisYearEnd)
}

func (impl RankPointStorerImpl) GetByCompositeKeyForLastYear(ctx context.Context, metricID primitive.ObjectID, function int8) (*RankPoint, error) {
	lastYearStart := timekit.FirstDayOfLastYear(time.Now)
	lastYearEnd := timekit.FirstDayOfThisYear(time.Now)
	return impl.GetByCompositeKey(ctx, metricID, function, PeriodYear, lastYearStart, lastYearEnd)
}
