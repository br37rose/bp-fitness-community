package datastore

import (
	"context"
	"log/slog"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (impl FitnessPlanStorerImpl) ListByFilter(ctx context.Context, f *FitnessPlanListFilter) (*FitnessPlanListResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 12*time.Second)
	defer cancel()

	// Create the filter based on the cursor
	filter := bson.M{}
	if !f.Cursor.IsZero() {
		filter["_id"] = bson.M{"$gt": f.Cursor} // Add the cursor condition to the filter
	}

	// Add filter conditions to the filter
	if f.OrganizationID != primitive.NilObjectID {
		filter["organization_id"] = f.OrganizationID
	}
	if f.UserID != primitive.NilObjectID {
		filter["user_id"] = f.UserID
	}
	if len(f.ExerciseNames) > 0 {
		// Use the $in operator to filter documents where ExerciseNames contains any of the provided names
		filter["exercise_names"] = bson.M{"$in": f.ExerciseNames}
	}
	if len(f.StatusList) > 0 {
		filter["status"] = bson.M{"$in": f.StatusList}
	}

	impl.Logger.Debug("fetching video categories list",
		slog.Any("Cursor", f.Cursor),
		slog.Int64("PageSize", f.PageSize),
		slog.String("SortField", f.SortField),
		slog.Any("SortOrder", f.SortOrder),
		slog.Any("OrganizationID", f.OrganizationID),
		slog.Any("Status", f.StatusList),
	)

	// Include additional filters for our cursor-based pagination pertaining to sorting and limit.
	options := options.Find().
		SetSort(bson.M{f.SortField: f.SortOrder}).
		SetLimit(f.PageSize)

	// Execute the query
	cursor, err := impl.Collection.Find(ctx, filter, options)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	// var results = []*FitnessPlan{}
	// if err = cursor.All(ctx, &results); err != nil {
	// 	panic(err)
	// }

	// Retrieve the documents and check if there is a next page
	results := []*FitnessPlan{}
	hasNextPage := false
	for cursor.Next(ctx) {
		document := &FitnessPlan{}
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

	return &FitnessPlanListResult{
		Results:     results,
		NextCursor:  nextCursor,
		HasNextPage: hasNextPage,
	}, nil
}
