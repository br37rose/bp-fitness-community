package controller

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/oauth2"

	gfa_ds "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/googlefitapp/datastore"
)

func (impl *GoogleFitAppControllerImpl) PullDataFromGoogle() error {
	ctx := context.Background()
	gfaIDs, err := impl.GoogleFitAppStorer.ListPhysicalIDsByStatus(ctx, gfa_ds.StatusActive)
	if err != nil {
		impl.Logger.Error("failed listing google fit apps by status",
			slog.Any("error", err))
		return err
	}
	for _, gfaID := range gfaIDs {
		if err := impl.pullDataFromGoogleWithGfaID(ctx, gfaID); err != nil {
			impl.Logger.Error("failed pulling data for google fit app",
				slog.Any("error", err))
			return err
		}
	}
	return nil
}

func (impl *GoogleFitAppControllerImpl) pullDataFromGoogleWithGfaID(ctx context.Context, gfaID primitive.ObjectID) error {
	// Lock this google fit app
	impl.DistributedMutex.Lockf(ctx, "googlefitapp_%v", gfaID.Hex())
	defer impl.DistributedMutex.Unlockf(ctx, "googlefitapp_%v", gfaID.Hex())

	// Get our database record.
	gfa, err := impl.GoogleFitAppStorer.GetByID(ctx, gfaID)
	if err != nil {
		impl.Logger.Error("failed getting google fit app from database",
			slog.String("gfa_id", gfaID.Hex()),
			slog.Any("error", err))
		return err
	}
	if gfa == nil {
		err := fmt.Errorf("google fit app does not exist for id: %s", gfaID.Hex())
		return err
	}

	impl.Logger.Debug("starting...",
		slog.String("gfa_id", gfaID.Hex()),
		slog.String("user_id", gfa.UserID.Hex()),
	)

	// Authenticated http client for a specific user's account. Note: No need
	// for refresh token handling as it's already handled!
	client, err := impl.GCP.NewHTTPClientFromToken(gfa.Token, func(newTok *oauth2.Token) {
		// Save the latest token provided by Google in our database and
		// make sure our status is set to running and having no problems.
		gfa.Token = newTok
		gfa.RequiresGoogleLoginAgain = false
		gfa.Status = gfa_ds.StatusActive
		gfa.Errors = ""
		if err := impl.GoogleFitAppStorer.UpdateByID(ctx, gfa); err != nil {
			impl.Logger.Error("failed updating google fit app in database",
				slog.String("gfa_id", gfaID.Hex()),
				slog.Any("error", err))
		}
		impl.Logger.Debug("updated google fit app with new token", slog.String("gfa_id", gfaID.Hex()))
	})
	if err != nil {
		//
		// If any errors occur let's force the user to log in again.
		//

		gfa.RequiresGoogleLoginAgain = true
		gfa.Status = gfa_ds.StatusError
		gfa.Errors = err.Error()
		if err := impl.GoogleFitAppStorer.UpdateByID(ctx, gfa); err != nil {
			impl.Logger.Error("failed updating google fit app in database",
				slog.String("gfa_id", gfaID.Hex()),
				slog.Any("error", err))
		}

		u, err := impl.UserStorer.GetByID(ctx, gfa.UserID)
		if err != nil {
			impl.Logger.Error("failed getting user from database",
				slog.String("gfa_id", gfaID.Hex()),
				slog.Any("error", err))
			return err
		}
		if u == nil {
			err := fmt.Errorf("user does not exist for id: %s", gfa.UserID.Hex())
			impl.Logger.Error("failed getting user from database",
				slog.String("gfa_id", gfaID.Hex()),
				slog.Any("error", err))
			return err
		}
		u.PrimaryHealthTrackingDeviceRequiresLoginAgain = true
		if err := impl.UserStorer.UpdateByID(ctx, u); err != nil {
			impl.Logger.Error("failed updating user to database",
				slog.String("gfa_id", gfaID.Hex()),
				slog.Any("error", err))
			return err
		}

		impl.Logger.Warn("access and refresh token expired, user must log in again for google fit app",
			slog.String("gfa_id", gfaID.Hex()),
			slog.String("user_id", gfa.UserID.Hex()),
		)

		return nil
	}
	if client == nil {
		err := fmt.Errorf("google fit app authenticated client does not exist for token: %v", gfa.Token)
		return err
	}

	return impl.pullDataFromGoogleWithGfaAndClient(ctx, gfa, client)
}

