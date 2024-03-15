package stripe

import (
	"net/http"

	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/utils/httperror"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CancelStripeSubscription Function will take the authenticated user and
// cancel their current subscription.
func (h *Handler) CancelStripeSubscription(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	// Here is where you extract url parameters.
	query := r.URL.Query()

	// Apply the subscription cancellation based on if the user submitted a
	// member ID or not. If no member ID was specified then business logic will
	// assume it is the logged in user making cancelation on their own account.
	memberID := primitive.NilObjectID
	memberIDStr := query.Get("member_id")
	if memberIDStr != "" {
		mid, err := primitive.ObjectIDFromHex(memberIDStr)
		if err != nil {
			httperror.ResponseError(w, err)
			return
		}
		memberID = mid
	}

	ctx := r.Context()

	err := h.Controller.CancelStripeSubscription(ctx, memberID)
	if err != nil {
		httperror.ResponseError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
