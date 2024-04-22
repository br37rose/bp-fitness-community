package scheduler

import (
	"log/slog"

	gcp_a "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/adapter/cloudprovider/google"
	dscheduler "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/adapter/distributedscheduler"
	dp_s "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/datapoint/datastore"
	googlefitdp_s "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/googlefitdatapoint/datastore"
	user_s "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/user/datastore"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/provider/kmutex"
)

type GoogleFitDataPointScheduler interface {
	RunEveryMinuteDeleteAllAnomalousData() error
}

// Handler Creates http request handler
type googleFitAppSchedulerImpl struct {
	Logger                   *slog.Logger
	Kmutex                   kmutex.Provider
	GCP                      gcp_a.GoogleCloudPlatformAdapter
	DistributedScheduler     dscheduler.DistributedSchedulerAdapter
	DataPointStorer          dp_s.DataPointStorer
	GoogleFitDataPointStorer googlefitdp_s.GoogleFitDataPointStorer
	UserStorer               user_s.UserStorer
}

// NewHandler Constructor
func NewScheduler(
	loggerp *slog.Logger,
	kmutexp kmutex.Provider,
	gcpa gcp_a.GoogleCloudPlatformAdapter,
	ds dscheduler.DistributedSchedulerAdapter,
	dp dp_s.DataPointStorer,
	gfdp googlefitdp_s.GoogleFitDataPointStorer,
	usr_storer user_s.UserStorer,
) GoogleFitDataPointScheduler {
	return &googleFitAppSchedulerImpl{
		Logger:                   loggerp,
		Kmutex:                   kmutexp,
		GCP:                      gcpa,
		DistributedScheduler:     ds,
		DataPointStorer:          dp,
		GoogleFitDataPointStorer: gfdp,
		UserStorer:               usr_storer,
	}
}
