package datastore

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (impl GoogleFitDataPointStorerImpl) ListByQueuedStatus(ctx context.Context) (*GoogleFitDataPointListResult, error) {
	f := &GoogleFitDataPointListFilter{
		Cursor:    primitive.NilObjectID,
		PageSize:  1_000_000_000,
		SortField: "_id",
		SortOrder: 1,
		Status:    StatusQueued,
	}
	return impl.ListByFilter(ctx, f)
}
