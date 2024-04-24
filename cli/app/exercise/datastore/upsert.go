package datastore

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (impl ExerciseStorerImpl) UpsertByID(ctx context.Context, e *Exercise) error {
	opts := options.Update().SetUpsert(true) // Use upsert option

	filter := bson.M{"_id": e.ID}

	update := bson.M{"$set": e}

	_, err := impl.Collection.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return err
	}

	return nil
}
