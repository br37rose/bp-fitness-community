package gateway

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/utils/httperror"
)

type PasswordResetRequestIDO struct {
	VerificationCode string `json:"verification_code"`
	Password         string `json:"password"`
	PasswordRepeated string `json:"password_repeated"`
}

func UnmarshalPasswordResetRequest(ctx context.Context, r *http.Request) (*PasswordResetRequestIDO, error) {
	// Initialize our array which will store all the results from the remote server.
	var requestData PasswordResetRequestIDO

	defer r.Body.Close()

	// Read the JSON string and convert it into our golang stuct else we need
	// to send a `400 Bad Request` errror message back to the client,
	err := json.NewDecoder(r.Body).Decode(&requestData) // [1]
	if err != nil {
		return nil, httperror.NewForSingleField(http.StatusBadRequest, "non_field_error", "payload structure is wrong")
	}

	// Perform our validation and return validation error on any issues detected.
	if err = ValidatePasswordResetRequest(&requestData); err != nil {
		return nil, err
	}
	return &requestData, nil
}

func ValidatePasswordResetRequest(dirtyData *PasswordResetRequestIDO) error {
	e := make(map[string]string)

	if dirtyData.VerificationCode == "" {
		e["verification_code"] = "missing value"
	}
	if dirtyData.Password == "" {
		e["password"] = "missing value"
	}
	if len(dirtyData.Password) > 255 {
		e["password"] = "too long"
	}
	if dirtyData.PasswordRepeated == "" {
		e["password_repeated"] = "missing value"
	}
	if len(dirtyData.PasswordRepeated) > 255 {
		e["password_repeated"] = "too long"
	}
	if dirtyData.Password != dirtyData.PasswordRepeated {
		e["password"] = "value does not match"
		e["password_repeated"] = "value does not match"
	}

	if len(e) != 0 {
		return httperror.NewForBadRequest(&e)
	}
	return nil
}

func (h *Handler) PasswordReset(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	data, err := UnmarshalPasswordResetRequest(ctx, r)
	if err != nil {
		httperror.ResponseError(w, err)
		return
	}

	if err := h.Controller.PasswordReset(ctx, data.VerificationCode, data.Password); err != nil {
		httperror.ResponseError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
