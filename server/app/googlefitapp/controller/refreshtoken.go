package controller

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/oauth2"

	gfa_ds "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/googlefitapp/datastore"
)

func (impl *GoogleFitAppControllerImpl) RefreshTokensFromGoogle() error {
	ctx := context.Background()
	gfaIDs, err := impl.GoogleFitAppStorer.ListPhysicalIDsByStatus(ctx, gfa_ds.StatusActive)
	if err != nil {
		impl.Logger.Error("failed listing google fit apps by status",
			slog.Any("error", err))
		return err
	}
	for _, gfaID := range gfaIDs {
		if err := impl.refreshTokenFromGoogle(ctx, gfaID); err != nil {
			impl.Logger.Error("failed pulling data for google fit app",
				slog.Any("error", err))
			return err
		}
	}
	return nil
}

func (impl *GoogleFitAppControllerImpl) refreshTokenFromGoogle(ctx context.Context, gfaID primitive.ObjectID) error {
	// Lock this google fit app
	impl.DistributedMutex.Lockf(ctx, "googlefitapp_%v", gfaID.Hex())
	defer impl.DistributedMutex.Unlockf(ctx, "googlefitapp_%v", gfaID.Hex())

	impl.Logger.Debug("checking gfa",
		slog.String("gfa_id", gfaID.Hex()))

	// Get our database record.
	gfa, err := impl.GoogleFitAppStorer.GetByID(ctx, gfaID)
	if err != nil {
		impl.Logger.Error("failed getting google fit app from database",
			slog.Any("error", err))
		return err
	}
	if gfa == nil {
		err := fmt.Errorf("google fit app does not exist for id: %s", gfaID.Hex())
		return err
	}

	// Authenticated http client for a specific user's account. Note: No need
	// for refresh token handling as it's already handled!
	client, err := impl.GCP.NewHTTPClientFromToken(gfa.Token, func(newTok *oauth2.Token) {
		// Save the latest token provided by Google in our database and
		// make sure our status is set to running and having no problems.
		gfa.Token = newTok
		gfa.RequiresGoogleLoginAgain = false
		gfa.Status = gfa_ds.StatusActive
		gfa.Errors = ""
		if err := impl.GoogleFitAppStorer.UpdateByID(ctx, gfa); err != nil {
			impl.Logger.Error("failed updating google fit app in database",
				slog.Any("error", err))
		}
		impl.Logger.Debug("updated google fit app with new token")
	})
	if err != nil {
		impl.Logger.Error("detected error when refreshing google fit token from oauth",
			slog.String("gfa_id", gfaID.Hex()),
			slog.String("user_id", gfa.UserID.Hex()),
			slog.Any("error", err),
		)

		//
		// If any errors occur let's force the user to log in again.
		//

		gfa.RequiresGoogleLoginAgain = true
		gfa.Status = gfa_ds.StatusError
		gfa.Errors = err.Error()
		if err := impl.GoogleFitAppStorer.UpdateByID(ctx, gfa); err != nil {
			impl.Logger.Error("failed updating google fit app in database",
				slog.Any("error", err))
		}

		u, err := impl.UserStorer.GetByID(ctx, gfa.UserID)
		if err != nil {
			impl.Logger.Error("failed getting user from database",
				slog.Any("error", err))
			return err
		}
		if u == nil {
			err := fmt.Errorf("user does not exist for id: %s", gfa.UserID.Hex())
			return err
		}
		u.PrimaryHealthTrackingDeviceRequiresLoginAgain = true
		if err := impl.UserStorer.UpdateByID(ctx, u); err != nil {
			impl.Logger.Error("failed updating user to database",
				slog.Any("error", err))
			return err
		}

		impl.Logger.Warn("access and refresh token expired, user must log in again for google fit app",
			slog.String("gfa_id", gfaID.Hex()),
			slog.String("user_id", gfa.UserID.Hex()),
		)

		return nil
	}
	if client == nil {
		err := fmt.Errorf("google fit app authenticated client does not exist for token: %v", gfa.Token)
		return err
	}

	expiryDur := time.Since(gfa.Token.Expiry)
	expiryDurInMins := expiryDur.Hours() * 60 * (-1)

	// Check to see if the token is approaching expiration date and if so then fetch a new one.
	// Note: This is a manual attempt by us and should be taken care of automatically by `oauth2` library.
	if expiryDurInMins <= 5 {
		impl.Logger.Debug("gfa token needs refreshing, attempting now...",
			slog.String("gfa_id", gfaID.Hex()),
			slog.String("user_id", gfa.UserID.Hex()),
			slog.Time("token_expiry", gfa.Token.Expiry),
			slog.Float64("token_mins_unitl_expiry", expiryDurInMins),
		)

		newTok, err := impl.GCP.NewTokenFromExistingToken(gfa.Token)
		if err != nil {
			impl.Logger.Error("failed getting new token from existing token",
				slog.Any("error", err))
			return err
		}
		gfa.Token = newTok
		if err := impl.GoogleFitAppStorer.UpdateByID(ctx, gfa); err != nil {
			impl.Logger.Error("failed updating google fit app in database",
				slog.Any("error", err))
		}

		expiryDur = time.Since(newTok.Expiry)
		expiryDurInMins = expiryDur.Hours() * 60 * (-1)
	}

	impl.Logger.Debug("checked gfa is ok",
		slog.String("gfa_id", gfaID.Hex()),
		slog.String("user_id", gfa.UserID.Hex()),
		slog.Time("token_expiry", gfa.Token.Expiry),
		slog.Float64("token_mins_unitl_expiry", expiryDurInMins),
	)
	return nil
}
