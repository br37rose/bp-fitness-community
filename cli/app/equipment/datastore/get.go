package datastore

import (
	"context"
	"log/slog"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func (impl EquipmentStorerImpl) GetByID(ctx context.Context, id primitive.ObjectID) (*Equipment, error) {
	filter := bson.M{"_id": id}

	var result Equipment
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

func (impl EquipmentStorerImpl) GetByName(ctx context.Context, name string) (*Equipment, error) {
	filter := bson.M{"name": name}

	var result Equipment
	err := impl.Collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// This error means your query did not match any documents.
			return nil, nil
		}
		impl.Logger.Error("database get by name error", slog.Any("error", err))
		return nil, err
	}
	return &result, nil
}
