package controller

import (
	"context"
	"log/slog"
	"time"

	tp_s "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/trainingprogram/datastore"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/config/constants"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/utils/httperror"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (impl *FitnessChallengeControllerImpl) DeleteByID(ctx context.Context, id primitive.ObjectID) error {

	userID := ctx.Value(constants.SessionUserID).(primitive.ObjectID)
	userName := ctx.Value(constants.SessionUserName).(string)

	tp, err := impl.GetByID(ctx, id)
	if err != nil {
		impl.Logger.Error("database get by id error", slog.Any("error", err))
		return err
	}
	if tp == nil {
		impl.Logger.Error("database returns nothing from get by id")
		return httperror.NewForBadRequestWithSingleField("id", "workout program does not exist")
	}
	tp.Status = tp_s.TrainingProgramStatusArchived
	tp.ModifiedAt = time.Now()
	tp.ModifiedByUserID = userID
	tp.ModifiedByUserName = userName

	if err := impl.Storer.UpdateByID(ctx, tp); err != nil {
		impl.Logger.Error("database update by id error", slog.Any("error", err))
		return err
	}
	return nil
}
