package controller

import (
	"context"

	"log/slog"

	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/utils/httperror"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (impl *OrganizationControllerImpl) DeleteByID(ctx context.Context, id primitive.ObjectID) error {
	organization, err := impl.GetByID(ctx, id)
	if err != nil {
		impl.Logger.Error("database get by id error", slog.Any("error", err))
		return err
	}
	if organization == nil {
		impl.Logger.Error("database returns nothing from get by id")
		return httperror.NewForBadRequestWithSingleField("id", "organization does not exist")
	}
	if err := impl.OrganizationStorer.DeleteByID(ctx, id); err != nil {
		impl.Logger.Error("database delete by id error", slog.Any("error", err))
		return err
	}
	return nil
}
