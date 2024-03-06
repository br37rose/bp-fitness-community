package controller

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"log/slog"

	gateway_s "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/gateway/datastore"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/utils/httperror"
)

func (impl *GatewayControllerImpl) Login(ctx context.Context, email, password string) (*gateway_s.LoginResponseIDO, error) {
	// Defensive Code: For security purposes we need to remove all whitespaces from the email and lower the characters.
	email = strings.ToLower(email)
	password = strings.ReplaceAll(password, " ", "")

	// Lookup the user in our database, else return a `400 Bad Request` error.
	u, err := impl.UserStorer.GetByEmail(ctx, email)
	if err != nil {
		impl.Logger.Error("database error", slog.Any("err", err))
		return nil, err
	}
	if u == nil {
		impl.Logger.Warn("user does not exist validation error", slog.String("email", email))
		return nil, httperror.NewForBadRequestWithSingleField("email", "does not exist")
	}

	// Verify the inputted password and hashed password match.
	passwordMatch, _ := impl.Password.ComparePasswordAndHash(password, u.PasswordHash)
	if passwordMatch == false {
		impl.Logger.Warn("password check validation error")
		return nil, httperror.NewForBadRequestWithSingleField("password", "password do not match with record")
	}

	// DEVELOPERS NOTE: DISABLED VERIFICATION.
	// // Enforce the verification code of the email.
	// if u.WasEmailVerified == false {
	// 	impl.Logger.Warn("email verification validation error", slog.Any("u", u))
	// 	return nil, httperror.NewForBadRequestWithSingleField("email", "was not verified")
	// }

	// Enforce 2FA if enabled.
	if u.OTPEnabled {
		// We need to reset the `otp_validated` status to be false to force
		// the user to use their `totp authenticator` application.
		u.OTPValidated = false
		u.ModifiedAt = time.Now()
		if err := impl.UserStorer.UpdateByID(ctx, u); err != nil {
			impl.Logger.Error("failed updating user during login", slog.Any("err", err))
			return nil, err
		}
	}

	uBin, err := json.Marshal(u)
	if err != nil {
		impl.Logger.Error("marshalling error", slog.Any("err", err))
		return nil, err
	}

	// Set expiry duration.
	atExpiry := 24 * time.Hour
	rtExpiry := 14 * 24 * time.Hour

	// Start our session using an access and refresh token.
	sessionUUID := impl.UUID.NewUUID()

	err = impl.Cache.SetWithExpiry(ctx, sessionUUID, string(uBin), rtExpiry)
	if err != nil {
		impl.Logger.Error("cache set with expiry error", slog.Any("err", err))
		return nil, err
	}

	// Generate our JWT token.
	accessToken, accessTokenExpiry, refreshToken, refreshTokenExpiry, err := impl.JWT.GenerateJWTTokenPair(sessionUUID, atExpiry, rtExpiry)
	if err != nil {
		impl.Logger.Error("jwt generate pairs error", slog.Any("err", err))
		return nil, err
	}

	// For security purposes.
	u.PasswordHashAlgorithm = ""
	u.PasswordHash = ""
	u.EmailVerificationCode = ""
	u.EmailVerificationExpiry = time.Now()

	// Return our auth keys.
	return &gateway_s.LoginResponseIDO{
		User:                   u,
		AccessToken:            accessToken,
		AccessTokenExpiryTime:  accessTokenExpiry,
		RefreshToken:           refreshToken,
		RefreshTokenExpiryTime: refreshTokenExpiry,
	}, nil
}
