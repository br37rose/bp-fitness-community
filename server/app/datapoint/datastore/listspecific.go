package datastore

import (
	"context"
	"time"

	"github.com/bartmika/timekit"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (impl DataPointStorerImpl) ListForToday(ctx context.Context, metricIDs []primitive.ObjectID) (*DataPointPaginationListResult, error) {
	thisDayStart := timekit.Midnight(time.Now)
	thisDayEnd := timekit.MidnightTomorrow(time.Now)

	dpf := &DataPointPaginationListFilter{
		MetricIDs: metricIDs,
		Cursor:    "",
		PageSize:  1_000_000_000,
		SortField: "timestamp",
		SortOrder: OrderAscending,
		GTE:       thisDayStart,
		LTE:       thisDayEnd,
	}
	return impl.ListByFilter(ctx, dpf)
}

func (impl DataPointStorerImpl) ListForYesterday(ctx context.Context, metricIDs []primitive.ObjectID) (*DataPointPaginationListResult, error) {
	lastDayStart := timekit.MidnightYesterday(time.Now)
	lastDayEnd := timekit.Midnight(time.Now)

	dpf := &DataPointPaginationListFilter{
		MetricIDs: metricIDs,
		Cursor:    "",
		PageSize:  1_000_000_000,
		SortField: "timestamp",
		SortOrder: OrderAscending,
		GTE:       lastDayStart,
		LTE:       lastDayEnd,
	}
	return impl.ListByFilter(ctx, dpf)
}

func (impl DataPointStorerImpl) ListForThisISOWeek(ctx context.Context, metricIDs []primitive.ObjectID) (*DataPointPaginationListResult, error) {
	thisWeekStart := timekit.FirstDayOfThisISOWeek(time.Now)
	thisWeekEnd := timekit.FirstDayOfNextISOWeek(time.Now)

	dpf := &DataPointPaginationListFilter{
		MetricIDs: metricIDs,
		Cursor:    "",
		PageSize:  1_000_000_000,
		SortField: "timestamp",
		SortOrder: OrderAscending,
		GTE:       thisWeekStart,
		LTE:       thisWeekEnd,
	}
	return impl.ListByFilter(ctx, dpf)
}

func (impl DataPointStorerImpl) ListForLastISOWeek(ctx context.Context, metricIDs []primitive.ObjectID) (*DataPointPaginationListResult, error) {
	lastWeekStart := timekit.FirstDayOfLastISOWeek(time.Now)
	lastWeekEnd := timekit.FirstDayOfThisISOWeek(time.Now)

	dpf := &DataPointPaginationListFilter{
		MetricIDs: metricIDs,
		Cursor:    "",
		PageSize:  1_000_000_000,
		SortField: "timestamp",
		SortOrder: OrderAscending,
		GTE:       lastWeekStart,
		LTE:       lastWeekEnd,
	}
	return impl.ListByFilter(ctx, dpf)
}

func (impl DataPointStorerImpl) ListForThisMonth(ctx context.Context, metricIDs []primitive.ObjectID) (*DataPointPaginationListResult, error) {
	thisMonthStart := timekit.FirstDayOfThisMonth(time.Now)
	thisMonthEnd := timekit.FirstDayOfNextMonth(time.Now)

	dpf := &DataPointPaginationListFilter{
		MetricIDs: metricIDs,
		Cursor:    "",
		PageSize:  1_000_000_000,
		SortField: "timestamp",
		SortOrder: OrderAscending,
		GTE:       thisMonthStart,
		LTE:       thisMonthEnd,
	}
	return impl.ListByFilter(ctx, dpf)
}

func (impl DataPointStorerImpl) ListForLastMonth(ctx context.Context, metricIDs []primitive.ObjectID) (*DataPointPaginationListResult, error) {
	lastMonthStart := timekit.FirstDayOfLastMonth(time.Now)
	lastMonthEnd := timekit.FirstDayOfThisMonth(time.Now)

	dpf := &DataPointPaginationListFilter{
		MetricIDs: metricIDs,
		Cursor:    "",
		PageSize:  1_000_000_000,
		SortField: "timestamp",
		SortOrder: OrderAscending,
		GTE:       lastMonthStart,
		LTE:       lastMonthEnd,
	}
	return impl.ListByFilter(ctx, dpf)
}

func (impl DataPointStorerImpl) ListForThisYear(ctx context.Context, metricIDs []primitive.ObjectID) (*DataPointPaginationListResult, error) {
	thisYearStart := timekit.FirstDayOfThisYear(time.Now)
	thisYearEnd := timekit.FirstDayOfNextYear(time.Now)
	dpf := &DataPointPaginationListFilter{
		MetricIDs: metricIDs,
		Cursor:    "",
		PageSize:  1_000_000_000,
		SortField: "timestamp",
		SortOrder: OrderAscending,
		GTE:       thisYearStart,
		LTE:       thisYearEnd,
	}
	return impl.ListByFilter(ctx, dpf)
}

func (impl DataPointStorerImpl) ListForLastYear(ctx context.Context, metricIDs []primitive.ObjectID) (*DataPointPaginationListResult, error) {
	lastYearStart := timekit.FirstDayOfLastYear(time.Now)
	lastYearEnd := timekit.FirstDayOfThisYear(time.Now)

	dpf := &DataPointPaginationListFilter{
		MetricIDs: metricIDs,
		Cursor:    "",
		PageSize:  1_000_000_000,
		SortField: "timestamp",
		SortOrder: OrderAscending,
		GTE:       lastYearStart,
		LTE:       lastYearEnd,
	}
	return impl.ListByFilter(ctx, dpf)
}
