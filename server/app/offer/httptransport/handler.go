package httptransport

import (
	offer_c "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/offer/controller"
	"log/slog"
)

// Handler Creates http request handler
type Handler struct {
	Logger     *slog.Logger
	Controller offer_c.Offerontroller
}

// NewHandler Constructor
func NewHandler(loggerp *slog.Logger, c offer_c.Offerontroller) *Handler {
	return &Handler{
		Logger:     loggerp,
		Controller: c,
	}
}
