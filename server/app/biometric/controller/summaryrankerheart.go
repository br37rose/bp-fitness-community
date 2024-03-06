package controller

import (
	"context"
	"log/slog"

	rp_d "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/rankpoint/datastore"
	u_d "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/user/datastore"
)

func (impl *BiometricControllerImpl) generateSummaryRankingsForHR(ctx context.Context, u *u_d.User, res *AggregatePointSummaryResponse) error {
	if err := impl.generateSummaryRankingsForHRToday(ctx, u, res); err != nil {
		return err
	}
	if err := impl.generateSummaryRankingsForHRISOWeek(ctx, u, res); err != nil {
		return err
	}
	if err := impl.generateSummaryRankingsForHRMonth(ctx, u, res); err != nil {
		return err
	}
	if err := impl.generateSummaryRankingsForHRYear(ctx, u, res); err != nil {
		return err
	}
	return nil
}

func (impl *BiometricControllerImpl) generateSummaryRankingsForHRToday(ctx context.Context, u *u_d.User, res *AggregatePointSummaryResponse) error {

	////
	//// Generate the rankings for TODAY.
	////

	rp, err := impl.RankPointStorer.GetByCompositeKeyForToday(ctx, u.PrimaryHealthTrackingDeviceHeartRateMetricID, rp_d.FunctionAverage)
	if err != nil {
		impl.Logger.Error("failed getting by composite key for today",
			slog.String("metric_id", u.PrimaryHealthTrackingDeviceHeartRateMetricID.Hex()),
			slog.Any("error", err))
		return err
	}
	if rp != nil {
		// The following code will take the current user's place and generate
		// a `startPlace` which is two rankings below and `endPlace` which is
		// two rankings above the current user's rank place. As a result,
		// this will be used to filter and get limited rank that user falls
		// under. Because we are using `uint64` datatype, we cannot have non-
		// negative numbers, hence the `if-else` conditions below.
		var startPlace uint64 = 0
		if int64(rp.Place)-2 <= 0 {
			startPlace = 0
		} else {
			startPlace = rp.Place - 2
		}
		var endPlace uint64 = rp.Place + 2

		impl.Logger.Debug("rank within place",
			slog.Any("metric_type", "hr"),
			slog.Any("date_rank", "today"),
			slog.Any("start_place", startPlace),
			slog.Any("curr_place", rp.Place),
			slog.Any("end_place", endPlace))

		rprp, err := impl.RankPointStorer.ListWithinPlaceAndToday(ctx, []int8{rp_d.MetricTypeHeartRate}, rp_d.FunctionAverage, rp_d.PeriodDay, startPlace, endPlace)
		if err != nil {
			impl.Logger.Error("failed getting list for today",
				slog.String("metric_id", u.PrimaryHealthTrackingDeviceHeartRateMetricID.Hex()),
				slog.Any("error", err))
			return err
		}

		if rprp != nil {
			res.HeartRateThisDayRanking = rprp.Results
		}
	}

	return nil
}

func (impl *BiometricControllerImpl) generateSummaryRankingsForHRISOWeek(ctx context.Context, u *u_d.User, res *AggregatePointSummaryResponse) error {

	////
	//// Generate the rankings for ISO WEEK.
	////

	rp, err := impl.RankPointStorer.GetByCompositeKeyForThisISOWeek(ctx, u.PrimaryHealthTrackingDeviceHeartRateMetricID, rp_d.FunctionAverage)
	if err != nil {
		impl.Logger.Error("failed getting by composite key for iso week",
			slog.String("metric_id", u.PrimaryHealthTrackingDeviceHeartRateMetricID.Hex()),
			slog.Any("error", err))
		return err
	}
	if rp != nil {
		// The following code will take the current user's place and generate
		// a `startPlace` which is two rankings below and `endPlace` which is
		// two rankings above the current user's rank place. As a result,
		// this will be used to filter and get limited rank that user falls
		// under. Because we are using `uint64` datatype, we cannot have non-
		// negative numbers, hence the `if-else` conditions below.
		var startPlace uint64 = 0
		if int64(rp.Place)-2 <= 0 {
			startPlace = 0
		} else {
			startPlace = rp.Place - 2
		}
		var endPlace uint64 = rp.Place + 2

		impl.Logger.Debug("rank within place",
			slog.Any("metric_type", "hr"),
			slog.Any("date_rank", "iso_week"),
			slog.Any("start_place", startPlace),
			slog.Any("curr_place", rp.Place),
			slog.Any("end_place", endPlace))

		rprp, err := impl.RankPointStorer.ListWithinPlaceAndISOWeek(ctx, []int8{rp_d.MetricTypeHeartRate}, rp_d.FunctionAverage, rp_d.PeriodWeek, startPlace, endPlace)
		if err != nil {
			impl.Logger.Error("failed getting list for iso week",
				slog.String("metric_id", u.PrimaryHealthTrackingDeviceHeartRateMetricID.Hex()),
				slog.Any("error", err))
			return err
		}

		if rprp != nil {
			res.HeartRateThisISOWeekRanking = rprp.Results
		}
	} else {
		impl.Logger.Debug("rank within place not found",
			slog.Any("metric_type", "hr"),
			slog.Any("function", rp_d.FunctionAverage),
			slog.Any("date_rank", "iso_week"))
	}
	return nil
}

