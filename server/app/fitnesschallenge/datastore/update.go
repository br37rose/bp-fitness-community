package datastore

import (
	"context"
	"log/slog"

	"go.mongodb.org/mongo-driver/bson"
)

func (impl FitnessChallengeStorerImpl) UpdateByID(ctx context.Context, tp *FitnessChallenge) error {
	filter := bson.M{"_id": tp.ID}
	update := bson.M{"$set": tp}
	_, err := impl.Collection.UpdateOne(ctx, filter, update)
	if err != nil {
		impl.Logger.Error("database update error", slog.Any("error", err))
		return err
	}
	return nil
}
