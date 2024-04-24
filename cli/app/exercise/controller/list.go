package controller

import (
	"context"

	domain "github.com/bci-innovation-labs/bp8fitnesscommunity-cli/app/exercise/datastore"
	"log/slog"
)

func (c *ExerciseControllerImpl) ListByFilter(ctx context.Context, f *domain.ExerciseListFilter) (*domain.ExerciseListResult, error) {
	// // Extract from our session the following data.
	// userID := ctx.Value(constants.SessionUserID).(primitive.ObjectID)
	// userRole := ctx.Value(constants.SessionUserRole).(int8)
	//
	// // Apply filtering based on ownership and role.
	// if userRole != user_d.UserRoleRoot {
	// 	f.UserID = userID
	// 	f.UserRole = userRole
	// }

	m, err := c.ExerciseStorer.ListByFilter(ctx, f)
	if err != nil {
		c.Logger.Error("database list by filter error", slog.Any("error", err))
		return nil, err
	}
	return m, err
}
