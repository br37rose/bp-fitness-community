package controller

import (
	"context"
	"log/slog"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	mg "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/adapter/emailer/mailgun"
	openai "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/adapter/openai"
	s3_storage "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/adapter/storage/s3"
	exercise_d "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/exercise/controller"
	exercise_s "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/exercise/datastore"
	qstn_c "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/question/controller"

	domain "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/fitnessplan/datastore"
	fitnessplan_s "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/fitnessplan/datastore"
	user_s "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/user/datastore"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/config"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/provider/kmutex"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/provider/uuid"
)

// FitnessPlanController Interface for fitnessplan business logic controller.
type FitnessPlanController interface {
	Create(ctx context.Context, req *FitnessPlanCreateRequestIDO) (*domain.FitnessPlan, error)
	GetByID(ctx context.Context, id primitive.ObjectID) (*domain.FitnessPlan, error)
	UpdateByID(ctx context.Context, ns *FitnessPlanUpdateRequestIDO) (*domain.FitnessPlan, error)
	ListByFilter(ctx context.Context, f *domain.FitnessPlanListFilter) (*domain.FitnessPlanListResult, error)
	ListAsSelectOptionByFilter(ctx context.Context, f *domain.FitnessPlanListFilter) ([]*domain.FitnessPlanAsSelectOption, error)
	DeleteByID(ctx context.Context, id primitive.ObjectID) error
	PermanentlyDeleteByID(ctx context.Context, id primitive.ObjectID) error
}

type FitnessPlanControllerImpl struct {
	Config             *config.Conf
	Logger             *slog.Logger
	UUID               uuid.Provider
	S3                 s3_storage.S3Storager
	Emailer            mg.Emailer
	DbClient           *mongo.Client
	Kmutex             kmutex.Provider
	OpenAI             openai.OpenAIConnector
	FitnessPlanStorer  fitnessplan_s.FitnessPlanStorer
	UserStorer         user_s.UserStorer
	ExerciseStorer     exercise_s.ExerciseStorer
	ExcerciseContr     exercise_d.ExerciseController
	QuestionController qstn_c.QuestionController
}

func NewController(
	appCfg *config.Conf,
	loggerp *slog.Logger,
	uuidp uuid.Provider,
	s3 s3_storage.S3Storager,
	emailer mg.Emailer,
	client *mongo.Client,
	kmutexp kmutex.Provider,
	ai openai.OpenAIConnector,
	org_storer fitnessplan_s.FitnessPlanStorer,
	ex_storer exercise_s.ExerciseStorer,
	usr_storer user_s.UserStorer,
	ex_contr exercise_d.ExerciseController,
	qstn_contr qstn_c.QuestionController,
) FitnessPlanController {
	s := &FitnessPlanControllerImpl{
		Config:             appCfg,
		Logger:             loggerp,
		UUID:               uuidp,
		S3:                 s3,
		Emailer:            emailer,
		DbClient:           client,
		Kmutex:             kmutexp,
		OpenAI:             ai,
		FitnessPlanStorer:  org_storer,
		UserStorer:         usr_storer,
		ExerciseStorer:     ex_storer,
		ExcerciseContr:     ex_contr,
		QuestionController: qstn_contr,
	}
	s.Logger.Debug("fitnessplan controller initialization started...")
	s.Logger.Debug("fitnessplan controller initialized")
	return s
}
