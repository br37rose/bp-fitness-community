package crontab

import (
	"context"
	"fmt"
	"log/slog"

	gfa_ds "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/googlefitapp/datastore"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/oauth2"
)

func (impl *googleFitAppCrontaberImpl) RefreshTokensFromGoogleJob() error {
	ctx := context.Background()
	impl.Logger.Debug("refreshing tokens from google")
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

func (impl *googleFitAppCrontaberImpl) refreshTokenFromGoogle(ctx context.Context, gfaID primitive.ObjectID) error {
	// Lock this google fit app
	impl.Kmutex.Lockf("googlefitapp_%v", gfaID.Hex())
	defer impl.Kmutex.Unlockf("googlefitapp_%v", gfaID.Hex())

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
	return nil
}
