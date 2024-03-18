package controller

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	gfa_ds "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/googlefitapp/datastore"
	u_s "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/user/datastore"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/utils/httperror"
)

type GoogleCallbackResponse struct {
	URL string `bson:"url" json:"url"`
}

func (impl *GoogleFitAppControllerImpl) GoogleCallback(ctx context.Context, state, code string) (*GoogleCallbackResponse, error) {
	////
	//// Start the transaction.
	////

	session, err := impl.DbClient.StartSession()
	if err != nil {
		impl.Logger.Error("start session error",
			slog.Any("error", err))
		return nil, err
	}
	defer session.EndSession(ctx)

	// Define a transaction function with a series of operations
	transactionFunc := func(sessCtx mongo.SessionContext) (interface{}, error) {
		impl.Logger.Debug("google callback to bp8 fitness community system",
			slog.Any("code", code),
			slog.Any("state", state))

		userIDs := make([]primitive.ObjectID, 0, len(impl.CodeVerifierMap))
		for k := range impl.CodeVerifierMap {
			userIDs = append(userIDs, k)
		}

		// Iterate through all the verification codes and try to match with our
		// `state` provided by Google. If match is made then proceed with process
		// it.
		for _, userID := range userIDs {
			codeVerifier := impl.CodeVerifierMap[userID]
			if state == codeVerifier {
				if err := impl.attemptAuthorizationForKey(sessCtx, userID, code); err != nil {
					impl.Logger.Error("google callback failed attempt authorization",
						slog.Any("user_id", userID),
						slog.Any("code", code),
						slog.Any("error", err))
					return nil, err
				}

				////
				//// End transaction with success.
				////
				return &GoogleCallbackResponse{URL: impl.Config.FitBitApp.RegistrationSuccessRedirectURL}, nil
			}
		}

		////
		//// End transaction with error.
		////

		// If the `state` provided by Google does not exist in our system then
		// we need to generate an error and do not proceed any further.
		err := httperror.NewForBadRequestWithSingleField("state", "was not verified with bp8 fitness community system")
		impl.Logger.Error("google callback failed verifying state",
			slog.Any("state", state),
			slog.Any("error", err))
		return nil, err
	}

	// Start a transaction
	res, err := session.WithTransaction(ctx, transactionFunc)
	if err != nil {
		impl.Logger.Error("session failed error",
			slog.Any("error", err))
		return nil, err
	}

	return res.(*GoogleCallbackResponse), nil
}

func (impl *GoogleFitAppControllerImpl) attemptAuthorizationForKey(sessCtx mongo.SessionContext, userID primitive.ObjectID, code string) error {
	token, err := impl.GCP.OAuth2ExchangeCode(code)
	if err != nil {
		impl.Logger.Error("google callback failed exchanging code",
			slog.String("user_id", userID.Hex()),
			slog.String("code", code),
			slog.Any("error", err))
		return err
	}
	if token == nil {
		err := fmt.Errorf("no token exchanged for code: %s", code)
		impl.Logger.Error("google callback exchanging code did not get token",
			slog.String("user_id", userID.Hex()),
			slog.Any("error", err))
		return err
	}

	impl.Logger.Debug("successfully exchanged authorization code with access token from google",
		slog.Any("code", code),
		slog.Any("token", token))

	u, err := impl.UserStorer.GetByID(sessCtx, userID)
	if err != nil {
		impl.Logger.Error("failed getting user by id",
			slog.String("user_id", userID.Hex()),
			slog.String("code", code),
			slog.Any("error", err))
		return err
	}
	if u == nil {
		err := fmt.Errorf("user does not exist for `user_id`: %s", userID)
		impl.Logger.Error("google callback failed getting user",
			slog.String("user_id", userID.Hex()),
			slog.Any("error", err))
		return err
	}

	// Get previous record to update or create a new record.
	gfa, err := impl.GoogleFitAppStorer.GetByUserID(sessCtx, u.ID)
	if err != nil {
		impl.Logger.Error("failed getting google fit app by user id",
			slog.String("user_id", userID.Hex()),
			slog.String("code", code),
			slog.Any("error", err))
		return err
	}
	if gfa == nil {
		gfa = &gfa_ds.GoogleFitApp{
			ID:                       primitive.NewObjectID(),
			UserFirstName:            u.FirstName,
			UserLastName:             u.LastName,
			UserName:                 u.Name,
			UserLexicalName:          u.LexicalName,
			UserID:                   u.ID,
			Status:                   gfa_ds.StatusActive,
			CreatedAt:                time.Now(),
			ModifiedAt:               time.Now(),
			OrganizationID:           u.OrganizationID,
			OrganizationName:         u.OrganizationName,
			AuthType:                 gfa_ds.AuthTypeOAuth2,
			Errors:                   "",
			Token:                    token,
			LastFetchedAt:            time.Date(2014, 1, 1, 00, 00, 00, 000000000, time.UTC), // 2013-01-01 00:00:00.00 UTC
			HeartRateMetricID:        primitive.NewObjectID(),
			StepsCountMetricID:       primitive.NewObjectID(),
			IsTestMode:               false,
			SimulatorAlgorithm:       "",
			RequiresGoogleLoginAgain: false,
		}
		if err := impl.GoogleFitAppStorer.Create(sessCtx, gfa); err != nil {
			impl.Logger.Error("failed creating google fit app in database",
				slog.String("user_id", userID.Hex()),
				slog.String("code", code),
				slog.Any("gfa", gfa),
				slog.Any("error", err))
			return err
		}

		impl.Logger.Debug("created google fit app")
	} else {
		gfa.Token = token
		gfa.RequiresGoogleLoginAgain = false
		if err := impl.GoogleFitAppStorer.UpdateByID(sessCtx, gfa); err != nil {
			impl.Logger.Error("failed updating google fit app in database",
				slog.String("user_id", userID.Hex()),
				slog.String("code", code),
				slog.Any("gfa", gfa),
				slog.Any("error", err))
			return err
		}
		impl.Logger.Debug("updated google fit app")
	}

	// Update our user with our new Google Fit registration / login.
	u.PrimaryHealthTrackingDeviceType = u_s.UserPrimaryHealthTrackingDeviceTypeFitBit
	u.PrimaryHealthTrackingDeviceRequiresLoginAgain = false
	u.PrimaryHealthTrackingDeviceHeartRateMetricID = gfa.HeartRateMetricID
	u.PrimaryHealthTrackingDeviceStepsCountMetricID = gfa.StepsCountMetricID
	u.GoogleFitAppID = gfa.ID
	u.ModifiedAt = time.Now()
	if err := impl.UserStorer.UpdateByID(sessCtx, u); err != nil {
		impl.Logger.Error("failed updating user by id",
			slog.Any("code", code),
			slog.Any("user_id", userID),
			slog.Any("error", err))
		return err
	}

	impl.Logger.Debug("updated user account with google fit app")
	return nil
}
