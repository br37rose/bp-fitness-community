package datastore

import (
	"context"
	"log/slog"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (impl ExerciseStorerImpl) ListByFilter(ctx context.Context, f *ExerciseListFilter) (*ExerciseListResult, error) {
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
	if f.ExcludeArchived {
		filter["status"] = bson.M{"$ne": ExerciseStatusArchived} // Do not list archived items! This code
	}
	if f.Gender != "" {
		filter["gender"] = f.Gender
	}
	if f.MovementType != 0 {
		filter["movement_type"] = f.MovementType
	}
	if f.Category != 0 {
		filter["category"] = f.Category
	}
	if f.Status != 0 {
		filter["status"] = f.Status
	}
	if f.VideoType != 0 {
		filter["video_type"] = f.VideoType
	}
	if !f.OfferID.IsZero() {
		filter["offer_id"] = f.OfferID
	}
	if len(f.Names) > 0 {
		// Use the $in operator to filter documents where `Names` contains any of the provided names
		filter["name"] = bson.M{"$in": f.Names}
	}
	if len(f.InTagIDs) > 0 {
		filter["tags.id"] = bson.M{"$in": f.InTagIDs}
	}
	if len(f.AllTagIDs) > 0 {
		filter["tags.id"] = bson.M{"$all": f.AllTagIDs}
	}

	impl.Logger.Debug("listing",
		slog.Any("names", f.Names),
		slog.Any("in_tags", f.InTagIDs),
		slog.Any("all_tags", f.AllTagIDs),
		slog.Int("category", int(f.Category)),
		slog.Int("video_type", int(f.VideoType)),
		slog.Any("offer_id", f.OfferID),
		slog.Int("status", int(f.Status)))

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
	results := []*Exercise{}
	hasNextPage := false
	for cursor.Next(ctx) {
		document := &Exercise{}
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

	return &ExerciseListResult{
		Results:     results,
		NextCursor:  nextCursor,
		HasNextPage: hasNextPage,
	}, nil
}
