package controller

import (
	"context"
	"log/slog"
	"time"

	q_s "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/question/datastore"
	u_d "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/user/datastore"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/config/constants"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/utils/httperror"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type QuestionUpdateRequest struct {
	ID primitive.ObjectID `json:"id"`
	QuestionRequest
}

func (c *QuestionControllerImpl) UpdateByID(ctx context.Context, req *QuestionUpdateRequest) (*q_s.Question, error) {

	userID := ctx.Value(constants.SessionUserID).(primitive.ObjectID)
	userName := ctx.Value(constants.SessionUserName).(string)
	urole, ok := ctx.Value(constants.SessionUserRole).(int8)
	if ok && urole != u_d.UserRoleAdmin {
		return nil, httperror.NewForBadRequestWithSingleField("message", "you role does not grant you access to this")
	}

	session, err := c.DbClient.StartSession()
	if err != nil {
		c.Logger.Error("start session error",
			slog.Any("error", err))
		return nil, err
	}
	defer session.EndSession(ctx)

	// Define a transaction function with a series of operations
	transactionFunc := func(sessCtx mongo.SessionContext) (interface{}, error) {

		os, err := c.QuestionStorer.GetByID(sessCtx, req.ID)
		if err != nil {
			c.Logger.Error("database get by id error", slog.Any("error", err), slog.Any("workout_id", req.ID))
			return nil, err
		}
		if os == nil {
			c.Logger.Error("question does not exist error", slog.Any("question_id", req.ID))
			return nil, httperror.NewForBadRequestWithSingleField("message", "workout does not exist")
		}

		// Create our record in the database.
		res := &q_s.Question{
			ID:                 req.ID,
			Title:              req.Title,
			IsMultiSelect:      req.IsMultiSelect,
			Subtitle:           req.Subtitle,
			Options:            req.Options,
			CreatedAt:          time.Now(),
			CreatedByUserID:    userID,
			CreatedByUserName:  userName,
			ModifiedAt:         time.Now(),
			ModifiedByUserName: userName,
			ModifiedByUserID:   userID,
			Status:             req.Status,
		}

		err = c.QuestionStorer.UpdateByID(sessCtx, res)
		if err != nil {
			c.Logger.Error("question create error", slog.Any("error", err))
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

	return result.(*q_s.Question), nil
}
