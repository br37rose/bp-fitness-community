package fitbitapp

import (
	"net/http"

	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/utils/httperror"
)

func (h *Handler) Subscriber(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	res, err := h.Controller.GetRegistrationURL(ctx)
	if err != nil {
		httperror.ResponseError(w, err)
		return
	}

	MarshalDetailResponse(res, w)
}
