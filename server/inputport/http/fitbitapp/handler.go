package fitbitapp

import (
	"log/slog"

	fitbitapp_c "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/fitbitapp/controller"
)

// Handler Creates http request handler
type Handler struct {
	Logger     *slog.Logger
	Controller fitbitapp_c.FitBitAppController
}

// NewHandler Constructor
func NewHandler(loggerp *slog.Logger, c fitbitapp_c.FitBitAppController) *Handler {
	return &Handler{
		Logger:     loggerp,
		Controller: c,
	}
}
