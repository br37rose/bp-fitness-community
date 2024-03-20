package datastore

import (
	"context"

	"log/slog"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func (impl GoogleFitDataPointStorerImpl) GetByID(ctx context.Context, id primitive.ObjectID) (*GoogleFitDataPoint, error) {
	filter := bson.M{"_id": id}

	var result GoogleFitDataPoint
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

func (impl GoogleFitDataPointStorerImpl) GetByUserID(ctx context.Context, userID primitive.ObjectID) (*GoogleFitDataPoint, error) {
	filter := bson.M{"user_id": userID}

	var result GoogleFitDataPoint
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

func (impl GoogleFitDataPointStorerImpl) GetByName(ctx context.Context, name string) (*GoogleFitDataPoint, error) {
	filter := bson.M{"name": name}

	var result GoogleFitDataPoint
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

func (impl GoogleFitDataPointStorerImpl) GetByPaymentProcessorGoogleFitDataPointID(ctx context.Context, paymentProcessorGoogleFitDataPointID string) (*GoogleFitDataPoint, error) {
	filter := bson.M{"payment_processor_googlefitapp_id": paymentProcessorGoogleFitDataPointID}

	var result GoogleFitDataPoint
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
