package controller

import (
	"log/slog"
	"sync"
	"time"

	"github.com/bartmika/timekit"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	ap_s "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/aggregatepoint/datastore"
	u_d "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/user/datastore"
)

func (impl *BiometricControllerImpl) generateSummaryDataForStepsDelta(sessCtx mongo.SessionContext, u *u_d.User, res *AggregatePointSummaryResponse, mu *sync.Mutex, wg *sync.WaitGroup) error {
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
	wg.Add(8)

	// --- Today --- //
	go func() {
		if err := impl.generateSummaryDataForStepsDeltaToday(sessCtx, u, res, mu, wg); err != nil {
			//TODO
		}
	}()
	// --- Yesterday --- //
	go func() {
		if err := impl.generateSummaryDataForStepsDeltaYesterday(sessCtx, u, res, mu, wg); err != nil {
			//TODO
		}
	}()
	// --- This ISO Week --- //
	go func() {
		if err := impl.generateSummaryDataForStepsDeltaThisISOWeek(sessCtx, u, res, mu, wg); err != nil {
			//TODO
		}
	}()
	// --- Last ISO Week --- //
	go func() {
		if err := impl.generateSummaryDataForStepsDeltaLastISOWeek(sessCtx, u, res, mu, wg); err != nil {
			//TODO
		}
	}()
	// --- This Month --- //
	go func() {
		if err := impl.generateSummaryDataForStepsDeltaThisMonth(sessCtx, u, res, mu, wg); err != nil {
			//TODO
		}
	}()
	// --- Last Month --- //
	go func() {
		if err := impl.generateSummaryDataForStepsDeltaLastMonth(sessCtx, u, res, mu, wg); err != nil {
			//TODO
		}
	}()
	// --- This Year --- //
	go func() {
		if err := impl.generateSummaryDataForStepsDeltaThisYear(sessCtx, u, res, mu, wg); err != nil {
			//TODO
		}
	}()
	// --- Last Year --- //
	go func() {
		if err := impl.generateSummaryDataForStepsDeltaLastYear(sessCtx, u, res, mu, wg); err != nil {
			//TODO
		}
	}()

	return nil
}

func (impl *BiometricControllerImpl) generateSummaryDataForStepsDeltaToday(sessCtx mongo.SessionContext, u *u_d.User, res *AggregatePointSummaryResponse, mu *sync.Mutex, wg *sync.WaitGroup) error {
	// Once this function has been completed (whether successfully or not) then
	// update the `WaitGroup` that this goroutine is finished.
	defer wg.Done()

	// Fetch all aggregate data points for the today based on the `per hour`
	// period. This is due to the fact that the frontend graph requires posting
	// latest values per hour for the day; as a result, this code will work
	// to meet that requirement.
	thisDayStart := timekit.Midnight(time.Now)
	thisDayEnd := timekit.MidnightTomorrow(time.Now)
	thisDayFilter := &ap_s.AggregatePointPaginationListFilter{
		Cursor:    "",
		PageSize:  1_000_000_000, // Unlimited
		MetricIDs: []primitive.ObjectID{u.PrimaryHealthTrackingDevice.StepCountDeltaMetricID},
		Period:    ap_s.PeriodHour,
		SortField: "start",
		SortOrder: ap_s.SortOrderDescending,
		StartGTE:  thisDayStart,
		EndLTE:    thisDayEnd,
	}
	thisDayList, err := impl.AggregatePointStorer.ListByFilter(sessCtx, thisDayFilter)
	if err != nil {
		impl.Logger.Error("failed listing aggregate points",
			slog.String("metric_id", u.PrimaryHealthTrackingDevice.StepCountDeltaMetricID.Hex()),
			slog.Any("period", thisDayFilter.Period),
			slog.Any("sort_field", thisDayFilter.SortField),
			slog.Any("sort_order", thisDayFilter.SortOrder),
			slog.Time("start", thisDayStart),
			slog.Time("end", thisDayEnd),
			slog.Any("error", err))
		return err
	}

	// Lock the mutex before accessing res
	mu.Lock()
	defer mu.Unlock()

	if thisDayList != nil {
		// impl.Logger.Debug("debugging purposes only",
		// 	slog.String("metric_id", u.PrimaryHealthTrackingDeviceStepsCounterMetricID.Hex()),
		// 	slog.Any("period", ap_s.PeriodHour),
		// 	slog.Any("sort_field", thisDayFilter.SortField),
		// 	slog.Any("sort_order", thisDayFilter.SortOrder),
		// 	slog.Time("start", thisDayStart),
		// 	slog.Time("end", thisDayEnd))
		res.StepsCounterThisDayData = thisDayList.Results
	}

	return nil
}

