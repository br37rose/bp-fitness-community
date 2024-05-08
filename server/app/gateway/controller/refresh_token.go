package controller

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"log/slog"

	user_s "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/user/datastore"
)

func (impl *GatewayControllerImpl) RefreshToken(ctx context.Context, value string) (*user_s.User, string, time.Time, string, time.Time, error) {
	////
	//// Extract the `sessionID` so we can process it.
	////

	sessionID, err := impl.JWT.ProcessJWTToken(value)
	if err != nil {
		impl.Logger.Warn("process jwt refresh token does not exist", slog.String("value", value))
		err := errors.New("jwt refresh token failed")
		return nil, "", time.Now(), "", time.Now(), err
	}

	////
	//// Lookup in our in-memory the user record for the `sessionID` or error.
	////

	uBin, err := impl.Cache.Get(ctx, sessionID)
	if err != nil {
		impl.Logger.Error("in-memory get error", slog.Any("err", err))
		return nil, "", time.Now(), "", time.Now(), err
	}

	var u *user_s.User
	err = json.Unmarshal([]byte(uBin), &u)
	if err != nil {
		impl.Logger.Error("unmarshal error", slog.Any("err", err))
		return nil, "", time.Now(), "", time.Now(), err
	}

	////
	//// Generate new access and refresh tokens and return them.
	////

	// Set expiry duration.
	atExpiry := 24 * time.Hour
	rtExpiry := 14 * 24 * time.Hour

	// Start our session using an access and refresh token.
	newSessionUUID := impl.UUID.NewUUID()

	err = impl.Cache.SetWithExpiry(ctx, newSessionUUID, uBin, rtExpiry)
	if err != nil {
		impl.Logger.Error("cache set with expiry error", slog.Any("err", err))
		return nil, "", time.Now(), "", time.Now(), err
	}

	// Generate our JWT token.
	accessToken, accessTokenExpiry, refreshToken, refreshTokenExpiry, err := impl.JWT.GenerateJWTTokenPair(newSessionUUID, atExpiry, rtExpiry)
	if err != nil {
		impl.Logger.Error("jwt generate pairs error", slog.Any("err", err))
		return nil, "", time.Now(), "", time.Now(), err
	}

	// Return our auth keys.
	return u, accessToken, accessTokenExpiry, refreshToken, refreshTokenExpiry, nil
}
