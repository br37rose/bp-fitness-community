package datastore

import (
	"context"
	"log/slog"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (impl WorkouStorerImpl) Create(ctx context.Context, w *Workout) error {

	if w.ID == primitive.NilObjectID {
		w.ID = primitive.NewObjectID()
	}

	_, err := impl.Collection.InsertOne(ctx, w)
	if err != nil {
		impl.Logger.Error("database insert error", slog.Any("error", err))
		return err
	}

	return nil
}
