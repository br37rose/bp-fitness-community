package controller

import (
	"log/slog"
	"sync"

	"go.mongodb.org/mongo-driver/mongo"

	rp_d "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/rankpoint/datastore"
	u_d "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/user/datastore"
)

func (impl *BiometricControllerImpl) generateSummaryRankingsForHR(sessCtx mongo.SessionContext, u *u_d.User, res *AggregatePointSummaryResponse, mu *sync.Mutex, wg *sync.WaitGroup) error {
	// Once this function has been completed (whether successfully or not) then
	// update the `WaitGroup` that this goroutine is finished.
	defer wg.Done()

	// Developers note:
	// Since we'll add more goroutines in this function, let's increase our
	// `WaitGroup` now.
	// - Today
	// - Yesterday
	// - This ISO Week
	// - Last ISO Week
	// - This Month
	// - Last Month
	// - This Year
	// - Last Year
	wg.Add(4)

	// --- Today --- //
	go func() {
		if err := impl.generateSummaryRankingsForHRToday(sessCtx, u, res, mu, wg); err != nil {
			// return err
		}
	}()

	go func() {
		if err := impl.generateSummaryRankingsForHRISOWeek(sessCtx, u, res, mu, wg); err != nil {
			//return err
		}
	}()

	go func() {
		if err := impl.generateSummaryRankingsForHRMonth(sessCtx, u, res, mu, wg); err != nil {
			// return err
		}
	}()

	go func() {
		if err := impl.generateSummaryRankingsForHRYear(sessCtx, u, res, mu, wg); err != nil {
			// return err
		}
	}()
	return nil
}

