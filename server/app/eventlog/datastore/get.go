package datastore

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"log/slog"
)

func (impl EventLogStorerImpl) GetByID(ctx context.Context, id primitive.ObjectID) (*EventLog, error) {
	filter := bson.M{"_id": id}

	var result EventLog
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

func (impl EventLogStorerImpl) GetByName(ctx context.Context, name string) (*EventLog, error) {
	filter := bson.M{"name": name}

	var result EventLog
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

func (impl EventLogStorerImpl) GetByPaymentProcessorEventLogID(ctx context.Context, paymentProcessorEventLogID string) (*EventLog, error) {
	filter := bson.M{"payment_processor_eventlog_id": paymentProcessorEventLogID}

	var result EventLog
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
