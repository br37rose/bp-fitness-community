package scheduler

import (
	"log/slog"

	dscheduler "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/adapter/distributedscheduler"
	ap_c "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/aggregatepoint/controller"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/provider/kmutex"
)

type AggregatePointScheduler interface {
	RunEveryMinuteDeleteAllAnomalousData() error
}

// Handler Creates http request handler
type aggregatePointSchedulerImpl struct {
	Logger               *slog.Logger
	Kmutex               kmutex.Provider
	DistributedScheduler dscheduler.DistributedSchedulerAdapter
	Controller           ap_c.AggregatePointController
}

// NewHandler Constructor
func NewScheduler(
	loggerp *slog.Logger,
	kmutexp kmutex.Provider,
	ds dscheduler.DistributedSchedulerAdapter,
	c ap_c.AggregatePointController,
) AggregatePointScheduler {
	return &aggregatePointSchedulerImpl{
		Logger:               loggerp,
		Kmutex:               kmutexp,
		DistributedScheduler: ds,
		Controller:           c,
	}
}