func (impl *BiometricControllerImpl) generateSummaryDataForStepsDeltaYesterday(sessCtx mongo.SessionContext, u *u_d.User, res *AggregatePointSummaryResponse, mu *sync.Mutex, wg *sync.WaitGroup) error {
	// Once this function has been completed (whether successfully or not) then
	// update the `WaitGroup` that this goroutine is finished.
	defer wg.Done()

	lastDayStart := timekit.MidnightYesterday(time.Now)
	lastDayEnd := timekit.Midnight(time.Now)
	lastDayFilter := &ap_s.AggregatePointPaginationListFilter{
		Cursor:    "",
		PageSize:  1_000_000_000, // Unlimited
		MetricIDs: []primitive.ObjectID{u.PrimaryHealthTrackingDevice.StepCountDeltaMetricID},
		Period:    ap_s.PeriodHour,
		SortField: "start",
		SortOrder: ap_s.SortOrderDescending,
		StartGTE:  lastDayStart,
		EndLTE:    lastDayEnd,
	}
	lastDayList, err := impl.AggregatePointStorer.ListByFilter(sessCtx, lastDayFilter)
	if err != nil {
		impl.Logger.Error("failed listing aggregate points",
			slog.String("metric_id", u.PrimaryHealthTrackingDevice.StepCountDeltaMetricID.Hex()),
			slog.Any("period", lastDayFilter.Period),
			slog.Any("sort_field", lastDayFilter.SortField),
			slog.Any("sort_order", lastDayFilter.SortOrder),
			slog.Time("start", lastDayStart),
			slog.Time("end", lastDayEnd),
			slog.Any("error", err))
		return err
	}

	// Lock the mutex before accessing res
	mu.Lock()
	defer mu.Unlock()

	if lastDayList != nil {
		// impl.Logger.Debug("debugging purposes only",
		// 	slog.String("metric_id", u.PrimaryHealthTrackingDeviceStepsCounterMetricID.Hex()),
		// 	slog.Any("period", ap_s.PeriodHour),
		// 	slog.Any("sort_field", lastDayList.SortField),
		// 	slog.Any("sort_order", lastDayList.SortOrder),
		// 	slog.Time("start", lastDayStart),
		// 	slog.Time("end", lastDayEnd))
		res.StepsCounterLastDayData = lastDayList.Results
	}
	return nil
}

func (impl *BiometricControllerImpl) generateSummaryDataForStepsDeltaThisISOWeek(sessCtx mongo.SessionContext, u *u_d.User, res *AggregatePointSummaryResponse, mu *sync.Mutex, wg *sync.WaitGroup) error {
	// Once this function has been completed (whether successfully or not) then
	// update the `WaitGroup` that this goroutine is finished.
	defer wg.Done()

	// --- This ISO Week --- //
	// For returning week aggregate points, we want to return the data for
	// the week in a `per day` basis. So for example the frontend graph will
	// show the values for Monday, then show all the values for Tuesday, etc.
	// We are doing this to meet the requirements of the GUI which wants per
	// day values.
	thisISOWeekStart := timekit.FirstDayOfThisISOWeek(time.Now)
	thisISOWeekEnd := timekit.FirstDayOfNextISOWeek(time.Now)
	thisISOWeekFilter := &ap_s.AggregatePointPaginationListFilter{
		Cursor:    "",
		PageSize:  1_000_000_000, // Unlimited
		MetricIDs: []primitive.ObjectID{u.PrimaryHealthTrackingDevice.StepCountDeltaMetricID},
		Period:    ap_s.PeriodDay,
		SortField: "start",
		SortOrder: ap_s.SortOrderDescending,
		StartGTE:  thisISOWeekStart,
		EndLTE:    thisISOWeekEnd,
	}
	thisISOWeekList, err := impl.AggregatePointStorer.ListByFilter(sessCtx, thisISOWeekFilter)
	if err != nil {
		impl.Logger.Error("failed listing aggregate points",
			slog.String("metric_id", u.PrimaryHealthTrackingDevice.StepCountDeltaMetricID.Hex()),
			slog.Any("period", thisISOWeekFilter.Period),
			slog.Any("sort_field", thisISOWeekFilter.SortField),
			slog.Any("sort_order", thisISOWeekFilter.SortOrder),
			slog.Time("start", thisISOWeekStart),
			slog.Time("end", thisISOWeekEnd),
			slog.Any("error", err))
		return err
	}

	// Lock the mutex before accessing res
	mu.Lock()
	defer mu.Unlock()

	if thisISOWeekList != nil {
		// impl.Logger.Debug("debugging purposes only",
		// 	slog.String("metric_id", u.PrimaryHealthTrackingDeviceStepsCounterMetricID.Hex()),
		// 	slog.Any("period", ap_s.PeriodHour),
		// 	slog.Any("sort_field", thisISOWeekFilter.SortField),
		// 	slog.Any("sort_order", thisISOWeekFilter.SortOrder),
		// 	slog.Time("start", thisISOWeekStart),
		// 	slog.Time("end", thisISOWeekEnd))
		res.StepsCounterThisISOWeekData = thisISOWeekList.Results
	}

	return nil
}

