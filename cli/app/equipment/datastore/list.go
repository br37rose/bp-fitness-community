package datastore

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log/slog"
)

func (impl EquipmentStorerImpl) ListByFilter(ctx context.Context, f *EquipmentListFilter) (*EquipmentListResult, error) {
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
	if f.Status != 0 {
		filter["status"] = f.Status
	}

	impl.Logger.Debug("fetching video categories list",
		slog.Any("Cursor", f.Cursor),
		slog.Int64("PageSize", f.PageSize),
		slog.String("SortField", f.SortField),
		slog.Any("SortOrder", f.SortOrder),
		slog.Any("OrganizationID", f.OrganizationID),
		slog.Any("Status", f.Status),
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

	// var results = []*Equipment{}
	// if err = cursor.All(ctx, &results); err != nil {
	// 	panic(err)
	// }

	// Retrieve the documents and check if there is a next page
	results := []*Equipment{}
	hasNextPage := false
	for cursor.Next(ctx) {
		document := &Equipment{}
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

	return &EquipmentListResult{
		Results:     results,
		NextCursor:  nextCursor,
		HasNextPage: hasNextPage,
	}, nil
}
