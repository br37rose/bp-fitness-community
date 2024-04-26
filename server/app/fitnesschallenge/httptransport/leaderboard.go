package httptransport

import (
	"encoding/json"
	"net/http"

	rank "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/rankpoint/datastore"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/utils/httperror"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (h *Handler) GetChallengeLeaderBoard(w http.ResponseWriter, r *http.Request, id string) {
	ctx := r.Context()

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		httperror.ResponseError(w, err)
		return
	}

	m, err := h.Controller.GetChallengeLeaderBoard(ctx, objectID)
	if err != nil {
		httperror.ResponseError(w, err)
		return
	}

	MarshalRankListResponse(m, w)
}

func MarshalRankListResponse(res *rank.RankPointPaginationListResult, w http.ResponseWriter) {
	if err := json.NewEncoder(w).Encode(&res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
