package datastore

import (
	"context"
	"log/slog"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (impl TrainingProgramStorerImpl) Create(ctx context.Context, tp *TrainingProgram) error {

	if tp.ID == primitive.NilObjectID {
		tp.ID = primitive.NewObjectID()
	}

	_, err := impl.Collection.InsertOne(ctx, tp)

	// check for errors in the insertion
	if err != nil {
		impl.Logger.Error("database insert error", slog.Any("error", err))
	}

	return nil
}