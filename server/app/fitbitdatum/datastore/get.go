package datastore

import (
	"context"

	"log/slog"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func (impl FitBitDatumStorerImpl) GetByID(ctx context.Context, id primitive.ObjectID) (*FitBitDatum, error) {
	filter := bson.M{"_id": id}

	var result FitBitDatum
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

func (impl FitBitDatumStorerImpl) GetByUserID(ctx context.Context, userID primitive.ObjectID) (*FitBitDatum, error) {
	filter := bson.M{"user_id": userID}

	var result FitBitDatum
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

func (impl FitBitDatumStorerImpl) GetByName(ctx context.Context, name string) (*FitBitDatum, error) {
	filter := bson.M{"name": name}

	var result FitBitDatum
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

func (impl FitBitDatumStorerImpl) GetByPaymentProcessorFitBitDatumID(ctx context.Context, paymentProcessorFitBitDatumID string) (*FitBitDatum, error) {
	filter := bson.M{"payment_processor_fitbitapp_id": paymentProcessorFitBitDatumID}

	var result FitBitDatum
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
