package datastore

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log/slog"
)

func (impl ExerciseStorerImpl) CheckIfExistsByID(ctx context.Context, id primitive.ObjectID) (bool, error) {
	filter := bson.D{{"_id", id}}
	count, err := impl.Collection.CountDocuments(ctx, filter)
	if err != nil {
		impl.Logger.Error("database check if exists by id error", slog.Any("error", err))
		return false, err
	}
	return count >= 1, nil
}
