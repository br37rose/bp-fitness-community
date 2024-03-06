package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	fitbitapp_s "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/fitbitapp/datastore"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/config/constants"
)

func (c *FitBitAppControllerImpl) refreshAccessTokenForFitBit(
	ctx context.Context,
	fba *fitbitapp_s.FitBitApp,
) error {
	duration := fba.ExpireTime.Sub(time.Now())
	durInHours := int64(duration.Hours())

	// DEVELOPERS NOTE:
	// ACCORDING TO THE DOCUMENTATION, THE REFRESH TOKEN DOES NOT EXPIRE SO
	// WE DO NOT NEED TO WORRY ABOUT THIS CASE. SEE:
	// https://dev.fitbit.com/build/reference/web-api/developer-guide/best-practices/#Using-Tokens-Effectively

	// Run the token refresh process if we are about to expire our access token.
	if durInHours <= 0 {
		//
		// STEP 1: Get New Access Token Using The Refresh Token
		//

		// Get our client ID and client secret.
		config := c.Config
		clientID := config.FitBitApp.ClientID
		clientSecret := config.FitBitApp.ClientSecret

		// Generate the bearer.
		bearerValueRaw := fmt.Sprintf("%v:%v", clientID, clientSecret)
		bearerValue := Base64EncodeStrippedFromString(bearerValueRaw)
		bearer := fmt.Sprintf("Basic %s", string(bearerValue))

		xxcurl := fmt.Sprintf("%s?grant_type=refresh_token&refresh_token=%s",
			constants.FitBitExchangeURL,
			fba.RefreshToken,
		)

		req, err := http.NewRequest("POST", xxcurl, nil)
		if err != nil {
			c.Logger.Error("new request error",
				slog.Any("fitbit_app_id", fba.ID),
				slog.Any("xxcurl", xxcurl),
				slog.Any("error", err))
			return err
		}

		// Set headers
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.Header.Add("Authorization", bearer)

		// Set client timeout
		client := &http.Client{Timeout: time.Second * 25}

		// Send req using http Client
		resp, err := client.Do(req)
		if err != nil {
			c.Logger.Error("fetch error",
				slog.Any("fitbit_app_id", fba.ID),
				slog.Any("req", req),
				slog.Any("error", err))
			return err
		}

		defer resp.Body.Close()

		//
		// STEP 2: Unmarshal.
		//

		data := FitBitAuthReponse{}

		json.NewDecoder(resp.Body).Decode(&data)
		if err != nil {
			c.Logger.Error("decode error",
				slog.Any("fitbit_app_id", fba.ID),
				slog.Any("resp.Body", resp.Body),
				slog.Any("error", err))
			return err
		}

		//
		// STEP 3: Update our record.
		//

		expireTime := time.Now().Add(time.Duration(data.ExpiresIn))

		// // For debugging purposes only.
		// bc.Logger.Info().Str("taskID", taskID).Str("taskName", "RefreshRemoteDeviceAccessTokens").
		// 	Str("OLD AccessToken", rd.AccessToken).
		// 	Str("OLD RefreshToken", rd.RefreshToken).
		// 	Time("OLD ExpireTime", rd.ExpireTime).
		// 	Str("NEW AccessToken", data.AccessToken).
		// 	Str("NEW RefreshToken", data.RefreshToken).
		// 	Time("NEW ExpireTime", expireTime).
		// 	Str("taskName", "RefreshRemoteDeviceAccessTokens").
		// 	Msg("")

		fba.TokenType = data.TokenType
		fba.AccessToken = data.AccessToken
		fba.ExpiresIn = data.ExpiresIn
		fba.RefreshToken = data.RefreshToken
		fba.ExpireTime = expireTime
		fba.LastFetchedAt = time.Now()
		fba.ModifiedAt = time.Now()

		err = c.FitBitAppStorer.UpdateByID(ctx, fba)
		if err != nil {
			c.Logger.Error("failed updating fitbit app",
				slog.Any("fitbit_app_id", fba.ID),
				slog.Any("error", err))
			return err
		}
	}

	return nil
}
