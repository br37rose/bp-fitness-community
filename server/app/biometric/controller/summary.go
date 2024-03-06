package controller

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"go.mongodb.org/mongo-driver/bson/primitive"

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

	HeartRateThisDayRanking []*rp_s.RankPoint `bson:"heart_rate_this_day_ranking" json:"heart_rate_this_day_ranking"`
	// HeartRateLastDayRanking     []*rp_s.RankPoint `bson:"heart_rate_last_day_ranking" json:"heart_rate_last_day_ranking"`
	HeartRateThisISOWeekRanking []*rp_s.RankPoint `bson:"heart_rate_this_iso_week_ranking" json:"heart_rate_this_iso_week_ranking"`
	// HeartRateLastISOWeekRanking []*rp_s.RankPoint `bson:"heart_rate_last_iso_week_ranking" json:"heart_rate_last_iso_week_ranking"`
	HeartRateThisMonthRanking []*rp_s.RankPoint `bson:"heart_rate_this_month_ranking" json:"heart_rate_this_month_ranking"`
	// HeartRateLastMonthRanking   []*rp_s.RankPoint `bson:"heart_rate_last_month_ranking" json:"heart_rate_last_month_ranking"`
	HeartRateThisYearRanking []*rp_s.RankPoint `bson:"heart_rate_this_year_ranking" json:"heart_rate_this_year_ranking"`
	// HeartRateLastYearRanking    []*rp_s.RankPoint `bson:"heart_rate_last_year_ranking" json:"heart_rate_last_year_ranking"`

	StepsCounterThisDaySummary     *ap_s.AggregatePoint `bson:"steps_counter_this_day_summary" json:"steps_counter_this_day_summary"`
	StepsCounterLastDaySummary     *ap_s.AggregatePoint `bson:"steps_counter_last_day_summary" json:"steps_counter_last_day_summary"`
	StepsCounterThisISOWeekSummary *ap_s.AggregatePoint `bson:"steps_counter_this_iso_week_summary" json:"steps_counter_this_iso_week_summary"`
	StepsCounterLastISOWeekSummary *ap_s.AggregatePoint `bson:"steps_counter_last_iso_week_summary" json:"steps_counter_last_iso_week_summary"`
	StepsCounterThisMonthSummary   *ap_s.AggregatePoint `bson:"steps_counter_this_month_summary" json:"steps_counter_this_month_summary"`
	StepsCounterLastMonthSummary   *ap_s.AggregatePoint `bson:"steps_counter_last_month_summary" json:"steps_counter_last_month_summary"`
	StepsCounterThisYearSummary    *ap_s.AggregatePoint `bson:"steps_counter_this_year_summary" json:"steps_counter_this_year_summary"`
	StepsCounterLastYearSummary    *ap_s.AggregatePoint `bson:"steps_counter_last_year_summary" json:"steps_counter_last_year_summary"`

	StepsCounterThisDayData     []*ap_s.AggregatePoint `bson:"steps_counter_this_day_data" json:"steps_counter_this_day_data"`
	StepsCounterLastDayData     []*ap_s.AggregatePoint `bson:"steps_counter_last_day_data" json:"steps_counter_last_day_data"`
	StepsCounterThisISOWeekData []*ap_s.AggregatePoint `bson:"steps_counter_this_iso_week_data" json:"steps_counter_this_iso_week_data"`
	StepsCounterLastISOWeekData []*ap_s.AggregatePoint `bson:"steps_counter_last_iso_week_data" json:"steps_counter_last_iso_week_data"`
	StepsCounterThisMonthData   []*ap_s.AggregatePoint `bson:"steps_counter_this_month_data" json:"steps_counter_this_month_data"`
	StepsCounterLastMonthData   []*ap_s.AggregatePoint `bson:"steps_counter_last_month_data" json:"steps_counter_last_month_data"`
	StepsCounterThisYearData    []*ap_s.AggregatePoint `bson:"steps_counter_this_year_data" json:"steps_counter_this_year_data"`
	StepsCounterLastYearData    []*ap_s.AggregatePoint `bson:"steps_counter_last_year_data" json:"steps_counter_last_year_data"`

	StepsCounterThisDayRanking []*rp_s.RankPoint `bson:"steps_counter_this_day_ranking" json:"steps_counter_this_day_ranking"`
	// StepsCounterLastDayRanking     []*rp_s.RankPoint `bson:"steps_counter_last_day_ranking" json:"steps_counter_last_day_ranking"`
	StepsCounterThisISOWeekRanking []*rp_s.RankPoint `bson:"steps_counter_this_iso_week_ranking" json:"steps_counter_this_iso_week_ranking"`
	// StepsCounterLastISOWeekRanking []*rp_s.RankPoint `bson:"steps_counter_last_iso_week_ranking" json:"steps_counter_last_iso_week_ranking"`
	StepsCounterThisMonthRanking []*rp_s.RankPoint `bson:"steps_counter_this_month_ranking" json:"steps_counter_this_month_ranking"`
	// StepsCounterLastMonthRanking   []*rp_s.RankPoint `bson:"steps_counter_last_month_ranking" json:"steps_counter_last_month_ranking"`
	StepsCounterThisYearRanking []*rp_s.RankPoint `bson:"steps_counter_this_year_ranking" json:"steps_counter_this_year_ranking"`
	// StepsCounterLastYearRanking    []*rp_s.RankPoint `bson:"steps_counter_last_year_ranking" json:"steps_counter_last_year_ranking"`
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
	case u_d.UserPrimaryHealthTrackingDeviceTypeFitBit:
		// Do nothing except continue execution of this function...
	default:
		impl.Logger.Error("user has unsupported health tracker", slog.String("user_id", uid.Hex()))
		return nil, httperror.NewForBadRequestWithSingleField("user_id", fmt.Sprintf("user has unsupported health tracker for type: %v", u.PrimaryHealthTrackingDeviceType))
	}

	// TODO:
	// DEVELOPERS NOTE:
	// In the future use golang goroutines to improve performance. This is the
	// techdebt we will incur for now.

	// Variable used to return a summary for all our data.
	res := &AggregatePointSummaryResponse{}

	////
	//// Heart Rate
	////

	if !u.PrimaryHealthTrackingDeviceHeartRateMetricID.IsZero() {
		if err := impl.generateSummaryForHR(ctx, u, res); err != nil {
			impl.Logger.Error("failed generating my summary data for heart rate",
				slog.String("metric_id", u.PrimaryHealthTrackingDeviceHeartRateMetricID.Hex()),
				slog.Any("error", err))
			return nil, err
		}
		if err := impl.generateSummaryDataForHR(ctx, u, res); err != nil {
			impl.Logger.Error("failed generating my summary data for heart rate",
				slog.String("metric_id", u.PrimaryHealthTrackingDeviceHeartRateMetricID.Hex()),
				slog.Any("error", err))
			return nil, err
		}
		if err := impl.generateSummaryRankingsForHR(ctx, u, res); err != nil {
			impl.Logger.Error("failed generating my summary rankings for heart rate",
				slog.String("metric_id", u.PrimaryHealthTrackingDeviceHeartRateMetricID.Hex()),
				slog.Any("error", err))
			return nil, err
		}
	}

	////
	//// Step Count
	////

	if !u.PrimaryHealthTrackingDeviceStepsCountMetricID.IsZero() {
		if err := impl.generateSummaryForStepCounter(ctx, u, res); err != nil {
			impl.Logger.Error("failed generating my summary rankings for steps counter",
				slog.String("metric_id", u.PrimaryHealthTrackingDeviceStepsCountMetricID.Hex()),
				slog.Any("error", err))
			return nil, err
		}
		if err := impl.generateSummaryDataForStepsCounter(ctx, u, res); err != nil {
			impl.Logger.Error("failed generating my summary rankings for steps counter",
				slog.String("metric_id", u.PrimaryHealthTrackingDeviceStepsCountMetricID.Hex()),
				slog.Any("error", err))
			return nil, err
		}
		if err := impl.generateSummaryRankingsForStepsCounter(ctx, u, res); err != nil {
			impl.Logger.Error("failed generating my summary rankings for steps counter",
				slog.String("metric_id", u.PrimaryHealthTrackingDeviceStepsCountMetricID.Hex()),
				slog.Any("error", err))
			return nil, err
		}
	}

	return res, nil
}
