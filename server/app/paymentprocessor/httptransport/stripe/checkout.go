package stripe

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/utils/httperror"
)

type CreateStripeCheckoutSessionRequestIDO struct {
	PriceID string `json:"price_id"`
}

func UnmarshalCreateStripeCheckoutSessionRequest(ctx context.Context, r *http.Request) (*CreateStripeCheckoutSessionRequestIDO, error) {
	// Initialize our array which will store all the results from the remote server.
	var requestData *CreateStripeCheckoutSessionRequestIDO

	defer r.Body.Close()

	// Read the JSON string and convert it into our golang stuct else we need
	// to send a `400 Bad Request` errror message back to the client,
	err := json.NewDecoder(r.Body).Decode(&requestData) // [1]
	if err != nil {
		return nil, httperror.NewForSingleField(http.StatusBadRequest, "non_field_error", "payload structure is wrong")
	}

	// Perform our validation and return validation error on any issues detected.
	if err := ValidateCreateStripeCheckoutSessionRequest(requestData); err != nil {
		return nil, err
	}

	return requestData, nil
}

func ValidateCreateStripeCheckoutSessionRequest(dirtyData *CreateStripeCheckoutSessionRequestIDO) error {
	e := make(map[string]string)

	if dirtyData.PriceID == "" {
		e["price_id"] = "missing value"
	}

	if len(e) != 0 {
		return httperror.NewForBadRequest(&e)
	}
	return nil
}

// CreateStripeCheckoutSession Function handles accepting the
// subscription choice from the authenticated user and create a subscription
// checkout session URL. Return this URL back to the user so the user can
// redirect to the checkout.
func (h *Handler) CreateStripeCheckoutSession(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	data, err := UnmarshalCreateStripeCheckoutSessionRequest(ctx, r)
	if err != nil {
		httperror.ResponseError(w, err)
		return
	}

	checkoutURL, err := h.Controller.CreateStripeCheckoutSessionURL(ctx, data.PriceID)
	if err != nil {
		httperror.ResponseError(w, err)
		return
	}

	m := make(map[string]string)
	m["checkout_url"] = checkoutURL

	if err := json.NewEncoder(w).Encode(&m); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// CompleteStripeCheckoutSession Function will take the session id
// that was returned by Stripe and process the subscription.
func (h *Handler) CompleteStripeCheckoutSession(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
	sessionID := r.URL.Query().Get("session_id")

	h.Logger.Debug("completing session", slog.String("sessionID", sessionID))

	ctx := r.Context()

	res, err := h.Controller.CompleteStripeCheckoutSession(ctx, sessionID)
	if err != nil {
		httperror.ResponseError(w, err)
		return
	}

	if err := json.NewEncoder(w).Encode(&res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
