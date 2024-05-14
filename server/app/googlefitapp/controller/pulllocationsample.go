package controller

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

func (impl *GoogleFitAppControllerImpl) pullLocationSampleDataFromGoogleWithGfaAndFitnessStore(ctx context.Context, gfa *gfa_ds.GoogleFitApp, svc *fitness.Service) error {
	// impl.Logger.Debug("pulling location sample dataset",
	// 	slog.String("gfa_id", gfa.ID.Hex()))

	////
	//// Get `Google Fit` data
	////

	maxTime := time.Now()
	minTime := gfa.LastFetchedAt
	dataset, err := impl.GCP.NotAggregatedDatasets(svc, minTime, maxTime, gcp_a.DataTypeShortNameLocationSample)
	if err != nil {
		impl.Logger.Error("failed listing location sample dataset",
			slog.Any("error", err))
		return err
	}

	if len(dataset) == 0 {
		// impl.Logger.Warn("pulled empty location sample dataset",
		// 	slog.String("gfa_id", gfa.ID.Hex()))
		return nil
	}

	// impl.Logger.Debug("pulled location sample dataset",
	// 	slog.String("gfa_id", gfa.ID.Hex()))

	////
	//// Convert from `Google Fit` format into our apps format.
	////

	locationSampleDataset := gcp_a.ParseLocationSample(dataset)

	////
	//// Save into our database.
	////

	for _, locationSampleSegmentDatapoint := range locationSampleDataset {
		exists, err := impl.GoogleFitDataPointStorer.CheckIfExistsByCompositeKey(ctx, gfa.UserID, gcp_a.DataTypeNameLocationSample, locationSampleSegmentDatapoint.StartTime, locationSampleSegmentDatapoint.EndTime)
		if err != nil {
			impl.Logger.Error("failed checking google fit datapoint by composite key",
				slog.Any("error", err))
			return err
		}
		if !exists {
			if locationSampleSegmentDatapoint.EndTime.Before(time.Now()) && locationSampleSegmentDatapoint.StartTime.After(time.Date(2000, 1, 1, 1, 0, 0, 0, time.UTC)) {
				dp := &dp_ds.GoogleFitDataPoint{
					ID:              primitive.NewObjectID(),
					DataTypeName:    gcp_a.DataTypeNameLocationSample,
					Status:          dp_ds.StatusQueued,
					UserID:          gfa.UserID,
					UserName:        gfa.UserName,
					UserLexicalName: gfa.UserLexicalName,
					GoogleFitAppID:  gfa.ID,
					MetricID:        gfa.LocationSampleMetricID,
					StartAt:         locationSampleSegmentDatapoint.StartTime,
					EndAt:           locationSampleSegmentDatapoint.EndTime,
					LocationSample:  &locationSampleSegmentDatapoint,
					Error:           "",
					CreatedAt:       time.Now(),
					ModifiedAt:      time.Now(),
					OrganizationID:  gfa.OrganizationID,
				}
				if err := impl.GoogleFitDataPointStorer.Create(ctx, dp); err != nil {
					impl.Logger.Error("failed inserting google fit data point for locationSample into database",
						slog.Any("error", err))
					return err
				}
				// impl.Logger.Debug("inserted location sample data point",
				// 	slog.Any("dp", dp))
			}
		}
	}

	return nil
}
