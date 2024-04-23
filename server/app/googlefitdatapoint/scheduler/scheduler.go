package scheduler

import (
	"log/slog"

	gcp_a "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/adapter/cloudprovider/google"
	dscheduler "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/adapter/distributedscheduler"
	dp_c "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/googlefitdatapoint/controller"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/provider/kmutex"
)

type GoogleFitDataPointScheduler interface {
	RunEveryMinuteDeleteAllAnomalousData() error
}

// Handler Creates http request handler
type googleFitDataPointSchedulerImpl struct {
	Logger               *slog.Logger
	Kmutex               kmutex.Provider
	DistributedScheduler dscheduler.DistributedSchedulerAdapter
	Controller           dp_c.GoogleFitDataPointController
}

// NewHandler Constructor
func NewScheduler(
	loggerp *slog.Logger,
	kmutexp kmutex.Provider,
	gcpa gcp_a.GoogleCloudPlatformAdapter,
	ds dscheduler.DistributedSchedulerAdapter,
	c dp_c.GoogleFitDataPointController,
) GoogleFitDataPointScheduler {
	return &googleFitDataPointSchedulerImpl{
		Logger:               loggerp,
		Kmutex:               kmutexp,
		DistributedScheduler: ds,
		Controller:           c,
	}
}

func (impl *googleFitDataPointSchedulerImpl) RunEveryMinuteDeleteAllAnomalousData() error {
	impl.Logger.Debug("scheduled every minute: delete all anomalous data")
	err := impl.DistributedScheduler.ScheduleEveryMinuteFunc(func() {
		impl.Logger.Debug("running delete all anomalous data...")
		if err := impl.Controller.DeleteAllAnomalousData(); err != nil {
			impl.Logger.Error("error with deleting all anomalous data from scheduler", slog.Any("err", err))
		}
		impl.Logger.Debug("finished running delete all anomalous data")
	})
	if err != nil {
		impl.Logger.Error("error with pinging scheduler", slog.Any("err", err))
	}
	return nil
}
