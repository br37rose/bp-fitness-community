package controller

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"log/slog"

	user_d "github.com/bci-innovation-labs/bp8fitnesscommunity-cli/app/user/datastore"
	attch_d "github.com/bci-innovation-labs/bp8fitnesscommunity-cli/app/videocategory/datastore"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-cli/config/constants"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-cli/utils/httperror"
)

func (impl *VideoCategoryControllerImpl) DeleteByID(ctx context.Context, id primitive.ObjectID) error {
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
	videocategory, err := impl.GetByID(ctx, id)
	videocategory.Status = attch_d.StatusArchived
	if err != nil {
		impl.Logger.Error("database get by id error", slog.Any("error", err))
		return err
	}
	if videocategory == nil {
		impl.Logger.Error("database returns nothing from get by id")
		return err
	}
	// // Security: Prevent deletion of root user(s).
	// if videocategory.Type == attch_d.RootType {
	// 	impl.Logger.Warn("root videocategory cannot be deleted error")
	// 	return httperror.NewForForbiddenWithSingleField("role", "root videocategory cannot be deleted")
	// }
	videocategory.Status = attch_d.StatusArchived

	// Save to the database the modified videocategory.
	if err := impl.VideoCategoryStorer.UpdateByID(ctx, videocategory); err != nil {
		impl.Logger.Error("database update by id error", slog.Any("error", err))
		return err
	}
	return nil
}

func (impl *VideoCategoryControllerImpl) PermanentlyDeleteByID(ctx context.Context, id primitive.ObjectID) error {
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
	videocategory, err := impl.GetByID(ctx, id)
	if err != nil {
		impl.Logger.Error("database get by id error", slog.Any("error", err))
		return err
	}
	if videocategory == nil {
		impl.Logger.Error("database returns nothing from get by id")
		return errors.New("does not exist")
	}

	if err := impl.VideoCategoryStorer.DeleteByID(ctx, videocategory.ID); err != nil {
		impl.Logger.Error("database delete by id error", slog.Any("error", err))
		return err
	}
	impl.Logger.Debug("deleted from database", slog.Any("videocategory_id", id))
	return nil
}

// Auto-generated comment for change 27
