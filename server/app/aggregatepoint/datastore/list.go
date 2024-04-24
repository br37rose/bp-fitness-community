package datastore

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

func (impl AggregatePointStorerImpl) ListByFilter(ctx context.Context, f *AggregatePointPaginationListFilter) (*AggregatePointPaginationListResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 12*time.Second)
	defer cancel()

	filter, err := impl.newPaginationFilter(f)
	if err != nil {
		return nil, err
	}

	// Add filter conditions to the filter
	if f.Period != 0 {
		filter["period"] = f.Period
	}

	// Create a slice to store conditions
	var conditions []bson.M

	// Add filter conditions to the slice
	if !f.StartGTE.IsZero() {
		conditions = append(conditions, bson.M{"start": bson.M{"$gte": f.StartGTE}})
	}
	if !f.StartGT.IsZero() {
		conditions = append(conditions, bson.M{"start": bson.M{"$gt": f.StartGT}})
	}
	if !f.EndLTE.IsZero() {
		conditions = append(conditions, bson.M{"end": bson.M{"$lte": f.EndLTE}})
	}
	if !f.EndLT.IsZero() {
		conditions = append(conditions, bson.M{"end": bson.M{"$lt": f.EndLT}})
	}

	// Combine conditions with $and operator
	if len(conditions) > 0 {
		filter["$and"] = conditions
	}

	// DEVELOPERS NOTE: We will restrict this list to whatever metric ID's were selected.
	// if len(f.MetricIDs) > 0 {
	// 	filter["metric_id"] = bson.M{"$in": f.MetricIDs}
	// }
	filter["metric_id"] = bson.M{"$in": f.MetricIDs}

	// impl.Logger.Debug("listing filter:",
	// 	slog.Any("filter", filter),
	// 	slog.Any("sort_field", f.SortField),
	// 	slog.Any("sort_order", f.SortOrder))

	// Include additional filters for our cursor-based pagination pertaining to sorting and limit.
	options, err := impl.newPaginationOptions(f)
	if err != nil {
		return nil, err
	}

	// Execute the query
	cursor, err := impl.Collection.Find(ctx, filter, options)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	// Retrieve the documents and check if there is a next page
	results := []*AggregatePoint{}
	hasNextPage := false
	for cursor.Next(ctx) {
		document := &AggregatePoint{}
		if err := cursor.Decode(document); err != nil {
			return nil, err
		}
		results = append(results, document)
		// Stop fetching documents if we have reached the desired page size
		if int64(len(results)) >= f.PageSize {
			hasNextPage = true
			break
		}
	}

	// Get the next cursor and encode it
	var nextCursor string
	if hasNextPage {
		nextCursor, err = impl.newPaginatorNextCursor(f, results)
		if err != nil {
			return nil, err
		}
	}

	return &AggregatePointPaginationListResult{
		Results:     results,
		NextCursor:  nextCursor,
		HasNextPage: hasNextPage,
	}, nil
}
