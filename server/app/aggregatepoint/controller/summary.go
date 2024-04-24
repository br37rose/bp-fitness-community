package controller

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"go.mongodb.org/mongo-driver/bson/primitive"

	ap_s "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/aggregatepoint/datastore"
	u_d "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/user/datastore"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/config/constants"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/utils/httperror"
)

type AggregatePointSummaryResponse struct {
	CaloriesBurnedThisDay     *ap_s.AggregatePoint `bson:"calories_burned_this_day" json:"calories_burned_this_day"`
	CaloriesBurnedLastDay     *ap_s.AggregatePoint `bson:"calories_burned_last_day" json:"calories_burned_last_day"`
	CaloriesBurnedThisISOWeek *ap_s.AggregatePoint `bson:"calories_burned_this_iso_week" json:"calories_burned_this_iso_week"`
	CaloriesBurnedLastISOWeek *ap_s.AggregatePoint `bson:"calories_burned_last_iso_week" json:"calories_burned_last_iso_week"`
	CaloriesBurnedThisMonth   *ap_s.AggregatePoint `bson:"calories_burned_this_month" json:"calories_burned_this_month"`
	CaloriesBurnedLastMonth   *ap_s.AggregatePoint `bson:"calories_burned_last_month" json:"calories_burned_last_month"`
	CaloriesBurnedThisYear    *ap_s.AggregatePoint `bson:"calories_burned_this_year" json:"calories_burned_this_year"`
	CaloriesBurnedLastYear    *ap_s.AggregatePoint `bson:"calories_burned_last_year" json:"calories_burned_last_year"`

	StepCountDeltaThisDay     *ap_s.AggregatePoint `bson:"step_count_delta_this_day" json:"step_count_delta_this_day"`
	StepCountDeltaLastDay     *ap_s.AggregatePoint `bson:"step_count_delta_last_day" json:"step_count_delta_last_day"`
	StepCountDeltaThisISOWeek *ap_s.AggregatePoint `bson:"step_count_delta_this_iso_week" json:"step_count_delta_this_iso_week"`
	StepCountDeltaLastISOWeek *ap_s.AggregatePoint `bson:"step_count_delta_last_iso_week" json:"step_count_delta_last_iso_week"`
	StepCountDeltaThisMonth   *ap_s.AggregatePoint `bson:"step_count_delta_this_month" json:"step_count_delta_this_month"`
	StepCountDeltaLastMonth   *ap_s.AggregatePoint `bson:"step_count_delta_last_month" json:"step_count_delta_last_month"`
	StepCountDeltaThisYear    *ap_s.AggregatePoint `bson:"step_count_delta_this_year" json:"step_count_delta_this_year"`
	StepCountDeltaLastYear    *ap_s.AggregatePoint `bson:"step_count_delta_last_year" json:"step_count_delta_last_year"`

	DistanceDeltaThisDay     *ap_s.AggregatePoint `bson:"distance_delta_this_day" json:"distance_delta_this_day"`
	DistanceDeltaLastDay     *ap_s.AggregatePoint `bson:"distance_delta_last_day" json:"distance_delta_last_day"`
	DistanceDeltaThisISOWeek *ap_s.AggregatePoint `bson:"distance_delta_this_iso_week" json:"distance_delta_this_iso_week"`
	DistanceDeltaLastISOWeek *ap_s.AggregatePoint `bson:"distance_delta_last_iso_week" json:"distance_delta_last_iso_week"`
	DistanceDeltaThisMonth   *ap_s.AggregatePoint `bson:"distance_delta_this_month" json:"distance_delta_this_month"`
	DistanceDeltaLastMonth   *ap_s.AggregatePoint `bson:"distance_delta_last_month" json:"distance_delta_last_month"`
	DistanceDeltaThisYear    *ap_s.AggregatePoint `bson:"distance_delta_this_year" json:"distance_delta_this_year"`
	DistanceDeltaLastYear    *ap_s.AggregatePoint `bson:"distance_delta_last_year" json:"distance_delta_last_year"`

	HeartRateThisDay     *ap_s.AggregatePoint `bson:"heart_rate_this_day" json:"heart_rate_this_day"`
	HeartRateLastDay     *ap_s.AggregatePoint `bson:"heart_rate_last_day" json:"heart_rate_last_day"`
	HeartRateThisISOWeek *ap_s.AggregatePoint `bson:"heart_rate_this_iso_week" json:"heart_rate_this_iso_week"`
	HeartRateLastISOWeek *ap_s.AggregatePoint `bson:"heart_rate_last_iso_week" json:"heart_rate_last_iso_week"`
	HeartRateThisMonth   *ap_s.AggregatePoint `bson:"heart_rate_this_month" json:"heart_rate_this_month"`
	HeartRateLastMonth   *ap_s.AggregatePoint `bson:"heart_rate_last_month" json:"heart_rate_last_month"`
	HeartRateThisYear    *ap_s.AggregatePoint `bson:"heart_rate_this_year" json:"heart_rate_this_year"`
	HeartRateLastYear    *ap_s.AggregatePoint `bson:"heart_rate_last_year" json:"heart_rate_last_year"`
}

