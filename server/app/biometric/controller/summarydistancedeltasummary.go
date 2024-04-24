package controller

import (
	"log/slog"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/mongo"

	"github.com/bartmika/timekit"
	ap_s "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/aggregatepoint/datastore"
	u_d "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/user/datastore"
)

func (impl *BiometricControllerImpl) generateSummarySummaryForDistanceDelta(sessCtx mongo.SessionContext, u *u_d.User, res *AggregatePointSummaryResponse, mu *sync.Mutex, wg *sync.WaitGroup) error {
	// Once this function has been completed (whether successfully or not) then
	// update the `WaitGroup` that this goroutine is finished.
	defer wg.Done()

	// Developers note:
	// Since we'll add more goroutines in this function, let's increase our
	// `WaitGroup` now.
	// - This Hour
	// - Last Hour
	// - Today
	// - Yesterday
	// - This ISO Week
	// - Last ISO Week
	// - This Month
	// - Last Month
	// - This Year
	// - Last Year
	wg.Add(10)

	// --- This Hour --- //
	go func() {
		if err := impl.generateSummarySummaryForDistanceDeltaThisHour(sessCtx, u, res, mu, wg); err != nil {
			//TODO
		}
	}()
	// --- Last Hour --- //
	go func() {
		if err := impl.generateSummarySummaryForDistanceDeltaLastHour(sessCtx, u, res, mu, wg); err != nil {
			//TODO
		}
	}()

	// --- Today --- //
	go func() {
		if err := impl.generateSummarySummaryForDistanceDeltaToday(sessCtx, u, res, mu, wg); err != nil {
			//TODO
		}
	}()
	// --- Yesterday --- //
	go func() {
		if err := impl.generateSummarySummaryForDistanceDeltaYesterday(sessCtx, u, res, mu, wg); err != nil {
			//TODO
		}
	}()
	// --- This ISO Week --- //
	go func() {
		if err := impl.generateSummarySummaryForDistanceDeltaThisISOWeek(sessCtx, u, res, mu, wg); err != nil {
			//TODO
		}
	}()
	// --- Last ISO Week --- //
	go func() {
		if err := impl.generateSummarySummaryForDistanceDeltaLastISOWeek(sessCtx, u, res, mu, wg); err != nil {
			//TODO
		}
	}()
	// --- This Month --- //
	go func() {
		if err := impl.generateSummarySummaryForDistanceDeltaThisMonth(sessCtx, u, res, mu, wg); err != nil {
			//TODO
		}
	}()
	// --- Last Month --- //
	go func() {
		if err := impl.generateSummarySummaryForDistanceDeltaLastMonth(sessCtx, u, res, mu, wg); err != nil {
			//TODO
		}
	}()
	// --- This Year --- //
	go func() {
		if err := impl.generateSummarySummaryForDistanceDeltaThisYear(sessCtx, u, res, mu, wg); err != nil {
			//TODO
		}
	}()
	// --- Last Year --- //
	go func() {
		if err := impl.generateSummarySummaryForDistanceDeltaLastYear(sessCtx, u, res, mu, wg); err != nil {
			//TODO
		}
	}()

	return nil
}

func (impl *BiometricControllerImpl) generateSummarySummaryForDistanceDeltaThisHour(sessCtx mongo.SessionContext, u *u_d.User, res *AggregatePointSummaryResponse, mu *sync.Mutex, wg *sync.WaitGroup) error {
	// Once this function has been completed (whether successfully or not) then
	// update the `WaitGroup` that this goroutine is finished.
	defer wg.Done()

	thisHourStart, thisHourEnd := timekit.HourRangeForNow(time.Now)
	thisHour, err := impl.AggregatePointStorer.GetByCompositeKey(sessCtx, u.PrimaryHealthTrackingDevice.DistanceDeltaMetricID, ap_s.PeriodHour, thisHourStart, thisHourEnd)
	if err != nil {
		impl.Logger.Error("failed getting aggregate point by composite key",
			slog.String("metric_id", u.PrimaryHealthTrackingDevice.DistanceDeltaMetricID.Hex()),
			slog.Time("start", thisHourStart),
			slog.Time("end", thisHourEnd),
			slog.Any("error", err))
		return err
	}

	// Lock the mutex before accessing res
	mu.Lock()
	defer mu.Unlock()

	if thisHour != nil {
		res.DistanceDeltaThisHourSummary = thisHour
	}

	return nil
}

