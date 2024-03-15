package httptransport

import (
	gateway_c "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/gateway/controller"
	"log/slog"
)

// Handler Creates http request handler
type Handler struct {
	Logger     *slog.Logger
	Controller gateway_c.GatewayController
}

// NewHandler Constructor
func NewHandler(loggerp *slog.Logger, c gateway_c.GatewayController) *Handler {
	return &Handler{
		Logger:     loggerp,
		Controller: c,
	}
}
