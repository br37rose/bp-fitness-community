package videocategory

import (
	videocategory_c "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/videocategory/controller"
	"log/slog"
)

// Handler Creates http request handler
type Handler struct {
	Logger     *slog.Logger
	Controller videocategory_c.VideoCategoryController
}

// NewHandler Constructor
func NewHandler(loggerp *slog.Logger, c videocategory_c.VideoCategoryController) *Handler {
	return &Handler{
		Logger:     loggerp,
		Controller: c,
	}
}