func (impl *BiometricControllerImpl) generateSummaryDataForStepsDeltaLastISOWeek(sessCtx mongo.SessionContext, u *u_d.User, res *AggregatePointSummaryResponse, mu *sync.Mutex, wg *sync.WaitGroup) error {
	// Once this function has been completed (whether successfully or not) then
	// update the `WaitGroup` that this goroutine is finished.
	defer wg.Done()

	lastISOWeekWeekStart := timekit.FirstDayOfLastISOWeek(time.Now)
	lastISOWeekWeekEnd := timekit.FirstDayOfThisISOWeek(time.Now)
	lastISOWeekWeekFilter := &ap_s.AggregatePointPaginationListFilter{
		Cursor:    "",
		PageSize:  1_000_000_000, // Unlimited
		MetricIDs: []primitive.ObjectID{u.PrimaryHealthTrackingDevice.StepCountDeltaMetricID},
		Period:    ap_s.PeriodDay,
		SortField: "start",
		SortOrder: ap_s.SortOrderDescending,
		StartGTE:  lastISOWeekWeekStart,
		EndLTE:    lastISOWeekWeekEnd,
	}
	lastISOWeekWeekList, err := impl.AggregatePointStorer.ListByFilter(sessCtx, lastISOWeekWeekFilter)
	if err != nil {
		impl.Logger.Error("failed listing aggregate points",
			slog.String("metric_id", u.PrimaryHealthTrackingDevice.StepCountDeltaMetricID.Hex()),
			slog.Any("period", lastISOWeekWeekFilter.Period),
			slog.Any("sort_field", lastISOWeekWeekFilter.SortField),
			slog.Any("sort_order", lastISOWeekWeekFilter.SortOrder),
			slog.Time("start", lastISOWeekWeekStart),
			slog.Time("end", lastISOWeekWeekEnd),
			slog.Any("error", err))
		return err
	}

	// Lock the mutex before accessing res
	mu.Lock()
	defer mu.Unlock()

	if lastISOWeekWeekList != nil {
		// impl.Logger.Debug("debugging purposes only",
		// 	slog.String("metric_id", u.PrimaryHealthTrackingDeviceStepsCounterMetricID.Hex()),
		// 	slog.Any("period", ap_s.PeriodHour),
		// 	slog.Any("sort_field", lastISOWeekWeekFilter.SortField),
		// 	slog.Any("sort_order", lastISOWeekWeekFilter.SortOrder),
		// 	slog.Time("start", lastISOWeekWeekStart),
		// 	slog.Time("end", lastISOWeekWeekEnd))
		res.StepsCounterLastISOWeekData = lastISOWeekWeekList.Results
	}

	return nil
}

