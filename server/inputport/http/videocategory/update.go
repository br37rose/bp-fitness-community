package videocategory

import (
	"context"
	"encoding/json"
	"net/http"

	sub_c "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/videocategory/controller"
	sub_s "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/videocategory/datastore"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/utils/httperror"
)

func UnmarshalUpdateRequest(ctx context.Context, r *http.Request) (*sub_c.VideoCategoryUpdateRequestIDO, error) {
	// Initialize our array which will store all the results from the remote server.
	var requestData sub_c.VideoCategoryUpdateRequestIDO

	defer r.Body.Close()

	// Read the JSON string and convert it into our golang stuct else we need
	// to send a `400 Bad Request` errror message back to the client,
	err := json.NewDecoder(r.Body).Decode(&requestData) // [1]
	if err != nil {
		return nil, httperror.NewForSingleField(http.StatusBadRequest, "non_field_error", "payload structure is wrong")
	}

	return &requestData, nil
}

func (h *Handler) UpdateByID(w http.ResponseWriter, r *http.Request, id string) {
	ctx := r.Context()

	data, err := UnmarshalUpdateRequest(ctx, r)
	if err != nil {
		httperror.ResponseError(w, err)
		return
	}

	org, err := h.Controller.UpdateByID(ctx, data)
	if err != nil {
		httperror.ResponseError(w, err)
		return
	}

	MarshalUpdateResponse(org, w)
}

func MarshalUpdateResponse(res *sub_s.VideoCategory, w http.ResponseWriter) {
	if err := json.NewEncoder(w).Encode(&res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
