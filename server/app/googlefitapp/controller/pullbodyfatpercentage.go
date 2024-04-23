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

func (impl *GoogleFitAppControllerImpl) pullBodyFatPercentageDataFromGoogleWithGfaAndFitnessStore(ctx context.Context, gfa *gfa_ds.GoogleFitApp, svc *fitness.Service) error {
	impl.Logger.Debug("pulling body fat percentage dataset",
		slog.String("gfa_id", gfa.ID.Hex()))

	////
	//// Get `Google Fit` data
	////

	maxTime := time.Now()
	minTime := gfa.LastFetchedAt
	dataset, err := impl.GCP.NotAggregatedDatasets(svc, minTime, maxTime, gcp_a.DataTypeShortNameBodyFatPercentage)
	if err != nil {
		impl.Logger.Error("failed listing body fat percentage dataset",
			slog.Any("error", err))
		return err
	}

	if len(dataset) == 0 {
		impl.Logger.Warn("pulled empty body fat percentage dataset",
			slog.String("gfa_id", gfa.ID.Hex()))
		return nil
	}

	impl.Logger.Debug("pulled body fat percentage dataset",
		slog.String("gfa_id", gfa.ID.Hex()))

	////
	//// Convert from `Google Fit` format into our apps format.
	////

	bodyFatPercentageDataset := gcp_a.ParseBodyFatPercentage(dataset)
	//
	// impl.Logger.Debug("",
	// 	slog.String("gfa_id", gfa.ID.Hex()),
	// 	slog.Any("dataset", dataset),
	// 	slog.Any("bodyFatPercentageDataset", bodyFatPercentageDataset),
	// )

	////
	//// Save into our database.
	////

	for _, bodyFatPercentageDatapoint := range bodyFatPercentageDataset {
		exists, err := impl.GoogleFitDataPointStorer.CheckIfExistsByCompositeKey(ctx, gfa.UserID, gcp_a.DataTypeNameBodyFatPercentage, bodyFatPercentageDatapoint.StartTime, bodyFatPercentageDatapoint.EndTime)
		if err != nil {
			impl.Logger.Error("failed checking google fit datapoint by composite key",
				slog.Any("error", err))
			return err
		}
		if !exists {
			if bodyFatPercentageDatapoint.EndTime.Before(time.Now()) && bodyFatPercentageDatapoint.StartTime.After(time.Date(2000, 1, 1, 1, 0, 0, 0, time.UTC)) {
				dp := &dp_ds.GoogleFitDataPoint{
					ID:                primitive.NewObjectID(),
					DataTypeName:      gcp_a.DataTypeNameBodyFatPercentage,
					Status:            dp_ds.StatusQueued,
					UserID:            gfa.UserID,
					UserName:          gfa.UserName,
					UserLexicalName:   gfa.UserLexicalName,
					GoogleFitAppID:    gfa.ID,
					MetricID:          gfa.BodyFatPercentageMetricID,
					StartAt:           bodyFatPercentageDatapoint.StartTime,
					EndAt:             bodyFatPercentageDatapoint.EndTime,
					BodyFatPercentage: &bodyFatPercentageDatapoint,
					Error:             "",
					CreatedAt:         time.Now(),
					ModifiedAt:        time.Now(),
					OrganizationID:    gfa.OrganizationID,
				}
				if err := impl.GoogleFitDataPointStorer.Create(ctx, dp); err != nil {
					impl.Logger.Error("failed inserting google fit data point for body fat percentage into database",
						slog.Any("error", err))
					return err
				}
				impl.Logger.Debug("inserted body fat percentage data point",
					slog.Any("dp", dp))
			}
		}
	}

	return nil
}
