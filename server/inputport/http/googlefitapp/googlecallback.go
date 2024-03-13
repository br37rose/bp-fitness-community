package googlefitapp

import (
	"log"
	"net/http"

	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/utils/httperror"
)

func (h *Handler) GoogleCallback(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	code := r.FormValue("code")
	state := r.FormValue("state")

	log.Println("state", state)
	log.Println("code", code)

	res, err := h.Controller.GoogleCallback(ctx, state, code)
	if err != nil {
		httperror.ResponseError(w, err)
		return
	}

	http.Redirect(w, r, res.URL, http.StatusFound)
}
