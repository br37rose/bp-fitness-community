package httptransport

import (
	"log/slog"

	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/workout/controller"
)

// Handler Creates http request handler
type Handler struct {
	Logger     *slog.Logger
	Controller controller.WorkoutController
}

// NewHandler Constructor
func NewHandler(loggerp *slog.Logger, c controller.WorkoutController) *Handler {
	return &Handler{
		Logger:     loggerp,
		Controller: c,
	}
}