func (impl *GoogleFitAppControllerImpl) pullDataFromGoogleWithGfaAndClient(ctx context.Context, gfa *gfa_ds.GoogleFitApp, client *http.Client) error {
	////
	//// Load up the Google Fitness Store.
	////

	svc, err := impl.GCP.NewFitnessStoreFromClient(client)
	if err != nil {
		impl.Logger.Error("failed creating new fitness store from client",
			slog.Any("error", err))
		return err
	}
	if svc == nil {
		err := fmt.Errorf("google fit app fitness store is empty for token: %v", gfa.Token)
		return err
	}

	////
	//// Get various data.
	////

	// Variable used to track the latest fetch time we've done. Take the time
	// at this moment and then hold onto it. There will be time delay for
	// running all of the calls below. Afterwords we will save this variables
	// time as the most recent time of our fetch.
	lastFetchedAtNow := time.Now()

	// --- Activity --- //

	if err := impl.pullActivitySegmentDataFromGoogleWithGfaAndFitnessStore(ctx, gfa, svc); err != nil {
		impl.Logger.Error("failed pulling activity dataset from google",
			slog.Any("error", err))
		return err
	}
	if err := impl.pullBasalMetabolicRateDataFromGoogleWithGfaAndFitnessStore(ctx, gfa, svc); err != nil {
		impl.Logger.Error("failed pulling activity dataset from google",
			slog.Any("error", err))
		return err
	}
	if err := impl.pullCaloriesBurnedDataFromGoogleWithGfaAndFitnessStore(ctx, gfa, svc); err != nil {
		impl.Logger.Error("failed pulling calories burned data from google",
			slog.Any("error", err))
		return err
	}
	if err := impl.pullCyclingPedalingCadenceDataFromGoogleWithGfaAndFitnessStore(ctx, gfa, svc); err != nil {
		impl.Logger.Error("failed pulling cycling pedaling cadence dataset from google",
			slog.Any("error", err))
		return err
	}
	if err := impl.pullCyclingPedalingCumulativeDataFromGoogleWithGfaAndFitnessStore(ctx, gfa, svc); err != nil {
		impl.Logger.Error("failed pulling cycling pedaling cumulative dataset from google",
			slog.Any("error", err))
		return err
	}
	if err := impl.pullHeartPointsDataFromGoogleWithGfaAndFitnessStore(ctx, gfa, svc); err != nil {
		impl.Logger.Error("failed pulling heart points dataset from google",
			slog.Any("error", err))
		return err
	}
	if err := impl.pullMoveMinutesDataFromGoogleWithGfaAndFitnessStore(ctx, gfa, svc); err != nil {
		impl.Logger.Error("failed pulling move minutes dataset from google",
			slog.Any("error", err))
		return err
	}
	if err := impl.pullPowerDataFromGoogleWithGfaAndFitnessStore(ctx, gfa, svc); err != nil {
		impl.Logger.Error("failed pulling power dataset from google",
			slog.Any("error", err))
		return err
	}
	if err := impl.pullStepCountCadenceDataFromGoogleWithGfaAndFitnessStore(ctx, gfa, svc); err != nil {
		impl.Logger.Error("failed pulling step count cadence dataset from google",
			slog.Any("error", err))
		return err
	}
	if err := impl.pullStepCountDeltaDataFromGoogleWithGfaAndFitnessStore(ctx, gfa, svc); err != nil {
		impl.Logger.Error("failed pulling step count delta data from google",
			slog.Any("error", err))
		return err
	}
	if err := impl.pullWorkoutDataFromGoogleWithGfaAndFitnessStore(ctx, gfa, svc); err != nil {
		impl.Logger.Error("failed pulling workout dataset from google",
			slog.Any("error", err))
		return err
	}

	// --- Location --- //

	if err := impl.pullCyclingWheelRevolutionRPMDataFromGoogleWithGfaAndFitnessStore(ctx, gfa, svc); err != nil {
		impl.Logger.Error("failed pulling cycling wheel revolution rpm dataset from google",
			slog.Any("error", err))
		return err
	}
	if err := impl.pullCyclingWheelRevolutionCumulativeDataFromGoogleWithGfaAndFitnessStore(ctx, gfa, svc); err != nil {
		impl.Logger.Error("failed pulling cycling wheel revolution cumulative dataset from google",
			slog.Any("error", err))
		return err
	}
	if err := impl.pullDistanceDeltaDataFromGoogleWithGfaAndFitnessStore(ctx, gfa, svc); err != nil {
		impl.Logger.Error("failed pulling distance delta dataset from google",
			slog.Any("error", err))
		return err
	}
	if err := impl.pullLocationSampleDataFromGoogleWithGfaAndFitnessStore(ctx, gfa, svc); err != nil {
		impl.Logger.Error("failed pulling location sample dataset from google",
			slog.Any("error", err))
		return err
	}
	if err := impl.pullSpeedDataFromGoogleWithGfaAndFitnessStore(ctx, gfa, svc); err != nil {
		impl.Logger.Error("failed pulling speed dataset from google",
			slog.Any("error", err))
		return err
	}

	// --- Nutrition --- //

	if err := impl.pullHydrationDataFromGoogleWithGfaAndFitnessStore(ctx, gfa, svc); err != nil {
		impl.Logger.Error("failed pulling hydration data from google",
			slog.Any("error", err))
		return err
	}
	if err := impl.pullNutritionDataFromGoogleWithGfaAndFitnessStore(ctx, gfa, svc); err != nil {
		impl.Logger.Error("failed pulling hydration data from google",
			slog.Any("error", err))
		return err
	}

	// --- Health --- //

	if err := impl.pullBloodGlucoseDataFromGoogleWithGfaAndFitnessStore(ctx, gfa, svc); err != nil {
		impl.Logger.Error("failed pulling blood glucose data from google",
			slog.Any("error", err))
		return err
	}
	if err := impl.pullBloodPressureDataFromGoogleWithGfaAndFitnessStore(ctx, gfa, svc); err != nil {
		impl.Logger.Error("failed pulling blood pressure data from google",
			slog.Any("error", err))
		return err
	}
	if err := impl.pullBodyFatPercentageDataFromGoogleWithGfaAndFitnessStore(ctx, gfa, svc); err != nil {
		impl.Logger.Error("failed pulling body fat percentage data from google",
			slog.Any("error", err))
		return err
	}
	if err := impl.pullBodyTemperatureDataFromGoogleWithGfaAndFitnessStore(ctx, gfa, svc); err != nil {
		impl.Logger.Error("failed pulling body temperature data from google",
			slog.Any("error", err))
		return err
	}
	if err := impl.pullHeartRateDataFromGoogleWithGfaAndFitnessStore(ctx, gfa, svc); err != nil {
		impl.Logger.Error("failed pulling heart rate dataset from google",
			slog.Any("error", err))
		return err
	}
	if err := impl.pullHeightDataFromGoogleWithGfaAndFitnessStore(ctx, gfa, svc); err != nil {
		impl.Logger.Error("failed pulling height data from google",
			slog.Any("error", err))
		return err
	}
	if err := impl.pullOxygenSaturationDataFromGoogleWithGfaAndFitnessStore(ctx, gfa, svc); err != nil {
		impl.Logger.Error("failed pulling oxygen saturation data from google",
			slog.Any("error", err))
		return err
	}
	if err := impl.pullSleepDataFromGoogleWithGfaAndFitnessStore(ctx, gfa, svc); err != nil {
		impl.Logger.Error("failed pulling sleep data from google",
			slog.Any("error", err))
		return err
	}
	if err := impl.pullWeightDataFromGoogleWithGfaAndFitnessStore(ctx, gfa, svc); err != nil {
		impl.Logger.Error("failed pulling weight data from google",
			slog.Any("error", err))
		return err
	}

	//
	// Keep track of last fetch time.
	//

	gfa.LastFetchedAt = lastFetchedAtNow
	if err := impl.GoogleFitAppStorer.UpdateByID(ctx, gfa); err != nil {
		impl.Logger.Error("failed pulling weight data from google",
			slog.Any("error", err))
		return err
	}

	return nil
}
