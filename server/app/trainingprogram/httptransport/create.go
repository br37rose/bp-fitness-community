package httptransport

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/trainingprogram/controller"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/trainingprogram/datastore"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/utils/httperror"
)

func UnmarshalCreateRequest(ctx context.Context, r *http.Request) (*controller.TrainingProgramCreateRequestIDO, error) {
	var requestData controller.TrainingProgramCreateRequestIDO

	defer r.Body.Close()
	err := json.NewDecoder(r.Body).Decode(&requestData) // [1]
	if err != nil {
		log.Println(err)
		return nil, httperror.NewForSingleField(http.StatusBadRequest, "non_field_error", "payload structure is wrong")
	}
	if err := ValidateCreateRequest(&requestData); err != nil {
		return nil, err
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

func MarshalCreateResponse(res *datastore.TrainingProgram, w http.ResponseWriter) {
	if err := json.NewEncoder(w).Encode(&res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func ValidateCreateRequest(dirtyData *controller.TrainingProgramCreateRequestIDO) error {
	e := make(map[string]string)

	if dirtyData.Name == "" {
		e["name"] = "missing value"
	}
	if dirtyData.Phases == 0 {
		e["phases"] = "missing value"
	}
	if dirtyData.Weeks == 0 {
		e["weeks"] = "missing value"
	}

	if len(e) != 0 {
		return httperror.NewForBadRequest(&e)
	}
	return nil
}
