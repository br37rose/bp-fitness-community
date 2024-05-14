package controller

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/bartmika/timekit"
	"go.mongodb.org/mongo-driver/bson/primitive"

	gcp_a "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/adapter/cloudprovider/google"
	ap_s "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/aggregatepoint/datastore"
	gfa_ds "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/googlefitapp/datastore"
)

func (impl *AggregatePointControllerImpl) AggregateForAllActiveGoogleFitApps(ctx context.Context) error {
	res, err := impl.GoogleFitAppStorer.ListIDsByStatus(ctx, gfa_ds.StatusActive)
	if err != nil {
		impl.Logger.Error("failed listing by active status",
			slog.Any("error", err))
		return err
	}
	for _, gfaID := range res {
		if err := impl.AggregateForGoogleFitAppID(ctx, gfaID); err != nil {
			impl.Logger.Error("failed aggregation",
				slog.Any("google_fit_app_id", gfaID),
				slog.Any("error", err))
			return err
		}

		return nil
	}

	impl.Logger.Debug("aggregation completed")
	return nil
}
func (impl *AggregatePointControllerImpl) AggregateForGoogleFitAppID(ctx context.Context, gfaID primitive.ObjectID) error {
	impl.DistributedMutex.Lockf(ctx, "gfa_%v", gfaID.Hex())
	defer impl.DistributedMutex.Unlockf(ctx, "gfa_%v", gfaID.Hex())

	gfa, err := impl.GoogleFitAppStorer.GetByID(ctx, gfaID)
	if err != nil {
		impl.Logger.Error("failed getting google fit app by id",
			slog.Any("google_fit_app_id", gfaID),
			slog.Any("error", err))
		return err
	}
	if gfa == nil {
		err := fmt.Errorf("google fit app does not exist for id %v", gfaID.Hex())
		impl.Logger.Error("google fit app does not exist",
			slog.Any("google_fit_app_id", gfaID),
			slog.Any("error", err))
		return err
	}

	if err := impl.aggregateForGoogleFitApp(ctx, gfa); err != nil {
		impl.Logger.Error("failed aggregation",
			slog.Any("google_fit_app_id", gfaID),
			slog.Any("error", err))
		return err
	}
	return nil
}

