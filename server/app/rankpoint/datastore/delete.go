package datastore

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (impl RankPointStorerImpl) DeleteByID(ctx context.Context, id primitive.ObjectID) error {
	_, err := impl.Collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		log.Fatal("DeleteOne() ERROR:", err)
	}
	return nil
}

func (impl RankPointStorerImpl) DeleteAll(ctx context.Context) error {
	_, err := impl.Collection.DeleteMany(ctx, bson.M{})
	if err != nil {
		return err
	}
	return nil
}
