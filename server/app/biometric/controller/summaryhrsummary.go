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

func (impl *BiometricControllerImpl) generateSummarySummaryForHR(sessCtx mongo.SessionContext, u *u_d.User, res *AggregatePointSummaryResponse, mu *sync.Mutex, wg *sync.WaitGroup) error {
	// Once this function has been completed (whether successfully or not) then
	// update the `WaitGroup` that this goroutine is finished.
	defer wg.Done()

	// --- This Hour --- //
	thisHourStart, thisHourEnd := timekit.HourRangeForNow(time.Now)
	thisHour, err := impl.AggregatePointStorer.GetByCompositeKey(sessCtx, u.PrimaryHealthTrackingDevice.HeartRateBPMMetricID, ap_s.PeriodHour, thisHourStart, thisHourEnd)
	if err != nil {
		impl.Logger.Error("failed getting aggregate point by composite key",
			slog.String("metric_id", u.PrimaryHealthTrackingDevice.HeartRateBPMMetricID.Hex()),
			slog.Time("start", thisHourStart),
			slog.Time("end", thisHourEnd),
			slog.Any("error", err))
		return err
	}
	// --- Last Hour --- //
	lastHourStart := thisHourStart.Add((-1) * time.Hour)
	lastHourEnd := thisHourEnd.Add((-1) * time.Hour)
	lastHour, err := impl.AggregatePointStorer.GetByCompositeKey(sessCtx, u.PrimaryHealthTrackingDevice.HeartRateBPMMetricID, ap_s.PeriodHour, lastHourStart, lastHourEnd)
	if err != nil {
		impl.Logger.Error("failed getting aggregate point by composite key",
			slog.String("metric_id", u.PrimaryHealthTrackingDevice.HeartRateBPMMetricID.Hex()),
			slog.Time("start", lastHourStart),
			slog.Time("end", lastHourEnd),
			slog.Any("error", err))
		return err
	}
	// --- Today --- //
	thisDayStart := timekit.Midnight(time.Now)
	thisDayEnd := timekit.MidnightTomorrow(time.Now)
	thisDay, err := impl.AggregatePointStorer.GetByCompositeKey(sessCtx, u.PrimaryHealthTrackingDevice.HeartRateBPMMetricID, ap_s.PeriodDay, thisDayStart, thisDayEnd)
	if err != nil {
		impl.Logger.Error("failed getting aggregate point by composite key",
			slog.String("metric_id", u.PrimaryHealthTrackingDevice.HeartRateBPMMetricID.Hex()),
			slog.Time("start", thisDayStart),
			slog.Time("end", thisDayEnd),
			slog.Any("error", err))
		return err
	}
	// --- Last Day --- //
	lastDayStart := timekit.MidnightYesterday(time.Now)
	lastDayEnd := timekit.Midnight(time.Now)
	lastDay, err := impl.AggregatePointStorer.GetByCompositeKey(sessCtx, u.PrimaryHealthTrackingDevice.HeartRateBPMMetricID, ap_s.PeriodDay, lastDayStart, lastDayEnd)
	if err != nil {
		impl.Logger.Error("failed getting aggregate point by composite key",
			slog.String("metric_id", u.PrimaryHealthTrackingDevice.HeartRateBPMMetricID.Hex()),
			slog.Time("start", lastDayStart),
			slog.Time("end", lastDayEnd),
			slog.Any("error", err))
		return err
	}
	// --- This ISO Week --- //
	thisWeekStart := timekit.FirstDayOfThisISOWeek(time.Now)
	thisWeekEnd := timekit.FirstDayOfNextISOWeek(time.Now)
	thisWeek, err := impl.AggregatePointStorer.GetByCompositeKey(sessCtx, u.PrimaryHealthTrackingDevice.HeartRateBPMMetricID, ap_s.PeriodWeek, thisWeekStart, thisWeekEnd)
	if err != nil {
		impl.Logger.Error("failed getting aggregate point by composite key",
			slog.String("metric_id", u.PrimaryHealthTrackingDevice.HeartRateBPMMetricID.Hex()),
			slog.Time("start", thisWeekStart),
			slog.Time("end", thisWeekEnd),
			slog.Any("error", err))
		return err
	}
	// --- Last ISO Week --- //
	lastWeekStart := timekit.FirstDayOfLastISOWeek(time.Now)
	lastWeekEnd := timekit.FirstDayOfThisISOWeek(time.Now)
	lastWeek, err := impl.AggregatePointStorer.GetByCompositeKey(sessCtx, u.PrimaryHealthTrackingDevice.HeartRateBPMMetricID, ap_s.PeriodWeek, lastWeekStart, lastWeekEnd)
	if err != nil {
		impl.Logger.Error("failed getting aggregate point by composite key",
			slog.String("metric_id", u.PrimaryHealthTrackingDevice.HeartRateBPMMetricID.Hex()),
			slog.Time("start", lastWeekStart),
			slog.Time("end", lastWeekEnd),
			slog.Any("error", err))
		return err
	}
	// --- This Month --- //
	thisMonthStart := timekit.FirstDayOfThisMonth(time.Now)
	thisMonthEnd := timekit.FirstDayOfNextMonth(time.Now)
	thisMonth, err := impl.AggregatePointStorer.GetByCompositeKey(sessCtx, u.PrimaryHealthTrackingDevice.HeartRateBPMMetricID, ap_s.PeriodMonth, thisMonthStart, thisMonthEnd)
	if err != nil {
		impl.Logger.Error("failed getting aggregate point by composite key",
			slog.String("metric_id", u.PrimaryHealthTrackingDevice.HeartRateBPMMetricID.Hex()),
			slog.Time("start", thisMonthStart),
			slog.Time("end", thisMonthEnd),
			slog.Any("error", err))
		return err
	}
	// --- Last Month --- //
	lastMonthStart := timekit.FirstDayOfLastMonth(time.Now)
	lastMonthEnd := timekit.FirstDayOfThisMonth(time.Now)
	lastMonth, err := impl.AggregatePointStorer.GetByCompositeKey(sessCtx, u.PrimaryHealthTrackingDevice.HeartRateBPMMetricID, ap_s.PeriodMonth, lastMonthStart, lastMonthEnd)
	if err != nil {
		impl.Logger.Error("failed getting aggregate point by composite key",
			slog.String("metric_id", u.PrimaryHealthTrackingDevice.HeartRateBPMMetricID.Hex()),
			slog.Time("start", lastMonthStart),
			slog.Time("end", lastMonthEnd),
			slog.Any("error", err))
		return err
	}
	// --- This Year --- //
	thisYearStart := timekit.FirstDayOfThisYear(time.Now)
	thisYearEnd := timekit.FirstDayOfNextYear(time.Now)
	thisYear, err := impl.AggregatePointStorer.GetByCompositeKey(sessCtx, u.PrimaryHealthTrackingDevice.HeartRateBPMMetricID, ap_s.PeriodYear, thisYearStart, thisYearEnd)
	if err != nil {
		impl.Logger.Error("failed getting aggregate point by composite key",
			slog.String("metric_id", u.PrimaryHealthTrackingDevice.HeartRateBPMMetricID.Hex()),
			slog.Time("start", thisYearStart),
			slog.Time("end", thisYearEnd),
			slog.Any("error", err))
		return err
	}
	// --- Last Year --- //
	lastYearStart := timekit.FirstDayOfLastYear(time.Now)
	lastYearEnd := timekit.FirstDayOfThisYear(time.Now)
	lastYear, err := impl.AggregatePointStorer.GetByCompositeKey(sessCtx, u.PrimaryHealthTrackingDevice.HeartRateBPMMetricID, ap_s.PeriodYear, lastYearStart, lastYearEnd)
	if err != nil {
		impl.Logger.Error("failed getting aggregate point by composite key",
			slog.String("metric_id", u.PrimaryHealthTrackingDevice.HeartRateBPMMetricID.Hex()),
			slog.Time("start", lastYearStart),
			slog.Time("end", lastYearEnd),
			slog.Any("error", err))
		return err
	}

	// Lock the mutex before accessing res
	mu.Lock()
	defer mu.Unlock()

	// --- This Hour --- //
	if thisHour != nil {
		res.HeartRateThisHourSummary = thisHour
	}
	// --- Last Hour --- //
	if lastHour != nil {
		res.HeartRateLastHourSummary = lastHour
	}
	// --- Today --- //
	if thisDay != nil {
		res.HeartRateThisDaySummary = thisDay
	}
	// --- Last Day --- //
	if lastDay != nil {
		res.HeartRateLastDaySummary = lastDay
	}
	// --- This ISO Week --- //
	if thisWeek != nil {
		res.HeartRateThisISOWeekSummary = thisWeek
	}
	// --- Last ISO Week --- //
	if lastWeek != nil {
		res.HeartRateLastISOWeekSummary = lastWeek
	}
	// --- This Month --- //
	if thisMonth != nil {
		res.HeartRateThisMonthSummary = thisMonth
	}
	// --- Last Month --- //
	if lastMonth != nil {
		res.HeartRateLastMonthSummary = lastMonth
	}
	// --- This Year --- //
	if thisYear != nil {
		res.HeartRateThisYearSummary = thisYear
	}
	// --- Last Year --- //
	if lastYear != nil {
		res.HeartRateLastYearSummary = lastYear
	}

	return nil
}
