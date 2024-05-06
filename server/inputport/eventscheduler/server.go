package eventscheduler

import (
	"log/slog"

	escheduler "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/adapter/eventscheduler"
	ap_task "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/aggregatepoint/scheduler"
	fp_task "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/fitnessplan/scheduler"
	googlefitapp_task "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/googlefitapp/scheduler"
	googlefitdp_task "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/googlefitdatapoint/scheduler"
	rp_task "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/rankpoint/scheduler"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/config"
)

type InputPortServer interface {
	Run()
	Shutdown()
}

type schedulerInputPort struct {
	Config                      *config.Conf
	Logger                      *slog.Logger
	EventScheduler              escheduler.EventSchedulerAdapter
	GoogleFitDataPointScheduler googlefitdp_task.GoogleFitDataPointScheduler
	GoogleFitAppScheduler       googlefitapp_task.GoogleFitAppScheduler
	AggregatePointScheduler     ap_task.AggregatePointScheduler
	RankPointScheduler          rp_task.RankPointScheduler
	FitnessPlanScheduler        fp_task.FitnessPlanScheduler
}

func NewInputPort(
	configp *config.Conf,
	loggerp *slog.Logger,
	es escheduler.EventSchedulerAdapter,
	gfdp googlefitdp_task.GoogleFitDataPointScheduler,
	gfa googlefitapp_task.GoogleFitAppScheduler,
	ap ap_task.AggregatePointScheduler,
	rp rp_task.RankPointScheduler,
	fp fp_task.FitnessPlanScheduler,
) InputPortServer {
	// Initialize.

	// Create our server controller.
	p := &schedulerInputPort{
		Config:                      configp,
		Logger:                      loggerp,
		EventScheduler:              es,
		GoogleFitDataPointScheduler: gfdp,
		GoogleFitAppScheduler:       gfa,
		AggregatePointScheduler:     ap,
		RankPointScheduler:          rp,
		FitnessPlanScheduler:        fp,
	}

	return p
}

func (port *schedulerInputPort) Run() {
	port.Logger.Info("scheduler server starting...")
	port.EventScheduler.Start()
	port.ping()

	defer func() {
		// Schedule one-time jobs.
		if err := port.GoogleFitAppScheduler.RunOnceAndStartImmediatelyProcessAllQueuedData(); err != nil {
			port.Logger.Error("scheduler failed for processing queued data ", slog.Any("err", err))
		}
		if err := port.AggregatePointScheduler.RunOnceAndStartImmediatelyAggregation(); err != nil {
			port.Logger.Error("scheduler has error", slog.Any("err", err))
		}
		if err := port.RankPointScheduler.RunOnceAndStartImmediatelyRanking(); err != nil {
			port.Logger.Error("scheduler has error", slog.Any("err", err))
		}
	}()

	// Schedule the following background tasks to run.
	if err := port.GoogleFitDataPointScheduler.RunEveryMinuteDeleteAllAnomalousData(); err != nil {
		port.Logger.Error("scheduler has error", slog.Any("err", err))
	}
	if err := port.GoogleFitAppScheduler.RunEveryMinuteRefreshTokensFromGoogle(); err != nil {
		port.Logger.Error("scheduler has error", slog.Any("err", err))
	}
	if err := port.GoogleFitAppScheduler.RunEveryMinuteProcessAllQueuedData(); err != nil {
		port.Logger.Error("scheduler has error", slog.Any("err", err))
	}
	if err := port.GoogleFitAppScheduler.RunEveryFifteenMinutesPullDataFromGoogle(); err != nil {
		port.Logger.Error("scheduler has error", slog.Any("err", err))
	}
	if err := port.AggregatePointScheduler.RunEveryMinuteAggregation(); err != nil {
		port.Logger.Error("scheduler has error", slog.Any("err", err))
	}
	if err := port.RankPointScheduler.RunEveryMinuteRanking(); err != nil {
		port.Logger.Error("scheduler has error", slog.Any("err", err))
	}
	if err := port.FitnessPlanScheduler.RunEveryMinuteUpdateFitnessPlans(); err != nil {
		port.Logger.Error("UpdateFitnessPlans scheduler has error", slog.Any("err", err))
	}

}

func (port *schedulerInputPort) Shutdown() {
	port.Logger.Info("scheduler server shutdown")
	port.EventScheduler.Shutdown()
}

// ping function will send out a one-time task to verify the cluster is
// successfullly connected.
func (port *schedulerInputPort) ping() {
	err := port.EventScheduler.ScheduleOneTimeFunc(func() {
		port.Logger.Info("scheduler pinged")
	})
	if err != nil {
		port.Logger.Error("error with pinging scheduler", slog.Any("err", err))
	}
}
