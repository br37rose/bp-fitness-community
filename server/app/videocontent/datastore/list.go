package datastore

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (impl VideoContentStorerImpl) ListByFilter(ctx context.Context, f *VideoContentListFilter) (*VideoContentListResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 12*time.Second)
	defer cancel()

	// Create the filter based on the cursor
	filter := bson.M{}
	if !f.Cursor.IsZero() {
		filter["_id"] = bson.M{"$gt": f.Cursor} // Add the cursor condition to the filter
	}

	// Apply the conditions to the filter
	if !f.OrganizationID.IsZero() {
		filter["organization_id"] = f.OrganizationID
	}
	if f.Gender != "" {
		filter["gender"] = f.Gender
	}
	if f.MovementType != 0 {
		filter["movement_type"] = f.MovementType
	}
	if !f.CategoryID.IsZero() {
		filter["category_id"] = f.CategoryID
	}
	if f.Status != 0 {
		filter["status"] = f.Status
	}
	if f.VideoType != 0 {
		filter["video_type"] = f.VideoType
	}
	if !f.CollectionID.IsZero() {
		filter["collection_id"] = f.CollectionID
	}
	if !f.OfferID.IsZero() {
		filter["offer_id"] = f.OfferID
	}

	// impl.Logger.Debug("listing",
	// 	slog.Any("collection_id", f.CollectionID),
	// 	slog.Any("category_id", f.CategoryID),
	// 	slog.Int("video_type", int(f.VideoType)),
	// 	slog.Int("status", int(f.Status)))

	// Include additional filters for our cursor-based pagination pertaining to sorting and limit.
	options := options.Find().
		SetSort(bson.M{f.SortField: f.SortOrder}).
		SetLimit(f.PageSize)

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

	// var results = []*ComicSubmission{}
	// if err = cursor.All(ctx, &results); err != nil {
	// 	panic(err)
	// }

	// Retrieve the documents and check if there is a next page
	results := []*VideoContent{}
	hasNextPage := false
	for cursor.Next(ctx) {
		document := &VideoContent{}
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
	nextCursor := primitive.NilObjectID
	if int64(len(results)) == f.PageSize {
		// Remove the extra document from the current page
		results = results[:len(results)]

		// Get the last document's _id as the next cursor
		nextCursor = results[len(results)-1].ID
	}

	return &VideoContentListResult{
		Results:     results,
		NextCursor:  nextCursor,
		HasNextPage: hasNextPage,
	}, nil
}

func (impl VideoContentStorerImpl) ListAsSelectOptionByFilter(ctx context.Context, f *VideoContentListFilter) ([]*VideoContentAsSelectOption, error) {
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
	if f.Gender != "" {
		query["gender"] = f.Gender
	}
	if f.MovementType != 0 {
		query["movement_type"] = f.MovementType
	}
	if !f.CategoryID.IsZero() {
		query["category_id"] = f.CategoryID
	}
	if f.Status != 0 {
		query["status"] = f.Status
	}
	if f.VideoType != 0 {
		query["video_type"] = f.VideoType
	}
	if !f.CollectionID.IsZero() {
		query["collection_id"] = f.CollectionID
	}
	if !f.OfferID.IsZero() {
		query["offer_id"] = f.OfferID
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

	var results = []*VideoContentAsSelectOption{}
	if err = cursor.All(ctx, &results); err != nil {
		panic(err)
	}

	return results, nil
}