func (impl *BiometricControllerImpl) generateSummaryDataForStepsDeltaThisMonth(sessCtx mongo.SessionContext, u *u_d.User, res *AggregatePointSummaryResponse, mu *sync.Mutex, wg *sync.WaitGroup) error {
	// Once this function has been completed (whether successfully or not) then
	// update the `WaitGroup` that this goroutine is finished.
	defer wg.Done()

	// --- This Month --- //
	thisMonthStart := timekit.FirstDayOfThisMonth(time.Now)
	thisMonthEnd := timekit.FirstDayOfNextMonth(time.Now)
	thisMonthFilter := &ap_s.AggregatePointPaginationListFilter{
		Cursor:    "",
		PageSize:  1_000_000_000, // Unlimited
		MetricIDs: []primitive.ObjectID{u.PrimaryHealthTrackingDevice.StepCountDeltaMetricID},
		Period:    ap_s.PeriodDay,
		SortField: "start",
		SortOrder: ap_s.SortOrderDescending,
		StartGTE:  thisMonthStart,
		EndLTE:    thisMonthEnd,
	}
	thisMonthList, err := impl.AggregatePointStorer.ListByFilter(sessCtx, thisMonthFilter)
	if err != nil {
		impl.Logger.Error("failed listing aggregate points",
			slog.String("metric_id", u.PrimaryHealthTrackingDevice.StepCountDeltaMetricID.Hex()),
			slog.Any("period", thisMonthFilter.Period),
			slog.Any("sort_field", thisMonthFilter.SortField),
			slog.Any("sort_order", thisMonthFilter.SortOrder),
			slog.Time("start", thisMonthStart),
			slog.Time("end", thisMonthEnd),
			slog.Any("error", err))
		return err
	}

	// Lock the mutex before accessing res
	mu.Lock()
	defer mu.Unlock()

	if thisMonthList != nil {
		// impl.Logger.Debug("debugging purposes only",
		// 	slog.String("metric_id", u.PrimaryHealthTrackingDeviceStepsCounterMetricID.Hex()),
		// 	slog.Any("period", thisMonthFilter.Period),
		// 	slog.Any("sort_field", thisMonthFilter.SortField),
		// 	slog.Any("sort_order", thisMonthFilter.SortOrder),
		// 	slog.Time("start", thisMonthStart),
		// 	slog.Time("end", thisMonthEnd))
		res.StepsCounterThisMonthData = thisMonthList.Results
	}

	return nil
}

func (impl *BiometricControllerImpl) generateSummaryDataForStepsDeltaLastMonth(sessCtx mongo.SessionContext, u *u_d.User, res *AggregatePointSummaryResponse, mu *sync.Mutex, wg *sync.WaitGroup) error {
	// Once this function has been completed (whether successfully or not) then
	// update the `WaitGroup` that this goroutine is finished.
	defer wg.Done()

	// --- Last Month --- //
	lastMonthStart := timekit.FirstDayOfLastMonth(time.Now)
	lastMonthEnd := timekit.FirstDayOfThisMonth(time.Now)
	lastMonthFilter := &ap_s.AggregatePointPaginationListFilter{
		Cursor:    "",
		PageSize:  1_000_000_000, // Unlimited
		MetricIDs: []primitive.ObjectID{u.PrimaryHealthTrackingDevice.StepCountDeltaMetricID},
		Period:    ap_s.PeriodDay,
		SortField: "start",
		SortOrder: ap_s.SortOrderDescending,
		StartGTE:  lastMonthStart,
		EndLTE:    lastMonthEnd,
	}
	lastMonthList, err := impl.AggregatePointStorer.ListByFilter(sessCtx, lastMonthFilter)
	if err != nil {
		impl.Logger.Error("failed listing aggregate points",
			slog.String("metric_id", u.PrimaryHealthTrackingDevice.StepCountDeltaMetricID.Hex()),
			slog.Any("period", lastMonthFilter.Period),
			slog.Any("sort_field", lastMonthFilter.SortField),
			slog.Any("sort_order", lastMonthFilter.SortOrder),
			slog.Time("start", lastMonthStart),
			slog.Time("end", lastMonthEnd),
			slog.Any("error", err))
		return err
	}

	// Lock the mutex before accessing res
	mu.Lock()
	defer mu.Unlock()

	if lastMonthList != nil {
		// impl.Logger.Debug("debugging purposes only",
		// 	slog.String("metric_id", u.PrimaryHealthTrackingDeviceStepsCounterMetricID.Hex()),
		// 	slog.Any("period", ap_s.PeriodHour),
		// 	slog.Any("sort_field", lastMonthFilter.SortField),
		// 	slog.Any("sort_order", lastMonthFilter.SortOrder),
		// 	slog.Time("start", lastMonthStart),
		// 	slog.Time("end", lastMonthEnd))
		res.StepsCounterLastMonthData = lastMonthList.Results
	}

	return nil
}

