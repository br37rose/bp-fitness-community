package controller

import (
	"context"
	"fmt"
	"log/slog"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type GoogleCallbackResponse struct {
	URL string `bson:"url" json:"url"`
}

func (impl *GoogleFitAppControllerImpl) GoogleCallback(ctx context.Context, state, code string) (*GoogleCallbackResponse, error) {
	impl.Logger.Debug("google callback to bp8 fitness community system",
		slog.Any("code", code),
		slog.Any("state", state),
		slog.Any("web_service", "google oauth2"))

	userIDs := make([]primitive.ObjectID, 0, len(impl.CodeVerifierMap))
	for k := range impl.CodeVerifierMap {
		userIDs = append(userIDs, k)
	}

	// Variable tracks whether the `state` provided by Google matches something
	// we have in our records. Assume nothing matches.
	var wasCodeVerified bool = false

	// Iterate through all the verification codes and try to match with our
	// `state` provided by Google. If match is made then proceed with process
	// it.
	for _, userID := range userIDs {
		codeVerifier := impl.CodeVerifierMap[userID]
		if state == codeVerifier {
			if err := impl.attemptAuthorizationForKey(ctx, userID, code); err != nil {
				impl.Logger.Error("google callback failed attempt authorization",
					slog.Any("web_service", "google oauth2"),
					slog.Any("user_id", userID),
					slog.Any("code", code),
					slog.Any("error", err))
				return nil, err
			}
			wasCodeVerified = true
		}
	}

	// If the `state` provided by Google does not exist in our system then
	// we need to generate an error and do not proceed any further.
	if wasCodeVerified == false {

		// For debugging purposes only:
		// The following code will use a test record to verify the code works.
		// Only run this for testing or developing reasons. When you are done
		// please comment out.
		tmpUserID, _ := primitive.ObjectIDFromHex("65b16a259c54a618ae7fd05e")
		tmpCode := "4/0AeaYSHAQYH5ySpc2zjxSqyLyQ7StvtUzQgUPzptpmylP4I9KwxZImw8uQzV4-pP0gio4Rg"
		if err := impl.attemptAuthorizationForKey(ctx, tmpUserID, tmpCode); err != nil {
			impl.Logger.Error("google callback failed attempt authorization",
				slog.Any("error", err))
			return nil, err
		}

		// err := httperror.NewForBadRequestWithSingleField("state", "was not verified with bp8 fitness community system")
		// impl.Logger.Error("google callback failed verifying state",
		// 	slog.Any("web_service", "google oauth2"),
		// 	slog.Any("state", state),
		// 	slog.Any("error", err))
		// return nil, err
	}

	return &GoogleCallbackResponse{URL: impl.Config.FitBitApp.RegistrationSuccessRedirectURL}, nil
}

func (impl *GoogleFitAppControllerImpl) attemptAuthorizationForKey(ctx context.Context, userID primitive.ObjectID, code string) error {

	token, err := impl.GCP.ExchangeCode(code)
	if err != nil {
		impl.Logger.Error("google callback failed exchanging code",
			slog.String("user_id", userID.Hex()),
			slog.String("code", code),
			slog.Any("error", err))
	}

	return fmt.Errorf("halt because of %s with response: %v", "programmer", token)
	return nil
}
