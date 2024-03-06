package datastore

import (
	"context"
	"log/slog"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func (impl FitnessPlanStorerImpl) GetByID(ctx context.Context, id primitive.ObjectID) (*FitnessPlan, error) {
	filter := bson.M{"_id": id}

	var result FitnessPlan
	err := impl.Collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// This error means your query did not match any documents.
			return nil, nil
		}
		impl.Logger.Error("database get by id error", slog.Any("error", err))
		return nil, err
	}
	return &result, nil
}
