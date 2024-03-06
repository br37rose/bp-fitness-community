package fitnessplan

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	e_c "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/fitnessplan/controller"
	e_d "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/fitnessplan/datastore"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/utils/httperror"
)

func UnmarshalCreateRequest(ctx context.Context, r *http.Request) (*e_c.FitnessPlanCreateRequestIDO, error) {
	// Initialize our array which will store all the results from the remote server.
	var requestData e_c.FitnessPlanCreateRequestIDO

	defer r.Body.Close()

	// Read the JSON string and convert it into our golang stuct else we need
	// to send a `400 Bad Request` errror message back to the client,
	err := json.NewDecoder(r.Body).Decode(&requestData) // [1]
	if err != nil {
		log.Println(err)
		return nil, httperror.NewForSingleField(http.StatusBadRequest, "non_field_error", "payload structure is wrong")
	}
	return &requestData, nil
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	req, err := UnmarshalCreateRequest(ctx, r)
	if err != nil {
		httperror.ResponseError(w, err)
		return
	}

	data, err := h.Controller.Create(ctx, req)
	if err != nil {
		httperror.ResponseError(w, err)
		return
	}

	MarshalCreateResponse(data, w)
}

func MarshalCreateResponse(res *e_d.FitnessPlan, w http.ResponseWriter) {
	if err := json.NewEncoder(w).Encode(&res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
