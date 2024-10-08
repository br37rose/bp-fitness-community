package controller

import (
	"context"
	"log/slog"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	mg "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/adapter/emailer/mailgun"
	openai "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/adapter/openai"
	s3_storage "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/adapter/storage/s3"
	exercise_s "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/exercise/datastore"
	domain "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/nutritionplan/datastore"
	nutritionplan_s "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/nutritionplan/datastore"
	user_s "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/user/datastore"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/config"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/provider/kmutex"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/provider/uuid"
)

// NutritionPlanController Interface for nutritionplan business logic controller.
type NutritionPlanController interface {
	Create(ctx context.Context, req *NutritionPlanCreateRequestIDO) (*domain.NutritionPlan, error)
	GetByID(ctx context.Context, id primitive.ObjectID) (*domain.NutritionPlan, error)
	UpdateByID(ctx context.Context, ns *NutritionPlanUpdateRequestIDO) (*domain.NutritionPlan, error)
	ListByFilter(ctx context.Context, f *domain.NutritionPlanListFilter) (*domain.NutritionPlanListResult, error)
	ListAsSelectOptionByFilter(ctx context.Context, f *domain.NutritionPlanListFilter) ([]*domain.NutritionPlanAsSelectOption, error)
	DeleteByID(ctx context.Context, id primitive.ObjectID) error
	PermanentlyDeleteByID(ctx context.Context, id primitive.ObjectID) error
}

type NutritionPlanControllerImpl struct {
	Config              *config.Conf
	Logger              *slog.Logger
	UUID                uuid.Provider
	S3                  s3_storage.S3Storager
	Emailer             mg.Emailer
	DbClient            *mongo.Client
	Kmutex              kmutex.Provider
	OpenAI              openai.OpenAIConnector
	NutritionPlanStorer nutritionplan_s.NutritionPlanStorer
	UserStorer          user_s.UserStorer
	ExerciseStorer      exercise_s.ExerciseStorer
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
	org_storer nutritionplan_s.NutritionPlanStorer,
	ex_storer exercise_s.ExerciseStorer,
	usr_storer user_s.UserStorer,
) NutritionPlanController {
	s := &NutritionPlanControllerImpl{
		Config:              appCfg,
		Logger:              loggerp,
		UUID:                uuidp,
		S3:                  s3,
		Emailer:             emailer,
		DbClient:            client,
		Kmutex:              kmutexp,
		OpenAI:              ai,
		NutritionPlanStorer: org_storer,
		UserStorer:          usr_storer,
		ExerciseStorer:      ex_storer,
	}
	s.Logger.Debug("nutritionplan controller initialization started...")
	s.Logger.Debug("nutritionplan controller initialized")
	return s
}
