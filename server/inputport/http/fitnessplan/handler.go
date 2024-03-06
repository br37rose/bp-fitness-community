package fitnessplan

import (
	fitnessplan_c "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/fitnessplan/controller"
	"log/slog"
)

// Handler Creates http request handler
type Handler struct {
	Logger     *slog.Logger
	Controller fitnessplan_c.FitnessPlanController
}

// NewHandler Constructor
func NewHandler(loggerp *slog.Logger, c fitnessplan_c.FitnessPlanController) *Handler {
	return &Handler{
		Logger:     loggerp,
		Controller: c,
	}
}
