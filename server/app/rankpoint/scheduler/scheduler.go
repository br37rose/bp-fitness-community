package scheduler

import (
	"context"
	"log/slog"

	dscheduler "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/adapter/distributedscheduler"
	ap_c "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/rankpoint/controller"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/provider/kmutex"
)

type RankPointScheduler interface {
	RunEveryMinuteRankToday() error
	RunEveryMinuteRankThisISOWeek() error
	RunEveryMinuteRankThisMonth() error
	RunEveryMinuteRankThisYear() error
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

func (impl *rankPointSchedulerImpl) RunEveryMinuteRankToday() error {
	impl.Logger.Debug("scheduled: rank today", slog.String("interval", "every minute"))
	err := impl.DistributedScheduler.ScheduleEveryMinuteFunc(func() {
		if err := impl.Controller.GenerateGlobalRankingForTodayUsingActiveGoogleFitApps(context.Background()); err != nil {
			impl.Logger.Error("failed ranking today",
				slog.Any("error", err))
			return
		}
	})
	if err != nil {
		impl.Logger.Error("rank today error with scheduler", slog.Any("err", err))
	}
	return nil
}

func (impl *rankPointSchedulerImpl) RunEveryMinuteRankThisISOWeek() error {
	impl.Logger.Debug("scheduled: rank this iso week", slog.String("interval", "every minute"))
	err := impl.DistributedScheduler.ScheduleEveryMinuteFunc(func() {
		if err := impl.Controller.GenerateGlobalRankingForThisISOWeekUsingActiveGoogleFitApps(context.Background()); err != nil {
			impl.Logger.Error("failed ranking this iso week",
				slog.Any("error", err))
			return
		}
	})
	if err != nil {
		impl.Logger.Error("rank this iso week error with scheduler", slog.Any("err", err))
	}
	return nil
}

func (impl *rankPointSchedulerImpl) RunEveryMinuteRankThisMonth() error {
	impl.Logger.Debug("scheduled: rank this month", slog.String("interval", "every minute"))
	err := impl.DistributedScheduler.ScheduleEveryMinuteFunc(func() {
		if err := impl.Controller.GenerateGlobalRankingForThisMonthUsingActiveGoogleFitApps(context.Background()); err != nil {
			impl.Logger.Error("failed ranking this month",
				slog.Any("error", err))
			return
		}
	})
	if err != nil {
		impl.Logger.Error("rank this month error with scheduler", slog.Any("err", err))
	}
	return nil
}

func (impl *rankPointSchedulerImpl) RunEveryMinuteRankThisYear() error {
	impl.Logger.Debug("scheduled: rank this year", slog.String("interval", "every minute"))
	err := impl.DistributedScheduler.ScheduleEveryMinuteFunc(func() {
		if err := impl.Controller.GenerateGlobalRankingForThisYearUsingActiveGoogleFitApps(context.Background()); err != nil {
			impl.Logger.Error("failed ranking this year",
				slog.Any("error", err))
			return
		}
	})
	if err != nil {
		impl.Logger.Error("rank this year error with scheduler", slog.Any("err", err))
	}
	return nil
}
