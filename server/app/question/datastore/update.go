package datastore

import (
	"context"
	"log/slog"

	"go.mongodb.org/mongo-driver/bson"
)

func (impl *QuestionStorerImpl) UpdateByID(ctx context.Context, q *Question) error {
	filter := bson.D{{Key: "_id", Value: q.ID}}

	update := bson.M{
		"$set": q,
	}
	_, err := impl.Collection.UpdateOne(ctx, filter, update)
	if err != nil {
		impl.Logger.Error("database update by user id error", slog.Any("error", err))
	}

	return nil
}
