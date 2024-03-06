package datastore

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (impl ExerciseStorerImpl) ListAsSelectOptionByFilter(ctx context.Context, f *ExerciseListFilter) ([]*ExerciseAsSelectOption, error) {
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
	if f.Category != 0 {
		query["category"] = f.Category
	}
	if f.Status != 0 {
		query["status"] = f.Status
	}
	if f.VideoType != 0 {
		query["video_type"] = f.VideoType
	}
	if len(f.Names) > 0 {
		// Use the $in operator to filter documents where `Names` contains any of the provided names
		query["name"] = bson.M{"$in": f.Names}
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
		query["status"] = bson.M{"$ne": ExerciseStatusArchived} // Do not list archived items! This code
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

	var results = []*ExerciseAsSelectOption{}
	if err = cursor.All(ctx, &results); err != nil {
		panic(err)
	}

	return results, nil
}

func (impl ExerciseStorerImpl) ListAllAsSelectOption(ctx context.Context, oid primitive.ObjectID) ([]*ExerciseAsSelectOption, error) {
	f := &ExerciseListFilter{
		Cursor:         primitive.NilObjectID,
		PageSize:       1_000_000,
		SortField:      "name",
		SortOrder:      1, // 1=ascending | -1=descending
		Status:         ExerciseStatusActive,
		OrganizationID: oid,
	}
	return impl.ListAsSelectOptionByFilter(ctx, f)
}
