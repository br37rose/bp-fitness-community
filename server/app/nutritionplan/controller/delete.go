package controller

import (
	"context"
	"errors"
	"log/slog"

	"go.mongodb.org/mongo-driver/bson/primitive"

	attch_d "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/nutritionplan/datastore"
	user_d "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/user/datastore"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/config/constants"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/utils/httperror"
)

func (impl *NutritionPlanControllerImpl) DeleteByID(ctx context.Context, id primitive.ObjectID) error {
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
	nutritionplan, err := impl.GetByID(ctx, id)
	nutritionplan.Status = attch_d.StatusArchived
	if err != nil {
		impl.Logger.Error("database get by id error", slog.Any("error", err))
		return err
	}
	if nutritionplan == nil {
		impl.Logger.Error("database returns nothing from get by id")
		return err
	}
	// // Security: Prevent deletion of root user(s).
	// if nutritionplan.Type == attch_d.RootType {
	// 	impl.Logger.Warn("root nutritionplan cannot be deleted error")
	// 	return httperror.NewForForbiddenWithSingleField("role", "root nutritionplan cannot be deleted")
	// }
	nutritionplan.Status = attch_d.StatusArchived

	// Save to the database the modified nutritionplan.
	if err := impl.NutritionPlanStorer.UpdateByID(ctx, nutritionplan); err != nil {
		impl.Logger.Error("database update by id error", slog.Any("error", err))
		return err
	}
	return nil
}

func (impl *NutritionPlanControllerImpl) PermanentlyDeleteByID(ctx context.Context, id primitive.ObjectID) error {
	// Extract from our session the following data.
	// userID := ctx.Value(constants.SessionUserID).(primitive.ObjectID)
	// userRole := ctx.Value(constants.SessionUserRole).(int8)

	// Apply protection based on ownership and role.
	// if userRole != user_d.UserRoleRoot && userRole != user_d.UserRoleAdmin {
	// 	impl.Logger.Error("authenticated user is not staff role error",
	// 		slog.Any("role", userRole),
	// 		slog.Any("userID", userID))
	// 	return httperror.NewForForbiddenWithSingleField("message", "you role does not grant you access to this")
	// }

	// Update the database.
	nutritionplan, err := impl.GetByID(ctx, id)
	if err != nil {
		impl.Logger.Error("database get by id error", slog.Any("error", err))
		return err
	}
	if nutritionplan == nil {
		impl.Logger.Error("database returns nothing from get by id")
		return errors.New("does not exist")
	}

	if err := impl.NutritionPlanStorer.DeleteByID(ctx, nutritionplan.ID); err != nil {
		impl.Logger.Error("database delete by id error", slog.Any("error", err))
		return err
	}
	impl.Logger.Debug("deleted from database", slog.Any("nutritionplan_id", id))
	return nil
}
