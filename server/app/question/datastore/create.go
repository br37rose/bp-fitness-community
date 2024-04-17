package datastore

import (
	"context"
	"log/slog"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (impl *QuestionStorerImpl) Create(ctx context.Context, q *Question) error {

	if q.ID == primitive.NilObjectID {
		q.ID = primitive.NewObjectID()
		impl.Logger.Warn("database insert user not included id value, created id now.", slog.Any("id", q.ID))
	}

	_, err := impl.Collection.InsertOne(ctx, q)

	// check for errors in the insertion
	if err != nil {
		impl.Logger.Error("database insert error", slog.Any("error", err))
	}

	return nil
}
