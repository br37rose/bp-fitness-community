package datastore

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (impl FitBitAppStorerImpl) UpsertByUserID(ctx context.Context, fba *FitBitApp) error {
	opts := options.Update().SetUpsert(true) // Use upsert option

	filter := bson.M{"user_id": fba.UserID}

	update := bson.M{"$set": fba}

	_, err := impl.Collection.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return err
	}

	return nil
}
