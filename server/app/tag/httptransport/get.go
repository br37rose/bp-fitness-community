package httptransport

import (
	"encoding/json"
	"net/http"

	"go.mongodb.org/mongo-driver/bson/primitive"

	sub_s "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/tag/datastore"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/utils/httperror"
)

func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request, id string) {
	ctx := r.Context()

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		httperror.ResponseError(w, err)
		return
	}

	m, err := h.Controller.GetByID(ctx, objectID)
	if err != nil {
		httperror.ResponseError(w, err)
		return
	}

	MarshalDetailResponse(m, w)
}

func MarshalDetailResponse(res *sub_s.Tag, w http.ResponseWriter) {
	if err := json.NewEncoder(w).Encode(&res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
