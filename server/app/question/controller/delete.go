package controller

import (
	"context"
	"log/slog"
	"time"

	u_d "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/user/datastore"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/config/constants"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/utils/httperror"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (impl *QuestionControllerImpl) DeleteByID(ctx context.Context, id primitive.ObjectID) error {
	userID := ctx.Value(constants.SessionUserID).(primitive.ObjectID)
	userName := ctx.Value(constants.SessionUserName).(string)
	urole, ok := ctx.Value(constants.SessionUserRole).(int8)
	if ok && urole != u_d.UserRoleAdmin {
		return httperror.NewForBadRequestWithSingleField("message", "you role does not grant you access to this")
	}
	q, err := impl.GetByID(ctx, id)
	if err != nil {
		impl.Logger.Error("database get by id error", slog.Any("error", err))
		return err
	}
	if q == nil {
		impl.Logger.Error("database returns nothing from get by id")
		return httperror.NewForBadRequestWithSingleField("id", "workout does not exist")
	}
	q.Status = false
	q.CreatedAt = time.Now()
	q.CreatedByUserID = userID
	q.CreatedByUserName = userName
	q.ModifiedAt = time.Now()
	q.ModifiedByUserName = userName
	q.ModifiedByUserID = userID

	if err := impl.QuestionStorer.UpdateByID(ctx, q); err != nil {
		impl.Logger.Error("database update by id error", slog.Any("error", err))
		return err
	}
	return nil
}