func (impl *BiometricControllerImpl) generateSummarySummaryForDistanceDeltaLastHour(sessCtx mongo.SessionContext, u *u_d.User, res *AggregatePointSummaryResponse, mu *sync.Mutex, wg *sync.WaitGroup) error {
	// Once this function has been completed (whether successfully or not) then
	// update the `WaitGroup` that this goroutine is finished.
	defer wg.Done()

	thisHourStart, thisHourEnd := timekit.HourRangeForNow(time.Now)
	lastHourStart := thisHourStart.Add((-1) * time.Hour)
	lastHourEnd := thisHourEnd.Add((-1) * time.Hour)
	lastHour, err := impl.AggregatePointStorer.GetByCompositeKey(sessCtx, u.PrimaryHealthTrackingDevice.DistanceDeltaMetricID, ap_s.PeriodHour, lastHourStart, lastHourEnd)
	if err != nil {
		impl.Logger.Error("failed getting aggregate point by composite key",
			slog.String("metric_id", u.PrimaryHealthTrackingDevice.DistanceDeltaMetricID.Hex()),
			slog.Time("start", lastHourStart),
			slog.Time("end", lastHourEnd),
			slog.Any("error", err))
		return err
	}

	// Lock the mutex before accessing res
	mu.Lock()
	defer mu.Unlock()

	if lastHour != nil {
		res.DistanceDeltaLastHourSummary = lastHour
	}

	return nil
}

func (impl *BiometricControllerImpl) generateSummarySummaryForDistanceDeltaToday(sessCtx mongo.SessionContext, u *u_d.User, res *AggregatePointSummaryResponse, mu *sync.Mutex, wg *sync.WaitGroup) error {
	// Once this function has been completed (whether successfully or not) then
	// update the `WaitGroup` that this goroutine is finished.
	defer wg.Done()

	thisDayStart := timekit.Midnight(time.Now)
	thisDayEnd := timekit.MidnightTomorrow(time.Now)
	thisDay, err := impl.AggregatePointStorer.GetByCompositeKey(sessCtx, u.PrimaryHealthTrackingDevice.DistanceDeltaMetricID, ap_s.PeriodDay, thisDayStart, thisDayEnd)
	if err != nil {
		impl.Logger.Error("failed getting aggregate point by composite key",
			slog.String("metric_id", u.PrimaryHealthTrackingDevice.DistanceDeltaMetricID.Hex()),
			slog.Time("start", thisDayStart),
			slog.Time("end", thisDayEnd),
			slog.Any("error", err))
		return err
	}

	// Lock the mutex before accessing res
	mu.Lock()
	defer mu.Unlock()

	if thisDay != nil {
		res.DistanceDeltaThisDaySummary = thisDay
	}

	return nil
}

func (impl *BiometricControllerImpl) generateSummarySummaryForDistanceDeltaYesterday(sessCtx mongo.SessionContext, u *u_d.User, res *AggregatePointSummaryResponse, mu *sync.Mutex, wg *sync.WaitGroup) error {
	// Once this function has been completed (whether successfully or not) then
	// update the `WaitGroup` that this goroutine is finished.
	defer wg.Done()

	lastDayStart := timekit.MidnightYesterday(time.Now)
	lastDayEnd := timekit.Midnight(time.Now)
	lastDay, err := impl.AggregatePointStorer.GetByCompositeKey(sessCtx, u.PrimaryHealthTrackingDevice.DistanceDeltaMetricID, ap_s.PeriodDay, lastDayStart, lastDayEnd)
	if err != nil {
		impl.Logger.Error("failed getting aggregate point by composite key",
			slog.String("metric_id", u.PrimaryHealthTrackingDevice.DistanceDeltaMetricID.Hex()),
			slog.Time("start", lastDayStart),
			slog.Time("end", lastDayEnd),
			slog.Any("error", err))
		return err
	}

	// Lock the mutex before accessing res
	mu.Lock()
	defer mu.Unlock()

	if lastDay != nil {
		res.DistanceDeltaLastDaySummary = lastDay
	}

	return nil
}

