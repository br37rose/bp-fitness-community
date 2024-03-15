package stripe

import (
	"log/slog"

	paymentprocessor_c "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/paymentprocessor/controller/stripe"
)

// Handler Creates http request handler
type Handler struct {
	Logger     *slog.Logger
	Controller paymentprocessor_c.StripePaymentProcessorController
}

// NewHandler Constructor
func NewHandler(loggerp *slog.Logger, c paymentprocessor_c.StripePaymentProcessorController) *Handler {
	return &Handler{
		Logger:     loggerp,
		Controller: c,
	}
}
