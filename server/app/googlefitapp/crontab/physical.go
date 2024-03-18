package crontab

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"time"

	gfa_ds "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/googlefitapp/datastore"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (impl *googleFitAppCrontaberImpl) PullDataFromGoogleJob() error {
	ctx := context.Background()
	impl.Logger.Debug("pulling data from google cloud platform...")
	gfaIDs, err := impl.GoogleFitAppStorer.ListPhysicalIDsByStatus(ctx, gfa_ds.StatusActive)
	if err != nil {
		impl.Logger.Error("failed listing google fit apps by status",
			slog.Any("error", err))
		return err
	}
	for _, gfaID := range gfaIDs {
		impl.Logger.Debug("processing google fit app",
			slog.String("gfa_id", gfaID.Hex()),
		)
		if err := impl.pullDataFromGoogle(ctx, gfaID); err != nil {
			impl.Logger.Error("failed pulling data for google fit app",
				slog.Any("error", err))
			return err
		}
	}
	return nil
}

func (impl *googleFitAppCrontaberImpl) pullDataFromGoogle(ctx context.Context, gfaID primitive.ObjectID) error {
	// Get our database record.
	gfa, err := impl.GoogleFitAppStorer.GetByID(ctx, gfaID)
	if err != nil {
		impl.Logger.Error("failed getting google fit app from database",
			slog.Any("error", err))
		return err
	}
	if gfa == nil {
		err := fmt.Errorf("google fit app does not exist for id: %s", gfaID.Hex())
		return err
	}

	// Authenticated http client for a specific user's account. Note: No need
	// for refresh token handling as it's already handled!
	client := impl.GCP.NewHTTPClientFromToken(gfa.Token)
	if client == nil {
		err := fmt.Errorf("google fit app authenticated client does not exist for token: %v", gfa.Token)
		return err
	}

	// Load up the Google Fitness Store.

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

	// Get data
	maxTime := time.Now()
	minTime := time.Now()
	dataType := "hydration"
	dataset, err := impl.GCP.NotAggregatedDatasets(svc, minTime, maxTime, dataType)
	if err != nil {
		impl.Logger.Error("failed listing hydration",
			slog.Any("error", err))
		return err
	}
	log.Println("-->", dataset)

	return nil
}
