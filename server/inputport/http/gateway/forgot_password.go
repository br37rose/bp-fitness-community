package gateway

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/utils/httperror"
)

type ForgotPasswordRequestIDO struct {
	Email string `json:"email"`
}

func UnmarshalForgotPasswordRequest(ctx context.Context, r *http.Request) (*ForgotPasswordRequestIDO, error) {
	// Initialize our array which will store all the results from the remote server.
	var requestData ForgotPasswordRequestIDO

	defer r.Body.Close()

	// Read the JSON string and convert it into our golang stuct else we need
	// to send a `400 Bad Request` errror message back to the client,
	err := json.NewDecoder(r.Body).Decode(&requestData) // [1]
	if err != nil {
		return nil, httperror.NewForSingleField(http.StatusBadRequest, "non_field_error", "payload structure is wrong")
	}

	// Defensive Code: For security purposes we need to remove all whitespaces from the email and lower the characters.
	requestData.Email = strings.ToLower(requestData.Email)
	requestData.Email = strings.ReplaceAll(requestData.Email, " ", "")

	// Perform our validation and return validation error on any issues detected.
	if err = ValidateForgotPasswordRequest(&requestData); err != nil {
		return nil, err
	}
	return &requestData, nil
}

func ValidateForgotPasswordRequest(dirtyData *ForgotPasswordRequestIDO) error {
	e := make(map[string]string)

	if dirtyData.Email == "" {
		e["email"] = "missing value"
	}
	if len(dirtyData.Email) > 255 {
		e["email"] = "too long"
	}

	if len(e) != 0 {
		return httperror.NewForBadRequest(&e)
	}
	return nil
}

func (h *Handler) ForgotPassword(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	data, err := UnmarshalForgotPasswordRequest(ctx, r)
	if err != nil {
		httperror.ResponseError(w, err)
		return
	}

	if err := h.Controller.ForgotPassword(ctx, data.Email); err != nil {
		httperror.ResponseError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
