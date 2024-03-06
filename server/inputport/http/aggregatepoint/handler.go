package aggregatepoint

import (
	"log/slog"

	ap_c "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/aggregatepoint/controller"
)

// Handler Creates http request handler
type Handler struct {
	Logger     *slog.Logger
	Controller ap_c.AggregatePointController
}

// NewHandler Constructor
func NewHandler(loggerp *slog.Logger, c ap_c.AggregatePointController) *Handler {
	return &Handler{
		Logger:     loggerp,
		Controller: c,
	}
}
