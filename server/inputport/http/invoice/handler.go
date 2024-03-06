package invoice

import (
	invoice_c "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/invoice/controller"
	"log/slog"
)

// Handler Creates http request handler
type Handler struct {
	Logger     *slog.Logger
	Controller invoice_c.InvoiceController
}

// NewHandler Constructor
func NewHandler(loggerp *slog.Logger, c invoice_c.InvoiceController) *Handler {
	return &Handler{
		Logger:     loggerp,
		Controller: c,
	}
}
