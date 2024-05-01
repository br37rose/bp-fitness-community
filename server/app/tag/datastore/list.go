package datastore

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (impl TagStorerImpl) ListByFilter(ctx context.Context, f *TagPaginationListFilter) (*TagPaginationListResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 12*time.Second)
	defer cancel()

	// Create the paginated filter based on the cursor
	filter, err := impl.newPaginationFilter(f)
	if err != nil {
		return nil, err
	}

	// Add filter conditions to the filter
	if !f.OrganizationID.IsZero() {
		filter["organization_id"] = f.OrganizationID
	}
	if f.Status != 0 {
		filter["status"] = f.Status
	}
	if f.SearchText != "" {
		filter["text"] = bson.M{"$regex": primitive.Regex{Pattern: f.SearchText, Options: "i"}}
	}

	// impl.Logger.Debug("listing filter:",
	// 	slog.Any("filter", filter))

	// Include additional filters for our cursor-based pagination pertaining to sorting and limit.
	options, err := impl.newPaginationOptions(f)
	if err != nil {
		return nil, err
	}

	// Include Full-text search
	// if f.SearchText != "" {
	// 	filter["$text"] = bson.M{"$search": f.SearchText}
	// 	options.SetProjection(bson.M{"score": bson.M{"$meta": "textScore"}})
	// 	options.SetSort(bson.D{{"score", bson.M{"$meta": "textScore"}}})
	// }

	// Execute the query
	cursor, err := impl.Collection.Find(ctx, filter, options)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	// var results = []*ComicSubmission{}
	// if err = cursor.All(ctx, &results); err != nil {
	// 	panic(err)
	// }

	// Retrieve the documents and check if there is a next page
	results := []*Tag{}
	hasNextPage := false
	for cursor.Next(ctx) {
		document := &Tag{}
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

	return &TagPaginationListResult{
		Results:     results,
		NextCursor:  nextCursor,
		HasNextPage: hasNextPage,
	}, nil
}

func (impl TagStorerImpl) ListAsSelectOptionByFilter(ctx context.Context, f *TagListFilter) ([]*TagAsSelectOption, error) {
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
	if !f.OrganizationID.IsZero() {
		query["organization_id"] = f.OrganizationID
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

	// Full-text search
	if f.SearchText != "" {
		query["$text"] = bson.M{"$search": f.SearchText}
		options.SetProjection(bson.M{"score": bson.M{"$meta": "textScore"}})
		options.SetSort(bson.D{{"score", bson.M{"$meta": "textScore"}}})
	}

	options.SetSort(bson.D{{sortField, 1}}) // Sort in ascending order based on the specified field

	// Retrieve the list of items from the collection
	cursor, err := collection.Find(ctx, query, options)
	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(ctx)

	var results = []*TagAsSelectOption{}
	if err = cursor.All(ctx, &results); err != nil {
		panic(err)
	}

	return results, nil
}
