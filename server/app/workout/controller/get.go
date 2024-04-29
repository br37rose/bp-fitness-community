package controller

import (
	"context"
	"log/slog"

	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/workout/datastore"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (c *WorkoutControllerImpl) GetByID(ctx context.Context, id primitive.ObjectID) (*datastore.Workout, error) {
	// Retrieve from our database the record for the specific id.
	m, err := c.WorkoutStorer.GetByID(ctx, id)
	if err != nil {
		c.Logger.Error("database get by id error", slog.Any("error", err))
		return nil, err
	}
	for i, e := range m.WorkoutExercises {
		if !e.IsRest && !e.ExerciseID.IsZero() {
			exc, err := c.ExcController.GetByID(ctx, e.ExerciseID)
			if err != nil {
				return nil, err
			}
			m.WorkoutExercises[i].Excercise = *exc
		}
	}
	return m, err
}
