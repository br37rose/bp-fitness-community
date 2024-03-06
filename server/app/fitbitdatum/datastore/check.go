package datastore

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log/slog"
)

func (impl FitBitDatumStorerImpl) CheckIfExistsByNameInOrgBranch(ctx context.Context, name string, orgID primitive.ObjectID, branchID primitive.ObjectID) (bool, error) {
	filter := bson.M{}
	filter["name"] = name
	filter["organization_id"] = orgID
	filter["branch_id"] = branchID
	count, err := impl.Collection.CountDocuments(ctx, filter)
	if err != nil {
		impl.Logger.Error("database check if exists by email error", slog.Any("error", err))
		return false, err
	}
	return count >= 1, nil
}
