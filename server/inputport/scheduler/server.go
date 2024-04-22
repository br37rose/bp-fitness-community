package scheduler

import (
	"log/slog"

	dscheduler "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/adapter/distributedscheduler"
	googlefitapp_task "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/googlefitapp/scheduler"
	googlefitdp_task "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/googlefitdatapoint/scheduler"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/config"
)

type InputPortServer interface {
	Run()
	Shutdown()
}

type schedulerInputPort struct {
	Config                      *config.Conf
	Logger                      *slog.Logger
	DistributedScheduler        dscheduler.DistributedSchedulerAdapter
	GoogleFitDataPointScheduler googlefitdp_task.GoogleFitDataPointScheduler
	GoogleFitAppScheduler       googlefitapp_task.GoogleFitAppScheduler
}

func NewInputPort(
	configp *config.Conf,
	loggerp *slog.Logger,
	ds dscheduler.DistributedSchedulerAdapter,
	gfdp googlefitdp_task.GoogleFitDataPointScheduler,
	gfa googlefitapp_task.GoogleFitAppScheduler,
) InputPortServer {
	// Initialize.

	// Create our server controller.
	p := &schedulerInputPort{
		Config:                      configp,
		Logger:                      loggerp,
		DistributedScheduler:        ds,
		GoogleFitDataPointScheduler: gfdp,
		GoogleFitAppScheduler:       gfa,
	}

	return p
}

func (port *schedulerInputPort) Run() {
	port.Logger.Info("Scheduler server starting...")
	port.DistributedScheduler.Start()
	port.ping()

	// Schedule the following background tasks to run.
	if err := port.GoogleFitDataPointScheduler.RunEveryMinuteDeleteAllAnomalousData(); err != nil {
		port.Logger.Error("scheduler has error", slog.Any("err", err))
	}
}

func (port *schedulerInputPort) Shutdown() {
	port.Logger.Info("Scheduler server shutdown")
	port.DistributedScheduler.Shutdown()
}

// ping function will send out a one-time task to verify the cluster is
// successfullly connected.
func (port *schedulerInputPort) ping() {
	err := port.DistributedScheduler.ScheduleOneTimeFunc(func() {
		port.Logger.Info("Scheduler pinged")
	})
	if err != nil {
		port.Logger.Error("error with pinging scheduler", slog.Any("err", err))
	}
}