func (impl *AggregatePointControllerImpl) aggregateForGoogleFitApp(ctx context.Context, gfa *gfa_ds.GoogleFitApp) error {
	// DEVELOPERS NOTE:
	// The following code below will aggregate health tracker sensor data
	// for the following data types our device supports.

	/*
	   ################
	   DEVELOPERS NOTE:
	   ################
	   THE FOLLOWING CONSTANTS ARE THE HEALTH TRACKER SENSORS OUR CODE SUPPORTS WHICH
	   ARE MARKED WITH THE `[DONE]` TEXT.

	   - - - - - - - - - - - - - - - - - - - - - - - - - - -
	   DataTypeNameActivitySegment
	   DataTypeNameBasalMetabolicRate
	   DataTypeNameCaloriesBurned        [DONE]
	   DataTypeNameCyclingPedalingCadence
	   DataTypeNameCyclingPedalingCumulative
	   DataTypeNameHeartPoints
	   DataTypeNameMoveMinutes
	   DataTypeNamePower
	   DataTypeNameStepCountCadence
	   DataTypeNameStepCountDelta        [DONE]
	   DataTypeNameWorkout
	   - - - - - - - - - - - - - - - - - - - - - - - - - - -
	   DataTypeNameCyclingWheelRevolutionRPM
	   DataTypeNameCyclingWheelRevolutionCumulative
	   DataTypeNameDistanceDelta         [DONE]
	   DataTypeNameLocationSample
	   DataTypeNameSpeed
	   - - - - - - - - - - - - - - - - - - - - - - - - - - -
	   DataTypeNameHydration
	   DataTypeNameNutrition
	   - - - - - - - - - - - - - - - - - - - - - - - - - - -
	   DataTypeNameBloodGlucose
	   DataTypeNameBloodPressure
	   DataTypeNameBodyFatPercentage
	   DataTypeNameBodyTemperature
	   DataTypeNameCervicalMucus
	   DataTypeNameCervicalPosition
	   DataTypeNameHeartRateBPM     [DONE]
	   DataTypeNameHeight
	   DataTypeNameMenstruation
	   DataTypeNameOvulationTest
	   DataTypeNameOxygenSaturation
	   DataTypeNameSleep
	   DataTypeNameVaginalSpotting
	   DataTypeNameWeight
	   - - -
	*/

	// Variable defines all the biometric sensors we want to process for
	// this aggregation function.
	metricIDs := []primitive.ObjectID{
		gfa.CaloriesBurnedMetricID,
		gfa.StepCountDeltaMetricID,
		gfa.DistanceDeltaMetricID,
		gfa.HeartRateBPMMetricID,
		//TODO: Add more health sensors here...
	}

	metricDataTypes := map[primitive.ObjectID]string{
		gfa.CaloriesBurnedMetricID: gcp_a.DataTypeNameCaloriesBurned,
		gfa.StepCountDeltaMetricID: gcp_a.DataTypeNameStepCountDelta,
		gfa.DistanceDeltaMetricID:  gcp_a.DataTypeNameDistanceDelta,
		gfa.HeartRateBPMMetricID:   gcp_a.DataTypeNameHeartRateBPM,
	}

	// impl.Logger.Debug("aggregation starting...",
	// 	slog.String("gfa_id", gfa.ID.Hex()))

	// Variable stores the number of goroutines we expect to wait for. We
	// set value of `1` because we have the following functions we want to
	// process in the background as goroutines:
	// - this hour
	// - last hour
	// - today
	// - yesterday
	// - this iso week
	// - last iso week
	// - this month
	// - last month
	// - this year
	// - last last
	numWorkers := 10

	// Create a channel to collect errors from goroutines.
	errCh := make(chan error, numWorkers)

	// Variable used to synchronize all the go routines running in
	// background outside of this function.
	var wg sync.WaitGroup

	// // Variable used to lock / unlock access when the goroutines want to
	// // perform writes to our output response.
	// var mu sync.Mutex

	// Load up the number of workers our waitgroup will need to handle.
	wg.Add(numWorkers)

	go func() {
		// // Lock the mutex before accessing res
		// mu.Lock()
		// defer mu.Unlock()
		// impl.Logger.Debug("processing this hour")
		start, end := timekit.HourRangeForNow(time.Now)
		for _, metricID := range metricIDs {
			if err := impl.aggregateForMetric(ctx, metricID, metricDataTypes[metricID], gfa.UserID, ap_s.PeriodHour, start, end); err != nil {
				impl.Logger.Error("failed aggregating this hour",
					slog.Any("google_fit_app_id", gfa.ID),
					slog.Any("metric_id", metricID),
					slog.Any("error", err))
			}
		}
		wg.Done() // We are done this background task.
	}()
	go func() {
		// // Lock the mutex before accessing res
		// mu.Lock()
		// defer mu.Unlock()
		// impl.Logger.Debug("processing last hour")
		start, end := timekit.HourRangeForNow(time.Now)

		// Calculate last hours.
		start = start.Add((-1) * time.Hour)
		end = end.Add((-1) * time.Hour)
		for _, metricID := range metricIDs {
			if err := impl.aggregateForMetric(ctx, metricID, metricDataTypes[metricID], gfa.UserID, ap_s.PeriodHour, start, end); err != nil {
				impl.Logger.Error("failed aggregating last hour",
					slog.Any("google_fit_app_id", gfa.ID),
					slog.Any("metric_id", metricID),
					slog.Any("error", err))
			}
		}
		wg.Done() // We are done this background task.
	}()
	go func() {
		// // Lock the mutex before accessing res
		// mu.Lock()
		// defer mu.Unlock()
		// impl.Logger.Debug("processing today")
		start := timekit.Midnight(time.Now)
		end := timekit.MidnightTomorrow(time.Now)
		for _, metricID := range metricIDs {
			if err := impl.aggregateForMetric(ctx, metricID, metricDataTypes[metricID], gfa.UserID, ap_s.PeriodDay, start, end); err != nil {
				impl.Logger.Error("failed aggregating today",
					slog.Any("google_fit_app_id", gfa.ID),
					slog.Any("metric_id", metricID),
					slog.Any("error", err))
			}
		}
		wg.Done() // We are done this background task.
	}()
	go func() {
		// // Lock the mutex before accessing res
		// mu.Lock()
		// defer mu.Unlock()
		// impl.Logger.Debug("processing yesterday")
		start := timekit.MidnightYesterday(time.Now)
		end := timekit.Midnight(time.Now)
		for _, metricID := range metricIDs {
			if err := impl.aggregateForMetric(ctx, metricID, metricDataTypes[metricID], gfa.UserID, ap_s.PeriodDay, start, end); err != nil {
				impl.Logger.Error("failed aggregating yesterday",
					slog.Any("google_fit_app_id", gfa.ID),
					slog.Any("metric_id", metricID),
					slog.Any("error", err))
			}
		}
		wg.Done() // We are done this background task.
	}()
	go func() {
		// // Lock the mutex before accessing res
		// mu.Lock()
		// defer mu.Unlock()
		// impl.Logger.Debug("processing this iso week")
		start := timekit.FirstDayOfThisISOWeek(time.Now)
		end := timekit.FirstDayOfNextISOWeek(time.Now)
		for _, metricID := range metricIDs {
			if err := impl.aggregateForMetric(ctx, metricID, metricDataTypes[metricID], gfa.UserID, ap_s.PeriodWeek, start, end); err != nil {
				impl.Logger.Error("failed aggregating this iso week",
					slog.Any("google_fit_app_id", gfa.ID),
					slog.Any("metric_id", metricID),
					slog.Any("error", err))
			}
		}
		wg.Done() // We are done this background task.
	}()
	go func() {
		// // Lock the mutex before accessing res
		// mu.Lock()
		// defer mu.Unlock()
		// impl.Logger.Debug("processing last iso week")
		start := timekit.FirstDayOfLastISOWeek(time.Now)
		end := timekit.FirstDayOfThisISOWeek(time.Now)
		for _, metricID := range metricIDs {
			if err := impl.aggregateForMetric(ctx, metricID, metricDataTypes[metricID], gfa.UserID, ap_s.PeriodWeek, start, end); err != nil {
				impl.Logger.Error("failed aggregating last iso week",
					slog.Any("google_fit_app_id", gfa.ID),
					slog.Any("metric_id", metricID),
					slog.Any("error", err))
			}
		}
		wg.Done() // We are done this background task.
	}()
	go func() {
		// // Lock the mutex before accessing res
		// mu.Lock()
		// defer mu.Unlock()
		// impl.Logger.Debug("processing this month")
		start := timekit.FirstDayOfThisMonth(time.Now)
		end := timekit.FirstDayOfNextMonth(time.Now)
		for _, metricID := range metricIDs {
			if err := impl.aggregateForMetric(ctx, metricID, metricDataTypes[metricID], gfa.UserID, ap_s.PeriodMonth, start, end); err != nil {
				impl.Logger.Error("failed aggregating this month",
					slog.Any("google_fit_app_id", gfa.ID),
					slog.Any("metric_id", metricID),
					slog.Any("error", err))
			}
		}
		wg.Done() // We are done this background task.
	}()
	go func() {
		// // Lock the mutex before accessing res
		// mu.Lock()
		// defer mu.Unlock()
		// impl.Logger.Debug("processing last month")
		start := timekit.FirstDayOfLastMonth(time.Now)
		end := timekit.FirstDayOfThisMonth(time.Now)
		for _, metricID := range metricIDs {
			if err := impl.aggregateForMetric(ctx, metricID, metricDataTypes[metricID], gfa.UserID, ap_s.PeriodMonth, start, end); err != nil {
				impl.Logger.Error("failed aggregating last month",
					slog.Any("google_fit_app_id", gfa.ID),
					slog.Any("metric_id", metricID),
					slog.Any("error", err))
			}
		}
		wg.Done() // We are done this background task.
	}()
	go func() {
		// // Lock the mutex before accessing res
		// mu.Lock()
		// defer mu.Unlock()
		// impl.Logger.Debug("processing this year")
		start := timekit.FirstDayOfThisYear(time.Now)
		end := timekit.FirstDayOfNextYear(time.Now)
		for _, metricID := range metricIDs {
			if err := impl.aggregateForMetric(ctx, metricID, metricDataTypes[metricID], gfa.UserID, ap_s.PeriodYear, start, end); err != nil {
				impl.Logger.Error("failed aggregating this year",
					slog.Any("google_fit_app_id", gfa.ID),
					slog.Any("metric_id", metricID),
					slog.Any("error", err))
			}
		}
		wg.Done() // We are done this background task.
	}()
	go func() {
		// // Lock the mutex before accessing res
		// mu.Lock()
		// defer mu.Unlock()
		// impl.Logger.Debug("processing last year")
		start := timekit.FirstDayOfLastYear(time.Now)
		end := timekit.FirstDayOfThisYear(time.Now)
		for _, metricID := range metricIDs {
			if err := impl.aggregateForMetric(ctx, metricID, metricDataTypes[metricID], gfa.UserID, ap_s.PeriodYear, start, end); err != nil {
				impl.Logger.Error("failed aggregating last year",
					slog.Any("google_fit_app_id", gfa.ID),
					slog.Any("metric_id", metricID),
					slog.Any("error", err))
			}
		}
		wg.Done() // We are done this background task.
	}()

	// Create a goroutine to close the error channel when all workers are done
	go func() {
		wg.Wait()
		close(errCh)
	}()

	// Iterate over the error channel to collect any errors from workers
	for err := range errCh {
		impl.Logger.Error("failed executing in goroutine",
			slog.Any("error", err))
		return err
	}

	// impl.Logger.Debug("aggregation completed", slog.String("gfa_id", gfa.ID.Hex()))
	return nil
}
