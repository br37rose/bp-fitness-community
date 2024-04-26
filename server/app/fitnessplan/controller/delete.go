package controller

import (
	"context"
	"errors"
	"log/slog"

	"go.mongodb.org/mongo-driver/bson/primitive"

	attch_d "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/fitnessplan/datastore"
	user_d "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/user/datastore"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/config/constants"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/utils/httperror"
)

func (impl *FitnessPlanControllerImpl) DeleteByID(ctx context.Context, id primitive.ObjectID) error {
	// Extract from our session the following data.
	// userID := ctx.Value(constants.SessionUserID).(primitive.ObjectID)
	// userRole := ctx.Value(constants.SessionUserRole).(int8)

	// Apply protection based on ownership and role.
	// if userRole != user_d.UserRoleRoot && userRole != user_d.UserRoleAdmin && user {
	// 	impl.Logger.Error("authenticated user is not staff role error",
	// 		slog.Any("role", userRole),
	// 		slog.Any("userID", userID))
	// 	return httperror.NewForForbiddenWithSingleField("message", "you role does not grant you access to this")
	// }

	// Update the database.
	fitnessplan, err := impl.GetByID(ctx, id)
	fitnessplan.Status = attch_d.StatusArchived
	if err != nil {
		impl.Logger.Error("database get by id error", slog.Any("error", err))
		return err
	}
	if fitnessplan == nil {
		impl.Logger.Error("database returns nothing from get by id")
		return err
	}
	// // Security: Prevent deletion of root user(s).
	// if fitnessplan.Type == attch_d.RootType {
	// 	impl.Logger.Warn("root fitnessplan cannot be deleted error")
	// 	return httperror.NewForForbiddenWithSingleField("role", "root fitnessplan cannot be deleted")
	// }
	fitnessplan.Status = attch_d.StatusArchived

	// Save to the database the modified fitnessplan.
	if err := impl.FitnessPlanStorer.UpdateByID(ctx, fitnessplan); err != nil {
		impl.Logger.Error("database update by id error", slog.Any("error", err))
		return err
	}
	return nil
}

func (impl *FitnessPlanControllerImpl) PermanentlyDeleteByID(ctx context.Context, id primitive.ObjectID) error {
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
	fitnessplan, err := impl.GetByID(ctx, id)
	if err != nil {
		impl.Logger.Error("database get by id error", slog.Any("error", err))
		return err
	}
	if fitnessplan == nil {
		impl.Logger.Error("database returns nothing from get by id")
		return errors.New("does not exist")
	}

	if err := impl.FitnessPlanStorer.DeleteByID(ctx, fitnessplan.ID); err != nil {
		impl.Logger.Error("database delete by id error", slog.Any("error", err))
		return err
	}
	impl.Logger.Debug("deleted from database", slog.Any("fitnessplan_id", id))
	return nil
}
