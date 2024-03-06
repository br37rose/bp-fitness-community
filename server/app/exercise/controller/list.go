package controller

import (
	"context"
	"errors"
	"log/slog"
	"time"

	domain "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/exercise/datastore"
	u_s "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/user/datastore"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/config/constants"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (c *ExerciseControllerImpl) ListByFilter(ctx context.Context, f *domain.ExerciseListFilter) (*domain.ExerciseListResult, error) {
	// // Extract from our session the following data.
	userID, _ := ctx.Value(constants.SessionUserID).(primitive.ObjectID)
	// userRole := ctx.Value(constants.SessionUserRole).(int8)
	orgID, _ := ctx.Value(constants.SessionUserOrganizationID).(primitive.ObjectID)

	// Apply filtering organization tenancy.
	f.OrganizationID = orgID

	listRes, err := c.ExerciseStorer.ListByFilter(ctx, f)
	if err != nil {
		c.Logger.Error("database list by filter error", slog.Any("error", err))
		return nil, err
	}

	// DEVELOPERS NOTE:
	// We need to restrict offerings if the user previously purchased it.
	u, err := c.UserStorer.GetByID(ctx, userID)
	if err != nil {
		c.Logger.Error("database get error", slog.Any("error", err))
		return nil, err
	}
	if u == nil {
		return nil, errors.New("user does not exist")
	}

	// Iterate through all the exercises and refresh the presigned URL if it
	// has expired. Furthermore we will restrict content if that particular
	// content requires offer purchase.
	for _, e := range listRes.Results {

		// Refresh the presigned URL if it has expired or else skip this step.
		if e.VideoType == domain.ExerciseVideoTypeS3 && e.VideoObjectKey != "" {
			if time.Now().After(e.VideoObjectExpiry) {
				c.Kmutex.Lockf("exercise-%v-for-video", e.ID.Hex())         // Step 1
				defer c.Kmutex.Unlockf("exercise-%v-for-video", e.ID.Hex()) // Step 2

				// Generate a presigned URL for today.
				expiryDur := time.Hour * 12
				videoObjectURL, presignErr := c.S3.GetPresignedURL(ctx, e.VideoObjectKey, expiryDur)
				if presignErr != nil {
					c.Logger.Error("video s3 presign url error", slog.Any("presignErr", presignErr))
					return nil, err
				}

				// Update the exercise.
				e.VideoObjectURL = videoObjectURL
				e.VideoObjectExpiry = time.Now().Add(expiryDur)
				if err := c.ExerciseStorer.UpdateByID(ctx, e); err != nil {
					c.Logger.Error("exercise database update by id error", slog.Any("error", err))
					return nil, err
				}
			}
		}

		// Refresh the presigned URL if it has expired or else skip this step.
		if e.ThumbnailType == domain.ExerciseThumbnailTypeS3 && e.ThumbnailObjectKey != "" {
			if time.Now().After(e.ThumbnailObjectExpiry) {
				c.Kmutex.Lockf("exercise-%v-for-thumbnail", e.ID.Hex()) // Step 1
				defer func() {
					c.Kmutex.Unlockf("exercise-%v-for-thumbnail", e.ID.Hex()) // Step 2
				}()

				// Generate a presigned URL for today.
				expiryDur := time.Hour * 12
				thumbnailObjectURL, presignErr := c.S3.GetPresignedURL(ctx, e.ThumbnailObjectKey, expiryDur)
				if presignErr != nil {
					c.Logger.Error("thumbnail s3 presign url error", slog.Any("presignErr", presignErr))
					return nil, err
				}

				// Update the exercise.
				e.ThumbnailObjectURL = thumbnailObjectURL
				e.ThumbnailObjectExpiry = time.Now().Add(expiryDur)
				if err := c.ExerciseStorer.UpdateByID(ctx, e); err != nil {
					c.Logger.Error("exercise database update by id error", slog.Any("error", err))
					return nil, err
				}
			}
		}

		c.Logger.Debug("exercise thumbnail",
			slog.Any("IsExpired", time.Now().After(e.ThumbnailObjectExpiry)),
			slog.Any("ThumbnailObjectURL", e.ThumbnailObjectURL),
			slog.Any("ThumbnailObjectKey", e.ThumbnailObjectKey),
			slog.Any("ThumbnailObjectExpiry", e.ThumbnailObjectExpiry))

		// Restrict access based on if user is member and whether the member
		// made the purchse to be granted access to this exercise.
		if e.HasMonetization && u.Role == u_s.UserRoleMember {
			// STEP 1: Do you have the correct membership rank? AKA. Did you
			//         make the correct purchase?
			if u.OfferMembershipRank < e.OfferMembershipRank {
				c.Logger.Warn("exercise access denied reason: did not purchase")
				e.CurrentUserHasAccessGranted = false
				e.VideoObjectURL = "[HIDDEN]"
				e.VideoURL = "[HIDDEN]"
				e.VideoObjectKey = "[HIDDEN]"

				// Do not continue this loop and just start loading the next
				// loop.
				continue
			}

			// STEP 2: Is this exercise time locked and do you have access?
			if e.HasTimedLock {
				// Create the date expected to unlock to the user.
				expectedUnlockDate := u.CreatedAt.Add(e.TimedLockDuration)

				// If today is not past the the expected unlock date then
				// let us lock this exercise to the user.
				if time.Now().Before(expectedUnlockDate) {
					c.Logger.Warn("exercise access denied reason: timed lock")
					e.CurrentUserHasAccessGranted = false
					e.VideoObjectURL = "[HIDDEN]"
					e.VideoURL = "[HIDDEN]"
					e.VideoObjectKey = "[HIDDEN]"

					// Do not continue this loop and just start loading the next
					// loop.
					continue
				}
			}
		}
	}
	return listRes, err
}
