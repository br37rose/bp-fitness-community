package scheduler

import (
	"log/slog"

	dscheduler "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/adapter/distributedscheduler"
	ap_c "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/rankpoint/controller"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/provider/kmutex"
)

type RankPointScheduler interface {
	// RunEveryMinuteDeleteAllAnomalousData() error
}

// Handler Creates http request handler
type rankPointSchedulerImpl struct {
	Logger               *slog.Logger
	Kmutex               kmutex.Provider
	DistributedScheduler dscheduler.DistributedSchedulerAdapter
	Controller           ap_c.RankPointController
}

// NewHandler Constructor
func NewScheduler(
	loggerp *slog.Logger,
	kmutexp kmutex.Provider,
	ds dscheduler.DistributedSchedulerAdapter,
	c ap_c.RankPointController,
) RankPointScheduler {
	return &rankPointSchedulerImpl{
		Logger:               loggerp,
		Kmutex:               kmutexp,
		DistributedScheduler: ds,
		Controller:           c,
	}
}
