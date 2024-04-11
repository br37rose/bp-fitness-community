package controller

import (
	"context"
	"log/slog"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	ud "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/user/datastore"
	domain "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/workout/datastore"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/config/constants"
	httperror "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/utils/httperror"
)

type WorkoutUpdateRequest struct {
	ID                        primitive.ObjectID        `json:"id,omitempty" bson:"_id,omitempty"`
	Name                      string                    `json:"name" bson:"name"`
	Description               string                    `json:"description"`
	Type                      int8                      `json:"type" bson:"type"`
	Status                    int8                      `json:"status" bson:"status"`
	Visibility                int8                      `json:"visibility"`
	WorkoutExercises          []*domain.WorkoutExercise `json:"workout_exercises,omitempty" bson:"workout_exercises,omitempty"`
	WorkoutExerciseTimeInMins int64                     `json:"workout_exercise_time_in_mins" bson:"workout_exercise_time_in_mins"`
}

func (c *WorkoutControllerImpl) UpdateByID(ctx context.Context, req *WorkoutUpdateRequest) (*domain.Workout, error) {
	userID := ctx.Value(constants.SessionUserID).(primitive.ObjectID)
	userName := ctx.Value(constants.SessionUserName).(string)
	urole, ok := ctx.Value(constants.SessionUserRole).(int8)
	if ok && urole == ud.UserRoleMember {
		req.Visibility = domain.WorkoutPersonalVisible
	}

	session, err := c.DbClient.StartSession()
	if err != nil {
		c.Logger.Error("start session error", slog.Any("error", err))
		return nil, err
	}
	defer session.EndSession(ctx)

	transactionFunc := func(sessCtx mongo.SessionContext) (interface{}, error) {
		os, err := c.WorkoutStorer.GetByID(sessCtx, req.ID)
		if err != nil {
			c.Logger.Error("database get by id error", slog.Any("error", err), slog.Any("workout_id", req.ID))
			return nil, err
		}
		if os == nil {
			c.Logger.Error("workout does not exist error", slog.Any("workout_id", req.ID))
			return nil, httperror.NewForBadRequestWithSingleField("message", "workout does not exist")
		}

		for i, v := range req.WorkoutExercises {
			if v.ID.IsZero() {
				v.ID = primitive.NewObjectID()
			}
			if v.CreatedAt.IsZero() {
				v.CreatedAt = time.Now().UTC()
			}
			v.ModifiedAt = time.Now().UTC()
			v.OrderNumber = int64(i + 1)
		}

		os.Name = req.Name
		os.Description = req.Description
		os.Type = req.Type
		os.Status = domain.WorkoutStatusActive
		os.Visibility = req.Visibility
		os.WorkoutExercises = req.WorkoutExercises
		os.WorkoutExerciseTimeInMins = req.WorkoutExerciseTimeInMins
		os.ModifiedAt = time.Now()
		os.ModifiedByUserID = userID
		os.ModifiedByUserName = userName

		if err := c.WorkoutStorer.UpdateByID(sessCtx, os); err != nil {
			c.Logger.Error("database update by id error", slog.Any("error", err))
			return nil, err
		}
		return os, nil
	}

	result, err := session.WithTransaction(ctx, transactionFunc)
	if err != nil {
		c.Logger.Error("session failed error", slog.Any("error", err))
		return nil, err
	}

	return result.(*domain.Workout), nil
}
