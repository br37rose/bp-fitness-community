package controller

import (
	"context"
	"log/slog"
	"time"

	w_d "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/workout/datastore"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/config/constants"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type WorkoutCreateRequestIDO struct {
	Name                      string                 `json:"name"`
	Type                      int8                   `json:"type"`
	Description               string                 `json:"description"`
	Status                    int8                   `json:"status"`
	WorkoutExercises          []*w_d.WorkoutExercise `json:"workout_exercises,omitempty"`
	WorkoutExerciseTimeInMins int64                  `json:"workout_exercise_time_in_mins"`
	Visibility                int8                   `json:"visibility"`
}

func (c *WorkoutControllerImpl) Create(ctx context.Context, req *WorkoutCreateRequestIDO) (*w_d.Workout, error) {
	userID := ctx.Value(constants.SessionUserID).(primitive.ObjectID)
	userName := ctx.Value(constants.SessionUserName).(string)

	session, err := c.DbClient.StartSession()
	if err != nil {
		c.Logger.Error("start session error",
			slog.Any("error", err))
		return nil, err
	}
	defer session.EndSession(ctx)

	// Define a transaction function with a series of operations
	transactionFunc := func(sessCtx mongo.SessionContext) (interface{}, error) {
		for i, v := range req.WorkoutExercises {
			if v.ID.IsZero() {
				v.ID = primitive.NewObjectID()
			}
			v.CreatedAt = time.Now().UTC()
			v.ModifiedAt = time.Now().UTC()
			v.OrderNumber = int64(i + 1)
		}

		// Create our record in the database.
		res := &w_d.Workout{
			ID:                        primitive.NewObjectID(),
			Name:                      req.Name,
			Description:               req.Description,
			Type:                      req.Type,
			Status:                    w_d.WorkoutStatusActive,
			WorkoutExercises:          req.WorkoutExercises,
			WorkoutExerciseTimeInMins: req.WorkoutExerciseTimeInMins,
			CreatedAt:                 time.Now(),
			CreatedByUserID:           userID,
			CreatedByUserName:         userName,
			ModifiedAt:                time.Now(),
			ModifiedByUserName:        userName,
			ModifiedByUserID:          userID,
			Visibility:                req.Visibility,
		}
		if req.Visibility == w_d.WorkoutPersonalVisible {
			res.UserId = userID
			res.UserName = userName
		}

		err := c.WorkoutStorer.Create(sessCtx, res)
		if err != nil {
			c.Logger.Error("fitnessplan create error", slog.Any("error", err))
			return nil, err
		}

		return res, nil
	}

	// Start a transaction
	result, err := session.WithTransaction(ctx, transactionFunc)
	if err != nil {
		c.Logger.Error("session failed error",
			slog.Any("error", err))
		return nil, err
	}

	return result.(*w_d.Workout), nil
}
