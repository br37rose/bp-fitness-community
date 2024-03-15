package httptransport

import (
	exercise_c "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/exercise/controller"
	"log/slog"
)

// Handler Creates http request handler
type Handler struct {
	Logger     *slog.Logger
	Controller exercise_c.ExerciseController
}

// NewHandler Constructor
func NewHandler(loggerp *slog.Logger, c exercise_c.ExerciseController) *Handler {
	return &Handler{
		Logger:     loggerp,
		Controller: c,
	}
}
