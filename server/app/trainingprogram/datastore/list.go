package datastore

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (impl TrainingProgramStorerImpl) ListByFilter(ctx context.Context, f *TrainingProgramListFilter) (*TrainingProgramListResult, error) {
	filter := bson.M{}

	if !f.Cursor.IsZero() {
		filter["_id"] = bson.M{"$gt": f.Cursor}
	}

	if !f.UserID.IsZero() {
		filter["user_id"] = f.UserID
	}

	if !f.OrganizationID.IsZero() {
		filter["organization_id"] = f.OrganizationID
	}

	if len(f.StatusList) > 0 {
		filter["status"] = bson.M{"$in": f.StatusList}
	}

	if f.DurationInWeeks != 0 {
		filter["duration_in_weeks"] = f.DurationInWeeks
	}

	if f.Phases != 0 {
		filter["phases"] = f.Phases
	}

	if f.Weeks != 0 {
		filter["weeks"] = f.Weeks
	}

	if !f.StartTime.IsZero() {
		filter["start_time"] = bson.M{"$gte": f.StartTime}
	}

	if !f.EndTime.IsZero() {
		filter["end_time"] = bson.M{"$lte": f.EndTime}
	}

	options := options.Find().
		SetSort(bson.M{f.SortField: f.SortOrder}).
		SetLimit(f.PageSize)

	if f.SearchText != "" {
		filter["$text"] = bson.M{"$search": f.SearchText}
		options.SetProjection(bson.M{"score": bson.M{"$meta": "textScore"}})
		options.SetSort(bson.D{{Key: "score", Value: bson.M{"$meta": "textScore"}}})
	}

	cursor, err := impl.Collection.Find(ctx, filter, options)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	results := []*TrainingProgram{}
	hasNextPage := false
	for cursor.Next(ctx) {
		var tp TrainingProgram
		if err := cursor.Decode(&tp); err != nil {
			return nil, err
		}
		results = append(results, &tp)
		if int64(len(results)) >= f.PageSize {
			hasNextPage = true
			break
		}
	}

	nextCursor := primitive.NilObjectID
	if len(results) == int(f.PageSize) {
		nextCursor = results[len(results)-1].ID
	}

	return &TrainingProgramListResult{
		Results:     results,
		NextCursor:  nextCursor,
		HasNextPage: hasNextPage,
	}, nil
}
