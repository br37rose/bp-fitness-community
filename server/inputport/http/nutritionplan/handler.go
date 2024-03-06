package nutritionplan

import (
	nutritionplan_c "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/nutritionplan/controller"
	"log/slog"
)

// Handler Creates http request handler
type Handler struct {
	Logger     *slog.Logger
	Controller nutritionplan_c.NutritionPlanController
}

// NewHandler Constructor
func NewHandler(loggerp *slog.Logger, c nutritionplan_c.NutritionPlanController) *Handler {
	return &Handler{
		Logger:     loggerp,
		Controller: c,
	}
}
