package datastore

import (
	"context"
	"log/slog"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (impl DataPointStorerImpl) GetByCompositeKey(ctx context.Context, metricID primitive.ObjectID, timestamp time.Time) (*DataPoint, error) {
	filter := bson.M{}
	filter["metric_id"] = metricID
	filter["timestamp"] = timestamp

	var result DataPoint
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

// GetByLatestTimestampAndMetricID retrieves the data point with the latest timestamp for a given metric.
func (impl DataPointStorerImpl) GetByLatestTimestampAndMetricID(ctx context.Context, metricID primitive.ObjectID) (*DataPoint, error) {
	filter := bson.M{"metric_id": metricID}

	// Specify options to sort by timestamp in descending order and limit to one result
	findOptions := options.FindOne().SetSort(bson.M{"timestamp": -1})

	var result DataPoint

	// Find the document with the latest timestamp for the given metric
	err := impl.Collection.FindOne(ctx, filter, findOptions).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// No data found for the specified metric and timestamp.
			return nil, nil
		}
		return nil, err
	}

	return &result, nil
}
