package datastore

import (
	"context"
	"log"
	"time"

	"log/slog"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (impl GoogleFitDataPointStorerImpl) ListByFilter(ctx context.Context, f *GoogleFitDataPointPaginationListFilter) (*GoogleFitDataPointPaginationListResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 12*time.Second)
	defer cancel()

	// Create the paginated filter based on the cursor
	filter, err := impl.newPaginationFilter(f)
	if err != nil {
		return nil, err
	}

	// Apply the conditions to the filter
	if len(f.MetricIDs) > 0 {
		filter["metric_id"] = bson.M{"$in": f.MetricIDs}
	}
	if !f.OrganizationID.IsZero() {
		filter["organization_id"] = f.OrganizationID
	}
	if !f.BranchID.IsZero() {
		filter["branch_id"] = f.BranchID
	}
	if f.ExcludeArchived {
		filter["status"] = bson.M{"$ne": StatusArchived} // Do not list archived items! This code
	}
	if f.Status > 0 {
		filter["status"] = f.Status
	}
	if len(f.DataTypeNames) > 0 {
		filter["data_type_name"] = bson.M{"$in": f.DataTypeNames}
	}

	// Create a slice to store conditions
	var conditions []bson.M

	// Add filter conditions to the slice
	if !f.StartAtGTE.IsZero() {
		conditions = append(conditions, bson.M{"start_at": bson.M{"$gte": f.StartAtGTE}})
	}
	if !f.StartAtGT.IsZero() {
		conditions = append(conditions, bson.M{"start_at": bson.M{"$gt": f.StartAtGT}})
	}
	if !f.StartAtLTE.IsZero() {
		conditions = append(conditions, bson.M{"start_at": bson.M{"$lte": f.StartAtLTE}})
	}
	if !f.StartAtLT.IsZero() {
		conditions = append(conditions, bson.M{"start_at": bson.M{"$lt": f.StartAtLT}})
	}

	// Combine conditions with $and operator
	if len(conditions) > 0 {
		filter["$and"] = conditions
	}

	impl.Logger.Debug("listing filter:",
		slog.Any("filter", filter))

	// Include additional filters for our cursor-based pagination pertaining to sorting and limit.
	options, err := impl.newPaginationOptions(f)
	if err != nil {
		return nil, err
	}

	// Include Full-text search
	if f.SearchText != "" {
		filter["$text"] = bson.M{"$search": f.SearchText}
		options.SetProjection(bson.M{"score": bson.M{"$meta": "textScore"}})
		options.SetSort(bson.D{{"score", bson.M{"$meta": "textScore"}}})
	}

	// Execute the query
	cursor, err := impl.Collection.Find(ctx, filter, options)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	// Retrieve the documents and check if there is a next page
	results := []*GoogleFitDataPoint{}
	hasNextPage := false
	for cursor.Next(ctx) {
		document := &GoogleFitDataPoint{}
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

	return &GoogleFitDataPointPaginationListResult{
		Results:     results,
		NextCursor:  nextCursor,
		HasNextPage: hasNextPage,
	}, nil
}

func (impl GoogleFitDataPointStorerImpl) ListAsSelectOptionByFilter(ctx context.Context, f *GoogleFitDataPointPaginationListFilter) ([]*GoogleFitDataPointAsSelectOption, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 12*time.Second)
	defer cancel()

	// Get a reference to the collection
	collection := impl.Collection

	// Pagination parameters
	pageSize := 10
	startAfter := "" // The ID to start after, initially empty for the first page

	// Sorting parameters
	sortField := "_id"
	sortOrder := 1 // 1=ascending | -1=descending

	// Pagination query
	query := bson.M{}
	options := options.Find().
		SetLimit(int64(pageSize)).
		SetSort(bson.D{{sortField, sortOrder}})

	// Add filter conditions to the query
	if !f.BranchID.IsZero() {
		query["branch_id"] = f.BranchID
	}

	if startAfter != "" {
		// Find the document with the given startAfter ID
		cursor, err := collection.FindOne(ctx, bson.M{"_id": startAfter}).DecodeBytes()
		if err != nil {
			log.Fatal(err)
		}
		options.SetSkip(1)
		query["_id"] = bson.M{"$gt": cursor.Lookup("_id").ObjectID()}
	}

	if f.ExcludeArchived {
		query["status"] = bson.M{"$ne": StatusArchived} // Do not list archived items! This code
	}

	options.SetSort(bson.D{{sortField, 1}}) // Sort in ascending order based on the specified field

	// Retrieve the list of items from the collection
	cursor, err := collection.Find(ctx, query, options)
	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(ctx)

	var results = []*GoogleFitDataPointAsSelectOption{}
	if err = cursor.All(ctx, &results); err != nil {
		panic(err)
	}

	return results, nil
}
