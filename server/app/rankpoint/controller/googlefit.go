package controller

import (
	"context"
	"log/slog"
	"time"

	"github.com/bartmika/timekit"
	"go.mongodb.org/mongo-driver/bson/primitive"

	gcp_a "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/adapter/cloudprovider/google"
	gfa_ds "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/googlefitapp/datastore"
	rp_s "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/rankpoint/datastore"
)

func (impl *RankPointControllerImpl) GenerateGlobalRankingForTodayUsingActiveGoogleFitApps(ctx context.Context) error {
	f := &gfa_ds.GoogleFitAppListFilter{
		Cursor:    primitive.NilObjectID,
		PageSize:  1_000_000,
		SortField: "_id",
		SortOrder: -1,
		Status:    gfa_ds.StatusActive,
	}
	gfas, err := impl.GoogleFitAppStorer.ListByFilter(ctx, f)
	if err != nil {
		impl.Logger.Error("failed listing by active status",
			slog.Any("date_range", "today"),
			slog.Any("error", err))
		return err
	}
	start := timekit.Midnight(time.Now)
	end := timekit.MidnightTomorrow(time.Now)

	/*
		DataTypeKeyActivitySegment           = 1  // https://developers.google.com/fit/datatypes/activity
		DataTypeKeyBasalMetabolicRate        = 2  // https://developers.google.com/fit/datatypes/activity#basal_metabolic_rate_bmr
		DataTypeKeyCaloriesBurned            [DONE]
		DataTypeKeyCyclingPedalingCadence    = 4  // https://developers.google.com/fit/datatypes/activity#cycling_pedaling_cadence
		DataTypeKeyCyclingPedalingCumulative = 5  // https://developers.google.com/fit/datatypes/activity#cycling_pedaling_cumulative
		DataTypeKeyHeartPoints               = 6  // https://developers.google.com/fit/datatypes/activity#heart_points
		DataTypeKeyMoveMinutes               = 7  // https://developers.google.com/fit/datatypes/activity#move_minutes
		DataTypeKeyPower                     = 8  // https://developers.google.com/fit/datatypes/activity#power
		DataTypeKeyStepCountCadence          = 9  //https://developers.google.com/fit/datatypes/activity#step_count_cadence
		DataTypeKeyStepCountDelta            [DONE]
		DataTypeKeyWorkout                   = 11 //https://developers.google.com/fit/datatypes/activity#workout

		DataTypeKeyCyclingWheelRevolutionRPM        = 12 // https://developers.google.com/fit/datatypes/location#cycling_wheel_revolution_rpm
		DataTypeKeyCyclingWheelRevolutionCumulative = 13 // https://developers.google.com/fit/datatypes/location#cycling_wheel_revolution_cumulative
		DataTypeKeyDistanceDelta                    = 14 // https://developers.google.com/fit/datatypes/location#distance_delta
		DataTypeKeyLocationSample                   = 15 // https://developers.google.com/fit/datatypes/location#location_sample
		DataTypeKeySpeed                            = 16 // https://developers.google.com/fit/datatypes/location#speed

		DataTypeKeyHydration = 17 // https://developers.google.com/fit/datatypes/nutrition
		DataTypeKeyNutrition = 18 // https://developers.google.com/fit/datatypes/nutrition

		DataTypeKeyBloodGlucose      = 19 // https://developers.google.com/fit/datatypes/health#blood_glucose
		DataTypeKeyBloodPressure     = 20 // https://developers.google.com/fit/datatypes/health#blood_pressure
		DataTypeKeyBodyFatPercentage = 21 // https://developers.google.com/fit/datatypes/health#body_fat_percentage
		DataTypeKeyBodyTemperature   = 22 // https://developers.google.com/fit/datatypes/health#body_temperature
		DataTypeKeyCervicalMucus     = 23 // https://developers.google.com/fit/datatypes/health#cervical_mucus
		DataTypeKeyCervicalPosition  = 24 // https://developers.google.com/fit/datatypes/health#cervical_position
		DataTypeKeyHeartRateBPM      [DONE]
		DataTypeKeyHeight            = 26 // https://developers.google.com/fit/datatypes/health#height
		DataTypeKeyMenstruation      = 27 // https://developers.google.com/fit/datatypes/health#menstruation
		DataTypeKeyOvulationTest     = 28 // https://developers.google.com/fit/datatypes/health#ovulation_test
		DataTypeKeyOxygenSaturation  = 29 // https://developers.google.com/fit/datatypes/health#oxygen_saturation
		DataTypeKeySleep             = 30 // https://developers.google.com/fit/datatypes/health#sleep
		DataTypeKeyVaginalSpotting   = 31 // https://developers.google.com/fit/datatypes/health#vaginal_spotting
		DataTypeKeyWeight            = 32 // https://developers.google.com/fit/datatypes/health#weight
	*/
	metricTypes := []int8{
		gcp_a.DataTypeKeyCaloriesBurned,
		gcp_a.DataTypeKeyStepCountDelta,
		gcp_a.DataTypeKeyHeartRateBPM,
		//TODO: Add more health sensors here...
	}
	funcTypes := []int8{
		rp_s.FunctionAverage,
		rp_s.FunctionSum,
		rp_s.FunctionCount,
		rp_s.FunctionMin,
		rp_s.FunctionMax,
	}
	for _, metricType := range metricTypes {
		for _, funcType := range funcTypes {
			go func(list []*gfa_ds.GoogleFitApp, mt int8, ft int8, startDT time.Time, endDT time.Time) {
				if err := impl.processGlobalRanksForGoogleFitApps(context.Background(), list, mt, ft, rp_s.PeriodDay, startDT, endDT); err != nil {
					impl.Logger.Error("failed generating global rate ranking for day",
						slog.Any("date_range", "today"),
						slog.Any("mt", mt),
						slog.Any("ft", ft),
						slog.Any("error", err))
					return
				}
			}(gfas.Results, metricType, funcType, start, end)
		}
	}
	return nil
}

