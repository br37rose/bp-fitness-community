package datastore

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type GoogleFitAppDevice struct {
	ID                 primitive.ObjectID `bson:"_id" json:"id,omitempty"`
	UserName           string             `bson:"user_name" json:"user_name"`
	UserLexicalName    string             `bson:"user_lexical_name" json:"user_lexical_name"`
	UserID             primitive.ObjectID `bson:"user_id" json:"user_id"`
	Status             int8               `bson:"status" json:"status"`
	HeartRateMetricID  primitive.ObjectID `bson:"heart_rate_metric_id" json:"heart_rate_metric_id,omitempty"`
	StepsCountMetricID primitive.ObjectID `bson:"steps_count_metric_id" json:"steps_count_metric_id,omitempty"`
}

func (impl GoogleFitAppStorerImpl) ListDevicesByStatus(ctx context.Context, status int8) ([]*GoogleFitAppDevice, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 12*time.Second)
	defer cancel()

	// Create the filter based on the cursor
	filter := bson.M{}

	if status != 0 {
		filter["status"] = status
	}

	results := make([]*GoogleFitAppDevice, 0)

	// Find documents that match the filter
	cursor, err := impl.Collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	// Iterate over the cursor and extract the "_id" field values
	for cursor.Next(ctx) {
		document := &GoogleFitAppDevice{}
		if err := cursor.Decode(document); err != nil {
			return nil, err
		}
		results = append(results, document)
	}

	return results, nil
}
