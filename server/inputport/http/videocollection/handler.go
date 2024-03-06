package videocollection

import (
	videocollection_c "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/videocollection/controller"
	"log/slog"
)

// Handler Creates http request handler
type Handler struct {
	Logger     *slog.Logger
	Controller videocollection_c.VideoCollectionController
}

// NewHandler Constructor
func NewHandler(loggerp *slog.Logger, c videocollection_c.VideoCollectionController) *Handler {
	return &Handler{
		Logger:     loggerp,
		Controller: c,
	}
}
