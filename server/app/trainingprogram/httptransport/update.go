package httptransport

import (
	"context"
	"encoding/json"
	"net/http"

	tp_c "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/trainingprogram/controller"
	tp_s "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/trainingprogram/datastore"

	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/utils/httperror"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (h *Handler) UpdateTrainingPhase(
	w http.ResponseWriter,
	r *http.Request,
	tpId string) {

	ctx := r.Context()
	tpObjectID, err := primitive.ObjectIDFromHex(tpId)
	if err != nil {
		httperror.ResponseError(w, err)
		return
	}
	data, err := UnmarshalUpdateRequest(ctx, r)
	if err != nil {
		httperror.ResponseError(w, err)
		return
	}

	m, err := h.Controller.UpdateTPPhase(ctx, tpObjectID, *data)
	if err != nil {
		httperror.ResponseError(w, err)
		return
	}

	MarshalDetailResponse(m, w)
}
func UnmarshalUpdateRequest(ctx context.Context, r *http.Request) (*tp_c.PhaseUpdateRequestIDO, error) {
	// Initialize our array which will store all the results from the remote server.
	var requestData tp_c.PhaseUpdateRequestIDO

	defer r.Body.Close()

	// Read the JSON string and convert it into our golang stuct else we need
	// to send a `400 Bad Request` errror message back to the client,
	err := json.NewDecoder(r.Body).Decode(&requestData) // [1]
	if err != nil {
		return nil, httperror.NewForSingleField(http.StatusBadRequest, "non_field_error", "payload structure is wrong")
	}
	return &requestData, nil
}

func MarshalDetailResponse(res *tp_s.TrainingProgram, w http.ResponseWriter) {
	if err := json.NewEncoder(w).Encode(&res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