func (impl *BiometricControllerImpl) generateSummaryDataForStepsDeltaThisYear(sessCtx mongo.SessionContext, u *u_d.User, res *AggregatePointSummaryResponse, mu *sync.Mutex, wg *sync.WaitGroup) error {
	// Once this function has been completed (whether successfully or not) then
	// update the `WaitGroup` that this goroutine is finished.
	defer wg.Done()

	// --- This Year --- //
	// When the user sees the year, we want to provide summary data on an only
	// `per month` basis.
	thisYearStart := timekit.FirstDayOfThisYear(time.Now)
	thisYearEnd := timekit.FirstDayOfNextYear(time.Now)
	thisYearFilter := &ap_s.AggregatePointPaginationListFilter{
		Cursor:    "",
		PageSize:  1_000_000_000, // Unlimited
		MetricIDs: []primitive.ObjectID{u.PrimaryHealthTrackingDevice.StepCountDeltaMetricID},
		Period:    ap_s.PeriodMonth,
		SortField: "start",
		SortOrder: ap_s.SortOrderDescending,
		StartGTE:  thisYearStart,
		EndLTE:    thisYearEnd,
	}
	thisYearList, err := impl.AggregatePointStorer.ListByFilter(sessCtx, thisYearFilter)
	if err != nil {
		impl.Logger.Error("failed listing aggregate points",
			slog.String("metric_id", u.PrimaryHealthTrackingDevice.StepCountDeltaMetricID.Hex()),
			slog.Any("period", thisYearFilter.Period),
			slog.Any("sort_field", thisYearFilter.SortField),
			slog.Any("sort_order", thisYearFilter.SortOrder),
			slog.Time("start", thisYearStart),
			slog.Time("end", thisYearEnd),
			slog.Any("error", err))
		return err
	}

	// Lock the mutex before accessing res
	mu.Lock()
	defer mu.Unlock()

	if thisYearList != nil {
		// impl.Logger.Debug("debugging purposes only",
		// 	slog.String("metric_id", u.PrimaryHealthTrackingDeviceStepsCounterMetricID.Hex()),
		// 	slog.Any("period", ap_s.PeriodHour),
		// 	slog.Any("sort_field", thisYearFilter.SortField),
		// 	slog.Any("sort_order", thisYearFilter.SortOrder),
		// 	slog.Time("start", thisYearStart),
		// 	slog.Time("end", thisYearEnd))
		res.StepsCounterThisYearData = thisYearList.Results
	}

	return nil
}

func (impl *BiometricControllerImpl) generateSummaryDataForStepsDeltaLastYear(sessCtx mongo.SessionContext, u *u_d.User, res *AggregatePointSummaryResponse, mu *sync.Mutex, wg *sync.WaitGroup) error {
	// Once this function has been completed (whether successfully or not) then
	// update the `WaitGroup` that this goroutine is finished.
	defer wg.Done()

	// --- Last Year --- //
	lastYearStart := timekit.FirstDayOfThisYear(time.Now)
	lastYearEnd := timekit.FirstDayOfNextYear(time.Now)
	lastYearFilter := &ap_s.AggregatePointPaginationListFilter{
		Cursor:    "",
		PageSize:  1_000_000_000, // Unlimited
		MetricIDs: []primitive.ObjectID{u.PrimaryHealthTrackingDevice.StepCountDeltaMetricID},
		Period:    ap_s.PeriodMonth,
		SortField: "start",
		SortOrder: ap_s.SortOrderDescending,
		StartGTE:  lastYearStart,
		EndLTE:    lastYearEnd,
	}
	lastYearList, err := impl.AggregatePointStorer.ListByFilter(sessCtx, lastYearFilter)
	if err != nil {
		impl.Logger.Error("failed listing aggregate points",
			slog.String("metric_id", u.PrimaryHealthTrackingDevice.StepCountDeltaMetricID.Hex()),
			slog.Any("period", lastYearFilter.Period),
			slog.Any("sort_field", lastYearFilter.SortField),
			slog.Any("sort_order", lastYearFilter.SortOrder),
			slog.Time("start", lastYearStart),
			slog.Time("end", lastYearEnd),
			slog.Any("error", err))
		return err
	}

	// Lock the mutex before accessing res
	mu.Lock()
	defer mu.Unlock()

	if lastYearList != nil {
		// impl.Logger.Debug("debugging purposes only",
		// 	slog.String("metric_id", u.PrimaryHealthTrackingDeviceStepsCounterMetricID.Hex()),
		// 	slog.Any("period", ap_s.PeriodHour),
		// 	slog.Any("sort_field", lastYearFilter.SortField),
		// 	slog.Any("sort_order", lastYearFilter.SortOrder),
		// 	slog.Time("start", lastYearStart),
		// 	slog.Time("end", lastYearEnd))
		res.StepsCounterLastYearData = lastYearList.Results
	}

	return nil
}