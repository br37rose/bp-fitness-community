package httptransport

import (
	"log/slog"

	tag_c "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/tag/controller"
)

// Handler Creates http request handler
type Handler struct {
	Logger     *slog.Logger
	Controller tag_c.TagController
}

// NewHandler Constructor
func NewHandler(loggerp *slog.Logger, c tag_c.TagController) *Handler {
	return &Handler{
		Logger:     loggerp,
		Controller: c,
	}
}
