package controller

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"sync"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	ap_s "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/aggregatepoint/datastore"
	rp_s "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/rankpoint/datastore"
	u_d "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/user/datastore"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/config/constants"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/utils/httperror"
)

type AggregatePointSummaryResponse struct {
	HeartRateThisHourSummary    *ap_s.AggregatePoint `bson:"heart_rate_this_hour_summary" json:"heart_rate_this_hour_summary"`
	HeartRateLastHourSummary    *ap_s.AggregatePoint `bson:"heart_rate_last_hour_summary" json:"heart_rate_last_hour_summary"`
	HeartRateThisDaySummary     *ap_s.AggregatePoint `bson:"heart_rate_this_day_summary" json:"heart_rate_this_day_summary"`
	HeartRateLastDaySummary     *ap_s.AggregatePoint `bson:"heart_rate_last_day_summary" json:"heart_rate_last_day_summary"`
	HeartRateThisISOWeekSummary *ap_s.AggregatePoint `bson:"heart_rate_this_iso_week_summary" json:"heart_rate_this_iso_week_summary"`
	HeartRateLastISOWeekSummary *ap_s.AggregatePoint `bson:"heart_rate_last_iso_week_summary" json:"heart_rate_last_iso_week_summary"`
	HeartRateThisMonthSummary   *ap_s.AggregatePoint `bson:"heart_rate_this_month_summary" json:"heart_rate_this_month_summary"`
	HeartRateLastMonthSummary   *ap_s.AggregatePoint `bson:"heart_rate_last_month_summary" json:"heart_rate_last_month_summary"`
	HeartRateThisYearSummary    *ap_s.AggregatePoint `bson:"heart_rate_this_year_summary" json:"heart_rate_this_year_summary"`
	HeartRateLastYearSummary    *ap_s.AggregatePoint `bson:"heart_rate_last_year_summary" json:"heart_rate_last_year_summary"`

	HeartRateThisDayData     []*ap_s.AggregatePoint `bson:"heart_rate_this_day_data" json:"heart_rate_this_day_data"`
	HeartRateLastDayData     []*ap_s.AggregatePoint `bson:"heart_rate_last_day_data" json:"heart_rate_last_day_data"`
	HeartRateThisISOWeekData []*ap_s.AggregatePoint `bson:"heart_rate_this_iso_week_data" json:"heart_rate_this_iso_week_data"`
	HeartRateLastISOWeekData []*ap_s.AggregatePoint `bson:"heart_rate_last_iso_week_data" json:"heart_rate_last_iso_week_data"`
	HeartRateThisMonthData   []*ap_s.AggregatePoint `bson:"heart_rate_this_month_data" json:"heart_rate_this_month_data"`
	HeartRateLastMonthData   []*ap_s.AggregatePoint `bson:"heart_rate_last_month_data" json:"heart_rate_last_month_data"`
	HeartRateThisYearData    []*ap_s.AggregatePoint `bson:"heart_rate_this_year_data" json:"heart_rate_this_year_data"`
	HeartRateLastYearData    []*ap_s.AggregatePoint `bson:"heart_rate_last_year_data" json:"heart_rate_last_year_data"`

	HeartRateThisDayRanking     []*rp_s.RankPoint `bson:"heart_rate_this_day_ranking" json:"heart_rate_this_day_ranking"`
	HeartRateThisISOWeekRanking []*rp_s.RankPoint `bson:"heart_rate_this_iso_week_ranking" json:"heart_rate_this_iso_week_ranking"`
	HeartRateThisMonthRanking   []*rp_s.RankPoint `bson:"heart_rate_this_month_ranking" json:"heart_rate_this_month_ranking"`
	HeartRateThisYearRanking    []*rp_s.RankPoint `bson:"heart_rate_this_year_ranking" json:"heart_rate_this_year_ranking"`
	// HeartRateLastDayRanking     []*rp_s.RankPoint `bson:"heart_rate_last_day_ranking" json:"heart_rate_last_day_ranking"`
	// HeartRateLastISOWeekRanking []*rp_s.RankPoint `bson:"heart_rate_last_iso_week_ranking" json:"heart_rate_last_iso_week_ranking"`
	// HeartRateLastMonthRanking   []*rp_s.RankPoint `bson:"heart_rate_last_month_ranking" json:"heart_rate_last_month_ranking"`
	// HeartRateLastYearRanking    []*rp_s.RankPoint `bson:"heart_rate_last_year_ranking" json:"heart_rate_last_year_ranking"`

	//-----------------------------------------------------------------------------------------------------------------------------------

	StepCountDeltaThisHourSummary    *ap_s.AggregatePoint `bson:"step_count_delta_this_hour_summary" json:"step_count_delta_this_hour_summary"`
	StepCountDeltaLastHourSummary    *ap_s.AggregatePoint `bson:"step_count_delta_last_hour_summary" json:"step_count_delta_last_hour_summary"`
	StepCountDeltaThisDaySummary     *ap_s.AggregatePoint `bson:"step_count_delta_this_day_summary" json:"step_count_delta_this_day_summary"`
	StepCountDeltaLastDaySummary     *ap_s.AggregatePoint `bson:"step_count_delta_last_day_summary" json:"step_count_delta_last_day_summary"`
	StepCountDeltaThisISOWeekSummary *ap_s.AggregatePoint `bson:"step_count_delta_this_iso_week_summary" json:"step_count_delta_this_iso_week_summary"`
	StepCountDeltaLastISOWeekSummary *ap_s.AggregatePoint `bson:"step_count_delta_last_iso_week_summary" json:"step_count_delta_last_iso_week_summary"`
	StepCountDeltaThisMonthSummary   *ap_s.AggregatePoint `bson:"step_count_delta_this_month_summary" json:"step_count_delta_this_month_summary"`
	StepCountDeltaLastMonthSummary   *ap_s.AggregatePoint `bson:"step_count_delta_last_month_summary" json:"step_count_delta_last_month_summary"`
	StepCountDeltaThisYearSummary    *ap_s.AggregatePoint `bson:"step_count_delta_this_year_summary" json:"step_count_delta_this_year_summary"`
	StepCountDeltaLastYearSummary    *ap_s.AggregatePoint `bson:"step_count_delta_last_year_summary" json:"step_count_delta_last_year_summary"`

	StepCountDeltaThisDayData     []*ap_s.AggregatePoint `bson:"step_count_delta_this_day_data" json:"step_count_delta_this_day_data"`
	StepCountDeltaLastDayData     []*ap_s.AggregatePoint `bson:"step_count_delta_last_day_data" json:"step_count_delta_last_day_data"`
	StepCountDeltaThisISOWeekData []*ap_s.AggregatePoint `bson:"step_count_delta_this_iso_week_data" json:"step_count_delta_this_iso_week_data"`
	StepCountDeltaLastISOWeekData []*ap_s.AggregatePoint `bson:"step_count_delta_last_iso_week_data" json:"step_count_delta_last_iso_week_data"`
	StepCountDeltaThisMonthData   []*ap_s.AggregatePoint `bson:"step_count_delta_this_month_data" json:"step_count_delta_this_month_data"`
	StepCountDeltaLastMonthData   []*ap_s.AggregatePoint `bson:"step_count_delta_last_month_data" json:"step_count_delta_last_month_data"`
	StepCountDeltaThisYearData    []*ap_s.AggregatePoint `bson:"step_count_delta_this_year_data" json:"step_count_delta_this_year_data"`
	StepCountDeltaLastYearData    []*ap_s.AggregatePoint `bson:"step_count_delta_last_year_data" json:"step_count_delta_last_year_data"`

	StepCountDeltaThisDayRanking     []*rp_s.RankPoint `bson:"step_count_delta_this_day_ranking" json:"step_count_delta_this_day_ranking"`
	StepCountDeltaThisISOWeekRanking []*rp_s.RankPoint `bson:"step_count_delta_this_iso_week_ranking" json:"step_count_delta_this_iso_week_ranking"`
	StepCountDeltaThisMonthRanking   []*rp_s.RankPoint `bson:"step_count_delta_this_month_ranking" json:"step_count_delta_this_month_ranking"`
	StepCountDeltaThisYearRanking    []*rp_s.RankPoint `bson:"step_count_delta_this_year_ranking" json:"step_count_delta_this_year_ranking"`
	// StepCountDeltaLastDayRanking     []*rp_s.RankPoint `bson:"step_count_delta_last_day_ranking" json:"step_count_delta_last_day_ranking"`
	// StepCountDeltaLastISOWeekRanking []*rp_s.RankPoint `bson:"step_count_delta_last_iso_week_ranking" json:"step_count_delta_last_iso_week_ranking"`
	// StepCountDeltaLastMonthRanking   []*rp_s.RankPoint `bson:"step_count_delta_last_month_ranking" json:"step_count_delta_last_month_ranking"`
	// StepCountDeltaLastYearRanking    []*rp_s.RankPoint `bson:"step_count_delta_last_year_ranking" json:"step_count_delta_last_year_ranking"`

	//-----------------------------------------------------------------------------------------------------------------------------------

	CaloriesBurnedThisHourSummary    *ap_s.AggregatePoint `bson:"calories_burned_this_hour_summary" json:"calories_burned_this_hour_summary"`
	CaloriesBurnedLastHourSummary    *ap_s.AggregatePoint `bson:"calories_burned_last_hour_summary" json:"calories_burned_last_hour_summary"`
	CaloriesBurnedThisDaySummary     *ap_s.AggregatePoint `bson:"calories_burned_this_day_summary" json:"calories_burned_this_day_summary"`
	CaloriesBurnedLastDaySummary     *ap_s.AggregatePoint `bson:"calories_burned_last_day_summary" json:"calories_burned_last_day_summary"`
	CaloriesBurnedThisISOWeekSummary *ap_s.AggregatePoint `bson:"calories_burned_this_iso_week_summary" json:"calories_burned_this_iso_week_summary"`
	CaloriesBurnedLastISOWeekSummary *ap_s.AggregatePoint `bson:"calories_burned_last_iso_week_summary" json:"calories_burned_last_iso_week_summary"`
	CaloriesBurnedThisMonthSummary   *ap_s.AggregatePoint `bson:"calories_burned_this_month_summary" json:"calories_burned_this_month_summary"`
	CaloriesBurnedLastMonthSummary   *ap_s.AggregatePoint `bson:"calories_burned_last_month_summary" json:"calories_burned_last_month_summary"`
	CaloriesBurnedThisYearSummary    *ap_s.AggregatePoint `bson:"calories_burned_this_year_summary" json:"calories_burned_this_year_summary"`
	CaloriesBurnedLastYearSummary    *ap_s.AggregatePoint `bson:"calories_burned_last_year_summary" json:"calories_burned_last_year_summary"`

	CaloriesBurnedThisDayData     []*ap_s.AggregatePoint `bson:"calories_burned_this_day_data" json:"calories_burned_this_day_data"`
	CaloriesBurnedLastDayData     []*ap_s.AggregatePoint `bson:"calories_burned_last_day_data" json:"calories_burned_last_day_data"`
	CaloriesBurnedThisISOWeekData []*ap_s.AggregatePoint `bson:"calories_burned_this_iso_week_data" json:"calories_burned_this_iso_week_data"`
	CaloriesBurnedLastISOWeekData []*ap_s.AggregatePoint `bson:"calories_burned_last_iso_week_data" json:"calories_burned_last_iso_week_data"`
	CaloriesBurnedThisMonthData   []*ap_s.AggregatePoint `bson:"calories_burned_this_month_data" json:"calories_burned_this_month_data"`
	CaloriesBurnedLastMonthData   []*ap_s.AggregatePoint `bson:"calories_burned_last_month_data" json:"calories_burned_last_month_data"`
	CaloriesBurnedThisYearData    []*ap_s.AggregatePoint `bson:"calories_burned_this_year_data" json:"calories_burned_this_year_data"`
	CaloriesBurnedLastYearData    []*ap_s.AggregatePoint `bson:"calories_burned_last_year_data" json:"calories_burned_last_year_data"`

	CaloriesBurnedThisDayRanking     []*rp_s.RankPoint `bson:"calories_burned_this_day_ranking" json:"calories_burned_this_day_ranking"`
	CaloriesBurnedThisISOWeekRanking []*rp_s.RankPoint `bson:"calories_burned_this_iso_week_ranking" json:"calories_burned_this_iso_week_ranking"`
	CaloriesBurnedThisMonthRanking   []*rp_s.RankPoint `bson:"calories_burned_this_month_ranking" json:"calories_burned_this_month_ranking"`
	CaloriesBurnedThisYearRanking    []*rp_s.RankPoint `bson:"calories_burned_this_year_ranking" json:"calories_burned_this_year_ranking"`
	// CaloriesBurnedLastDayRanking     []*rp_s.RankPoint `bson:"calories_burned_last_day_ranking" json:"calories_burned_last_day_ranking"`
	// CaloriesBurnedLastISOWeekRanking []*rp_s.RankPoint `bson:"calories_burned_last_iso_week_ranking" json:"calories_burned_last_iso_week_ranking"`
	// CaloriesBurnedLastMonthRanking   []*rp_s.RankPoint `bson:"calories_burned_last_month_ranking" json:"calories_burned_last_month_ranking"`
	// CaloriesBurnedLastYearRanking    []*rp_s.RankPoint `bson:"calories_burned_last_year_ranking" json:"calories_burned_last_year_ranking"`

	//-----------------------------------------------------------------------------------------------------------------------------------

	DistanceDeltaThisHourSummary    *ap_s.AggregatePoint `bson:"distance_delta_this_hour_summary" json:"distance_delta_this_hour_summary"`
	DistanceDeltaLastHourSummary    *ap_s.AggregatePoint `bson:"distance_delta_last_hour_summary" json:"distance_delta_last_hour_summary"`
	DistanceDeltaThisDaySummary     *ap_s.AggregatePoint `bson:"distance_delta_this_day_summary" json:"distance_delta_this_day_summary"`
	DistanceDeltaLastDaySummary     *ap_s.AggregatePoint `bson:"distance_delta_last_day_summary" json:"distance_delta_last_day_summary"`
	DistanceDeltaThisISOWeekSummary *ap_s.AggregatePoint `bson:"distance_delta_this_iso_week_summary" json:"distance_delta_this_iso_week_summary"`
	DistanceDeltaLastISOWeekSummary *ap_s.AggregatePoint `bson:"distance_delta_last_iso_week_summary" json:"distance_delta_last_iso_week_summary"`
	DistanceDeltaThisMonthSummary   *ap_s.AggregatePoint `bson:"distance_delta_this_month_summary" json:"distance_delta_this_month_summary"`
	DistanceDeltaLastMonthSummary   *ap_s.AggregatePoint `bson:"distance_delta_last_month_summary" json:"distance_delta_last_month_summary"`
	DistanceDeltaThisYearSummary    *ap_s.AggregatePoint `bson:"distance_delta_this_year_summary" json:"distance_delta_this_year_summary"`
	DistanceDeltaLastYearSummary    *ap_s.AggregatePoint `bson:"distance_delta_last_year_summary" json:"distance_delta_last_year_summary"`

	DistanceDeltaThisDayData     []*ap_s.AggregatePoint `bson:"distance_delta_this_day_data" json:"distance_delta_this_day_data"`
	DistanceDeltaLastDayData     []*ap_s.AggregatePoint `bson:"distance_delta_last_day_data" json:"distance_delta_last_day_data"`
	DistanceDeltaThisISOWeekData []*ap_s.AggregatePoint `bson:"distance_delta_this_iso_week_data" json:"distance_delta_this_iso_week_data"`
	DistanceDeltaLastISOWeekData []*ap_s.AggregatePoint `bson:"distance_delta_last_iso_week_data" json:"distance_delta_last_iso_week_data"`
	DistanceDeltaThisMonthData   []*ap_s.AggregatePoint `bson:"distance_delta_this_month_data" json:"distance_delta_this_month_data"`
	DistanceDeltaLastMonthData   []*ap_s.AggregatePoint `bson:"distance_delta_last_month_data" json:"distance_delta_last_month_data"`
	DistanceDeltaThisYearData    []*ap_s.AggregatePoint `bson:"distance_delta_this_year_data" json:"distance_delta_this_year_data"`
	DistanceDeltaLastYearData    []*ap_s.AggregatePoint `bson:"distance_delta_last_year_data" json:"distance_delta_last_year_data"`

	DistanceDeltaThisDayRanking     []*rp_s.RankPoint `bson:"distance_delta_this_day_ranking" json:"distance_delta_this_day_ranking"`
	DistanceDeltaThisISOWeekRanking []*rp_s.RankPoint `bson:"distance_delta_this_iso_week_ranking" json:"distance_delta_this_iso_week_ranking"`
	DistanceDeltaThisMonthRanking   []*rp_s.RankPoint `bson:"distance_delta_this_month_ranking" json:"distance_delta_this_month_ranking"`
	DistanceDeltaThisYearRanking    []*rp_s.RankPoint `bson:"distance_delta_this_year_ranking" json:"distance_delta_this_year_ranking"`
	// DistanceDeltaLastDayRanking     []*rp_s.RankPoint `bson:"distance_delta_last_day_ranking" json:"distance_delta_last_day_ranking"`
	// DistanceDeltaLastISOWeekRanking []*rp_s.RankPoint `bson:"distance_delta_last_iso_week_ranking" json:"distance_delta_last_iso_week_ranking"`
	// DistanceDeltaLastMonthRanking   []*rp_s.RankPoint `bson:"distance_delta_last_month_ranking" json:"distance_delta_last_month_ranking"`
	// DistanceDeltaLastYearRanking    []*rp_s.RankPoint `bson:"distance_delta_last_year_ranking" json:"distance_delta_last_year_ranking"`
}

