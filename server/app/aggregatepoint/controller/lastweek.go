package controller

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/bartmika/timekit"
	ap_s "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/aggregatepoint/datastore"
	gfa_ds "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/googlefitapp/datastore"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (impl *AggregatePointControllerImpl) AggregateLastISOWeekForAllActiveGoogleFitApps(ctx context.Context) error {
	res, err := impl.GoogleFitAppStorer.ListIDsByStatus(ctx, gfa_ds.StatusActive)
	if err != nil {
		impl.Logger.Error("failed listing by active status",
			slog.Any("error", err))
		return err
	}
	for _, gfaID := range res {
		// Lock this Google Fit App for modification.
		impl.Kmutex.Lockf("gfa_%v", gfaID.Hex())
		defer impl.Kmutex.Unlockf("gfa_%v", gfaID.Hex())

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

		start := timekit.FirstDayOfLastISOWeek(time.Now)
		end := timekit.FirstDayOfThisISOWeek(time.Now)

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
			DataTypeNameDistanceDelta
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
		metricIDs := []primitive.ObjectID{
			gfa.CaloriesBurnedMetricID,
			gfa.StepCountDeltaMetricID,
			gfa.HeartRateBPMMetricID,
		}
		for _, metricID := range metricIDs {
			if err := impl.aggregateForMetric(ctx, metricID, ap_s.PeriodWeek, start, end); err != nil {
				impl.Logger.Error("failed aggregating",
					slog.Any("google_fit_app_id", gfaID),
					slog.Any("metric_id", metricID),
					slog.Any("error", err))
			}
		}
	}
	return nil
}
