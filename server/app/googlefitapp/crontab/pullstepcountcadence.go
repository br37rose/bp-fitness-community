package crontab

import (
	"context"
	"log/slog"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/api/fitness/v1"

	gcp_a "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/adapter/cloudprovider/google"
	gfa_ds "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/googlefitapp/datastore"
	dp_ds "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/googlefitdatapoint/datastore"
)

func (impl *googleFitAppCrontaberImpl) pullStepCountCadenceDataFromGoogleWithGfaAndFitnessStore(ctx context.Context, gfa *gfa_ds.GoogleFitApp, svc *fitness.Service) error {
	impl.Logger.Debug("pulling step count delta dataset",
		slog.String("gfa_id", gfa.ID.Hex()))

	////
	//// Get `Google Fit` data
	////

	maxTime := time.Now()
	minTime := gfa.LastFetchedAt
	dataset, err := impl.GCP.NotAggregatedDatasets(svc, minTime, maxTime, gcp_a.DataTypeShortNameStepCountCadence)
	if err != nil {
		impl.Logger.Error("failed listing step count delta dataset",
			slog.Any("error", err))
		return err
	}

	if len(dataset) == 0 {
		impl.Logger.Warn("pulled empty step count delta dataset",
			slog.String("gfa_id", gfa.ID.Hex()))
		return nil
	}

	impl.Logger.Debug("pulled step count delta dataset",
		slog.String("gfa_id", gfa.ID.Hex()))

	////
	//// Convert from `Google Fit` format into our apps format.
	////

	stepCountCadenceDataset := gcp_a.ParseStepCountCadence(dataset)

	////
	//// Save into our database.
	////

	for _, stepCountCadenceDatapoint := range stepCountCadenceDataset {
		exists, err := impl.GoogleFitDataPointStorer.CheckIfExistsByCompositeKey(ctx, gfa.UserID, gcp_a.DataTypeNameStepCountCadence, stepCountCadenceDatapoint.StartTime, stepCountCadenceDatapoint.EndTime)
		if err != nil {
			impl.Logger.Error("failed checking google fit datapoint by composite key",
				slog.Any("error", err))
			return err
		}
		if !exists {
			if stepCountCadenceDatapoint.EndTime.Before(time.Now()) && stepCountCadenceDatapoint.StartTime.After(time.Date(2000, 1, 1, 1, 0, 0, 0, time.UTC)) {
				dp := &dp_ds.GoogleFitDataPoint{
					ID:               primitive.NewObjectID(),
					DataTypeName:     gcp_a.DataTypeNameStepCountCadence,
					Status:           dp_ds.StatusQueued,
					UserID:           gfa.UserID,
					UserName:         gfa.UserName,
					UserLexicalName:  gfa.UserLexicalName,
					GoogleFitAppID:   gfa.ID,
					MetricID:         gfa.StepCountCadenceMetricID,
					StartAt:          stepCountCadenceDatapoint.StartTime,
					EndAt:            stepCountCadenceDatapoint.EndTime,
					StepCountCadence: &stepCountCadenceDatapoint,
					Error:            "",
					CreatedAt:        time.Now(),
					ModifiedAt:       time.Now(),
					OrganizationID:   gfa.OrganizationID,
				}
				if err := impl.GoogleFitDataPointStorer.Create(ctx, dp); err != nil {
					impl.Logger.Error("failed inserting google fit data point for step count delta into database",
						slog.Any("error", err))
					return err
				}
				impl.Logger.Debug("inserted step count delta data point",
					slog.Any("dp", dp))
			}
		}
	}

	return nil
}
