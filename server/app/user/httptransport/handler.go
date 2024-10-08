package httptransport

import (
	"log/slog"

	user_c "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/user/controller"
)

// Handler Creates http request handler
type Handler struct {
	Logger     *slog.Logger
	Controller user_c.UserController
}

// NewHandler Constructor
func NewHandler(loggerp *slog.Logger, c user_c.UserController) *Handler {
	return &Handler{
		Logger:     loggerp,
		Controller: c,
	}
}
