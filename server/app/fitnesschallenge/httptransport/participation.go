package httptransport

import (
	"net/http"

	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/utils/httperror"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (h *Handler) ChangeParticipationStatus(
	w http.ResponseWriter,
	r *http.Request,
	tpId string) {

	ctx := r.Context()
	id, err := primitive.ObjectIDFromHex(tpId)
	if err != nil {
		httperror.ResponseError(w, err)
		return
	}
	if id.IsZero() {
		httperror.ResponseError(w, httperror.NewForSingleField(http.StatusBadRequest, "id", "incorrect value "))
		return
	}

	m, err := h.Controller.ChangeParticipationStatus(ctx, id)
	if err != nil {
		httperror.ResponseError(w, err)
		return
	}

	MarshalDetailResponse(m, w)
}
