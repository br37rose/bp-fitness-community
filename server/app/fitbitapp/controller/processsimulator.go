package controller

import (
	"context"
	"fmt"
	"log/slog"
	"math/rand"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	dp_s "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/datapoint/datastore"
	fitbitapp_s "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/fitbitapp/datastore"
)

func (c *FitBitAppControllerImpl) ProcessAllActiveSimulators(ctx context.Context) error {
	ids, err := c.FitBitAppStorer.ListSimulatorIDsByStatus(ctx, fitbitapp_s.StatusActive)
	if err != nil {
		c.Logger.Error("failed listing ids by status",
			slog.Any("error", err))
		return err
	}
	// // For debugging purposes only.
	// c.Logger.Debug("processing all active fitbit simulator devices", slog.Any("fba_ids", ids))
	for _, fitbitAppID := range ids {
		if err := c.ProcessSimulator(ctx, fitbitAppID); err != nil {
			c.Logger.Error("failed simulating fitbit app",
				slog.Any("fitbit_app_id", fitbitAppID),
				slog.Any("error", err))
			return err
		}
	}
	return nil
}

func (c *FitBitAppControllerImpl) ProcessSimulator(ctx context.Context, fitbitAppID primitive.ObjectID) error {
	// // For debugging purposes only.
	// c.Logger.Debug(fmt.Sprintf("processing fitbit simulator ID %v", fitbitAppID.Hex()))
	// defer c.Logger.Debug(fmt.Sprintf("processed fitbit simulator ID %v", fitbitAppID.Hex()))

	// // Lock this fitbit device for modification.
	// c.Kmutex.Lockf("fitbitapp_%v", fitbitAppID.Hex())
	// defer c.Kmutex.Unlockf("fitbitapp_%v", fitbitAppID.Hex())

	fba, err := c.FitBitAppStorer.GetByID(ctx, fitbitAppID)
	if err != nil {
		c.Logger.Error("failed getting fitbit simulator by id",
			slog.Any("fitbit_app_id", fitbitAppID),
			slog.Any("error", err))
		return err
	}
	if fba == nil {
		err := fmt.Errorf("fitbit simulator does not exist for id %v", fitbitAppID.Hex())
		c.Logger.Error("fitbit does not exist",
			slog.Any("fitbit_app_id", fitbitAppID),
			slog.Any("error", err))
		return err
	}

	// Use this to keep track when our pulling has started so we can compare
	// when the pulling stopped and find the duration of the task.
	startAt := time.Now()

	////
	//// ACTIVITY
	////

	// Process fitbit device activity data if user granted access to this data.
	if strings.Contains(fba.Scope, "activity") {
		var latestValue float64 = randomBetween(0, 100)

		// Generate our time-series datum.
		dp := &dp_s.DataPoint{
			ID:        primitive.NewObjectID(),
			MetricID:  fba.StepsCountMetricID,
			Timestamp: time.Now(),
			Value:     latestValue,
		}
		if err := c.DataPointStorer.Create(ctx, dp); err != nil {
			c.Logger.Error("failed creating simulated data point",
				slog.Any("fitbit_app_id", fitbitAppID),
				slog.Any("error", err))
			return err
		}

		// // For debugging purposes only.
		// c.Logger.Debug(fmt.Sprintf("created simulated data point #%s", dp.ID.Hex()),
		// 	slog.Any("metric_id", fba.StepsCountMetricID))
	}

	////
	//// HEART RATE
	////

	if strings.Contains(fba.Scope, "heartrate") {
		var latestValue float64 = randomBetween(50, 190)

		// Generate our time-series datum.
		dp := &dp_s.DataPoint{
			ID:        primitive.NewObjectID(),
			MetricID:  fba.HeartRateMetricID,
			Timestamp: time.Now(),
			Value:     latestValue,
		}
		if err := c.DataPointStorer.Create(ctx, dp); err != nil {
			c.Logger.Error("failed creating simulated data point",
				slog.Any("fitbit_app_id", fitbitAppID),
				slog.Any("error", err))
			return err
		}
		// // For debugging purposes only.
		// c.Logger.Debug(fmt.Sprintf("created simulated data point #%s", dp.ID.Hex()),
		// 	slog.Any("metric_id", fba.HeartRateMetricID))
	}

	// Record the current datetime to the remote device to let it know we
	// fetched it at this time and therefore not pull the service for a while.
	// fda.LastFetchedAt = nowDT
	fba.ModifiedAt = startAt
	if err := c.FitBitAppStorer.UpdateByID(ctx, fba); err != nil {
		return err
	}
	return nil
}

func randomBetween(min, max int) float64 {
	rand.Seed(time.Now().UnixNano())
	return float64(rand.Intn(max-min+1) + min)
}
