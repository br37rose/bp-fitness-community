package datastore

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
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

func (impl GoogleFitDataPointStorerImpl) ListByQueuedStatusAndGfaID(ctx context.Context, googleFitAppID primitive.ObjectID) (*GoogleFitDataPointPaginationListResult, error) {
	f := &GoogleFitDataPointPaginationListFilter{
		Cursor:         "",
		PageSize:       1_000_000_000,
		SortField:      "created_at",
		SortOrder:      1,
		Status:         StatusQueued,
		GoogleFitAppID: googleFitAppID,
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

func (impl GoogleFitDataPointStorerImpl) ListByAnomalousDetection(ctx context.Context) (*GoogleFitDataPointPaginationListResult, error) {
	// STEP 1: List all values which are greater then today as it's not possible
	// to get future health data if it hasn't happened yet!
	nowDT := time.Now()
	futuref := &GoogleFitDataPointPaginationListFilter{
		Cursor:    "",
		PageSize:  1_000_000_000,
		SortField: "created_at",
		SortOrder: 1,
		StartAtGT: nowDT,
	}
	futureres, err := impl.ListByFilter(ctx, futuref)
	if err != nil {
		return nil, err
	}

	res := futureres

	// STEP 2: List all values in the distant past before consumer biometrics
	// tackers were available.
	distantPastDT := time.Date(2000, 1, 2, 22, 00, 22, 380000000, time.UTC) // 2000-01-02 22:00:22.38 -0500 EST
	distantPastF := &GoogleFitDataPointPaginationListFilter{
		Cursor:    "",
		PageSize:  1_000_000_000,
		SortField: "created_at",
		SortOrder: 1,
		StartAtLT: distantPastDT,
	}
	distantPastRes, err := impl.ListByFilter(ctx, distantPastF)
	if err != nil {
		return nil, err
	}
	if distantPastRes == nil {
		err := fmt.Errorf("no results for values less then: %v", distantPastDT)
		return nil, err
	}

	// STEP 3: Append the distant past results to the future results and thus
	// return all the anomalous data.
	for _, v := range distantPastRes.Results {
		res.Results = append(res.Results, v)
	}

	return res, nil
}