func (impl *BiometricControllerImpl) generateSummaryRankingsForHRMonth(ctx context.Context, u *u_d.User, res *AggregatePointSummaryResponse) error {

	////
	//// Generate the rankings for ISO WEEK.
	////

	rp, err := impl.RankPointStorer.GetByCompositeKeyForThisMonth(ctx, u.PrimaryHealthTrackingDeviceHeartRateMetricID, rp_d.FunctionAverage)
	if err != nil {
		impl.Logger.Error("failed getting by composite key for month",
			slog.String("metric_id", u.PrimaryHealthTrackingDeviceHeartRateMetricID.Hex()),
			slog.Any("error", err))
		return err
	}
	if rp != nil {
		// The following code will take the current user's place and generate
		// a `startPlace` which is two rankings below and `endPlace` which is
		// two rankings above the current user's rank place. As a result,
		// this will be used to filter and get limited rank that user falls
		// under. Because we are using `uint64` datatype, we cannot have non-
		// negative numbers, hence the `if-else` conditions below.
		var startPlace uint64 = 0
		if int64(rp.Place)-2 <= 0 {
			startPlace = 0
		} else {
			startPlace = rp.Place - 2
		}
		var endPlace uint64 = rp.Place + 2

		impl.Logger.Debug("rank within place",
			slog.Any("metric_type", "hr"),
			slog.Any("date_rank", "month"),
			slog.Any("start_place", startPlace),
			slog.Any("curr_place", rp.Place),
			slog.Any("end_place", endPlace))

		rprp, err := impl.RankPointStorer.ListWithinPlaceAndMonth(ctx, []int8{rp_d.MetricTypeHeartRate}, rp_d.FunctionAverage, rp_d.PeriodMonth, startPlace, endPlace)
		if err != nil {
			impl.Logger.Error("failed getting list for month",
				slog.String("metric_id", u.PrimaryHealthTrackingDeviceHeartRateMetricID.Hex()),
				slog.Any("error", err))
			return err
		}

		if rprp != nil {
			res.HeartRateThisMonthRanking = rprp.Results
		}
	} else {
		impl.Logger.Debug("rank within place not found",
			slog.Any("metric_type", "hr"),
			slog.Any("date_rank", "month"))
	}

	return nil
}

func (impl *BiometricControllerImpl) generateSummaryRankingsForHRYear(ctx context.Context, u *u_d.User, res *AggregatePointSummaryResponse) error {
	rp, err := impl.RankPointStorer.GetByCompositeKeyForThisYear(ctx, u.PrimaryHealthTrackingDeviceHeartRateMetricID, rp_d.FunctionAverage)
	if err != nil {
		impl.Logger.Error("failed getting by composite key for year",
			slog.String("metric_id", u.PrimaryHealthTrackingDeviceHeartRateMetricID.Hex()),
			slog.Any("error", err))
		return err
	}
	if rp != nil {
		// The following code will take the current user's place and generate
		// a `startPlace` which is two rankings below and `endPlace` which is
		// two rankings above the current user's rank place. As a result,
		// this will be used to filter and get limited rank that user falls
		// under. Because we are using `uint64` datatype, we cannot have non-
		// negative numbers, hence the `if-else` conditions below.
		var startPlace uint64 = 0
		if int64(rp.Place)-2 <= 0 {
			startPlace = 0
		} else {
			startPlace = rp.Place - 2
		}
		var endPlace uint64 = rp.Place + 2

		impl.Logger.Debug("rank within place",
			slog.Any("metric_type", "hr"),
			slog.Any("date_rank", "year"),
			slog.Any("start_place", startPlace),
			slog.Any("curr_place", rp.Place),
			slog.Any("end_place", endPlace))

		rprp, err := impl.RankPointStorer.ListWithinPlaceAndYear(ctx, []int8{rp_d.MetricTypeHeartRate}, rp_d.FunctionAverage, rp_d.PeriodYear, startPlace, endPlace)
		if err != nil {
			impl.Logger.Error("failed getting list for year",
				slog.String("metric_id", u.PrimaryHealthTrackingDeviceHeartRateMetricID.Hex()),
				slog.Any("error", err))
			return err
		}

		if rprp != nil {
			res.HeartRateThisYearRanking = rprp.Results
		}
	} else {
		impl.Logger.Debug("rank within place not found",
			slog.Any("metric_type", "hr"),
			slog.Any("date_rank", "year"))
	}

	return nil
}
