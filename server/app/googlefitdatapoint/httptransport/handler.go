package httptransport

import (
	"log/slog"

	dp_c "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/googlefitdatapoint/controller"
)

// Handler Creates http request handler
type Handler struct {
	Logger     *slog.Logger
	Controller dp_c.GoogleFitDataPointController
}

// NewHandler Constructor
func NewHandler(loggerp *slog.Logger, c dp_c.GoogleFitDataPointController) *Handler {
	return &Handler{
		Logger:     loggerp,
		Controller: c,
	}
}
