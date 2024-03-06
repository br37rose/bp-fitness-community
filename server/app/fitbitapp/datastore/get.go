package datastore

import (
	"context"

	"log/slog"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func (impl FitBitAppStorerImpl) GetByID(ctx context.Context, id primitive.ObjectID) (*FitBitApp, error) {
	filter := bson.M{"_id": id}

	var result FitBitApp
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

func (impl FitBitAppStorerImpl) GetByUserID(ctx context.Context, userID primitive.ObjectID) (*FitBitApp, error) {
	filter := bson.M{"user_id": userID}

	var result FitBitApp
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

func (impl FitBitAppStorerImpl) GetByName(ctx context.Context, name string) (*FitBitApp, error) {
	filter := bson.M{"name": name}

	var result FitBitApp
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

func (impl FitBitAppStorerImpl) GetByPaymentProcessorFitBitAppID(ctx context.Context, paymentProcessorFitBitAppID string) (*FitBitApp, error) {
	filter := bson.M{"payment_processor_fitbitapp_id": paymentProcessorFitBitAppID}

	var result FitBitApp
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
