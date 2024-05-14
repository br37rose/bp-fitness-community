package scheduler

import (
	"context"
	"log/slog"

	escheduler "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/adapter/eventscheduler"
	ap_c "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/aggregatepoint/controller"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/provider/kmutex"
)

type AggregatePointScheduler interface {
	RunEveryFiveMinutesAggregation() error
	RunOnceAndStartImmediatelyAggregation() error
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

func (impl *aggregatePointSchedulerImpl) RunEveryFiveMinutesAggregation() error {
	impl.Logger.Debug("scheduled: aggregation", slog.String("interval", "every minute"))
	err := impl.EventScheduler.ScheduleEveryFiveMinutesFunc(func() {
		if err := impl.Controller.AggregateForAllActiveGoogleFitApps(context.Background()); err != nil {
			impl.Logger.Error("failed aggregation",
				slog.Any("error", err))
		}
	})
	if err != nil {
		impl.Logger.Error("aggregation error with scheduler", slog.Any("err", err))
		return err
	}
	return nil
}

func (impl *aggregatePointSchedulerImpl) RunOnceAndStartImmediatelyAggregation() error {
	impl.Logger.Debug("scheduled: aggregation", slog.String("interval", "once"))
	err := impl.EventScheduler.ScheduleOneTimeFunc(func() {
		if err := impl.Controller.AggregateForAllActiveGoogleFitApps(context.Background()); err != nil {
			impl.Logger.Error("failed aggregation",
				slog.Any("error", err))
		}
	})
	if err != nil {
		impl.Logger.Error("aggregation error with scheduler", slog.Any("err", err))
		return err
	}
	return nil
}
