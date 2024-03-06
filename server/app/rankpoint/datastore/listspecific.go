package datastore

import (
	"context"
	"time"

	"github.com/bartmika/timekit"
)

func (impl RankPointStorerImpl) ListWithinPlace(ctx context.Context, metricTypes []int8, function int8, period int8, start, end uint64) (*RankPointPaginationListResult, error) {
	f := &RankPointPaginationListFilter{
		Cursor:      "",
		PageSize:    1_000_000_000, // Unlimited
		SortField:   "place",
		SortOrder:   OrderAscending,
		MetricTypes: metricTypes,
		Function:    function,
		Period:      period,
		PlaceGTE:    start,
		PlaceLTE:    end,
	}
	return impl.ListByFilter(ctx, f)
}

func (impl RankPointStorerImpl) ListWithinPlaceAndToday(ctx context.Context, metricTypes []int8, function int8, period int8, start, end uint64) (*RankPointPaginationListResult, error) {
	thisDayStart := timekit.Midnight(time.Now)
	thisDayEnd := timekit.MidnightTomorrow(time.Now)
	f := &RankPointPaginationListFilter{
		Cursor:      "",
		PageSize:    1_000_000_000, // Unlimited
		SortField:   "place",
		SortOrder:   OrderAscending,
		MetricTypes: metricTypes,
		Function:    function,
		Period:      period,
		PlaceGTE:    start,
		PlaceLTE:    end,
		StartGTE:    thisDayStart,
		EndLTE:      thisDayEnd,
	}
	return impl.ListByFilter(ctx, f)
}

func (impl RankPointStorerImpl) ListWithinPlaceAndISOWeek(ctx context.Context, metricTypes []int8, function int8, period int8, start, end uint64) (*RankPointPaginationListResult, error) {
	thisWeekStart := timekit.FirstDayOfThisISOWeek(time.Now)
	thisWeekEnd := timekit.FirstDayOfNextISOWeek(time.Now)
	f := &RankPointPaginationListFilter{
		Cursor:      "",
		PageSize:    1_000_000_000, // Unlimited
		SortField:   "place",
		SortOrder:   OrderAscending,
		MetricTypes: metricTypes,
		Function:    function,
		Period:      period,
		PlaceGTE:    start,
		PlaceLTE:    end,
		StartGTE:    thisWeekStart,
		EndLTE:      thisWeekEnd,
	}
	return impl.ListByFilter(ctx, f)
}

func (impl RankPointStorerImpl) ListWithinPlaceAndMonth(ctx context.Context, metricTypes []int8, function int8, period int8, start, end uint64) (*RankPointPaginationListResult, error) {
	thisMonthStart := timekit.FirstDayOfThisMonth(time.Now)
	thisMonthEnd := timekit.FirstDayOfNextMonth(time.Now)
	f := &RankPointPaginationListFilter{
		Cursor:      "",
		PageSize:    1_000_000_000, // Unlimited
		SortField:   "place",
		SortOrder:   OrderAscending,
		MetricTypes: metricTypes,
		Function:    function,
		Period:      period,
		PlaceGTE:    start,
		PlaceLTE:    end,
		StartGTE:    thisMonthStart,
		EndLTE:      thisMonthEnd,
	}
	return impl.ListByFilter(ctx, f)
}

func (impl RankPointStorerImpl) ListWithinPlaceAndYear(ctx context.Context, metricTypes []int8, function int8, period int8, start, end uint64) (*RankPointPaginationListResult, error) {
	thisYearStart := timekit.FirstDayOfThisYear(time.Now)
	thisYearEnd := timekit.FirstDayOfNextYear(time.Now)
	f := &RankPointPaginationListFilter{
		Cursor:      "",
		PageSize:    1_000_000_000, // Unlimited
		SortField:   "place",
		SortOrder:   OrderAscending,
		MetricTypes: metricTypes,
		Function:    function,
		Period:      period,
		PlaceGTE:    start,
		PlaceLTE:    end,
		StartGTE:    thisYearStart,
		EndLTE:      thisYearEnd,
	}
	return impl.ListByFilter(ctx, f)
}
