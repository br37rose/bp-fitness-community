package datastore

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"log/slog"
)

func (impl VideoCategoryStorerImpl) UpdateByID(ctx context.Context, m *VideoCategory) error {
	filter := bson.D{{"_id", m.ID}}

	update := bson.M{ // DEVELOPERS NOTE: https://stackoverflow.com/a/60946010
		"$set": m,
	}

	// execute the UpdateOne() function to update the first matching document
	_, err := impl.Collection.UpdateOne(ctx, filter, update)
	if err != nil {
		impl.Logger.Error("database update by id error", slog.Any("error", err))
	}

	return nil
}