func (impl *BiometricControllerImpl) generateSummaryRankingsForHRToday(sessCtx mongo.SessionContext, u *u_d.User, res *AggregatePointSummaryResponse, mu *sync.Mutex, wg *sync.WaitGroup) error {
	// Once this function has been completed (whether successfully or not) then
	// update the `WaitGroup` that this goroutine is finished.
	defer wg.Done()

	////
	//// Generate the rankings for TODAY.
	////

	rp, err := impl.RankPointStorer.GetByCompositeKeyForToday(sessCtx, u.PrimaryHealthTrackingDevice.HeartRateBPMMetricID, rp_d.FunctionAverage)
	if err != nil {
		impl.Logger.Error("failed getting by composite key for today",
			slog.String("metric_id", u.PrimaryHealthTrackingDevice.HeartRateBPMMetricID.Hex()),
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

		rprp, err := impl.RankPointStorer.ListWithinPlaceAndToday(sessCtx, []int8{rp_d.MetricTypeHeartRate}, rp_d.FunctionAverage, rp_d.PeriodDay, startPlace, endPlace)
		if err != nil {
			impl.Logger.Error("failed getting list for today",
				slog.String("metric_id", u.PrimaryHealthTrackingDevice.HeartRateBPMMetricID.Hex()),
				slog.Any("error", err))
			return err
		}

		if rprp != nil {
			// Lock the mutex before accessing res
			mu.Lock()
			defer mu.Unlock()

			res.HeartRateThisDayRanking = rprp.Results
		}
	}

	return nil
}

func (impl *BiometricControllerImpl) generateSummaryRankingsForHRISOWeek(sessCtx mongo.SessionContext, u *u_d.User, res *AggregatePointSummaryResponse, mu *sync.Mutex, wg *sync.WaitGroup) error {
	// Once this function has been completed (whether successfully or not) then
	// update the `WaitGroup` that this goroutine is finished.
	defer wg.Done()

	////
	//// Generate the rankings for ISO WEEK.
	////

	rp, err := impl.RankPointStorer.GetByCompositeKeyForThisISOWeek(sessCtx, u.PrimaryHealthTrackingDevice.HeartRateBPMMetricID, rp_d.FunctionAverage)
	if err != nil {
		impl.Logger.Error("failed getting by composite key for iso week",
			slog.String("metric_id", u.PrimaryHealthTrackingDevice.HeartRateBPMMetricID.Hex()),
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

		rprp, err := impl.RankPointStorer.ListWithinPlaceAndISOWeek(sessCtx, []int8{rp_d.MetricTypeHeartRate}, rp_d.FunctionAverage, rp_d.PeriodWeek, startPlace, endPlace)
		if err != nil {
			impl.Logger.Error("failed getting list for iso week",
				slog.String("metric_id", u.PrimaryHealthTrackingDevice.HeartRateBPMMetricID.Hex()),
				slog.Any("error", err))
			return err
		}

		if rprp != nil {
			// Lock the mutex before accessing res
			mu.Lock()
			defer mu.Unlock()

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

func (impl *BiometricControllerImpl) generateSummaryRankingsForHRMonth(sessCtx mongo.SessionContext, u *u_d.User, res *AggregatePointSummaryResponse, mu *sync.Mutex, wg *sync.WaitGroup) error {
	// Once this function has been completed (whether successfully or not) then
	// update the `WaitGroup` that this goroutine is finished.
	defer wg.Done()

	////
	//// Generate the rankings for ISO WEEK.
	////

	rp, err := impl.RankPointStorer.GetByCompositeKeyForThisMonth(sessCtx, u.PrimaryHealthTrackingDevice.HeartRateBPMMetricID, rp_d.FunctionAverage)
	if err != nil {
		impl.Logger.Error("failed getting by composite key for month",
			slog.String("metric_id", u.PrimaryHealthTrackingDevice.HeartRateBPMMetricID.Hex()),
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

		rprp, err := impl.RankPointStorer.ListWithinPlaceAndMonth(sessCtx, []int8{rp_d.MetricTypeHeartRate}, rp_d.FunctionAverage, rp_d.PeriodMonth, startPlace, endPlace)
		if err != nil {
			impl.Logger.Error("failed getting list for month",
				slog.String("metric_id", u.PrimaryHealthTrackingDevice.HeartRateBPMMetricID.Hex()),
				slog.Any("error", err))
			return err
		}

		if rprp != nil {
			// Lock the mutex before accessing res
			mu.Lock()
			defer mu.Unlock()

			res.HeartRateThisMonthRanking = rprp.Results
		}
	} else {
		impl.Logger.Debug("rank within place not found",
			slog.Any("metric_type", "hr"),
			slog.Any("date_rank", "month"))
	}

	return nil
}

func (impl *BiometricControllerImpl) generateSummaryRankingsForHRYear(sessCtx mongo.SessionContext, u *u_d.User, res *AggregatePointSummaryResponse, mu *sync.Mutex, wg *sync.WaitGroup) error {
	// Once this function has been completed (whether successfully or not) then
	// update the `WaitGroup` that this goroutine is finished.
	defer wg.Done()

	rp, err := impl.RankPointStorer.GetByCompositeKeyForThisYear(sessCtx, u.PrimaryHealthTrackingDevice.HeartRateBPMMetricID, rp_d.FunctionAverage)
	if err != nil {
		impl.Logger.Error("failed getting by composite key for year",
			slog.String("metric_id", u.PrimaryHealthTrackingDevice.HeartRateBPMMetricID.Hex()),
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

		rprp, err := impl.RankPointStorer.ListWithinPlaceAndYear(sessCtx, []int8{rp_d.MetricTypeHeartRate}, rp_d.FunctionAverage, rp_d.PeriodYear, startPlace, endPlace)
		if err != nil {
			impl.Logger.Error("failed getting list for year",
				slog.String("metric_id", u.PrimaryHealthTrackingDevice.HeartRateBPMMetricID.Hex()),
				slog.Any("error", err))
			return err
		}

		if rprp != nil {
			// Lock the mutex before accessing res
			mu.Lock()
			defer mu.Unlock()

			res.HeartRateThisYearRanking = rprp.Results
		}
	} else {
		impl.Logger.Debug("rank within place not found",
			slog.Any("metric_type", "hr"),
			slog.Any("date_rank", "year"))
	}

	return nil
}
