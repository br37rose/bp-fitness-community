package datastore

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

func (impl DataPointStorerImpl) ListByFilter(ctx context.Context, f *DataPointPaginationListFilter) (*DataPointPaginationListResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 12*time.Second)
	defer cancel()

	filter, err := impl.newPaginationFilter(f)
	if err != nil {
		return nil, err
	}

	// Create a slice to store conditions
	var conditions []bson.M

	// Add filter conditions to the slice
	if !f.GTE.IsZero() {
		conditions = append(conditions, bson.M{"timestamp": bson.M{"$gte": f.GTE}})
	}
	if !f.GT.IsZero() {
		conditions = append(conditions, bson.M{"timestamp": bson.M{"$gt": f.GT}})
	}
	if !f.LTE.IsZero() {
		conditions = append(conditions, bson.M{"timestamp": bson.M{"$lte": f.LTE}})
	}
	if !f.LT.IsZero() {
		conditions = append(conditions, bson.M{"timestamp": bson.M{"$lt": f.LT}})
	}

	// Combine conditions with $and operator
	if len(conditions) > 0 {
		filter["$and"] = conditions
	}

	// DEVELOPERS NOTE: We will restrict this list to whatever metric ID's were selected.
	filter["metric_id"] = bson.M{"$in": f.MetricIDs}

	// impl.Logger.Debug("listing filter:",
	// 	slog.Any("filter", filter))

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
	results := []*DataPoint{}
	hasNextPage := false
	for cursor.Next(ctx) {
		document := &DataPoint{}
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

	return &DataPointPaginationListResult{
		Results:     results,
		NextCursor:  nextCursor,
		HasNextPage: hasNextPage,
	}, nil
}
