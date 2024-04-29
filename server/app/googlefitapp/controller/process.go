package controller

import (
	"context"
	"log/slog"

	gcp_a "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/adapter/cloudprovider/google"
	dp_s "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/datapoint/datastore"
	gfdp_ds "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/googlefitdatapoint/datastore"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

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
DataTypeNameDistanceDelta        [DONE]
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
- - - - - - - - - - - - - - - - - - - - - - - - - - -
*/

func (impl *GoogleFitAppControllerImpl) ProcessAllQueuedData() error {
	impl.Logger.Debug("starting task...")

	// DEVELOPERS NOTE:
	// Load up only the data we want to process in our application. In future
	// if we want more data processed then add below:
	dataTypeNames := []string{
		gcp_a.DataTypeNameCaloriesBurned,
		gcp_a.DataTypeNameStepCountDelta,
		gcp_a.DataTypeNameDistanceDelta,
		gcp_a.DataTypeNameHeartRateBPM,
		//TODO: Add more health sensors here...
	}

	ctx := context.Background()
	dpdp, err := impl.GoogleFitDataPointStorer.ListByQueuedStatusInDataTypeNames(ctx, dataTypeNames)
	if err != nil {
		impl.Logger.Error("failed listing queued google fit data points",
			slog.Any("data_type_names", dataTypeNames),
			slog.Any("error", err))
		return err
	}
	for _, dp := range dpdp.Results {
		if err := impl.processForQueuedData(ctx, dp); err != nil {
			impl.Logger.Error("failed transform queued google fit data point",
				slog.Any("dp", dp),
				slog.Any("error", err))
			return err
		}
	}

	impl.Logger.Debug("finished task")
	return nil
}

func (impl *GoogleFitAppControllerImpl) processForQueuedDataWithGfaID(ctx context.Context, gfaID primitive.ObjectID) error {
	dpdp, err := impl.GoogleFitDataPointStorer.ListByQueuedStatusAndGfaID(ctx, gfaID)
	if err != nil {
		impl.Logger.Error("failed listing queued google fit data points",
			slog.Any("google_fit_app_id", gfaID),
			slog.Any("error", err))
		return err
	}

	impl.Logger.Debug("processing queued datapoint",
		slog.Int("count", int(len(dpdp.Results))),
	)

	for _, dp := range dpdp.Results {
		if err := impl.processForQueuedData(ctx, dp); err != nil {
			impl.Logger.Error("failed transform queued google fit data point",
				slog.Any("dp", dp),
				slog.Any("error", err))
			return err
		}
	}

	impl.Logger.Debug("finished task")
	return nil
}

func (impl *GoogleFitAppControllerImpl) processForQueuedData(ctx context.Context, dp *gfdp_ds.GoogleFitDataPoint) error {
	// DEVELOPERS NOTE:
	// The following code will extract the data-point `value` based on the
	// datatype we are importing into our system.
	var val float64
	switch dp.DataTypeName {
	case gcp_a.DataTypeNameCaloriesBurned:
		if dp.CaloriesBurned != nil {
			val = dp.CaloriesBurned.Calories
		}
	case gcp_a.DataTypeNameStepCountDelta:
		if dp.StepCountDelta != nil {
			val = float64(dp.StepCountDelta.Steps)
		}
	case gcp_a.DataTypeNameDistanceDelta:
		if dp.DistanceDelta != nil {
			val = float64(dp.DistanceDelta.Distance)
		}
	case gcp_a.DataTypeNameHeartRateBPM:
		if dp.HeartRateBPM != nil {
			val = float64(dp.HeartRateBPM.BPM)
		}
		//TODO: Add more health sensors here...
	default:
		// DEVELOPERS NOTE:
		// FOR NOW WE DO NOT WANT TO ERROR, WE JUST WANT TO IGNORE THE
		// RECORD AND MOVE ON... UNCOMMENT THE CODE BELOW WHEN YOU ARE
		// READY TO IMPLEMENT DIFFERENT DATA TYPES.
		return nil

		// err := fmt.Errorf("unsupported data type name: %s", dp.DataTypeName)
		// impl.Logger.Error("",
		// 	slog.Any("error", err))
		// return err
	}

	// DEVELOPERS NOTE:
	// Create our record and insert it if it doesn't exist. We will also need
	// to update the Google Fit data point to no longer be `queued` but
	// `active` state b/c we've processed it.

	dataPoint := &dp_s.DataPoint{
		ID:        primitive.NewObjectID(),
		MetricID:  dp.MetricID,
		Timestamp: dp.EndAt,
		Value:     val,
		IsNull:    false,
	}
	exists, err := impl.DataPointStorer.CheckIfExistsByCompositeKey(ctx, dataPoint.MetricID, dataPoint.Timestamp)
	if err != nil {
		impl.Logger.Error("failed checking by datapoint by composite key",
			slog.Any("error", err))
		return err
	}
	if !exists {
		// STEP 1: Create our record.
		if err := impl.DataPointStorer.Create(ctx, dataPoint); err != nil {
			impl.Logger.Error("failed updating data point",
				slog.Any("error", err))
			return err
		}

		// STEP 2: Update our Google Fit datapoint to be `active` status.
		dp.Status = gfdp_ds.StatusActive
		if err := impl.GoogleFitDataPointStorer.UpdateByID(ctx, dp); err != nil {
			impl.Logger.Error("failed updating google fit data point",
				slog.Any("error", err))
			return err
		}

		// STEP 3: For debugging purposes update the console log.
		impl.Logger.Debug("created datapoint",
			slog.String("data_type_name", dp.DataTypeName),
			slog.String("metric_id", dp.MetricID.Hex()),
			slog.Time("start_at", dp.StartAt),
			slog.Time("end_at", dp.EndAt),
			slog.Int("status", int(dp.Status)),
		)
	}
	return nil
}
