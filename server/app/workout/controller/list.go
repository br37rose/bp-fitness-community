package controller

import (
	"context"
	"log/slog"

	ud "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/user/datastore"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/workout/datastore"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/config/constants"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (c *WorkoutControllerImpl) ListByFilter(ctx context.Context, f *datastore.WorkoutListFilter) (*datastore.WorkoutistResult, error) {
	urole, ok := ctx.Value(constants.SessionUserRole).(int8)
	if ok && urole == ud.UserRoleMember {
		userID, ok := ctx.Value(constants.SessionUserID).(primitive.ObjectID)
		if ok {
			f.UserId = userID
		}
	}

	workpouts, err := c.WorkoutStorer.ListByFilter(ctx, f)
	if err != nil {
		c.Logger.Error("database list by filter error", slog.Any("error", err))
		return nil, err
	}
	return workpouts, err
}
