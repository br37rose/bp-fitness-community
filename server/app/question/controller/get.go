package controller

import (
	"context"
	"log/slog"

	q_s "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/question/datastore"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (c *QuestionControllerImpl) GetByID(ctx context.Context, id primitive.ObjectID) (*q_s.Question, error) {
	m, err := c.QuestionStorer.GetByID(ctx, id)
	if err != nil {
		c.Logger.Error("database get by id error", slog.Any("error", err))
		return nil, err
	}
	return m, err
}
