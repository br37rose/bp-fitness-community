package scheduler

import (
	"context"
	"log/slog"

	escheduler "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/adapter/eventscheduler"
	ap_c "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/aggregatepoint/controller"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/provider/kmutex"
)

type AggregatePointScheduler interface {
	RunEveryMinuteAggregateThisHour() error
	RunEveryMinuteAggregateLastHour() error
	RunEveryMinuteAggregateToday() error
	RunEveryMinuteAggregateYesterday() error
	RunEveryMinuteAggregateThisISOWeek() error
	RunEveryMinuteAggregateLastISOWeek() error
	RunEveryMinuteAggregateThisMonth() error
	RunEveryMinuteAggregateLastMonth() error
	RunEveryMinuteAggregateThisYear() error
	RunEveryMinuteAggregateLastYear() error
}

// Handler Creates http request handler
type aggregatePointSchedulerImpl struct {
	Logger         *slog.Logger
	Kmutex         kmutex.Provider
	EventScheduler escheduler.EventSchedulerAdapter
	Controller     ap_c.AggregatePointController
}

// NewHandler Constructor
func NewScheduler(
	loggerp *slog.Logger,
	kmutexp kmutex.Provider,
	es escheduler.EventSchedulerAdapter,
	c ap_c.AggregatePointController,
) AggregatePointScheduler {
	return &aggregatePointSchedulerImpl{
		Logger:         loggerp,
		Kmutex:         kmutexp,
		EventScheduler: es,
		Controller:     c,
	}
}

func (impl *aggregatePointSchedulerImpl) RunEveryMinuteAggregateThisHour() error {
	impl.Logger.Debug("scheduled: aggregate this hour", slog.String("interval", "every minute"))
	err := impl.EventScheduler.ScheduleEveryMinuteFunc(func() {
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
	err := impl.EventScheduler.ScheduleEveryMinuteFunc(func() {
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
	err := impl.EventScheduler.ScheduleEveryMinuteFunc(func() {
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
	err := impl.EventScheduler.ScheduleEveryMinuteFunc(func() {
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

func (impl *aggregatePointSchedulerImpl) RunEveryMinuteAggregateThisISOWeek() error {
	impl.Logger.Debug("scheduled: aggregate this iso week", slog.String("interval", "every minute"))
	err := impl.EventScheduler.ScheduleEveryMinuteFunc(func() {
		if err := impl.Controller.AggregateThisISOWeekForAllActiveGoogleFitApps(context.Background()); err != nil {
			impl.Logger.Error("failed aggregating this iso week ",
				slog.Any("error", err))
		}
	})
	if err != nil {
		impl.Logger.Error("aggregate this iso week error with scheduler", slog.Any("err", err))
		return err
	}
	return nil
}

func (impl *aggregatePointSchedulerImpl) RunEveryMinuteAggregateLastISOWeek() error {
	impl.Logger.Debug("scheduled: aggregate last iso week", slog.String("interval", "every minute"))
	err := impl.EventScheduler.ScheduleEveryMinuteFunc(func() {
		if err := impl.Controller.AggregateLastISOWeekForAllActiveGoogleFitApps(context.Background()); err != nil {
			impl.Logger.Error("failed aggregating last iso week",
				slog.Any("error", err))
		}
	})
	if err != nil {
		impl.Logger.Error("aggregate last iso week error with scheduler", slog.Any("err", err))
		return err
	}
	return nil
}

func (impl *aggregatePointSchedulerImpl) RunEveryMinuteAggregateThisMonth() error {
	impl.Logger.Debug("scheduled: aggregate this month", slog.String("interval", "every minute"))
	err := impl.EventScheduler.ScheduleEveryMinuteFunc(func() {
		if err := impl.Controller.AggregateThisMonthForAllActiveGoogleFitApps(context.Background()); err != nil {
			impl.Logger.Error("failed aggregating this month",
				slog.Any("error", err))
		}
	})
	if err != nil {
		impl.Logger.Error("aggregate this month error with scheduler", slog.Any("err", err))
		return err
	}
	return nil
}

func (impl *aggregatePointSchedulerImpl) RunEveryMinuteAggregateLastMonth() error {
	impl.Logger.Debug("scheduled: aggregate last month", slog.String("interval", "every minute"))
	err := impl.EventScheduler.ScheduleEveryMinuteFunc(func() {
		if err := impl.Controller.AggregateLastMonthForAllActiveGoogleFitApps(context.Background()); err != nil {
			impl.Logger.Error("failed aggregating last month",
				slog.Any("error", err))
		}
	})
	if err != nil {
		impl.Logger.Error("aggregate last month error with scheduler", slog.Any("err", err))
		return err
	}
	return nil
}

func (impl *aggregatePointSchedulerImpl) RunEveryMinuteAggregateThisYear() error {
	impl.Logger.Debug("scheduled: aggregate this year", slog.String("interval", "every minute"))
	err := impl.EventScheduler.ScheduleEveryMinuteFunc(func() {
		if err := impl.Controller.AggregateThisYearForAllActiveGoogleFitApps(context.Background()); err != nil {
			impl.Logger.Error("failed aggregating this year",
				slog.Any("error", err))
		}
	})
	if err != nil {
		impl.Logger.Error("aggregate this year error with scheduler", slog.Any("err", err))
		return err
	}
	return nil
}

func (impl *aggregatePointSchedulerImpl) RunEveryMinuteAggregateLastYear() error {
	impl.Logger.Debug("scheduled: aggregate last year", slog.String("interval", "every minute"))
	err := impl.EventScheduler.ScheduleEveryMinuteFunc(func() {
		if err := impl.Controller.AggregateLastYearForAllActiveGoogleFitApps(context.Background()); err != nil {
			impl.Logger.Error("failed aggregating last year",
				slog.Any("error", err))
		}
	})
	if err != nil {
		impl.Logger.Error("aggregate last year error with scheduler", slog.Any("err", err))
		return err
	}
	return nil
}

// func (impl *aggregatePointSchedulerImpl) RunEveryMinuteAggregateToday() error {
// 	impl.Logger.Debug("scheduled: aggregate today", slog.String("interval", "every minute"))
// 	err := impl.EventScheduler.ScheduleEveryMinuteFunc(func() {
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
//
// func (impl *aggregatePointSchedulerImpl) RunEveryMinuteAggregateToday() error {
// 	impl.Logger.Debug("scheduled: aggregate today", slog.String("interval", "every minute"))
// 	err := impl.EventScheduler.ScheduleEveryMinuteFunc(func() {
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
