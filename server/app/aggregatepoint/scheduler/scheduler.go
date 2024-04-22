package scheduler

import (
	"context"
	"log/slog"

	dscheduler "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/adapter/distributedscheduler"
	ap_c "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/aggregatepoint/controller"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/provider/kmutex"
)

type AggregatePointScheduler interface {
	RunEveryMinuteAggregateThisHour() error
	RunEveryMinuteAggregateLastHour() error
	RunEveryMinuteAggregateToday() error
	RunEveryMinuteAggregateYesterday() error
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

func (impl *aggregatePointSchedulerImpl) RunEveryMinuteAggregateThisHour() error {
	impl.Logger.Debug("scheduled: aggregate this hour", slog.String("interval", "every minute"))
	err := impl.DistributedScheduler.ScheduleEveryMinuteFunc(func() {
		if err := impl.Controller.AggregateThisHourForAllActiveGoogleFitApps(context.Background()); err != nil {
			impl.Logger.Error("failed aggregating this hour",
				slog.Any("error", err))
		}
	})
	if err != nil {
		impl.Logger.Error("aggregate this hour error with scheduler", slog.Any("err", err))
		return err
	}
	return nil
}

func (impl *aggregatePointSchedulerImpl) RunEveryMinuteAggregateLastHour() error {
	impl.Logger.Debug("scheduled: aggregate last hour", slog.String("interval", "every minute"))
	err := impl.DistributedScheduler.ScheduleEveryMinuteFunc(func() {
		if err := impl.Controller.AggregateLastHourForAllActiveGoogleFitApps(context.Background()); err != nil {
			impl.Logger.Error("failed aggregating last hour",
				slog.Any("error", err))
		}
	})
	if err != nil {
		impl.Logger.Error("aggregate last hour error with scheduler", slog.Any("err", err))
		return err
	}
	return nil
}

func (impl *aggregatePointSchedulerImpl) RunEveryMinuteAggregateToday() error {
	impl.Logger.Debug("scheduled: aggregate today", slog.String("interval", "every minute"))
	err := impl.DistributedScheduler.ScheduleEveryMinuteFunc(func() {
		if err := impl.Controller.AggregateTodayForAllActiveGoogleFitApps(context.Background()); err != nil {
			impl.Logger.Error("failed aggregating today",
				slog.Any("error", err))
		}
	})
	if err != nil {
		impl.Logger.Error("aggregate today error with scheduler", slog.Any("err", err))
		return err
	}
	return nil
}

func (impl *aggregatePointSchedulerImpl) RunEveryMinuteAggregateYesterday() error {
	impl.Logger.Debug("scheduled: aggregate yesterday", slog.String("interval", "every minute"))
	err := impl.DistributedScheduler.ScheduleEveryMinuteFunc(func() {
		if err := impl.Controller.AggregateYesterdayForAllActiveGoogleFitApps(context.Background()); err != nil {
			impl.Logger.Error("failed aggregating yesterday",
				slog.Any("error", err))
		}
	})
	if err != nil {
		impl.Logger.Error("aggregate yesterday error with scheduler", slog.Any("err", err))
		return err
	}
	return nil
}

// func (impl *aggregatePointSchedulerImpl) RunEveryMinuteAggregateToday() error {
// 	impl.Logger.Debug("scheduled: aggregate today", slog.String("interval", "every minute"))
// 	err := impl.DistributedScheduler.ScheduleEveryMinuteFunc(func() {
// 		if err := impl.Controller.AggregateTodayForAllActiveGoogleFitApps(context.Background()); err != nil {
// 			impl.Logger.Error("failed aggregating today",
// 				slog.Any("error", err))
// 		}
// 	})
// 	if err != nil {
// 		impl.Logger.Error("aggregate today error with scheduler", slog.Any("err", err))
// 		return err
// 	}
// 	return nil
// }
