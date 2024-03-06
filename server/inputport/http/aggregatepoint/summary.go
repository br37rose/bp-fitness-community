package aggregatepoint

import (
	"encoding/json"
	"net/http"

	ap_c "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/aggregatepoint/controller"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/utils/httperror"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (h *Handler) GetSummary(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Here is where you extract url parameters.
	query := r.URL.Query()

	userIDStr := query.Get("user_id")
	userID, err := primitive.ObjectIDFromHex(userIDStr)
	if err != nil {
		httperror.ResponseError(w, err)
		return
	}

	// Perform our database operation.
	res, err := h.Controller.GetSummary(ctx, userID)
	if err != nil {
		httperror.ResponseError(w, err)
		return
	}
	MarshalAggregatePointSummaryResponse(res, w)
}

func MarshalAggregatePointSummaryResponse(res *ap_c.AggregatePointSummaryResponse, w http.ResponseWriter) {
	if err := json.NewEncoder(w).Encode(&res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
