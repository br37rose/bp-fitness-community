package member

import (
	member_c "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/member/controller"
	"log/slog"
)

// Handler Creates http request handler
type Handler struct {
	Logger     *slog.Logger
	Controller member_c.MemberController
}

// NewHandler Constructor
func NewHandler(loggerp *slog.Logger, c member_c.MemberController) *Handler {
	return &Handler{
		Logger:     loggerp,
		Controller: c,
	}
}
