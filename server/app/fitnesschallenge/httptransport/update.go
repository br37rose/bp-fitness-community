package httptransport

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/fitnesschallenge/controller"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/fitnesschallenge/datastore"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/utils/httperror"
)

func (h *Handler) UpdateByID(
	w http.ResponseWriter,
	r *http.Request,
	tpId string) {

	ctx := r.Context()

	data, err := UnmarshalUpdateRequest(ctx, r)
	if err != nil {
		httperror.ResponseError(w, err)
		return
	}

	m, err := h.Controller.UpdateByID(ctx, data)
	if err != nil {
		httperror.ResponseError(w, err)
		return
	}

	MarshalDetailResponse(m, w)
}
func UnmarshalUpdateRequest(ctx context.Context, r *http.Request) (*controller.FitnessChallengeUpdateRequestIDO, error) {
	// Initialize our array which will store all the results from the remote server.
	var requestData controller.FitnessChallengeUpdateRequestIDO

	defer r.Body.Close()

	// Read the JSON string and convert it into our golang stuct else we need
	// to send a `400 Bad Request` errror message back to the client,
	err := json.NewDecoder(r.Body).Decode(&requestData) // [1]
	if err != nil {
		return nil, httperror.NewForSingleField(http.StatusBadRequest, "non_field_error", "payload structure is wrong")
	}
	return &requestData, nil
}

func MarshalDetailResponse(res *datastore.FitnessChallenge, w http.ResponseWriter) {
	if err := json.NewEncoder(w).Encode(&res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
