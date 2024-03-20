package crontab

import (
	"context"
	"log/slog"
	"time"

	"google.golang.org/api/fitness/v1"

	gcp_a "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/adapter/cloudprovider/google"
	gfa_ds "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/googlefitapp/datastore"
)

func (impl *googleFitAppCrontaberImpl) pullHydrationDataFromGoogleWithGfaAndFitnessStore(ctx context.Context, gfa *gfa_ds.GoogleFitApp, svc *fitness.Service) error {
	impl.Logger.Debug("pulling hydration data",
		slog.String("gfa_id", gfa.ID.Hex()))

	////
	//// Get `Google Fit` data
	////

	maxTime := time.Now()
	minTime := gfa.LastFetchedAt
	dataType := "hydration" // Note: This is a `Google Fit` specific constant.
	dataset, err := impl.GCP.NotAggregatedDatasets(svc, minTime, maxTime, dataType)
	if err != nil {
		impl.Logger.Error("failed listing hydration",
			slog.Any("error", err))
		return err
	}

	impl.Logger.Debug("pulled hydration data",
		slog.String("gfa_id", gfa.ID.Hex()))

	if len(dataset) == 0 {
		return nil
	}

	////
	//// Convert from `Google Fit` format into our apps format.
	////

	hydration := gcp_a.ParseHydration(dataset)
	impl.Logger.Debug("parsed hydration data",
		slog.Any("data", hydration))

	////
	//// Save into our database.
	////

	return nil
}
