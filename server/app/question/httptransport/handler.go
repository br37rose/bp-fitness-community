package httptransport

import (
	"log/slog"

	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/question/controller"
)

// Handler Creates http request handler
type Handler struct {
	Logger     *slog.Logger
	Controller controller.QuestionController
}

// NewHandler Constructor
func NewHandler(loggerp *slog.Logger, c controller.QuestionController) *Handler {
	return &Handler{
		Logger:     loggerp,
		Controller: c,
	}
}
