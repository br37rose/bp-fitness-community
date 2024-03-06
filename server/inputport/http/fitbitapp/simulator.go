package fitbitapp

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	fitbitapp_c "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/fitbitapp/controller"
	fitbitapp_s "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/fitbitapp/datastore"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/utils/httperror"
)

func UnmarshalCreateRequest(ctx context.Context, r *http.Request) (*fitbitapp_c.FitBitSimulatorCreateRequestIDO, error) {
	// Initialize our array which will store all the results from the remote server.
	var requestData fitbitapp_c.FitBitSimulatorCreateRequestIDO

	defer r.Body.Close()

	// Read the JSON string and convert it into our golang stuct else we need
	// to send a `400 Bad Request` errror message back to the client,
	err := json.NewDecoder(r.Body).Decode(&requestData) // [1]
	if err != nil {
		log.Println(err)
		return nil, httperror.NewForSingleField(http.StatusBadRequest, "non_field_error", "payload structure is wrong")
	}

	// Perform our validation and return validation error on any issues detected.
	if err := ValidateCreateRequest(&requestData); err != nil {
		return nil, err
	}
	return &requestData, nil
}

func ValidateCreateRequest(dirtyData *fitbitapp_c.FitBitSimulatorCreateRequestIDO) error {
	e := make(map[string]string)

	if dirtyData.UserID.IsZero() {
		e["user_id"] = "missing value"
	}
	if dirtyData.SimulatorAlgorithm == "" {
		e["simulator_algorithm"] = "missing value"
	}

	if len(e) != 0 {
		return httperror.NewForBadRequest(&e)
	}
	return nil
}

func (h *Handler) CreateSimulator(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	data, err := UnmarshalCreateRequest(ctx, r)
	if err != nil {
		httperror.ResponseError(w, err)
		return
	}

	res, err := h.Controller.CreateSimulator(ctx, data)
	if err != nil {
		httperror.ResponseError(w, err)
		return
	}

	MarshalCreateSimulatorResponse(res, w)
}

func MarshalCreateSimulatorResponse(res *fitbitapp_s.FitBitApp, w http.ResponseWriter) {
	if err := json.NewEncoder(w).Encode(&res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
