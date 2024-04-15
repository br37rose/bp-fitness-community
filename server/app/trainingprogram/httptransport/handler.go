package httptransport

import (
	"log/slog"

	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/trainingprogram/controller"
)

// Handler Creates http request handler
type Handler struct {
	Logger     *slog.Logger
	Controller controller.TrainingprogramController
}

// NewHandler Constructor
func NewHandler(loggerp *slog.Logger, c controller.TrainingprogramController) *Handler {
	return &Handler{
		Logger:     loggerp,
		Controller: c,
	}
}
