package scheduler

import (
	"log/slog"

	escheduler "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/adapter/eventscheduler"
	googlefitapp_c "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/googlefitapp/controller"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/provider/kmutex"
)

type GoogleFitAppScheduler interface {
	RunEveryMinuteRefreshTokensFromGoogle() error
	RunEveryProcessAllQueuedData() error
	RunEveryFifteenMinutesPullDataFromGoogle() error
}

// Handler Creates http request handler
type googleFitAppSchedulerImpl struct {
	Logger         *slog.Logger
	Kmutex         kmutex.Provider
	EventScheduler escheduler.EventSchedulerAdapter
	Controller     googlefitapp_c.GoogleFitAppController
}

// NewHandler Constructor
func NewScheduler(
	loggerp *slog.Logger,
	kmutexp kmutex.Provider,
	es escheduler.EventSchedulerAdapter,
	c googlefitapp_c.GoogleFitAppController,
) GoogleFitAppScheduler {
	return &googleFitAppSchedulerImpl{
		Logger:         loggerp,
		Kmutex:         kmutexp,
		EventScheduler: es,
		Controller:     c,
	}
}

func (impl *googleFitAppSchedulerImpl) RunEveryMinuteRefreshTokensFromGoogle() error {
	impl.Logger.Debug("scheduled: refresh token", slog.String("interval", "every minute"))
	err := impl.EventScheduler.ScheduleEveryMinuteFunc(func() {
		impl.Logger.Debug("running refresh token...")
		if err := impl.Controller.RefreshTokensFromGoogle; err != nil {
			impl.Logger.Error("refresh token error with scheduler", slog.Any("err", err))
		}
		impl.Logger.Debug("finished refresh token")
	})
	if err != nil {
		impl.Logger.Error("error with scheduler", slog.Any("err", err))
	}
	return nil
}

func (impl *googleFitAppSchedulerImpl) RunEveryProcessAllQueuedData() error {
	impl.Logger.Debug("scheduled: process queued data", slog.String("interval", "every minute"))
	err := impl.EventScheduler.ScheduleEveryMinuteFunc(func() {
		impl.Logger.Debug("running process queued data...")
		if err := impl.Controller.ProcessAllQueuedData; err != nil {
			impl.Logger.Error("process queued data error with scheduler", slog.Any("err", err))
		}
		impl.Logger.Debug("finished process queued data")
	})
	if err != nil {
		impl.Logger.Error("error with scheduler", slog.Any("err", err))
	}
	return nil
}

func (impl *googleFitAppSchedulerImpl) RunEveryFifteenMinutesPullDataFromGoogle() error {
	impl.Logger.Debug("scheduled: pull data from google", slog.String("interval", "every minute"))
	err := impl.EventScheduler.ScheduleEveryMinuteFunc(func() {
		impl.Logger.Debug("running pull data from google...")
		if err := impl.Controller.PullDataFromGoogle; err != nil {
			impl.Logger.Error("pull data from google error with scheduler", slog.Any("err", err))
		}
		impl.Logger.Debug("finished pull data from google")
	})
	if err != nil {
		impl.Logger.Error("error with scheduler", slog.Any("err", err))
	}
	return nil
}
