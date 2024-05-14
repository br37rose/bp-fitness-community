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
	// impl.Logger.Debug("processing all queued data...")

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

	// DEVELOPERS NOTE:
	// We are only going to pull a limited number of data by using pagination
	// and hence will require us to iterate over each page. This is done to
	// help prvent the server from erroring with a "context exceed" error.
	cursor := ""
	isRunning := true
	for isRunning {
		f := &gfdp_ds.GoogleFitDataPointPaginationListFilter{
			Cursor:        cursor,
			PageSize:      250,
			SortField:     "created_at",
			SortOrder:     gfdp_ds.OrderDescending,
			Status:        gfdp_ds.StatusQueued,
			DataTypeNames: dataTypeNames,
		}
		dpdp, err := impl.GoogleFitDataPointStorer.ListByFilter(ctx, f)
		if err != nil {
			impl.Logger.Error("failed listing queued google fit data points",
				slog.Any("data_type_names", dataTypeNames),
				slog.Any("error", err))
			return err
		}
		// impl.Logger.Debug("part 1 of 2 - processing 250 queued data...",
		// 	slog.String("cursor", cursor),
		// 	slog.Any("f", f),
		// 	slog.Any("any", dpdp.Results),
		// 	slog.Int("count", int(len(dpdp.Results))))

		for _, dp := range dpdp.Results {
			// impl.Logger.Debug("processing...",
			// 	slog.Any("id", dp.ID),
			// 	slog.Any("dtn", dp.DataTypeName),
			// 	slog.Any("status", dp.Status),
			// 	slog.Any("start_at", dp.StartAt),
			// )
			if err := impl.processForQueuedData(ctx, dp); err != nil {
				impl.Logger.Error("failed transform queued google fit data point",
					slog.Any("dp", dp),
					slog.Any("error", err))
				return err
			}
		}

		// impl.Logger.Debug("part 2 of 2 - processing 250 queued data...",
		// 	slog.String("NextCursor", dpdp.NextCursor),
		// 	slog.Any("HasNextPage", dpdp.HasNextPage),
		// )

		if dpdp.HasNextPage {
			cursor = dpdp.NextCursor
			// impl.Logger.Debug("processing next page of all queued data...")
		} else {
			// If there is no more cursors from the list results then that means
			// all our pagination have been completed so we can exit while loop.
			isRunning = false
			// impl.Logger.Debug("processing finished of all queued data")
		}
	}

	// impl.Logger.Debug("finished task")
	return nil
}

func (impl *GoogleFitAppControllerImpl) processForQueuedDataWithGfaID(ctx context.Context, gfaID primitive.ObjectID) error {
	// impl.Logger.Debug("processing queued data for gfa...", slog.String("gfa_id", gfaID.Hex()))

	// DEVELOPERS NOTE:
	// We are only going to pull a limited number of data by using pagination
	// and hence will require us to iterate over each page. This is done to
	// help prvent the server from erroring with a "context exceed" error.
	var cursor string = ""
	var isRunning = true
	for isRunning {
		f := &gfdp_ds.GoogleFitDataPointPaginationListFilter{
			Cursor:         cursor,
			PageSize:       250,
			SortField:      "created_at",
			SortOrder:      gfdp_ds.OrderDescending,
			Status:         gfdp_ds.StatusQueued,
			GoogleFitAppID: gfaID,
		}
		dpdp, err := impl.GoogleFitDataPointStorer.ListByFilter(ctx, f)
		if err != nil {
			impl.Logger.Error("failed listing queued google fit data points",
				slog.String("gfa_id", gfaID.Hex()),
				slog.Any("error", err))
			return err
		}
		// impl.Logger.Debug("processing 250 queued data...",
		// 	slog.String("cursor", cursor),
		// 	slog.String("gfa_id", gfaID.Hex()),
		// 	slog.Int("count", int(len(dpdp.Results))))

		for _, dp := range dpdp.Results {
			if err := impl.processForQueuedData(ctx, dp); err != nil {
				impl.Logger.Error("failed transform queued google fit data point",
					slog.Any("dp", dp),
					slog.Any("error", err))
				return err
			}
		}

		if dpdp.HasNextPage {
			cursor = dpdp.NextCursor
			// impl.Logger.Debug("processing next page of queued data for gfa...",
			// 	slog.String("gfa_id", gfaID.Hex()),
			// )
		} else {
			// If there is no more cursors from the list results then that means
			// all our pagination have been completed so we can exit while loop.
			isRunning = false
			// impl.Logger.Debug("processing finished of queued data for gfa...",
			// 	slog.String("gfa_id", gfaID.Hex()),
			// )
		}
	}

	// impl.Logger.Debug("processed queued data for gfa", slog.String("gfa_id", gfaID.Hex()))
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
		ID:                 primitive.NewObjectID(),
		MetricID:           dp.MetricID,
		MetricDataTypeName: dp.DataTypeName,
		UserID:             dp.UserID,
		Timestamp:          dp.EndAt,
		Value:              val,
		IsNull:             false,
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

		// // STEP 3: For debugging purposes update the console log.
		// impl.Logger.Debug("created datapoint",
		// 	slog.String("data_type_name", dp.DataTypeName),
		// 	slog.String("metric_id", dp.MetricID.Hex()),
		// 	slog.String("user_id", dp.UserID.Hex()),
		// 	slog.Time("start_at", dp.StartAt),
		// 	slog.Time("end_at", dp.EndAt),
		// 	slog.Int("status", int(dp.Status)),
		// )
	} else {
		if dp.Status == gfdp_ds.StatusQueued {
			// impl.Logger.Debug("datapoint already exists",
			// 	slog.String("data_type_name", dp.DataTypeName),
			// 	slog.String("metric_id", dp.MetricID.Hex()),
			// 	slog.Time("start_at", dp.StartAt),
			// 	slog.Time("end_at", dp.EndAt),
			// 	slog.Int("status", int(dp.Status)),
			// )

			// Update our Google Fit datapoint to be `active` status.
			dp.Status = gfdp_ds.StatusActive
			if err := impl.GoogleFitDataPointStorer.UpdateByID(ctx, dp); err != nil {
				impl.Logger.Error("failed updating google fit data point",
					slog.Any("error", err))
				return err
			}

			// // STEP 3: For debugging purposes update the console log.
			// impl.Logger.Debug("updated datapoint",
			// 	slog.String("data_type_name", dp.DataTypeName),
			// 	slog.String("metric_id", dp.MetricID.Hex()),
			// 	slog.Time("start_at", dp.StartAt),
			// 	slog.Time("end_at", dp.EndAt),
			// 	slog.Int("status", int(dp.Status)),
			// )
		}

	}

	return nil
}
