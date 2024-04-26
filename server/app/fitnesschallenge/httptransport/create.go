package httptransport

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/fitnesschallenge/controller"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/fitnesschallenge/datastore"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/utils/httperror"
)

func UnmarshalCreateRequest(ctx context.Context, r *http.Request) (*controller.FitnessChallengeCreateRequestIDO, error) {
	var requestData controller.FitnessChallengeCreateRequestIDO

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

func MarshalCreateResponse(res *datastore.FitnessChallenge, w http.ResponseWriter) {
	if err := json.NewEncoder(w).Encode(&res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func ValidateCreateRequest(dirtyData *controller.FitnessChallengeCreateRequestIDO) error {
	e := make(map[string]string)

	if len(e) != 0 {
		return httperror.NewForBadRequest(&e)
	}
	if dirtyData.OrganizationID.IsZero() {
		e["organization_id"] = "missing value"
	}
	if dirtyData.Name == "" {
		e["name"] = "missing value"
	}
	if dirtyData.Description == "" {
		e["description"] = "missing value"
	}
	if dirtyData.StartOn.IsZero() {
		e["startOn"] = "missing value"
	}
	if len(dirtyData.Rules) == 0 {
		e["rules"] = "missing value"
	}
	return nil
}
