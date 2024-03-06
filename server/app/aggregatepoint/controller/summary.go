package controller

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/bartmika/timekit"
	ap_s "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/aggregatepoint/datastore"
	u_d "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/user/datastore"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/config/constants"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/utils/httperror"
)

type AggregatePointSummaryResponse struct {
	HeartRateThisDay        *ap_s.AggregatePoint `bson:"heart_rate_this_day" json:"heart_rate_this_day"`
	HeartRateLastDay        *ap_s.AggregatePoint `bson:"heart_rate_last_day" json:"heart_rate_last_day"`
	HeartRateThisISOWeek    *ap_s.AggregatePoint `bson:"heart_rate_this_iso_week" json:"heart_rate_this_iso_week"`
	HeartRateLastISOWeek    *ap_s.AggregatePoint `bson:"heart_rate_last_iso_week" json:"heart_rate_last_iso_week"`
	HeartRateThisMonth      *ap_s.AggregatePoint `bson:"heart_rate_this_month" json:"heart_rate_this_month"`
	HeartRateLastMonth      *ap_s.AggregatePoint `bson:"heart_rate_last_month" json:"heart_rate_last_month"`
	HeartRateThisYear       *ap_s.AggregatePoint `bson:"heart_rate_this_year" json:"heart_rate_this_year"`
	HeartRateLastYear       *ap_s.AggregatePoint `bson:"heart_rate_last_year" json:"heart_rate_last_year"`
	StepsCounterThisDay     *ap_s.AggregatePoint `bson:"steps_counter_this_day" json:"steps_counter_this_day"`
	StepsCounterLastDay     *ap_s.AggregatePoint `bson:"steps_counter_last_day" json:"steps_counter_last_day"`
	StepsCounterThisISOWeek *ap_s.AggregatePoint `bson:"steps_counter_this_iso_week" json:"steps_counter_this_iso_week"`
	StepsCounterLastISOWeek *ap_s.AggregatePoint `bson:"steps_counter_last_iso_week" json:"steps_counter_last_iso_week"`
	StepsCounterThisMonth   *ap_s.AggregatePoint `bson:"steps_counter_this_month" json:"steps_counter_this_month"`
	StepsCounterLastMonth   *ap_s.AggregatePoint `bson:"steps_counter_last_month" json:"steps_counter_last_month"`
	StepsCounterThisYear    *ap_s.AggregatePoint `bson:"steps_counter_this_year" json:"steps_counter_this_year"`
	StepsCounterLastYear    *ap_s.AggregatePoint `bson:"steps_counter_last_year" json:"steps_counter_last_year"`
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
	case u_d.UserPrimaryHealthTrackingDeviceTypeFitBit:
		// Do nothing except continue execution of this function...
	default:
		impl.Logger.Error("user has unsupported health tracker", slog.String("user_id", uid.Hex()))
		return nil, httperror.NewForBadRequestWithSingleField("user_id", fmt.Sprintf("user has unsupported health tracker for type: %v", u.PrimaryHealthTrackingDeviceType))
	}

	// DEVELOPERS NOTE:
	// In the future use golang goroutines to improve performance.

	// Variable used to return a summary for all our data.
	res := &AggregatePointSummaryResponse{}

	////
	//// Heart Rate
	////

	if !u.PrimaryHealthTrackingDeviceHeartRateMetricID.IsZero() {
		// --- Today --- //
		thisDayStart := timekit.Midnight(time.Now)
		thisDayEnd := timekit.MidnightTomorrow(time.Now)
		thisDay, err := impl.AggregatePointStorer.GetByCompositeKey(ctx, u.PrimaryHealthTrackingDeviceHeartRateMetricID, ap_s.PeriodDay, thisDayStart, thisDayEnd)
		if err != nil {
			impl.Logger.Error("failed getting aggregate point by composite key",
				slog.String("metric_id", u.PrimaryHealthTrackingDeviceHeartRateMetricID.Hex()),
				slog.Time("start", thisDayStart),
				slog.Time("end", thisDayEnd),
				slog.Any("error", err))
			return nil, err
		}
		if thisDay != nil {
			res.HeartRateThisDay = thisDay
		}

		// --- Last Day --- //
		lastDayStart := timekit.MidnightYesterday(time.Now)
		lastDayEnd := timekit.Midnight(time.Now)
		lastDay, err := impl.AggregatePointStorer.GetByCompositeKey(ctx, u.PrimaryHealthTrackingDeviceHeartRateMetricID, ap_s.PeriodDay, lastDayStart, lastDayEnd)
		if err != nil {
			impl.Logger.Error("failed getting aggregate point by composite key",
				slog.String("metric_id", u.PrimaryHealthTrackingDeviceHeartRateMetricID.Hex()),
				slog.Time("start", lastDayStart),
				slog.Time("end", lastDayEnd),
				slog.Any("error", err))
			return nil, err
		}
		if lastDay != nil {
			res.HeartRateLastDay = lastDay
		}

		// --- This ISO Week --- //
		thisWeekStart := timekit.FirstDayOfThisISOWeek(time.Now)
		thisWeekEnd := timekit.LastDayOfThisISOWeek(time.Now)
		thisWeek, err := impl.AggregatePointStorer.GetByCompositeKey(ctx, u.PrimaryHealthTrackingDeviceHeartRateMetricID, ap_s.PeriodWeek, thisWeekStart, thisWeekEnd)
		if err != nil {
			impl.Logger.Error("failed getting aggregate point by composite key",
				slog.String("metric_id", u.PrimaryHealthTrackingDeviceHeartRateMetricID.Hex()),
				slog.Time("start", thisWeekStart),
				slog.Time("end", thisWeekEnd),
				slog.Any("error", err))
			return nil, err
		}
		if thisWeek != nil {
			res.HeartRateThisISOWeek = thisWeek
		}

		// --- Last ISO Week --- //
		lastWeekStart := timekit.FirstDayOfLastISOWeek(time.Now)
		lastWeekEnd := timekit.FirstDayOfThisISOWeek(time.Now)
		lastWeek, err := impl.AggregatePointStorer.GetByCompositeKey(ctx, u.PrimaryHealthTrackingDeviceHeartRateMetricID, ap_s.PeriodWeek, lastWeekStart, lastWeekEnd)
		if err != nil {
			impl.Logger.Error("failed getting aggregate point by composite key",
				slog.String("metric_id", u.PrimaryHealthTrackingDeviceHeartRateMetricID.Hex()),
				slog.Time("start", lastWeekStart),
				slog.Time("end", lastWeekEnd),
				slog.Any("error", err))
			return nil, err
		}
		if lastWeek != nil {
			res.HeartRateLastISOWeek = lastWeek
		}

		// --- This Month --- //
		thisMonthStart := timekit.FirstDayOfThisMonth(time.Now)
		thisMonthEnd := timekit.FirstDayOfNextMonth(time.Now)
		thisMonth, err := impl.AggregatePointStorer.GetByCompositeKey(ctx, u.PrimaryHealthTrackingDeviceHeartRateMetricID, ap_s.PeriodMonth, thisMonthStart, thisMonthEnd)
		if err != nil {
			impl.Logger.Error("failed getting aggregate point by composite key",
				slog.String("metric_id", u.PrimaryHealthTrackingDeviceHeartRateMetricID.Hex()),
				slog.Time("start", thisMonthStart),
				slog.Time("end", thisMonthEnd),
				slog.Any("error", err))
			return nil, err
		}
		if thisMonth != nil {
			res.HeartRateThisMonth = thisMonth
		}

		// --- Last Month --- //
		lastMonthStart := timekit.FirstDayOfLastMonth(time.Now)
		lastMonthEnd := timekit.FirstDayOfThisMonth(time.Now)
		lastMonth, err := impl.AggregatePointStorer.GetByCompositeKey(ctx, u.PrimaryHealthTrackingDeviceHeartRateMetricID, ap_s.PeriodMonth, lastMonthStart, lastMonthEnd)
		if err != nil {
			impl.Logger.Error("failed getting aggregate point by composite key",
				slog.String("metric_id", u.PrimaryHealthTrackingDeviceHeartRateMetricID.Hex()),
				slog.Time("start", lastMonthStart),
				slog.Time("end", lastMonthEnd),
				slog.Any("error", err))
			return nil, err
		}
		if lastMonth != nil {
			res.HeartRateLastMonth = lastMonth
		}

		// --- This Year --- //
		thisYearStart := timekit.FirstDayOfThisYear(time.Now)
		thisYearEnd := timekit.FirstDayOfNextYear(time.Now)
		thisYear, err := impl.AggregatePointStorer.GetByCompositeKey(ctx, u.PrimaryHealthTrackingDeviceHeartRateMetricID, ap_s.PeriodYear, thisYearStart, thisYearEnd)
		if err != nil {
			impl.Logger.Error("failed getting aggregate point by composite key",
				slog.String("metric_id", u.PrimaryHealthTrackingDeviceHeartRateMetricID.Hex()),
				slog.Time("start", thisYearStart),
				slog.Time("end", thisYearEnd),
				slog.Any("error", err))
			return nil, err
		}
		if thisYear != nil {
			res.HeartRateThisYear = thisYear
		}

		// --- Last Year --- //
		lastYearStart := timekit.FirstDayOfLastYear(time.Now)
		lastYearEnd := timekit.FirstDayOfThisYear(time.Now)
		lastYear, err := impl.AggregatePointStorer.GetByCompositeKey(ctx, u.PrimaryHealthTrackingDeviceHeartRateMetricID, ap_s.PeriodYear, lastYearStart, lastYearEnd)
		if err != nil {
			impl.Logger.Error("failed getting aggregate point by composite key",
				slog.String("metric_id", u.PrimaryHealthTrackingDeviceHeartRateMetricID.Hex()),
				slog.Time("start", lastYearStart),
				slog.Time("end", lastYearEnd),
				slog.Any("error", err))
			return nil, err
		}
		if lastYear != nil {
			res.HeartRateLastYear = lastYear
		}
	}

	////
	//// Step Count
	////

	if !u.PrimaryHealthTrackingDeviceStepsCountMetricID.IsZero() {
		// --- Today --- //
		thisDayStart := timekit.Midnight(time.Now)
		thisDayEnd := timekit.MidnightTomorrow(time.Now)
		thisDay, err := impl.AggregatePointStorer.GetByCompositeKey(ctx, u.PrimaryHealthTrackingDeviceStepsCountMetricID, ap_s.PeriodDay, thisDayStart, thisDayEnd)
		if err != nil {
			impl.Logger.Error("failed getting aggregate point by composite key",
				slog.String("metric_id", u.PrimaryHealthTrackingDeviceStepsCountMetricID.Hex()),
				slog.Time("start", thisDayStart),
				slog.Time("end", thisDayEnd),
				slog.Any("error", err))
			return nil, err
		}
		if thisDay != nil {
			res.StepsCounterThisDay = thisDay
		}

		// --- Last Day --- //
		lastDayStart := timekit.MidnightYesterday(time.Now)
		lastDayEnd := timekit.Midnight(time.Now)
		lastDay, err := impl.AggregatePointStorer.GetByCompositeKey(ctx, u.PrimaryHealthTrackingDeviceStepsCountMetricID, ap_s.PeriodDay, lastDayStart, lastDayEnd)
		if err != nil {
			impl.Logger.Error("failed getting aggregate point by composite key",
				slog.String("metric_id", u.PrimaryHealthTrackingDeviceStepsCountMetricID.Hex()),
				slog.Time("start", lastDayStart),
				slog.Time("end", lastDayEnd),
				slog.Any("error", err))
			return nil, err
		}
		if lastDay != nil {
			res.StepsCounterLastDay = lastDay
		}

		// --- This ISO Week --- //
		thisWeekStart := timekit.FirstDayOfThisISOWeek(time.Now)
		thisWeekEnd := timekit.FirstDayOfNextISOWeek(time.Now)
		thisWeek, err := impl.AggregatePointStorer.GetByCompositeKey(ctx, u.PrimaryHealthTrackingDeviceStepsCountMetricID, ap_s.PeriodWeek, thisWeekStart, thisWeekEnd)
		if err != nil {
			impl.Logger.Error("failed getting aggregate point by composite key",
				slog.String("metric_id", u.PrimaryHealthTrackingDeviceStepsCountMetricID.Hex()),
				slog.Time("start", thisWeekStart),
				slog.Time("end", thisWeekEnd),
				slog.Any("error", err))
			return nil, err
		}
		if thisWeek != nil {
			res.StepsCounterThisISOWeek = thisWeek
		}

		// --- Last ISO Week --- //
		lastWeekStart := timekit.FirstDayOfLastISOWeek(time.Now)
		lastWeekEnd := timekit.FirstDayOfThisISOWeek(time.Now)
		lastWeek, err := impl.AggregatePointStorer.GetByCompositeKey(ctx, u.PrimaryHealthTrackingDeviceStepsCountMetricID, ap_s.PeriodWeek, lastWeekStart, lastWeekEnd)
		if err != nil {
			impl.Logger.Error("failed getting aggregate point by composite key",
				slog.String("metric_id", u.PrimaryHealthTrackingDeviceStepsCountMetricID.Hex()),
				slog.Time("start", lastWeekStart),
				slog.Time("end", lastWeekEnd),
				slog.Any("error", err))
			return nil, err
		}
		if lastWeek != nil {
			res.StepsCounterLastISOWeek = lastWeek
		}

		// --- This Month --- //
		thisMonthStart := timekit.FirstDayOfThisMonth(time.Now)
		thisMonthEnd := timekit.FirstDayOfNextMonth(time.Now)
		thisMonth, err := impl.AggregatePointStorer.GetByCompositeKey(ctx, u.PrimaryHealthTrackingDeviceStepsCountMetricID, ap_s.PeriodMonth, thisMonthStart, thisMonthEnd)
		if err != nil {
			impl.Logger.Error("failed getting aggregate point by composite key",
				slog.String("metric_id", u.PrimaryHealthTrackingDeviceStepsCountMetricID.Hex()),
				slog.Time("start", thisMonthStart),
				slog.Time("end", thisMonthEnd),
				slog.Any("error", err))
			return nil, err
		}
		if thisMonth != nil {
			res.StepsCounterThisMonth = thisMonth
		}

		// --- Last Month --- //
		lastMonthStart := timekit.FirstDayOfLastMonth(time.Now)
		lastMonthEnd := timekit.FirstDayOfThisMonth(time.Now)
		lastMonth, err := impl.AggregatePointStorer.GetByCompositeKey(ctx, u.PrimaryHealthTrackingDeviceStepsCountMetricID, ap_s.PeriodMonth, lastMonthStart, lastMonthEnd)
		if err != nil {
			impl.Logger.Error("failed getting aggregate point by composite key",
				slog.String("metric_id", u.PrimaryHealthTrackingDeviceStepsCountMetricID.Hex()),
				slog.Time("start", lastMonthStart),
				slog.Time("end", lastMonthEnd),
				slog.Any("error", err))
			return nil, err
		}
		if lastMonth != nil {
			res.StepsCounterLastMonth = lastMonth
		}

		// --- This Year --- //
		thisYearStart := timekit.FirstDayOfThisYear(time.Now)
		thisYearEnd := timekit.FirstDayOfNextYear(time.Now)
		thisYear, err := impl.AggregatePointStorer.GetByCompositeKey(ctx, u.PrimaryHealthTrackingDeviceStepsCountMetricID, ap_s.PeriodYear, thisYearStart, thisYearEnd)
		if err != nil {
			impl.Logger.Error("failed getting aggregate point by composite key",
				slog.String("metric_id", u.PrimaryHealthTrackingDeviceStepsCountMetricID.Hex()),
				slog.Time("start", thisYearStart),
				slog.Time("end", thisYearEnd),
				slog.Any("error", err))
			return nil, err
		}
		if thisYear != nil {
			res.StepsCounterThisYear = thisYear
		}

		// --- Last Year --- //
		lastYearStart := timekit.FirstDayOfLastYear(time.Now)
		lastYearEnd := timekit.FirstDayOfThisYear(time.Now)
		lastYear, err := impl.AggregatePointStorer.GetByCompositeKey(ctx, u.PrimaryHealthTrackingDeviceStepsCountMetricID, ap_s.PeriodYear, lastYearStart, lastYearEnd)
		if err != nil {
			impl.Logger.Error("failed getting aggregate point by composite key",
				slog.String("metric_id", u.PrimaryHealthTrackingDeviceStepsCountMetricID.Hex()),
				slog.Time("start", lastYearStart),
				slog.Time("end", lastYearEnd),
				slog.Any("error", err))
			return nil, err
		}
		if lastYear != nil {
			res.StepsCounterLastYear = lastYear
		}
	}

	return res, nil
}
