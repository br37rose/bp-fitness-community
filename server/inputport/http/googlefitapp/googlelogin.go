package googlefitapp

import (
	"encoding/json"
	"net/http"

	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/utils/httperror"
)

func (h *Handler) GetGoogleLoginURL(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	res, err := h.Controller.GetGoogleLoginURL(ctx)
	if err != nil {
		httperror.ResponseError(w, err)
		return
	}

	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(false)
	if err := enc.Encode(&res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
