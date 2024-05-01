package controller

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"log/slog"
	"time"

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

	// DEVELOPERS NOTE: We need to temporary store this value for the oAuth 2.0
	// grant authorization when the user successfully authorizes.
	// (Please see `remote_devices_googlefit_utils` to see how this is used later on.)
	// Please note this function will save the code challenge if it was not
	// previousl saved, else do nothing.
	if err := impl.setCodeVerifier(ctx, userID, oauthState); err != nil {
		return nil, err
	}

	googleFitURL := impl.GCP.OAuth2GenerateAuthCodeURL(oauthState)

	impl.Logger.DebugContext(ctx, "generated authorization url",
		slog.Any("authorization_url", googleFitURL),
		slog.Any("code_verifier", oauthState),
		slog.Any("web_service", "google oauth2"))

	res := &GoogleLoginURLResponse{
		URL: googleFitURL,
	}
	return res, nil
}

func (impl *GoogleFitAppControllerImpl) setCodeVerifier(ctx context.Context, userID primitive.ObjectID, oauthState string) error {
	impl.Logger.DebugContext(ctx, "locking code verifier", slog.String("func", "setCodeVerifier"))
	impl.Kmutex.Lock("google-code-verifier")
	defer impl.Kmutex.Unlock("google-code-verifier")
	defer impl.Logger.DebugContext(ctx, "unlocking code verifier", slog.String("func", "setCodeVerifier"))

	var codeVerifierMap map[string]primitive.ObjectID
	str, err := impl.Cache.Get(ctx, "google-code-verifier")
	if err != nil {
		impl.Logger.WarnContext(ctx, "failed getting code verifier from cache", slog.Any("err", err))
		codeVerifierMap = make(map[string]primitive.ObjectID, 0)
	}
	if str != "" {
		if err := json.Unmarshal([]byte(str), &codeVerifierMap); err != nil {
			impl.Logger.WarnContext(ctx, "failed unmarshalling code verifier", slog.Any("err", err))
			codeVerifierMap = make(map[string]primitive.ObjectID, 0)
		}
		impl.Logger.DebugContext(ctx, "unmarshalled code verifier successfully")
	}

	codeVerifierMap[oauthState] = userID
	bin, err := json.Marshal(codeVerifierMap)
	if err != nil {
		impl.Logger.WarnContext(ctx, "failed marshalling code verifier", slog.Any("err", err))
		return err
	}
	impl.Logger.DebugContext(ctx, "marshalled code verifier successfully",
		slog.Any("code_verifier_map", codeVerifierMap))
	return impl.Cache.SetWithExpiry(ctx, "google-code-verifier", string(bin), 15*time.Minute)
}

func (impl *GoogleFitAppControllerImpl) searchForUserIdFromCodeVerifier(ctx context.Context, oauthState string) (primitive.ObjectID, error) {
	impl.Logger.DebugContext(ctx, "locking code verifier", slog.String("func", "searchForUserIdFromCodeVerifier"))
	impl.Kmutex.Lock("google-code-verifier")
	defer impl.Kmutex.Unlock("google-code-verifier")
	defer impl.Logger.DebugContext(ctx, "unlocked code verifier", slog.String("func", "searchForUserIdFromCodeVerifier"))

	var codeVerifierMap map[string]primitive.ObjectID
	str, err := impl.Cache.Get(ctx, "google-code-verifier")
	if err != nil {
		impl.Logger.WarnContext(ctx, "failed getting code verifier from cache", slog.Any("err", err))
		codeVerifierMap = make(map[string]primitive.ObjectID, 0)
	}
	if str != "" {
		if err := json.Unmarshal([]byte(str), &codeVerifierMap); err != nil {
			impl.Logger.WarnContext(ctx, "failed unmarshalling code verifier", slog.Any("err", err))
			codeVerifierMap = make(map[string]primitive.ObjectID, 0)
		}
	}

	impl.Logger.DebugContext(ctx, "successfully unmarshalled code verifier, preparing to lookup `oauth_state` ...",
		slog.Any("code_verifier_map", codeVerifierMap),
		slog.Any("oauth_state", oauthState))

	userID := codeVerifierMap[oauthState]
	if !userID.IsZero() {
		impl.Logger.DebugContext(ctx, "successfully found user_id in code verifier",
			slog.Any("user_id", userID),
			slog.Any("code_verifier_map", codeVerifierMap))
		return userID, nil
	}

	impl.Logger.WarnContext(ctx, "failled finding user_id in code verifier",
		slog.Any("oauth_state", oauthState),
		slog.Any("code_verifier_map", codeVerifierMap))
	return primitive.NilObjectID, nil
}
