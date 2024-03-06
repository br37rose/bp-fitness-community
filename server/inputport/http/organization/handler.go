package organization

import (
	organization_c "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/organization/controller"
	"log/slog"
)

// Handler Creates http request handler
type Handler struct {
	Logger     *slog.Logger
	Controller organization_c.OrganizationController
}

// NewHandler Constructor
func NewHandler(loggerp *slog.Logger, c organization_c.OrganizationController) *Handler {
	return &Handler{
		Logger:     loggerp,
		Controller: c,
	}
}
