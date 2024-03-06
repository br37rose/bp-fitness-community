package controller

import (
	"context"

	vcon_s "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/videocontent/datastore"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/utils/httperror"
	"log/slog"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (impl *VideoContentControllerImpl) DeleteByID(ctx context.Context, id primitive.ObjectID) error {
	videocontent, err := impl.GetByID(ctx, id)
	if err != nil {
		impl.Logger.Error("database get by id error", slog.Any("error", err))
		return err
	}
	if videocontent == nil {
		impl.Logger.Error("database returns nothing from get by id")
		return err
	}

	// Defensive code: Do not edit system.
	if videocontent.Type == vcon_s.VideoContentTypeSystem {
		impl.Logger.Error("system videocontent delete error")
		return httperror.NewForBadRequestWithSingleField("message", "the `system` videocontents are read-only and thus access is denied")
	}

	if videocontent.VideoType == vcon_s.VideoContentVideoTypeS3 {
		keys := []string{videocontent.VideoObjectKey}
		if err := impl.S3.DeleteByKeys(ctx, keys); err != nil {
			impl.Logger.Error("s3 delete error", slog.Any("error", err))
			// Do not return an error, simply continue this function as there might
			// be a case were the file was removed on the s3 bucket by ourselves
			// or some other reason.
		}
		impl.Logger.Error("s3 deleted video", slog.Any("keys", keys))
	}

	if videocontent.ThumbnailType == vcon_s.VideoContentThumbnailTypeS3 {
		keys := []string{videocontent.ThumbnailObjectKey}
		if err := impl.S3.DeleteByKeys(ctx, keys); err != nil {
			impl.Logger.Error("s3 delete error", slog.Any("error", err))
			// Do not return an error, simply continue this function as there might
			// be a case were the file was removed on the s3 bucket by ourselves
			// or some other reason.
		}
		impl.Logger.Error("s3 deleted thumbnail", slog.Any("keys", keys))
	}

	if err := impl.VideoContentStorer.DeleteByID(ctx, id); err != nil {
		impl.Logger.Error("database delete by id error", slog.Any("error", err))
		return err
	}
	return nil
}
