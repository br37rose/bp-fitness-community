package datastore

import (
	"context"
	"log/slog"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func (impl AggregatePointStorerImpl) GetByID(ctx context.Context, id primitive.ObjectID) (*AggregatePoint, error) {
	filter := bson.M{"_id": id}

	var result AggregatePoint
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

func (impl AggregatePointStorerImpl) GetByCompositeKey(ctx context.Context, metricID primitive.ObjectID, period int8, start time.Time, end time.Time) (*AggregatePoint, error) {
	filter := bson.M{
		"metric_id": metricID,
		"period":    period,
		"start":     start,
		"end":       end,
	}

	var result AggregatePoint
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
