package scheduler

import (
	"context"
	"log/slog"

	escheduler "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/adapter/eventscheduler"
	ap_c "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/rankpoint/controller"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/provider/kmutex"
)

type RankPointScheduler interface {
	RunEveryFiveMinutesRanking() error
	RunOnceAndStartImmediatelyRanking() error
}

// Handler Creates http request handler
type rankPointSchedulerImpl struct {
	Logger         *slog.Logger
	Kmutex         kmutex.Provider
	EventScheduler escheduler.EventSchedulerAdapter
	Controller     ap_c.RankPointController
}

// NewHandler Constructor
func NewScheduler(
	loggerp *slog.Logger,
	kmutexp kmutex.Provider,
	es escheduler.EventSchedulerAdapter,
	c ap_c.RankPointController,
) RankPointScheduler {
	return &rankPointSchedulerImpl{
		Logger:         loggerp,
		Kmutex:         kmutexp,
		EventScheduler: es,
		Controller:     c,
	}
}

func (impl *rankPointSchedulerImpl) RunEveryFiveMinutesRanking() error {
	impl.Logger.Debug("scheduled: ranking", slog.String("interval", "every minute"))
	err := impl.EventScheduler.ScheduleEveryFiveMinutesFunc(func() {
		if err := impl.Controller.GenerateGlobalRankingForActiveGoogleFitApps(context.Background()); err != nil {
			impl.Logger.Error("failed ranking",
				slog.Any("error", err))
			return
		}
	})
	if err != nil {
		impl.Logger.Error("ranking error with scheduler", slog.Any("err", err))
	}
	return nil
}

func (impl *rankPointSchedulerImpl) RunOnceAndStartImmediatelyRanking() error {
	impl.Logger.Debug("scheduled: ranking", slog.String("interval", "once"))
	err := impl.EventScheduler.ScheduleOneTimeFunc(func() {
		if err := impl.Controller.GenerateGlobalRankingForActiveGoogleFitApps(context.Background()); err != nil {
			impl.Logger.Error("failed ranking",
				slog.Any("error", err))
			return
		}
	})
	if err != nil {
		impl.Logger.Error("ranking error with scheduler", slog.Any("err", err))
	}
	return nil
}
