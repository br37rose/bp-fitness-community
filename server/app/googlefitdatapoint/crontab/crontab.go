package crontab

import (
	"log/slog"

	gcp_a "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/adapter/cloudprovider/google"
	dp_s "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/datapoint/datastore"
	googlefitdp_c "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/googlefitdatapoint/controller"
	googlefitdp_s "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/googlefitdatapoint/datastore"
	user_s "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/user/datastore"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/provider/kmutex"
)

type GoogleFitDataPointCrontaber interface {
	DeleteAllAnomalousData() error
}

// Handler Creates http request handler
type googleFitAppCrontaberImpl struct {
	Logger                       *slog.Logger
	Kmutex                       kmutex.Provider
	GCP                          gcp_a.GoogleCloudPlatformAdapter
	DataPointStorer              dp_s.DataPointStorer
	GoogleFitDataPointStorer     googlefitdp_s.GoogleFitDataPointStorer
	GoogleFitDataPointController googlefitdp_c.GoogleFitDataPointController
	UserStorer                   user_s.UserStorer
}

// NewHandler Constructor
func NewCrontab(
	loggerp *slog.Logger,
	kmutexp kmutex.Provider,
	gcpa gcp_a.GoogleCloudPlatformAdapter,
	dp dp_s.DataPointStorer,
	gfdp googlefitdp_s.GoogleFitDataPointStorer,
	c googlefitdp_c.GoogleFitDataPointController,
	usr_storer user_s.UserStorer,
) GoogleFitDataPointCrontaber {
	return &googleFitAppCrontaberImpl{
		Logger:                       loggerp,
		Kmutex:                       kmutexp,
		GCP:                          gcpa,
		DataPointStorer:              dp,
		GoogleFitDataPointStorer:     gfdp,
		GoogleFitDataPointController: c,
		UserStorer:                   usr_storer,
	}
}
