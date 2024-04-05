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

func (impl *googleFitAppCrontaberImpl) pullHeartPointsDataFromGoogleWithGfaAndFitnessStore(ctx context.Context, gfa *gfa_ds.GoogleFitApp, svc *fitness.Service) error {
	impl.Logger.Debug("pulling heart point dataset",
		slog.String("gfa_id", gfa.ID.Hex()))

	////
	//// Get `Google Fit` data
	////

	maxTime := time.Now()
	minTime := gfa.LastFetchedAt
	dataset, err := impl.GCP.NotAggregatedDatasets(svc, minTime, maxTime, gcp_a.DataTypeShortNameHeartPoints)
	if err != nil {
		impl.Logger.Error("failed listing heart point dataset",
			slog.Any("error", err))
		return err
	}

	if len(dataset) == 0 {
		impl.Logger.Warn("pulled empty heart point dataset",
			slog.String("gfa_id", gfa.ID.Hex()))
		return nil
	}

	impl.Logger.Debug("pulled power dataset",
		slog.String("gfa_id", gfa.ID.Hex()))

	////
	//// Convert from `Google Fit` format into our apps format.
	////

	heartPointDataset := gcp_a.ParseHeartPoints(dataset)
	//
	// impl.Logger.Debug("",
	// 	slog.String("gfa_id", gfa.ID.Hex()),
	// 	slog.Any("dataset", dataset),
	// 	slog.Any("heartPointDataset", heartPointDataset),
	// )

	////
	//// Save into our database.
	////

	for _, heartPointDatapoint := range heartPointDataset {
		exists, err := impl.GoogleFitDataPointStorer.CheckIfExistsByCompositeKey(ctx, gfa.UserID, gcp_a.DataTypeNameHeartPoints, heartPointDatapoint.StartTime, heartPointDatapoint.EndTime)
		if err != nil {
			impl.Logger.Error("failed checking google fit datapoint by composite key",
				slog.Any("error", err))
			return err
		}
		if !exists {
			dp := &dp_ds.GoogleFitDataPoint{
				ID:              primitive.NewObjectID(),
				DataTypeName:    gcp_a.DataTypeNameHeartPoints,
				Status:          dp_ds.StatusQueued,
				UserID:          gfa.UserID,
				UserName:        gfa.UserName,
				UserLexicalName: gfa.UserLexicalName,
				GoogleFitAppID:  gfa.ID,
				MetricID:        gfa.HeartPointsMetricID,
				StartAt:         heartPointDatapoint.StartTime,
				EndAt:           heartPointDatapoint.EndTime,
				HeartPoints:     &heartPointDatapoint,
				Error:           "",
				CreatedAt:       time.Now(),
				ModifiedAt:      time.Now(),
				OrganizationID:  gfa.OrganizationID,
			}
			if err := impl.GoogleFitDataPointStorer.Create(ctx, dp); err != nil {
				impl.Logger.Error("failed inserting google fit data point for heart point into database",
					slog.Any("error", err))
				return err
			}
			impl.Logger.Debug("inserted heart point data point",
				slog.Any("dp", dp))
		}
	}

	return nil
}
