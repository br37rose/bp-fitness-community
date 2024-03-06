package datastore

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DataPointAggregateResult struct {
	Count   float64 `bson:"count"`
	Average float64 `bson:"average"`
	Min     float64 `bson:"min"`
	Max     float64 `bson:"max"`
	Sum     float64 `bson:"sum"`
}

func (impl DataPointStorerImpl) Aggregate(ctx context.Context, metricID primitive.ObjectID, startAt, endAt time.Time) (*DataPointAggregateResult, error) {
	// Add your date range conditions
	dateMatch := bson.M{
		"timestamp": bson.M{
			"$gte": startAt,
			"$lte": endAt,
		},
	}

	pipeline := []bson.M{
		{
			"$match": bson.M{
				"metric_id": metricID, // Add any other filters as needed
				"$and":      []bson.M{dateMatch},
			},
		},
		{
			"$group": bson.M{
				"_id":     nil,               // Group by all documents
				"count":   bson.M{"$sum": 1}, // Count of documents in the group
				"average": bson.M{"$avg": "$value"},
				"min":     bson.M{"$min": "$value"},
				"max":     bson.M{"$max": "$value"},
				"sum":     bson.M{"$sum": "$value"},
			},
		},
	}

	cursor, err := impl.Collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var result struct {
		Count   float64 `bson:"count"`
		Average float64 `bson:"average"`
		Min     float64 `bson:"min"`
		Max     float64 `bson:"max"`
		Sum     float64 `bson:"sum"`
	}

	if cursor.Next(ctx) {
		if err := cursor.Decode(&result); err != nil {
			return nil, err
		}
	}
	return &DataPointAggregateResult{
		Count:   result.Count,
		Average: result.Average,
		Min:     result.Min,
		Max:     result.Max,
		Sum:     result.Sum,
	}, nil
}
