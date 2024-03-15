package stripe

import (
	"encoding/json"
	"net/http"
	"strconv"

	u_d "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/user/datastore"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/utils/httperror"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

//	func (c *Handler) StripeWebhook(w http.ResponseWriter, r *http.Request) {
//		if r.Method != "POST" {
//			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
//			return
//		}
//
//		// https://stackoverflow.com/questions/28073395/limiting-file-size-in-formfile
//		const MaxBodyBytes = int64(65536)
//		r.Body = http.MaxBytesReader(w, r.Body, MaxBodyBytes)
//
//		payload, err := ioutil.ReadAll(r.Body)
//		if err != nil {
//			c.Logger.Error("stripe webhook", slog.String("payload", string(payload)), slog.Any("err", err))
//			http.Error(w, err.Error(), http.StatusBadRequest)
//			httperror.ResponseError(w, err)
//			return
//		}
//
//		ctx := r.Context()
//
//		// Defensive code: To prevent hackers from taking advantage of this public
//		// API endpoint then we need verify the code from `Stripe, Inc`.
//		header := r.Header.Get("Stripe-Signature")
//		if err := c.Controller.StripeWebhook(ctx, header, payload); err != nil {
//			httperror.ResponseError(w, err)
//			return
//		}
//
// }
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
