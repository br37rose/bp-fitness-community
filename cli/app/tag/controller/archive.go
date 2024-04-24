package controller

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"log/slog"

	tag_s "github.com/bci-innovation-labs/bp8fitnesscommunity-cli/app/tag/datastore"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-cli/utils/httperror"
)

func (impl *TagControllerImpl) ArchiveByID(ctx context.Context, id primitive.ObjectID) (*tag_s.Tag, error) {
	// // Extract from our session the following data.
	// userID := ctx.Value(constants.SessionUserID).(primitive.ObjectID)

	// Lookup the tag in our database, else return a `400 Bad Request` error.
	ou, err := impl.TagStorer.GetByID(ctx, id)
	if err != nil {
		impl.Logger.Error("database error", slog.Any("err", err))
		return nil, err
	}
	if ou == nil {
		impl.Logger.Warn("tag does not exist validation error")
		return nil, httperror.NewForBadRequestWithSingleField("id", "does not exist")
	}

	ou.Status = tag_s.TagStatusArchived

	if err := impl.TagStorer.UpdateByID(ctx, ou); err != nil {
		impl.Logger.Error("tag update by id error", slog.Any("error", err))
		return nil, err
	}
	return ou, nil
}
