package datastore

import (
	"context"

	"log/slog"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func (impl GoogleFitAppStorerImpl) GetByID(ctx context.Context, id primitive.ObjectID) (*GoogleFitApp, error) {
	filter := bson.M{"_id": id}

	var result GoogleFitApp
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

func (impl GoogleFitAppStorerImpl) GetByUserID(ctx context.Context, userID primitive.ObjectID) (*GoogleFitApp, error) {
	filter := bson.M{"user_id": userID}

	var result GoogleFitApp
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

func (impl GoogleFitAppStorerImpl) GetByName(ctx context.Context, name string) (*GoogleFitApp, error) {
	filter := bson.M{"name": name}

	var result GoogleFitApp
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

func (impl GoogleFitAppStorerImpl) GetByPaymentProcessorGoogleFitAppID(ctx context.Context, paymentProcessorGoogleFitAppID string) (*GoogleFitApp, error) {
	filter := bson.M{"payment_processor_googlefitapp_id": paymentProcessorGoogleFitAppID}

	var result GoogleFitApp
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
