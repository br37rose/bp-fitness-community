package datastore

import (
	"context"
	"time"

	"log/slog"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (impl GoogleFitDataPointStorerImpl) CheckIfExistsByCompositeKey(ctx context.Context, userID primitive.ObjectID, dataTypeName string, startAt time.Time, endAt time.Time) (bool, error) {
	filter := bson.M{}
	filter["user_id"] = userID
	filter["data_type_name"] = dataTypeName
	filter["start_at"] = startAt
	filter["end_at"] = endAt
	count, err := impl.Collection.CountDocuments(ctx, filter)
	if err != nil {
		impl.Logger.Error("database check if exists by pk error", slog.Any("error", err))
		return false, err
	}
	return count >= 1, nil
}
