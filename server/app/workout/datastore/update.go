package datastore

import (
	"context"
	"log/slog"

	"go.mongodb.org/mongo-driver/bson"
)

func (impl WorkouStorerImpl) UpdateByID(ctx context.Context, m *Workout) error {
	filter := bson.D{{Key: "_id", Value: m.ID}}

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
