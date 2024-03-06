package datastore

import (
	"context"
	"time"

	"log/slog"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (impl RankPointStorerImpl) CheckIfExistsByCompositeKey(ctx context.Context, metricID primitive.ObjectID, timestamp time.Time) (bool, error) {
	filter := bson.M{}
	filter["metric_id"] = metricID
	filter["timestamp"] = timestamp
	count, err := impl.Collection.CountDocuments(ctx, filter)
	if err != nil {
		impl.Logger.Error("database check if exists by composite key error", slog.Any("error", err))
		return false, err
	}
	return count >= 1, nil
}
