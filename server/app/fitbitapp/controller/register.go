package controller

import (
	"context"
	"fmt"
	"log/slog"

	xxx "github.com/nirasan/go-oauth-pkce-code-verifier"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/config/constants"
)

type RegistrationURLResponse struct {
	URL string `bson:"url" json:"url"`
}

// GetRegistrationURL function will generate a URL which is required for the
// user to visit in their browser to begin registering the user's FitBit device
// with our application.
func (c *FitBitAppControllerImpl) GetRegistrationURL(ctx context.Context) (*RegistrationURLResponse, error) {
	// Extract from our session the following data.
	// orgID := ctx.Value(constants.SessionUserOrganizationID).(primitive.ObjectID)
	// orgName := ctx.Value(constants.SessionUserOrganizationName).(string)
	userID := ctx.Value(constants.SessionUserID).(primitive.ObjectID)
	// userName := ctx.Value(constants.SessionUserName).(string)

	// DEVELOPERS NOTE:
	// For more information on how to process authorization then please read the following:
	// ---> https://dev.fitbit.com/build/reference/web-api/developer-guide/authorization/
	v, err := xxx.CreateCodeVerifier() // Create code_verifier
	if err != nil {
		return nil, err
	}
	codeVerifier := v.String()
	codeChallenge := v.CodeChallengeS256() // Create code_challenge with S256 method
	codeChallengeMethod := "S256"

	// DEVELOPERS NOTE: We need to temporary store this value for the oAuth 2.0
	// grant authorization when the user successfully authorizes.
	// (Please see `remote_devices_fitbit_utils` to see how this is used later on.)
	// Please note this function will save the code challenge if it was not
	// previousl saved, else do nothing.
	c.CodeVerifierMap[userID] = codeVerifier

	fitBitURL := fmt.Sprintf("%s?client_id=%s&code_challenge=%s&code_challenge_method=%s&response_type=code",
		constants.FitBitAuthorizationURL,
		c.Config.FitBitApp.ClientID,
		codeChallenge,
		codeChallengeMethod) + "&scope=activity%20cardio_fitness%20electrocardiogram%20heartrate%20location%20nutrition%20oxygen_saturation%20respiratory_rate%20sleep%20temperature%20weight"

	c.Logger.Debug("generated authorization url",
		slog.Any("authorization_url", fitBitURL),
		slog.Any("code_verifier", codeVerifier),
		slog.Any("code_challenge", codeChallenge),
		slog.Any("device", "FitBit"))

	res := &RegistrationURLResponse{
		URL: fitBitURL,
	}
	return res, nil
}
