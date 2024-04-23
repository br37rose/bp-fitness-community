package controller

import (
	"context"
	"log/slog"

	u_d "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/user/datastore"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/config/constants"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/utils/httperror"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (impl *QuestionControllerImpl) DeleteByID(ctx context.Context, id primitive.ObjectID) error {
	urole, ok := ctx.Value(constants.SessionUserRole).(int8)
	if ok && urole != u_d.UserRoleAdmin {
		return httperror.NewForBadRequestWithSingleField("message", "you role does not grant you access to this")
	}
	err := impl.QuestionStorer.DeleteByID(ctx, id)
	if err != nil {
		impl.Logger.Error("database get by id error", slog.Any("error", err))
		return err
	}

	return nil
}
