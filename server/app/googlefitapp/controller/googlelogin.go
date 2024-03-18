package controller

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"log/slog"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/config/constants"
)

type GoogleLoginURLResponse struct {
	URL string `bson:"url" json:"url"`
}

// GetRegistrationURL function will generate a URL which is required for the
// user to visit in their browser to begin registering the user's GoogleFit device
// with our application.
func (impl *GoogleFitAppControllerImpl) GetGoogleLoginURL(ctx context.Context) (*GoogleLoginURLResponse, error) {

	// // Extract from our session the following data.
	// // orgID := ctx.Value(constants.SessionUserOrganizationID).(primitive.ObjectID)
	// // orgName := ctx.Value(constants.SessionUserOrganizationName).(string)
	userID := ctx.Value(constants.SessionUserID).(primitive.ObjectID)

	b := make([]byte, 16)
	rand.Read(b)
	oauthState := base64.URLEncoding.EncodeToString(b)

	// // DEVELOPERS NOTE: We need to temporary store this value for the oAuth 2.0
	// // grant authorization when the user successfully authorizes.
	// // (Please see `remote_devices_googlefit_utils` to see how this is used later on.)
	// // Please note this function will save the code challenge if it was not
	// // previousl saved, else do nothing.
	impl.CodeVerifierMap[userID] = oauthState

	googleFitURL := impl.GCP.OAuth2GenerateAuthCodeURL(oauthState)

	impl.Logger.Debug("generated authorization url",
		slog.Any("authorization_url", googleFitURL),
		slog.Any("code_verifier", oauthState),
		slog.Any("web_service", "google oauth2"))

	res := &GoogleLoginURLResponse{
		URL: googleFitURL,
	}
	return res, nil
}
