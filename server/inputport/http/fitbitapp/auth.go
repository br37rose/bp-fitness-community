package fitbitapp

import (
	"encoding/json"
	"net/http"

	fba_c "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/fitbitapp/controller"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/utils/httperror"
)

func (h *Handler) Auth(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Extract the query parameters.
	code := r.URL.Query().Get("code")

	res, err := h.Controller.Auth(ctx, code)
	if err != nil {
		httperror.ResponseError(w, err)
		return
	}

	http.Redirect(w, r, res.URL, http.StatusFound)
}

func MarshalDetailResponse(res *fba_c.RegistrationURLResponse, w http.ResponseWriter) {
	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(false)
	if err := enc.Encode(&res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
