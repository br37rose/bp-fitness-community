package videocontent

import (
	videocontent_c "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/videocontent/controller"
	"log/slog"
)

// Handler Creates http request handler
type Handler struct {
	Logger     *slog.Logger
	Controller videocontent_c.VideoContentController
}

// NewHandler Constructor
func NewHandler(loggerp *slog.Logger, c videocontent_c.VideoContentController) *Handler {
	return &Handler{
		Logger:     loggerp,
		Controller: c,
	}
}
