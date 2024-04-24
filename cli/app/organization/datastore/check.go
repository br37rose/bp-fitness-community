package datastore

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"log/slog"
)

func (impl OrganizationStorerImpl) CheckIfExistsByName(ctx context.Context, name string) (bool, error) {
	filter := bson.D{{"name", name}}
	count, err := impl.Collection.CountDocuments(ctx, filter)
	if err != nil {
		impl.Logger.Error("database check if exists by email error", slog.Any("error", err))
		return false, err
	}
	return count >= 1, nil
}
