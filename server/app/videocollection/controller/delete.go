package controller

import (
	"context"

	vcol_d "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/videocollection/datastore"
	"log/slog"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (impl *VideoCollectionControllerImpl) DeleteByID(ctx context.Context, id primitive.ObjectID) error {
	videocollection, err := impl.GetByID(ctx, id)
	if err != nil {
		impl.Logger.Error("database get by id error", slog.Any("error", err))
		return err
	}
	if videocollection == nil {
		impl.Logger.Error("database returns nothing from get by id")
		return err
	}

	// // Defensive code: Do not edit system.
	// if videocollection.Type == vcol_d.VideoCollectionTypeSystem {
	// 	impl.Logger.Error("system videocollection delete error")
	// 	return httperror.NewForBadRequestWithSingleField("message", "the `system` videocollections are read-only and thus access is denied")
	// }

	// if videocollection.VideoType == vcol_d.VideoCollectionVideoTypeS3 {
	// 	keys := []string{videocollection.VideoObjectKey}
	// 	if err := impl.S3.DeleteByKeys(ctx, keys); err != nil {
	// 		impl.Logger.Error("s3 delete error", slog.Any("error", err))
	// 		// Do not return an error, simply continue this function as there might
	// 		// be a case were the file was removed on the s3 bucket by ourselves
	// 		// or some other reason.
	// 	}
	// 	impl.Logger.Error("s3 deleted video", slog.Any("keys", keys))
	// }

	if videocollection.ThumbnailType == vcol_d.VideoCollectionThumbnailTypeS3 {
		keys := []string{videocollection.ThumbnailObjectKey}
		if err := impl.S3.DeleteByKeys(ctx, keys); err != nil {
			impl.Logger.Error("s3 delete error", slog.Any("error", err))
			// Do not return an error, simply continue this function as there might
			// be a case were the file was removed on the s3 bucket by ourselves
			// or some other reason.
		}
		impl.Logger.Error("s3 deleted thumbnail", slog.Any("keys", keys))
	}

	if err := impl.VideoCollectionStorer.DeleteByID(ctx, id); err != nil {
		impl.Logger.Error("database delete by id error", slog.Any("error", err))
		return err
	}
	return nil
}
