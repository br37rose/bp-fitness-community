package crontab

import (
	"log/slog"

	gcp_a "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/adapter/cloudprovider/google"
	dp_s "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/datapoint/datastore"
	googlefitapp_c "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/googlefitapp/controller"
	gfa_ds "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/googlefitapp/datastore"
	googlefitdp_s "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/googlefitdatapoint/datastore"
	user_s "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/user/datastore"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/provider/kmutex"
)

type GoogleFitAppCrontaber interface {
	// RefreshTokensFromGoogleJob() error
	PullDataFromGoogleJob() error
}

// Handler Creates http request handler
type googleFitAppCrontaberImpl struct {
	Logger                   *slog.Logger
	Kmutex                   kmutex.Provider
	GCP                      gcp_a.GoogleCloudPlatformAdapter
	DataPointStorer          dp_s.DataPointStorer
	GoogleFitDataPointStorer googlefitdp_s.GoogleFitDataPointStorer
	GoogleFitAppStorer       gfa_ds.GoogleFitAppStorer
	GoogleFitAppController   googlefitapp_c.GoogleFitAppController
	UserStorer               user_s.UserStorer
}

// NewHandler Constructor
func NewCrontab(
	loggerp *slog.Logger,
	kmutexp kmutex.Provider,
	gcpa gcp_a.GoogleCloudPlatformAdapter,
	dp dp_s.DataPointStorer,
	gfdp googlefitdp_s.GoogleFitDataPointStorer,
	gfa_storer gfa_ds.GoogleFitAppStorer,
	c googlefitapp_c.GoogleFitAppController,
	usr_storer user_s.UserStorer,
) GoogleFitAppCrontaber {
	return &googleFitAppCrontaberImpl{
		Logger:                   loggerp,
		Kmutex:                   kmutexp,
		GCP:                      gcpa,
		DataPointStorer:          dp,
		GoogleFitDataPointStorer: gfdp,
		GoogleFitAppStorer:       gfa_storer,
		GoogleFitAppController:   c,
		UserStorer:               usr_storer,
	}
}