func (impl *BiometricControllerImpl) GetSummary(ctx context.Context, userID primitive.ObjectID) (*AggregatePointSummaryResponse, error) {
	// Extract from our session the following data.
	urole := ctx.Value(constants.SessionUserRole).(int8)
	uid := ctx.Value(constants.SessionUserID).(primitive.ObjectID)
	// uname := ctx.Value(constants.SessionUserName).(string)
	// oid := ctx.Value(constants.SessionUserOrganizationID).(primitive.ObjectID)
	// oname := ctx.Value(constants.SessionUserOrganizationName).(string)

	switch urole { // Security.
	case u_d.UserRoleRoot:
		return nil, httperror.NewForForbiddenWithSingleField("message", "you did not saasify offer")
	case u_d.UserRoleTrainer:
		return nil, httperror.NewForForbiddenWithSingleField("message", "you do not have permission")
	case u_d.UserRoleMember:
		if uid != userID {
			return nil, httperror.NewForForbiddenWithSingleField("message", "you do not have permission")
		}
	}

	////
	//// Start the transaction.
	////

	session, err := impl.DbClient.StartSession()
	if err != nil {
		impl.Logger.Error("start session error",
			slog.Any("error", err))
		return nil, err
	}
	defer session.EndSession(ctx)

	// Define a transaction function with a series of operations
	transactionFunc := func(sessCtx mongo.SessionContext) (interface{}, error) {

		if userID.IsZero() {
			impl.Logger.Error("user_id missing value")
			return nil, httperror.NewForBadRequestWithSingleField("user_id", "missing value")
		}

		u, err := impl.UserStorer.GetByID(sessCtx, uid)
		if err != nil {
			impl.Logger.Error("failed getting user",
				slog.String("user_id", uid.Hex()),
				slog.Any("error", err))
			return nil, err
		}
		if u == nil {
			impl.Logger.Error("user does not exist", slog.String("user_id", uid.Hex()))
			return nil, httperror.NewForBadRequestWithSingleField("user_id", fmt.Sprintf("user does not exist for ID: %v", uid.Hex()))
		}
		switch u.PrimaryHealthTrackingDeviceType {
		case u_d.UserPrimaryHealthTrackingDeviceTypeNone:
			err := errors.New("no health tracker attached")
			impl.Logger.Error("no health tracker attached",
				slog.String("user_id", uid.Hex()),
				slog.Any("error", err))
			return nil, err
		case u_d.UserPrimaryHealthTrackingDeviceTypeGoogleFit:
			// Do nothing except continue execution of this function...
		default:
			impl.Logger.Error("user has unsupported health tracker", slog.String("user_id", uid.Hex()))
			return nil, httperror.NewForBadRequestWithSingleField("user_id", fmt.Sprintf("user has unsupported health tracker for type: %v", u.PrimaryHealthTrackingDeviceType))
		}

		// Variable stores the number of goroutines we expect to wait for. We
		// set value of `7` because we have the following sensors we want to
		// process in the background as goroutines:
		// - generateSummarySummaryForHR
		// - generateSummaryDataForHR
		// - generateSummaryRankingsForHR
		// - generateSummarySummaryForStepCountDelta
		// - generateSummaryDataForStepsDelta
		// - generateSummarySummaryForCaloriesBurned
		// - generateSummaryDataForCaloriesBurned
		// - generateSummarySummaryForDistanceDelta
		// - generateSummaryDataForDistanceDelta
		numWorkers := 9

		// Create a channel to collect errors from goroutines.
		errCh := make(chan error, numWorkers)

		// Variable used to synchronize all the go routines running in
		// background outside of this function.
		var wg sync.WaitGroup

		// Variable used to lock / unlock access when the goroutines want to
		// perform writes to our output response.
		var mu sync.Mutex

		// Load up the number of workers our waitgroup will need to handle.
		wg.Add(numWorkers)

		// Variable used to return a summary for all our data.
		res := &AggregatePointSummaryResponse{}

		// Execute the following functions:
		// ---> Heart Rate:
		go func() {
			if err := impl.generateSummarySummaryForHR(sessCtx, u, res, &mu, &wg); err != nil {
				errCh <- err
			}
		}()
		go func() {
			if err := impl.generateSummaryDataForHR(sessCtx, u, res, &mu, &wg); err != nil {
				errCh <- err
			}
		}()
		go func() {
			if err := impl.generateSummaryRankingsForHR(sessCtx, u, res, &mu, &wg); err != nil {
				errCh <- err
			}
		}()
		// ---> Step Counter:
		go func() {
			if err := impl.generateSummarySummaryForStepCountDelta(sessCtx, u, res, &mu, &wg); err != nil {
				errCh <- err
			}
		}()
		go func() {
			if err := impl.generateSummaryDataForStepsDelta(sessCtx, u, res, &mu, &wg); err != nil {
				errCh <- err
			}
		}()
		// ---> Calories Burned:
		go func() {
			if err := impl.generateSummarySummaryForCaloriesBurned(sessCtx, u, res, &mu, &wg); err != nil {
				errCh <- err
			}
		}()
		go func() {
			if err := impl.generateSummaryDataForCaloriesBurned(sessCtx, u, res, &mu, &wg); err != nil {
				errCh <- err
			}
		}()
		// go func() {
		// 	if err := impl.generateSummaryRankingsForHR(sessCtx, u, res, &mu, &wg); err != nil {
		// 		errCh <- err
		// 	}
		// }()
		// ---> Distance Delta:
		go func() {
			if err := impl.generateSummarySummaryForDistanceDelta(sessCtx, u, res, &mu, &wg); err != nil {
				errCh <- err
			}
		}()
		go func() {
			if err := impl.generateSummaryDataForDistanceDelta(sessCtx, u, res, &mu, &wg); err != nil {
				errCh <- err
			}
		}()
		// go func() {
		// 	if err := impl.generateSummaryRankingsForHR(sessCtx, u, res, &mu, &wg); err != nil {
		// 		errCh <- err
		// 	}
		// }()

		// Create a goroutine to close the error channel when all workers are done
		go func() {
			wg.Wait()
			close(errCh)
		}()

		// Iterate over the error channel to collect any errors from workers
		for err := range errCh {
			impl.Logger.Error("failed executing in goroutine",
				slog.Any("error", err))
			return nil, err
		}

		return res, nil
	}

	// Start a transaction
	res, err := session.WithTransaction(ctx, transactionFunc)
	if err != nil {
		impl.Logger.Error("session failed error",
			slog.Any("error", err))
		return nil, err
	}

	return res.(*AggregatePointSummaryResponse), nil
}
