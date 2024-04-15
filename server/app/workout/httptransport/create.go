package httptransport

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	w_c "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/workout/controller"
	w_d "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/workout/datastore"

	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/utils/httperror"
)

func UnmarshalCreateRequest(ctx context.Context, r *http.Request) (*w_c.WorkoutCreateRequestIDO, error) {
	// Initialize our array which will store all the results from the remote server.
	var requestData *w_c.WorkoutCreateRequestIDO

	defer r.Body.Close()

	// Read the JSON string and convert it into our golang stuct else we need
	// to send a `400 Bad Request` errror message back to the client,
	err := json.NewDecoder(r.Body).Decode(&requestData) // [1]
	if err != nil {
		return nil, httperror.NewForSingleField(http.StatusBadRequest, "non_field_error", "payload structure is wrong")
	}

	// Perform our validation and return validation error on any issues detected.
	if err := ValidateCreateRequest(requestData); err != nil {
		return nil, err
	}

	return requestData, nil
}

func ValidateCreateRequest(dirtyData *w_c.WorkoutCreateRequestIDO) error {
	e := make(map[string]string)

	if dirtyData.Name == "" {
		e["name"] = "missing value"
	}

	if dirtyData.WorkoutExercises == nil {
		e["workout_exercises"] = "missing value"
	}

	if len(e) != 0 {
		return httperror.NewForBadRequest(&e)
	}
	return nil
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	data, err := UnmarshalCreateRequest(ctx, r)
	if err != nil {
		log.Println("Create | member | err:", err)
		httperror.ResponseError(w, err)
		return
	}

	workout, err := h.Controller.Create(ctx, data)
	if err != nil {
		httperror.ResponseError(w, err)
		return
	}

	MarshalCreateResponse(workout, w)
}

func MarshalCreateResponse(res *w_d.Workout, w http.ResponseWriter) {
	if err := json.NewEncoder(w).Encode(&res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