func (impl *RankPointControllerImpl) GenerateGlobalRankingForThisISOWeekUsingActiveGoogleFitApps(ctx context.Context) error {
	f := &gfa_ds.GoogleFitAppListFilter{
		Cursor:    primitive.NilObjectID,
		PageSize:  1_000_000,
		SortField: "_id",
		SortOrder: -1,
		Status:    gfa_ds.StatusActive,
	}
	gfas, err := impl.GoogleFitAppStorer.ListByFilter(ctx, f)
	if err != nil {
		impl.Logger.Error("failed listing by active status",
			slog.Any("date_range", "iso_week"),
			slog.Any("error", err))
		return err
	}
	start := timekit.FirstDayOfThisISOWeek(time.Now)
	end := timekit.FirstDayOfNextISOWeek(time.Now)

	metricTypes := []int8{
		gcp_a.DataTypeKeyCaloriesBurned,
		gcp_a.DataTypeKeyStepCountDelta,
		gcp_a.DataTypeKeyHeartRateBPM,
		//TODO: Add more health sensors here...
	}
	funcTypes := []int8{
		rp_s.FunctionAverage,
		rp_s.FunctionSum,
		rp_s.FunctionCount,
		rp_s.FunctionMin,
		rp_s.FunctionMax,
	}
	for _, metricType := range metricTypes {
		for _, funcType := range funcTypes {
			go func(list []*gfa_ds.GoogleFitApp, mt int8, ft int8, startDT time.Time, endDT time.Time) {
				if err := impl.processGlobalRanksForGoogleFitApps(context.Background(), list, mt, ft, rp_s.PeriodWeek, startDT, endDT); err != nil {
					impl.Logger.Error("failed generating global rate ranking for this iso week",
						slog.Any("date_range", "iso week"),
						slog.Any("mt", mt),
						slog.Any("ft", ft),
						slog.Any("error", err))
					return
				}
			}(gfas.Results, metricType, funcType, start, end)
		}
	}

	return nil
}

