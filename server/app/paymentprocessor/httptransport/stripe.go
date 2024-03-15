package httptransport

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	"log/slog"

	"go.mongodb.org/mongo-driver/bson/primitive"

	u_d "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/user/datastore"
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

// CancelStripeSubscription Function will take the authenticated user and
// cancel their current subscription.
func (h *Handler) CancelStripeSubscription(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	ctx := r.Context()

	err := h.Controller.CancelStripeSubscription(ctx)
	if err != nil {
		httperror.ResponseError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (c *Handler) StripeWebhook(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	// https://stackoverflow.com/questions/28073395/limiting-file-size-in-formfile
	const MaxBodyBytes = int64(65536)
	r.Body = http.MaxBytesReader(w, r.Body, MaxBodyBytes)

	payload, err := ioutil.ReadAll(r.Body)
	if err != nil {
		c.Logger.Error("stripe webhook", slog.String("payload", string(payload)), slog.Any("err", err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		httperror.ResponseError(w, err)
		return
	}

	ctx := r.Context()

	// Defensive code: To prevent hackers from taking advantage of this public
	// API endpoint then we need verify the code from `Stripe, Inc`.
	header := r.Header.Get("Stripe-Signature")
	if err := c.Controller.StripeWebhook(ctx, header, payload); err != nil {
		httperror.ResponseError(w, err)
		return
	}

}

func (h *Handler) ListLatestStripeInvoices(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Here is where you extract url parameters.
	query := r.URL.Query()

	var userID primitive.ObjectID
	userIDStr := query.Get("user_id")
	if userIDStr != "" {
		userIDRes, err := primitive.ObjectIDFromHex(userIDStr)
		if err != nil {
			httperror.ResponseError(w, err)
			return
		}
		userID = userIDRes
	}

	var cursor int64 = 0
	cursorStr := query.Get("cursor")
	if cursorStr != "" {
		cursorInt64, _ := strconv.ParseInt(cursorStr, 10, 64)
		cursor = cursorInt64
	}

	var limit int64 = 25
	limitStr := query.Get("page_size")
	if limitStr != "" {
		pageSize, _ := strconv.ParseInt(limitStr, 10, 64)
		if pageSize == 0 || pageSize > 250 {
			limit = 25
		}
		limit = pageSize
	}

	m, err := h.Controller.ListLatestStripeInvoices(ctx, userID, cursor, limit)
	if err != nil {
		httperror.ResponseError(w, err)
		return
	}

	MarshalPublicListResponse(m, w)
}

func MarshalPublicListResponse(res *u_d.StripeInvoiceListResult, w http.ResponseWriter) {
	if err := json.NewEncoder(w).Encode(&res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
