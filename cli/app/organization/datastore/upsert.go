package datastore

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (impl OrganizationStorerImpl) UpsertByID(ctx context.Context, o *Organization) error {
	opts := options.Update().SetUpsert(true) // Use upsert option

	filter := bson.M{"_id": o.ID}

	update := bson.M{"$set": o}

	_, err := impl.Collection.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return err
	}

	return nil
}
