package httptransport

import (
	"log/slog"

	googlefitapp_c "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/googlefitapp/controller"
)

// Handler Creates http request handler
type Handler struct {
	Logger     *slog.Logger
	Controller googlefitapp_c.GoogleFitAppController
}

// NewHandler Constructor
func NewHandler(loggerp *slog.Logger, c googlefitapp_c.GoogleFitAppController) *Handler {
	return &Handler{
		Logger:     loggerp,
		Controller: c,
	}
}
