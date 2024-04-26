package httptransport

import (
	"log/slog"

	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/fitnesschallenge/controller"
)

// Handler Creates http request handler
type Handler struct {
	Logger     *slog.Logger
	Controller controller.FitnessChallengeController
}

// NewHandler Constructor
func NewHandler(loggerp *slog.Logger, c controller.FitnessChallengeController) *Handler {
	return &Handler{
		Logger:     loggerp,
		Controller: c,
	}
}
