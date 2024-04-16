package controller

import (
	"context"
	"log/slog"

	q_s "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/question/datastore"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/config"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/provider/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type QuestionController interface {
	Create(ctx context.Context, req *QuestionRequest) (*q_s.Question, error)
	GetByID(ctx context.Context, id primitive.ObjectID) (*q_s.Question, error)
	UpdateByID(ctx context.Context, req *QuestionUpdateRequest) (*q_s.Question, error)
	DeleteByID(ctx context.Context, id primitive.ObjectID) error
	ListByFilter(ctx context.Context, f *q_s.QuestionListFilter) (*q_s.QuestionListResult, error)
}

type QuestionRequest struct {
	Question      string                `json:"question"`
	IsMultiSelect bool                  `json:"isMultiSelect"`
	Content       []q_s.QuestionContent `json:"content"`
	Status        bool                  `json:"status"`
}

type QuestionControllerImpl struct {
	Config         *config.Conf
	Logger         *slog.Logger
	UUID           uuid.Provider
	DbClient       *mongo.Client
	QuestionStorer q_s.QuestionStorer
}

func NewController(
	appCfg *config.Conf,
	loggerp *slog.Logger,
	uuid uuid.Provider,
	client *mongo.Client,
	questionStorer q_s.QuestionStorer,
) QuestionController {
	return &QuestionControllerImpl{
		Config:         appCfg,
		Logger:         loggerp,
		UUID:           uuid,
		DbClient:       client,
		QuestionStorer: questionStorer,
	}
}