func (impl *BiometricControllerImpl) generateSummarySummaryForDistanceDeltaThisISOWeek(sessCtx mongo.SessionContext, u *u_d.User, res *AggregatePointSummaryResponse, mu *sync.Mutex, wg *sync.WaitGroup) error {
	// Once this function has been completed (whether successfully or not) then
	// update the `WaitGroup` that this goroutine is finished.
	defer wg.Done()

	thisWeekStart := timekit.FirstDayOfThisISOWeek(time.Now)
	thisWeekEnd := timekit.FirstDayOfNextISOWeek(time.Now)
	thisWeek, err := impl.AggregatePointStorer.GetByCompositeKey(sessCtx, u.PrimaryHealthTrackingDevice.DistanceDeltaMetricID, ap_s.PeriodWeek, thisWeekStart, thisWeekEnd)
	if err != nil {
		impl.Logger.Error("failed getting aggregate point by composite key",
			slog.String("metric_id", u.PrimaryHealthTrackingDevice.DistanceDeltaMetricID.Hex()),
			slog.Time("start", thisWeekStart),
			slog.Time("end", thisWeekEnd),
			slog.Any("error", err))
		return err
	}

	// Lock the mutex before accessing res
	mu.Lock()
	defer mu.Unlock()

	if thisWeek != nil {
		res.DistanceDeltaThisISOWeekSummary = thisWeek
	}

	return nil
}

func (impl *BiometricControllerImpl) generateSummarySummaryForDistanceDeltaLastISOWeek(sessCtx mongo.SessionContext, u *u_d.User, res *AggregatePointSummaryResponse, mu *sync.Mutex, wg *sync.WaitGroup) error {
	// Once this function has been completed (whether successfully or not) then
	// update the `WaitGroup` that this goroutine is finished.
	defer wg.Done()

	lastWeekStart := timekit.FirstDayOfLastISOWeek(time.Now)
	lastWeekEnd := timekit.FirstDayOfThisISOWeek(time.Now)
	lastWeek, err := impl.AggregatePointStorer.GetByCompositeKey(sessCtx, u.PrimaryHealthTrackingDevice.DistanceDeltaMetricID, ap_s.PeriodWeek, lastWeekStart, lastWeekEnd)
	if err != nil {
		impl.Logger.Error("failed getting aggregate point by composite key",
			slog.String("metric_id", u.PrimaryHealthTrackingDevice.DistanceDeltaMetricID.Hex()),
			slog.Time("start", lastWeekStart),
			slog.Time("end", lastWeekEnd),
			slog.Any("error", err))
		return err
	}

	// Lock the mutex before accessing res
	mu.Lock()
	defer mu.Unlock()

	if lastWeek != nil {
		res.DistanceDeltaLastISOWeekSummary = lastWeek
	}

	return nil
}

func (impl *BiometricControllerImpl) generateSummarySummaryForDistanceDeltaThisMonth(sessCtx mongo.SessionContext, u *u_d.User, res *AggregatePointSummaryResponse, mu *sync.Mutex, wg *sync.WaitGroup) error {
	// Once this function has been completed (whether successfully or not) then
	// update the `WaitGroup` that this goroutine is finished.
	defer wg.Done()

	thisMonthStart := timekit.FirstDayOfThisMonth(time.Now)
	thisMonthEnd := timekit.FirstDayOfNextMonth(time.Now)
	thisMonth, err := impl.AggregatePointStorer.GetByCompositeKey(sessCtx, u.PrimaryHealthTrackingDevice.DistanceDeltaMetricID, ap_s.PeriodMonth, thisMonthStart, thisMonthEnd)
	if err != nil {
		impl.Logger.Error("failed getting aggregate point by composite key",
			slog.String("metric_id", u.PrimaryHealthTrackingDevice.DistanceDeltaMetricID.Hex()),
			slog.Time("start", thisMonthStart),
			slog.Time("end", thisMonthEnd),
			slog.Any("error", err))
		return err
	}

	// Lock the mutex before accessing res
	mu.Lock()
	defer mu.Unlock()

	if thisMonth != nil {
		res.DistanceDeltaThisMonthSummary = thisMonth
	}

	return nil
}

