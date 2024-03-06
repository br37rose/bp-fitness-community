package biometric

import (
	"log/slog"

	bio_c "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/biometric/controller"
)

// Handler Creates http request handler
type Handler struct {
	Logger     *slog.Logger
	Controller bio_c.BiometricController
}

// NewHandler Constructor
func NewHandler(loggerp *slog.Logger, c bio_c.BiometricController) *Handler {
	return &Handler{
		Logger:     loggerp,
		Controller: c,
	}
}
