package gateway

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	user_s "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/user/datastore"
)

type RefreshTokenRequestIDO struct {
	Value string `json:"value"`
}

// RefreshTokenResponseIDO struct used to represent the system's response when the `login` POST request was a success.
type RefreshTokenResponseIDO struct {
	Email                  string    `json:"username"`
	AccessToken            string    `json:"access_token"`
	AccessTokenExpiryDate  time.Time `json:"access_token_expiry_date"`
	RefreshToken           string    `json:"refresh_token"`
	RefreshTokenExpiryDate time.Time `json:"refresh_token_expiry_date"`
}

func UnmarshalRefreshTokenRequest(ctx context.Context, r *http.Request) (*RefreshTokenRequestIDO, error, int) {
	// Initialize our array which will store all the results from the remote server.
	var requestData RefreshTokenRequestIDO

	defer r.Body.Close()

	// Read the JSON string and convert it into our golang stuct else we need
	// to send a `400 Bad Request` errror message back to the client,
	err := json.NewDecoder(r.Body).Decode(&requestData) // [1]
	if err != nil {
		return nil, err, http.StatusBadRequest
	}

	// Perform our validation and return validation error on any issues detected.
	isValid, errStr := ValidateRefreshTokenRequest(&requestData)
	if isValid == false {
		return nil, errors.New(errStr), http.StatusBadRequest
	}

	return &requestData, nil, http.StatusOK
}

func ValidateRefreshTokenRequest(dirtyData *RefreshTokenRequestIDO) (bool, string) {
	e := make(map[string]string)

	if dirtyData.Value == "" {
		e["value"] = "missing value"
	}

	if len(e) != 0 {
		b, err := json.Marshal(e)
		if err != nil { // Defensive code
			return false, err.Error()
		}
		return false, string(b)
	}
	return true, ""
}

func (h *Handler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	requestData, err, errStatusCode := UnmarshalRefreshTokenRequest(ctx, r)
	if err != nil {
		http.Error(w, err.Error(), errStatusCode)
		return
	}

	user, accessToken, accessTokenExpiryDate, refreshToken, refreshTokenExpiryDate, err := h.Controller.RefreshToken(ctx, requestData.Value)
	if user == nil {
		http.Error(w, "{'non_field_error':'user does not exist'}", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusOK)
		return
	}

	MarshalRefreshTokenResponse(accessToken, accessTokenExpiryDate, refreshToken, refreshTokenExpiryDate, user, w)
}

func MarshalRefreshTokenResponse(accessToken string, accessTokenExpiryDate time.Time, refreshToken string, refreshTokenExpiryDate time.Time, u *user_s.User, w http.ResponseWriter) {
	responseData := RefreshTokenResponseIDO{
		Email:                  u.Email,
		AccessToken:            accessToken,
		AccessTokenExpiryDate:  accessTokenExpiryDate,
		RefreshToken:           refreshToken,
		RefreshTokenExpiryDate: refreshTokenExpiryDate,
	}
	if err := json.NewEncoder(w).Encode(&responseData); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
