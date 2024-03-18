package crontab

import (
	"log/slog"

	gcp_a "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/adapter/cloudprovider/google"
	googlefitapp_c "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/googlefitapp/controller"
	gfa_ds "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/googlefitapp/datastore"
)

type GoogleFitAppCrontaber interface {
	PullDataFromGoogleJob() error
}

// Handler Creates http request handler
type googleFitAppCrontaberImpl struct {
	Logger                 *slog.Logger
	GCP                    gcp_a.GoogleCloudPlatformAdapter
	GoogleFitAppStorer     gfa_ds.GoogleFitAppStorer
	GoogleFitAppController googlefitapp_c.GoogleFitAppController
}

// NewHandler Constructor
func NewCrontab(
	loggerp *slog.Logger,
	gcpa gcp_a.GoogleCloudPlatformAdapter,
	gfa_storer gfa_ds.GoogleFitAppStorer,
	c googlefitapp_c.GoogleFitAppController,
) GoogleFitAppCrontaber {
	return &googleFitAppCrontaberImpl{
		Logger:                 loggerp,
		GCP:                    gcpa,
		GoogleFitAppStorer:     gfa_storer,
		GoogleFitAppController: c,
	}
}
