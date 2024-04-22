package scheduler

import (
	"log/slog"

	dscheduler "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/adapter/distributedscheduler"
	fp_c "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/fitnessplan/controller"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/provider/kmutex"
)

type FitnessPlanScheduler interface {
	// RunEveryMinuteDeleteAllAnomalousData() error
}

// Handler Creates http request handler
type fitnessPlanSchedulerImpl struct {
	Logger               *slog.Logger
	Kmutex               kmutex.Provider
	DistributedScheduler dscheduler.DistributedSchedulerAdapter
	Controller           fp_c.FitnessPlanController
}

// NewHandler Constructor
func NewScheduler(
	loggerp *slog.Logger,
	kmutexp kmutex.Provider,
	ds dscheduler.DistributedSchedulerAdapter,
	c fp_c.FitnessPlanController,
) FitnessPlanScheduler {
	return &fitnessPlanSchedulerImpl{
		Logger:               loggerp,
		Kmutex:               kmutexp,
		DistributedScheduler: ds,
		Controller:           c,
	}
}