func (impl *AggregatePointControllerImpl) GetSummary(ctx context.Context, userID primitive.ObjectID) (*AggregatePointSummaryResponse, error) {
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

	if userID.IsZero() {
		impl.Logger.Error("user_id missing value")
		return nil, httperror.NewForBadRequestWithSingleField("user_id", "missing value")
	}

	u, err := impl.UserStorer.GetByID(ctx, uid)
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

	// Defensive Code: If no primary device has been set then error.
	if u.PrimaryHealthTrackingDevice == nil {
		err := errors.New("no health tracker attached")
		impl.Logger.Error("no health tracker attached",
			slog.String("user_id", uid.Hex()),
			slog.Any("error", err))
		return nil, err
	}

	// Variable used to return a summary for all our data.
	res := &AggregatePointSummaryResponse{}

	////
	//// Summarization.
	////

	// The following if-conditionals will look into our database for the
	// specific records and return them to the user.

	// TODO: In the future use golang goroutines to improve performance.

	if !u.PrimaryHealthTrackingDevice.CaloriesBurnedMetricID.IsZero() {
		if err := impl.summarizeCaloriesBurned(ctx, u, res); err != nil {
			impl.Logger.Error("failed summarizing calories burned",
				slog.String("metric_id", u.PrimaryHealthTrackingDevice.CaloriesBurnedMetricID.Hex()),
				slog.Any("error", err))
			return nil, err
		}
	}

	if !u.PrimaryHealthTrackingDevice.StepCountDeltaMetricID.IsZero() {
		if err := impl.summarizeStepCountDelta(ctx, u, res); err != nil {
			impl.Logger.Error("failed summarizing step counter delta",
				slog.String("metric_id", u.PrimaryHealthTrackingDevice.StepCountDeltaMetricID.Hex()),
				slog.Any("error", err))
			return nil, err
		}
	}

	if !u.PrimaryHealthTrackingDevice.DistanceDeltaMetricID.IsZero() {
		if err := impl.summarizeDistanceDelta(ctx, u, res); err != nil {
			impl.Logger.Error("failed summarizing distance delta",
				slog.String("metric_id", u.PrimaryHealthTrackingDevice.DistanceDeltaMetricID.Hex()),
				slog.Any("error", err))
			return nil, err
		}
	}

	if !u.PrimaryHealthTrackingDevice.HeartRateBPMMetricID.IsZero() {
		if err := impl.summarizeHeartRateBPM(ctx, u, res); err != nil {
			impl.Logger.Error("failed summarizing heart rate (bpm)",
				slog.String("metric_id", u.PrimaryHealthTrackingDevice.HeartRateBPMMetricID.Hex()),
				slog.Any("error", err))
			return nil, err
		}
	}

	//TODO: Add more health sensors here...

	return res, nil
}
