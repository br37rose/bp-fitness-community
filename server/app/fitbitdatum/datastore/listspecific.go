package datastore

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (impl FitBitDatumStorerImpl) ListByQueuedStatus(ctx context.Context) (*FitBitDatumListResult, error) {
	f := &FitBitDatumListFilter{
		Cursor:    primitive.NilObjectID,
		PageSize:  1_000_000_000,
		SortField: "_id",
		SortOrder: 1,
		Status:    StatusQueued,
	}
	return impl.ListByFilter(ctx, f)
}
