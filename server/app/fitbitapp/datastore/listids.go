package datastore

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (impl FitBitAppStorerImpl) ListIDsByStatus(ctx context.Context, status int8) ([]primitive.ObjectID, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 12*time.Second)
	defer cancel()

	// Create the filter based on the cursor
	filter := bson.M{}

	if status != 0 {
		filter["status"] = status
	}

	ids := make([]primitive.ObjectID, 0)

	// Find documents that match the filter
	cursor, err := impl.Collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	// Iterate over the cursor and extract the "_id" field values
	for cursor.Next(ctx) {
		var result bson.M
		if err := cursor.Decode(&result); err != nil {
			return nil, err
		}

		id, ok := result["_id"].(primitive.ObjectID)
		if !ok {
			return nil, errors.New("failed to convert _id to ObjectID")
		}

		ids = append(ids, id)
	}

	return ids, nil
}

func (impl FitBitAppStorerImpl) ListSimulatorIDsByStatus(ctx context.Context, status int8) ([]primitive.ObjectID, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 12*time.Second)
	defer cancel()

	// Create the filter based on the cursor
	filter := bson.M{}

	if status != 0 {
		filter["status"] = status
	}
	filter["is_test_mode"] = true

	ids := make([]primitive.ObjectID, 0)

	// Find documents that match the filter
	cursor, err := impl.Collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	// Iterate over the cursor and extract the "_id" field values
	for cursor.Next(ctx) {
		var result bson.M
		if err := cursor.Decode(&result); err != nil {
			return nil, err
		}

		id, ok := result["_id"].(primitive.ObjectID)
		if !ok {
			return nil, errors.New("failed to convert _id to ObjectID")
		}

		ids = append(ids, id)
	}

	return ids, nil
}

func (impl FitBitAppStorerImpl) ListPhysicalIDsByStatus(ctx context.Context, status int8) ([]primitive.ObjectID, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 12*time.Second)
	defer cancel()

	// Create the filter based on the cursor
	filter := bson.M{}

	if status != 0 {
		filter["status"] = status
	}
	filter["is_test_mode"] = false

	ids := make([]primitive.ObjectID, 0)

	// Find documents that match the filter
	cursor, err := impl.Collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	// Iterate over the cursor and extract the "_id" field values
	for cursor.Next(ctx) {
		var result bson.M
		if err := cursor.Decode(&result); err != nil {
			return nil, err
		}

		id, ok := result["_id"].(primitive.ObjectID)
		if !ok {
			return nil, errors.New("failed to convert _id to ObjectID")
		}

		ids = append(ids, id)
	}

	return ids, nil
}
