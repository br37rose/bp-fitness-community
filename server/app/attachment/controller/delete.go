package controller

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"log/slog"

	attch_d "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/attachment/datastore"
	user_d "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/user/datastore"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/config/constants"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/utils/httperror"
)

func (impl *AttachmentControllerImpl) DeleteByID(ctx context.Context, id primitive.ObjectID) error {
	// Extract from our session the following data.
	userID := ctx.Value(constants.SessionUserID).(primitive.ObjectID)
	userRole := ctx.Value(constants.SessionUserRole).(int8)

	// Apply protection based on ownership and role.
	if userRole != user_d.UserRoleRoot && userRole != user_d.UserRoleAdmin {
		impl.Logger.Error("authenticated user is not staff role error",
			slog.Any("role", userRole),
			slog.Any("userID", userID))
		return httperror.NewForForbiddenWithSingleField("message", "you role does not grant you access to this")
	}

	// Update the database.
	attachment, err := impl.GetByID(ctx, id)
	attachment.Status = attch_d.StatusArchived
	if err != nil {
		impl.Logger.Error("database get by id error", slog.Any("error", err))
		return err
	}
	if attachment == nil {
		impl.Logger.Error("database returns nothing from get by id")
		return err
	}
	// // Security: Prevent deletion of root user(s).
	// if attachment.Type == attch_d.RootType {
	// 	impl.Logger.Warn("root attachment cannot be deleted error")
	// 	return httperror.NewForForbiddenWithSingleField("role", "root attachment cannot be deleted")
	// }

	// Save to the database the modified attachment.
	if err := impl.AttachmentStorer.UpdateByID(ctx, attachment); err != nil {
		impl.Logger.Error("database update by id error", slog.Any("error", err))
		return err
	}
	return nil
}

func (impl *AttachmentControllerImpl) PermanentlyDeleteByID(ctx context.Context, id primitive.ObjectID) error {
	// Extract from our session the following data.
	userID := ctx.Value(constants.SessionUserID).(primitive.ObjectID)
	userRole := ctx.Value(constants.SessionUserRole).(int8)

	// Apply protection based on ownership and role.
	if userRole != user_d.UserRoleRoot && userRole != user_d.UserRoleAdmin {
		impl.Logger.Error("authenticated user is not staff role error",
			slog.Any("role", userRole),
			slog.Any("userID", userID))
		return httperror.NewForForbiddenWithSingleField("message", "you role does not grant you access to this")
	}

	// Update the database.
	attachment, err := impl.GetByID(ctx, id)
	if err != nil {
		impl.Logger.Error("database get by id error", slog.Any("error", err))
		return err
	}
	if attachment == nil {
		impl.Logger.Error("database returns nothing from get by id")
		return errors.New("does not exist")
	}

	// Proceed to delete the physical files from AWS s3.
	if err := impl.S3.DeleteByKeys(ctx, []string{attachment.ObjectKey}); err != nil {
		impl.Logger.Warn("s3 delete by keys error", slog.Any("error", err))
		// Do not return an error, simply continue this function as there might
		// be a case were the file was removed on the s3 bucket by ourselves
		// or some other reason.
	}
	impl.Logger.Debug("deleted from s3", slog.Any("attachment_id", id))

	if err := impl.AttachmentStorer.DeleteByID(ctx, attachment.ID); err != nil {
		impl.Logger.Error("database delete by id error", slog.Any("error", err))
		return err
	}
	impl.Logger.Debug("deleted from database", slog.Any("attachment_id", id))

	// Update exercise.
	if !attachment.OwnershipID.IsZero() {
		if attachment.OwnershipType == attch_d.OwnershipTypeExerciseVideo {
			impl.Logger.Debug("started updating exercise...", slog.Any("ownership_id", attachment.OwnershipID))

			e, err := impl.ExerciseStorer.GetByID(ctx, attachment.OwnershipID)
			if err != nil {
				impl.Logger.Error("database get by id error", slog.Any("error", err))
				return err
			}
			if e == nil {
				impl.Logger.Error("database returns nothing from get by id")
				return errors.New("does not exist")
			}
			e.VideoAttachmentID = primitive.NilObjectID
			e.VideoAttachmentFilename = ""
			e.VideoObjectKey = ""
			e.VideoObjectURL = ""
			e.VideoObjectExpiry = time.Now()
			e.VideoName = ""
			e.VideoURL = ""
			if err := impl.ExerciseStorer.UpdateByID(ctx, e); err != nil {
				impl.Logger.Error("database delete by id error", slog.Any("error", err))
				return err
			}
			impl.Logger.Debug("finished updating exercise", slog.Any("ownership_id", attachment.OwnershipID))
		}

		if attachment.OwnershipType == attch_d.OwnershipTypeExerciseThumbnail {
			impl.Logger.Debug("started updating exercise...", slog.Any("ownership_id", attachment.OwnershipID))

			e, err := impl.ExerciseStorer.GetByID(ctx, attachment.OwnershipID)
			if err != nil {
				impl.Logger.Error("database get by id error", slog.Any("error", err))
				return err
			}
			if e == nil {
				impl.Logger.Error("database returns nothing from get by id")
				return errors.New("does not exist")
			}
			e.ThumbnailAttachmentID = primitive.NilObjectID
			e.ThumbnailAttachmentFilename = ""
			e.ThumbnailObjectKey = ""
			e.ThumbnailObjectURL = ""
			e.ThumbnailObjectExpiry = time.Now()
			e.ThumbnailName = ""
			e.ThumbnailURL = ""
			if err := impl.ExerciseStorer.UpdateByID(ctx, e); err != nil {
				impl.Logger.Error("database delete by id error", slog.Any("error", err))
				return err
			}
			impl.Logger.Debug("finished updating exercise", slog.Any("ownership_id", attachment.OwnershipID))
		}
	}
	return nil
}
