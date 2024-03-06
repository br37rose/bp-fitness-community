package controller

import (
	"context"
	"log/slog"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/bartmika/timekit"
	dp_s "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/datapoint/datastore"
	fitbitapp_s "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/fitbitapp/datastore"
	u_d "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/user/datastore"
	u_s "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/user/datastore"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/config/constants"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/utils/httperror"
)

type FitBitSimulatorCreateRequestIDO struct {
	UserID             primitive.ObjectID `json:"user_id"`
	SimulatorAlgorithm string             `json:"simulator_algorithm"`
}

func (impl *FitBitAppControllerImpl) CreateSimulator(ctx context.Context, requestData *FitBitSimulatorCreateRequestIDO) (*fitbitapp_s.FitBitApp, error) {
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
		if uid != requestData.UserID {
			return nil, httperror.NewForForbiddenWithSingleField("message", "you do not have permission")
		}
	}

	u, err := impl.UserStorer.GetByID(ctx, requestData.UserID)
	if err != nil {
		impl.Logger.Error("failed getting user by id",
			// slog.Any("code", code),
			// slog.Any("fba", fba),
			// slog.Any("user_id", userID),
			slog.Any("error", err))
		return nil, err
	}

	// Create our simulator.
	fba := &fitbitapp_s.FitBitApp{
		ID:                 primitive.NewObjectID(),
		UserName:           u.Name,
		UserLexicalName:    u.LexicalName,
		UserID:             u.ID,
		Status:             fitbitapp_s.StatusActive,
		CreatedAt:          time.Now(),
		ModifiedAt:         time.Now(),
		OrganizationID:     u.OrganizationID,
		FitBitUserID:       "-",
		AuthType:           fitbitapp_s.AuthTypeOAuth2,
		Errors:             "",
		Scope:              "activity oxygen_saturation respiratory_rate cardio_fitness location temperature nutrition sleep heartrate electrocardiogram weight",
		TokenType:          "-",
		AccessToken:        "-",
		ExpiresIn:          9999999999,
		RefreshToken:       "-",
		ExpireTime:         time.Now().Add(time.Hour * time.Duration(999999)),
		LastFetchedAt:      time.Date(2014, 1, 1, 00, 00, 00, 000000000, time.UTC), // 2014-01-01 00:00:00.00 UTC
		HeartRateMetricID:  primitive.NewObjectID(),
		StepsCountMetricID: primitive.NewObjectID(),
		IsTestMode:         true,
		SimulatorAlgorithm: "random",
	}

	// Essentially run create or update function.
	if err := impl.FitBitAppStorer.Create(ctx, fba); err != nil {
		impl.Logger.Error("database create error",
			// slog.Any("code", code),
			// slog.Any("fba", fba),
			// slog.Any("user_id", userID),
			slog.Any("error", err))
		return nil, err
	}

	// Update our user with our new device.
	u.PrimaryHealthTrackingDeviceType = u_s.UserPrimaryHealthTrackingDeviceTypeFitBit
	u.PrimaryHealthTrackingDeviceHeartRateMetricID = fba.HeartRateMetricID
	u.PrimaryHealthTrackingDeviceStepsCountMetricID = fba.StepsCountMetricID
	u.FitBitAppID = fba.ID
	u.ModifiedAt = time.Now()
	if err := impl.UserStorer.UpdateByID(ctx, u); err != nil {
		impl.Logger.Error("failed updating user by id",
			// slog.Any("code", code),
			// slog.Any("user_id", userID),
			slog.Any("error", err))
		return nil, err
	}

	go func(fbapp *fitbitapp_s.FitBitApp) {
		// Populate fake data for simulator.
		start := time.Now().Add(time.Duration(-1) * time.Hour * 24 * 368)
		end := time.Now()
		ts := timekit.NewTimeStepper(start, end, 0, 0, 0, 0, 5, 0)

		var actual time.Time
		running := true
		for running {
			// Get the timestamp value we are on in the timestepper.
			actual = ts.Get()

			////
			//// ACTIVITY
			////

			// Process fitbit device activity data if user granted access to this data.
			if strings.Contains(fbapp.Scope, "activity") {
				var latestValue float64 = randomBetween(0, 100)

				// Generate our time-series datum.
				dp := &dp_s.DataPoint{
					ID:        primitive.NewObjectID(),
					MetricID:  fbapp.StepsCountMetricID,
					Timestamp: actual,
					Value:     latestValue,
				}
				if err := impl.DataPointStorer.Create(context.Background(), dp); err != nil {
					impl.Logger.Error("failed creating simulated data point",
						slog.Any("fitbit_app_id", fbapp.ID),
						slog.Any("error", err))
					return
				}
				// impl.Logger.Debug(fmt.Sprintf("created simulated data point #%s", dp.ID.Hex()),
				// 	slog.Any("metric_id", fbapp.StepsCountMetricID),
				// 	slog.Any("timestamp", actual),
				// 	slog.Any("value", latestValue),
				// )
			}

			////
			//// HEART RATE
			////

			if strings.Contains(fbapp.Scope, "heartrate") {
				var latestValue float64 = randomBetween(50, 190)

				// Generate our time-series datum.
				dp := &dp_s.DataPoint{
					ID:        primitive.NewObjectID(),
					MetricID:  fbapp.HeartRateMetricID,
					Timestamp: actual,
					Value:     latestValue,
				}
				if err := impl.DataPointStorer.Create(context.Background(), dp); err != nil {
					impl.Logger.Error("failed creating simulated data point",
						slog.Any("fitbit_app_id", fbapp.ID),
						slog.Any("error", err))
					return
				}
				// impl.Logger.Debug(fmt.Sprintf("created simulated data point #%s", dp.ID.Hex()),
				// 	slog.Any("metric_id", fbapp.HeartRateMetricID),
				// 	slog.Any("timestamp", actual),
				// 	slog.Any("value", latestValue),
				// )
			}

			// DEVELOPERS NOTE: Add more sensors here...

			// Run our timestepper to get our next value.
			ts.Next()

			running = ts.Done() == false
		}
		//
	}(fba)

	return fba, nil
}
