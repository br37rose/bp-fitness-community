package datastore

import (
	"context"
	"log/slog"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (impl TrainingProgramStorerImpl) UpdateByID(ctx context.Context, tp *TrainingProgram) error {
	filter := bson.M{"_id": tp.ID}
	update := bson.M{"$set": tp}
	_, err := impl.Collection.UpdateOne(ctx, filter, update)
	if err != nil {
		impl.Logger.Error("database update error", slog.Any("error", err))
		return err
	}
	return nil
}

func (impl TrainingProgramStorerImpl) UpdatePhase(ctx context.Context, tpId primitive.ObjectID, phases []*TrainingPhase) error {
	filter := bson.M{
		"_id": tpId,
	}

	update := bson.M{
		"$set": bson.M{
			"training_phases": phases,
		},
	}

	// Update the specific object within the array
	_, err := impl.Collection.UpdateOne(ctx, filter, update)
	if err != nil {
		// Handle error
		return err
	}

	return nil
}
