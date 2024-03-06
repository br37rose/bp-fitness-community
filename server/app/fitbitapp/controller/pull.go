package controller

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	fitbitapp_s "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/fitbitapp/datastore"
	fitbitdatum_s "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/fitbitdatum/datastore"
)

func (c *FitBitAppControllerImpl) PullAllActiveDevices(ctx context.Context) error {
	c.Logger.Debug("pulling all active fitbit devices")
	ids, err := c.FitBitAppStorer.ListPhysicalIDsByStatus(ctx, fitbitapp_s.StatusActive)
	if err != nil {
		c.Logger.Error("failed listing ids (for physical devices) by status",
			slog.Any("error", err))
		return err
	}
	for _, fitbitAppID := range ids {
		if err := c.PullDevice(ctx, fitbitAppID); err != nil {
			c.Logger.Error("failed pulling data for fitbit device",
				slog.Any("fitbit_app_id", fitbitAppID),
				slog.Any("error", err))
			return err
		}
	}
	return nil
}

// func (c *FitBitAppControllerImpl) Pull(ctx context.Context, fitBitAppID primitive.ObjectID) error {
// function will get the latest raw data for the particular FitBit device
// we have successfully registered in our application.
func (c *FitBitAppControllerImpl) PullDevice(ctx context.Context, fitbitAppID primitive.ObjectID) error {
	c.Logger.Debug(fmt.Sprintf("pulling fitbit device ID %v", fitbitAppID.Hex()))
	defer c.Logger.Debug(fmt.Sprintf("pulled fitbit device ID %v", fitbitAppID.Hex()))

	// Lock this fitbit device for modification.
	c.Kmutex.Lockf("fitbitapp_%v", fitbitAppID.Hex())
	defer c.Kmutex.Unlockf("fitbitapp_%v", fitbitAppID.Hex())

	////
	//// Start the transaction.
	////

	session, err := c.DbClient.StartSession()
	if err != nil {
		c.Logger.Error("start session error",
			slog.Any("fitbit_app_id", fitbitAppID),
			slog.Any("error", err))
		return err
	}
	defer session.EndSession(ctx)

	// Define a transaction function with a series of operations
	transactionFunc := func(sessCtx mongo.SessionContext) (interface{}, error) {
		fba, err := c.FitBitAppStorer.GetByID(sessCtx, fitbitAppID)
		if err != nil {
			c.Logger.Error("failed getting fitbit by id",
				slog.Any("fitbit_app_id", fitbitAppID),
				slog.Any("error", err))
			return nil, err
		}
		if fba == nil {
			err := fmt.Errorf("fitbit does not exist for id %v", fitbitAppID.Hex())
			c.Logger.Error("fitbit does not exist",
				slog.Any("fitbit_app_id", fitbitAppID),
				slog.Any("error", err))
			return nil, err
		}

		if fba.IsTestMode {
			c.Logger.Debug("skip pulling fitbit device in test mode",
				slog.Any("fitbit_app_id", fitbitAppID))
			return nil, nil
		}

		// Defensive code: Make sure the device is properly authenticated.
		if fba.AccessToken == "" || fba.RefreshToken == "" || fba.TokenType == "" {
			errMsg := "fitbit device missing access and or refresh token"
			err := errors.New(errMsg)
			c.Logger.Error("missing fitbit credentials",
				slog.Any("token_type", fba.TokenType),
				slog.Any("refresh_token", fba.RefreshToken),
				slog.Any("access_token", fba.AccessToken),
				slog.Any("err", err))

			// Update the state of the remote device.
			fba.Errors = errMsg
			fba.Status = fitbitapp_s.StatusError
			fba.ModifiedAt = time.Now()
			if err := c.FitBitAppStorer.UpdateByID(sessCtx, fba); err != nil {
				c.Logger.Error("failed updating fitbit app",
					slog.Any("err", err))
				return nil, err
			}
		}

		// Use this to keep track when our pulling has started so we can compare
		// when the pulling stopped and find the duration of the task.
		startAt := time.Now()

		// Pull fitbit device activity data if user granted access to this data.
		if strings.Contains(fba.Scope, "activity") {
			// Developers note:
			// - Ideal scope is as follows: response_type=code&client_id=&scope=activity+cardio_fitness+electrocardiogram+heartrate+location+nutrition+oxygen_saturation+respiratory_rate+sleep+temperature+weight
			// - See: https://dev.fitbit.com/build/reference/web-api/troubleshooting-guide/oauth2-tutorial/
			// - See: https://dev.fitbit.com/build/reference/web-api/developer-guide/authorization/#Authorization-Code-Grant-Flow-with-PKCE

			// var err error

			////
			//// https://dev.fitbit.com/build/reference/web-api/intraday/get-activity-intraday-by-date/
			////

			// resources := []string{"calories", "distance", "elevation", "floors", "steps"} // full list.
			resources := []string{"steps"} // lite list.
			for _, resource := range resources {
				// Lookup the scope in our map and if it exists then process it.
				typeID, ok := fitbitdatum_s.FitBitResourceToType[resource]
				if ok {
					if err := c.fetchActivityDataForFitBit(sessCtx, fba, startAt, resource, typeID); err != nil {
						c.Logger.Error("error with interacting with fitbit web-service for activity",
							slog.Any("error", err))
					}
				} else {
					c.Logger.Error("scope does not exist",
						slog.String("scope", resource))
				}
			}
		}

		if strings.Contains(fba.Scope, "heartrate") {
			if err := c.fetchHeartRateDataForFitBit(sessCtx, fba, startAt, "heartrate", fitbitdatum_s.TypeHeartRate); err != nil {
				c.Logger.Error("error with interacting with fitbit web-service for heartrate",
					slog.Any("error", err))
			}
		}

		if strings.Contains(fba.Scope, "respiratory_rate") {
			// TODO: Implement.
		}

		if strings.Contains(fba.Scope, "oxygen_saturation") {
			// TODO: Implement.
		}

		if strings.Contains(fba.Scope, "sleep") {
			// TODO: Implement.
		}

		if strings.Contains(fba.Scope, "temperature") {
			// TODO: Implement.
		}

		if strings.Contains(fba.Scope, "electrocardiogram") {
			// TODO: Implement.
		}

		if strings.Contains(fba.Scope, "nutrition") {
			// NOTE: Potential to work on in the future.
		}

		if strings.Contains(fba.Scope, "location") {
			// NOTE: Potential to work on in the future.
		}

		if strings.Contains(fba.Scope, "cardio_fitness") {
			// TODO: Implement.
		}

		if strings.Contains(fba.Scope, "weight") {
			// NOTE: Potential to work on in the future.
		}

		// Record the current datetime to the remote device to let it know we
		// fetched it at this time and therefore not pull the service for a while.
		// fda.LastFetchedAt = nowDT
		fba.ModifiedAt = startAt
		if err := c.FitBitAppStorer.UpdateByID(sessCtx, fba); err != nil {
			return nil, err
		}

		return nil, nil
	}

	// Start a transaction
	if _, err := session.WithTransaction(ctx, transactionFunc); err != nil {
		c.Logger.Error("session failed error",
			slog.Any("fitbit_app_id", fitbitAppID),
			slog.Any("error", err))
		return err
	}

	return nil
}
