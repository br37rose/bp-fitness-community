package datastore

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (impl OfferStorerImpl) Upsert(ctx context.Context, offer *Offer) error {
	opts := options.Update().SetUpsert(true) // Use upsert option

	filter := bson.M{"_id": offer.ID}

	update := bson.M{"$set": offer}

	_, err := impl.Collection.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return err
	}

	return nil
}
