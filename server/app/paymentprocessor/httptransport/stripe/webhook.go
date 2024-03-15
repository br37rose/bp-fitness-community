package stripe

import (
	"io/ioutil"
	"net/http"

	"log/slog"

	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/utils/httperror"
)

func (c *Handler) Webhook(w http.ResponseWriter, r *http.Request) {
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
	if err := c.Controller.Webhook(ctx, header, payload); err != nil {
		httperror.ResponseError(w, err)
		return
	}
}
