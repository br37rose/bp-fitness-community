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

func (impl *GoogleFitAppControllerImpl) pullBasalMetabolicRateDataFromGoogleWithGfaAndFitnessStore(ctx context.Context, gfa *gfa_ds.GoogleFitApp, svc *fitness.Service) error {
	impl.Logger.Debug("pulling basal metabolic rate dataset",
		slog.String("gfa_id", gfa.ID.Hex()))

	////
	//// Get `Google Fit` data
	////

	maxTime := time.Now()
	minTime := gfa.LastFetchedAt
	dataset, err := impl.GCP.NotAggregatedDatasets(svc, minTime, maxTime, gcp_a.DataTypeShortNameBasalMetabolicRate)
	if err != nil {
		impl.Logger.Error("failed listing basal metabolic rate dataset",
			slog.Any("error", err))
		return err
	}

	if len(dataset) == 0 {
		impl.Logger.Warn("pulled empty basal metabolic rate dataset",
			slog.String("gfa_id", gfa.ID.Hex()))
		return nil
	}

	impl.Logger.Debug("pulled basal metabolic rate dataset",
		slog.String("gfa_id", gfa.ID.Hex()))

	////
	//// Convert from `Google Fit` format into our apps format.
	////

	basalMetabolicRateDataset := gcp_a.ParseBasalMetabolicRate(dataset)
	//
	// impl.Logger.Debug("",
	// 	slog.String("gfa_id", gfa.ID.Hex()),
	// 	slog.Any("dataset", dataset),
	// 	slog.Any("basalMetabolicRateDataset", basalMetabolicRateDataset),
	// )

	////
	//// Save into our database.
	////

	for _, basalMetabolicRateDatapoint := range basalMetabolicRateDataset {
		exists, err := impl.GoogleFitDataPointStorer.CheckIfExistsByCompositeKey(ctx, gfa.UserID, gcp_a.DataTypeNameBasalMetabolicRate, basalMetabolicRateDatapoint.StartTime, basalMetabolicRateDatapoint.EndTime)
		if err != nil {
			impl.Logger.Error("failed checking google fit datapoint by composite key",
				slog.Any("error", err))
			return err
		}
		if !exists {
			if basalMetabolicRateDatapoint.EndTime.Before(time.Now()) && basalMetabolicRateDatapoint.StartTime.After(time.Date(2000, 1, 1, 1, 0, 0, 0, time.UTC)) {
				dp := &dp_ds.GoogleFitDataPoint{
					ID:                 primitive.NewObjectID(),
					DataTypeName:       gcp_a.DataTypeNameBasalMetabolicRate,
					Status:             dp_ds.StatusQueued,
					UserID:             gfa.UserID,
					UserName:           gfa.UserName,
					UserLexicalName:    gfa.UserLexicalName,
					GoogleFitAppID:     gfa.ID,
					MetricID:           gfa.BasalMetabolicRateMetricID,
					StartAt:            basalMetabolicRateDatapoint.StartTime,
					EndAt:              basalMetabolicRateDatapoint.EndTime,
					BasalMetabolicRate: &basalMetabolicRateDatapoint,
					Error:              "",
					CreatedAt:          time.Now(),
					ModifiedAt:         time.Now(),
					OrganizationID:     gfa.OrganizationID,
				}
				if err := impl.GoogleFitDataPointStorer.Create(ctx, dp); err != nil {
					impl.Logger.Error("failed inserting google fit data point for basal metabolic rate into database",
						slog.Any("error", err))
					return err
				}
				impl.Logger.Debug("inserted basal metabolic rate data point",
					slog.Any("dp", dp))
			}
		}
	}

	return nil
}
