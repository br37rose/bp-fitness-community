package datastore

import (
	"context"
)

func (impl GoogleFitDataPointStorerImpl) ListByQueuedStatus(ctx context.Context) (*GoogleFitDataPointPaginationListResult, error) {
	f := &GoogleFitDataPointPaginationListFilter{
		Cursor:    "",
		PageSize:  1_000_000_000,
		SortField: "created_at",
		SortOrder: 1,
		Status:    StatusQueued,
	}
	return impl.ListByFilter(ctx, f)
}

func (impl GoogleFitDataPointStorerImpl) ListByQueuedStatusInDataTypeNames(ctx context.Context, dataTypeNames []string) (*GoogleFitDataPointPaginationListResult, error) {
	f := &GoogleFitDataPointPaginationListFilter{
		Cursor:        "",
		PageSize:      1_000_000_000,
		SortField:     "created_at",
		SortOrder:     1,
		Status:        StatusQueued,
		DataTypeNames: dataTypeNames,
	}
	return impl.ListByFilter(ctx, f)
}

func (impl GoogleFitDataPointStorerImpl) ListByActiveStatusInDataTypeNames(ctx context.Context, dataTypeNames []string) (*GoogleFitDataPointPaginationListResult, error) {
	f := &GoogleFitDataPointPaginationListFilter{
		Cursor:        "",
		PageSize:      1_000_000_000,
		SortField:     "created_at",
		SortOrder:     1,
		Status:        StatusActive,
		DataTypeNames: dataTypeNames,
	}
	return impl.ListByFilter(ctx, f)
}
