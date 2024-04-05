package datastore

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

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
