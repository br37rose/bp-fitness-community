package datapoint

import (
	"log/slog"

	dp_c "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/datapoint/controller"
)

// Handler Creates http request handler
type Handler struct {
	Logger     *slog.Logger
	Controller dp_c.DataPointController
}

// NewHandler Constructor
func NewHandler(loggerp *slog.Logger, c dp_c.DataPointController) *Handler {
	return &Handler{
		Logger:     loggerp,
		Controller: c,
	}
}
