package datastore

import (
	"context"
	"log/slog"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (impl WorkouStorerImpl) ListByFilter(ctx context.Context, f *WorkoutListFilter) (*WorkoutistResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 12*time.Second)
	defer cancel()

	// Create the filter based on the cursor
	filter := bson.M{}
	if !f.Cursor.IsZero() {
		filter["_id"] = bson.M{"$gt": f.Cursor} // Add the cursor condition to the filter
	}

	if len(f.StatusList) > 0 {
		filter["status"] = bson.M{"$in": f.StatusList}
	}
	if f.ExcludeArchived {
		filter["status"] = bson.M{"$ne": WorkoutStatusArchived} // Do not list archived items! This code
	}
	if f.Visibility != 0 {
		if f.Visibility == 2 && !f.UserId.IsZero() {
			filter["$or"] = []bson.M{
				{"user_id": f.UserId},
				{"visibility": bson.M{"$eq": WorkoutVisibileToAll}},
			}
		}
	}

	if !f.CreatedByUserID.IsZero() {
		filter["created_by_user_id"] = f.CreatedByUserID
	}

	if len(f.Types) > 0 {
		filter["type"] = bson.M{"$in": f.Types}
	}
	if !f.UserId.IsZero() && f.Visibility != 2 {
		filter["user_id"] = bson.M{"$eq": f.UserId}
	}

	impl.Logger.Debug("listing",
		slog.Any("statusList", f.StatusList),
		slog.Any("visibility", f.Visibility),
		slog.Any("created_by_user_id", f.CreatedByUserID),
		slog.String("SearchText", f.SearchText),
	)
	// Include additional filters for our cursor-based pagination pertaining to sorting and limit.
	options := options.Find().
		SetSort(bson.M{f.SortField: f.SortOrder}).
		SetLimit(f.PageSize)

	// Include Full-text search
	if f.SearchText != "" {
		filter["$text"] = bson.M{"$search": f.SearchText}
		options.SetProjection(bson.M{"score": bson.M{"$meta": "textScore"}})
		options.SetSort(bson.D{{Key: "score", Value: bson.M{"$meta": "textScore"}}})
	}

	// Execute the query
	cursor, err := impl.Collection.Find(ctx, filter, options)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	// Retrieve the documents and check if there is a next page
	results := []*Workout{}
	hasNextPage := false
	for cursor.Next(ctx) {
		document := &Workout{}
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
		results = results[:len(results)]
		nextCursor = results[len(results)-1].ID
	}

	return &WorkoutistResult{
		Results:     results,
		NextCursor:  nextCursor,
		HasNextPage: hasNextPage,
	}, nil
}
