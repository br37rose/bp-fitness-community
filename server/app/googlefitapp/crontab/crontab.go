package crontab

import (
	"log/slog"

	googlefitapp_c "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/googlefitapp/controller"
)

type GoogleFitAppCrontaber interface {
	PullJob() error
}

// Handler Creates http request handler
type googleFitAppCrontaberImpl struct {
	Logger                 *slog.Logger
	GoogleFitAppController googlefitapp_c.GoogleFitAppController
}

// NewHandler Constructor
func NewCrontab(loggerp *slog.Logger, c googlefitapp_c.GoogleFitAppController) GoogleFitAppCrontaber {
	return &googleFitAppCrontaberImpl{
		Logger:                 loggerp,
		GoogleFitAppController: c,
	}
}

func (impl *googleFitAppCrontaberImpl) PullJob() error {
	impl.Logger.Error("pull")
	return nil
}
