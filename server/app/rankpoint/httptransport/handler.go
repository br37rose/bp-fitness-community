package httptransport

import (
	"log/slog"

	rp_c "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/rankpoint/controller"
)

// Handler Creates http request handler
type Handler struct {
	Logger     *slog.Logger
	Controller rp_c.RankPointController
}

// NewHandler Constructor
func NewHandler(loggerp *slog.Logger, c rp_c.RankPointController) *Handler {
	return &Handler{
		Logger:     loggerp,
		Controller: c,
	}
}