func (impl *BiometricControllerImpl) generateSummarySummaryForDistanceDeltaLastMonth(sessCtx mongo.SessionContext, u *u_d.User, res *AggregatePointSummaryResponse, mu *sync.Mutex, wg *sync.WaitGroup) error {
	// Once this function has been completed (whether successfully or not) then
	// update the `WaitGroup` that this goroutine is finished.
	defer wg.Done()

	lastMonthStart := timekit.FirstDayOfLastMonth(time.Now)
	lastMonthEnd := timekit.FirstDayOfThisMonth(time.Now)
	lastMonth, err := impl.AggregatePointStorer.GetByCompositeKey(sessCtx, u.PrimaryHealthTrackingDevice.DistanceDeltaMetricID, ap_s.PeriodMonth, lastMonthStart, lastMonthEnd)
	if err != nil {
		impl.Logger.Error("failed getting aggregate point by composite key",
			slog.String("metric_id", u.PrimaryHealthTrackingDevice.DistanceDeltaMetricID.Hex()),
			slog.Time("start", lastMonthStart),
			slog.Time("end", lastMonthEnd),
			slog.Any("error", err))
		return err
	}

	// Lock the mutex before accessing res
	mu.Lock()
	defer mu.Unlock()

	if lastMonth != nil {
		res.DistanceDeltaLastMonthSummary = lastMonth
	}

	return nil
}

func (impl *BiometricControllerImpl) generateSummarySummaryForDistanceDeltaThisYear(sessCtx mongo.SessionContext, u *u_d.User, res *AggregatePointSummaryResponse, mu *sync.Mutex, wg *sync.WaitGroup) error {
	// Once this function has been completed (whether successfully or not) then
	// update the `WaitGroup` that this goroutine is finished.
	defer wg.Done()

	thisYearStart := timekit.FirstDayOfThisYear(time.Now)
	thisYearEnd := timekit.FirstDayOfNextYear(time.Now)
	thisYear, err := impl.AggregatePointStorer.GetByCompositeKey(sessCtx, u.PrimaryHealthTrackingDevice.DistanceDeltaMetricID, ap_s.PeriodYear, thisYearStart, thisYearEnd)
	if err != nil {
		impl.Logger.Error("failed getting aggregate point by composite key",
			slog.String("metric_id", u.PrimaryHealthTrackingDevice.DistanceDeltaMetricID.Hex()),
			slog.Time("start", thisYearStart),
			slog.Time("end", thisYearEnd),
			slog.Any("error", err))
		return err
	}

	// Lock the mutex before accessing res
	mu.Lock()
	defer mu.Unlock()

	if thisYear != nil {
		res.DistanceDeltaThisYearSummary = thisYear
	}

	return nil
}

func (impl *BiometricControllerImpl) generateSummarySummaryForDistanceDeltaLastYear(sessCtx mongo.SessionContext, u *u_d.User, res *AggregatePointSummaryResponse, mu *sync.Mutex, wg *sync.WaitGroup) error {
	// Once this function has been completed (whether successfully or not) then
	// update the `WaitGroup` that this goroutine is finished.
	defer wg.Done()

	lastYearStart := timekit.FirstDayOfLastYear(time.Now)
	lastYearEnd := timekit.FirstDayOfThisYear(time.Now)
	lastYear, err := impl.AggregatePointStorer.GetByCompositeKey(sessCtx, u.PrimaryHealthTrackingDevice.DistanceDeltaMetricID, ap_s.PeriodYear, lastYearStart, lastYearEnd)
	if err != nil {
		impl.Logger.Error("failed getting aggregate point by composite key",
			slog.String("metric_id", u.PrimaryHealthTrackingDevice.DistanceDeltaMetricID.Hex()),
			slog.Time("start", lastYearStart),
			slog.Time("end", lastYearEnd),
			slog.Any("error", err))
		return err
	}

	// Lock the mutex before accessing res
	mu.Lock()
	defer mu.Unlock()

	if lastYear != nil {
		res.DistanceDeltaLastYearSummary = lastYear
	}

	return nil
}
