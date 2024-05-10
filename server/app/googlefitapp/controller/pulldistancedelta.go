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

func (impl *GoogleFitAppControllerImpl) pullDistanceDeltaDataFromGoogleWithGfaAndFitnessStore(ctx context.Context, gfa *gfa_ds.GoogleFitApp, svc *fitness.Service) error {
	// impl.Logger.Debug("pulling distance delta dataset",
	// 	slog.String("gfa_id", gfa.ID.Hex()))

	////
	//// Get `Google Fit` data
	////

	maxTime := time.Now()
	minTime := gfa.LastFetchedAt
	dataset, err := impl.GCP.NotAggregatedDatasets(svc, minTime, maxTime, gcp_a.DataTypeShortNameDistanceDelta)
	if err != nil {
		impl.Logger.Error("failed listing distance delta dataset",
			slog.Any("error", err))
		return err
	}

	if len(dataset) == 0 {
		// impl.Logger.Warn("pulled empty distance delta dataset",
		// 	slog.String("gfa_id", gfa.ID.Hex()))
		return nil
	}

	// impl.Logger.Debug("pulled distance delta dataset",
	// 	slog.String("gfa_id", gfa.ID.Hex()))

	////
	//// Convert from `Google Fit` format into our apps format.
	////

	distanceDeltaDataset := gcp_a.ParseDistanceDelta(dataset)

	////
	//// Save into our database.
	////

	for _, distanceDeltaDatapoint := range distanceDeltaDataset {
		exists, err := impl.GoogleFitDataPointStorer.CheckIfExistsByCompositeKey(ctx, gfa.UserID, gcp_a.DataTypeNameDistanceDelta, distanceDeltaDatapoint.StartTime, distanceDeltaDatapoint.EndTime)
		if err != nil {
			impl.Logger.Error("failed checking google fit datapoint by composite key",
				slog.Any("error", err))
			return err
		}
		if !exists {
			if distanceDeltaDatapoint.EndTime.Before(time.Now()) && distanceDeltaDatapoint.StartTime.After(time.Date(2000, 1, 1, 1, 0, 0, 0, time.UTC)) {
				dp := &dp_ds.GoogleFitDataPoint{
					ID:              primitive.NewObjectID(),
					DataTypeName:    gcp_a.DataTypeNameDistanceDelta,
					Status:          dp_ds.StatusQueued,
					UserID:          gfa.UserID,
					UserName:        gfa.UserName,
					UserLexicalName: gfa.UserLexicalName,
					GoogleFitAppID:  gfa.ID,
					MetricID:        gfa.DistanceDeltaMetricID,
					StartAt:         distanceDeltaDatapoint.StartTime,
					EndAt:           distanceDeltaDatapoint.EndTime,
					DistanceDelta:   &distanceDeltaDatapoint,
					Error:           "",
					CreatedAt:       time.Now(),
					ModifiedAt:      time.Now(),
					OrganizationID:  gfa.OrganizationID,
				}
				if err := impl.GoogleFitDataPointStorer.Create(ctx, dp); err != nil {
					impl.Logger.Error("failed inserting google fit data point for distance delta into database",
						slog.Any("error", err))
					return err
				}
				// impl.Logger.Debug("inserted distance delta data point",
				// 	slog.Any("dp", dp))
			}
		}
	}

	return nil
}