func (impl *RankPointControllerImpl) GenerateGlobalRankingForThisMonthUsingActiveGoogleFitApps(ctx context.Context) error {
	f := &gfa_ds.GoogleFitAppListFilter{
		Cursor:    primitive.NilObjectID,
		PageSize:  1_000_000,
		SortField: "_id",
		SortOrder: -1,
		Status:    gfa_ds.StatusActive,
	}
	gfas, err := impl.GoogleFitAppStorer.ListByFilter(ctx, f)
	if err != nil {
		impl.Logger.Error("failed listing by active status",
			slog.Any("date_range", "month"),
			slog.Any("error", err))
		return err
	}
	start := timekit.FirstDayOfThisMonth(time.Now)
	end := timekit.FirstDayOfNextMonth(time.Now)

	metricTypes := []int8{
		gcp_a.DataTypeKeyCaloriesBurned,
		gcp_a.DataTypeKeyStepCountDelta,
		gcp_a.DataTypeKeyHeartRateBPM,
		//TODO: Add more health sensors here...
	}
	funcTypes := []int8{
		rp_s.FunctionAverage,
		rp_s.FunctionSum,
		rp_s.FunctionCount,
		rp_s.FunctionMin,
		rp_s.FunctionMax,
	}
	for _, metricType := range metricTypes {
		for _, funcType := range funcTypes {
			go func(list []*gfa_ds.GoogleFitApp, mt int8, ft int8, startDT time.Time, endDT time.Time) {
				if err := impl.processGlobalRanksForGoogleFitApps(context.Background(), list, mt, ft, rp_s.PeriodMonth, startDT, endDT); err != nil {
					impl.Logger.Error("failed generating global rate ranking for this month",
						slog.Any("date_range", "month"),
						slog.Any("mt", mt),
						slog.Any("ft", ft),
						slog.Any("error", err))
					return
				}
			}(gfas.Results, metricType, funcType, start, end)
		}
	}

	return nil
}

func (impl *RankPointControllerImpl) GenerateGlobalRankingForThisYearUsingActiveGoogleFitApps(ctx context.Context) error {
	f := &gfa_ds.GoogleFitAppListFilter{
		Cursor:    primitive.NilObjectID,
		PageSize:  1_000_000,
		SortField: "_id",
		SortOrder: -1,
		Status:    gfa_ds.StatusActive,
	}
	gfas, err := impl.GoogleFitAppStorer.ListByFilter(ctx, f)
	if err != nil {
		impl.Logger.Error("failed listing by active status",
			slog.Any("date_range", "year"),
			slog.Any("error", err))
		return err
	}
	start := timekit.FirstDayOfThisYear(time.Now)
	end := timekit.FirstDayOfNextYear(time.Now)

	metricTypes := []int8{
		gcp_a.DataTypeKeyCaloriesBurned,
		gcp_a.DataTypeKeyStepCountDelta,
		gcp_a.DataTypeKeyHeartRateBPM,
		//TODO: Add more health sensors here...
	}
	funcTypes := []int8{
		rp_s.FunctionAverage,
		rp_s.FunctionSum,
		rp_s.FunctionCount,
		rp_s.FunctionMin,
		rp_s.FunctionMax,
	}
	for _, metricType := range metricTypes {
		for _, funcType := range funcTypes {
			go func(list []*gfa_ds.GoogleFitApp, mt int8, ft int8, startDT time.Time, endDT time.Time) {
				if err := impl.processGlobalRanksForGoogleFitApps(context.Background(), list, mt, ft, rp_s.PeriodYear, startDT, endDT); err != nil {
					impl.Logger.Error("failed generating global rate ranking for this year",
						slog.Any("date_range", "year"),
						slog.Any("mt", mt),
						slog.Any("ft", ft),
						slog.Any("error", err))
					return
				}
			}(gfas.Results, metricType, funcType, start, end)
		}
	}

	return nil
}
