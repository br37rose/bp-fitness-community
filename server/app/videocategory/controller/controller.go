package controller

import (
	"context"
	"log/slog"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	mg "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/adapter/emailer/mailgun"
	s3_storage "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/adapter/storage/s3"
	exercise_s "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/exercise/datastore"
	user_s "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/user/datastore"
	domain "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/videocategory/datastore"
	videocategory_s "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/videocategory/datastore"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/config"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/provider/uuid"
)

// VideoCategoryController Interface for videocategory business logic controller.
type VideoCategoryController interface {
	Create(ctx context.Context, req *VideoCategoryCreateRequestIDO) (*domain.VideoCategory, error)
	GetByID(ctx context.Context, id primitive.ObjectID) (*domain.VideoCategory, error)
	UpdateByID(ctx context.Context, ns *VideoCategoryUpdateRequestIDO) (*domain.VideoCategory, error)
	ListByFilter(ctx context.Context, f *domain.VideoCategoryListFilter) (*domain.VideoCategoryListResult, error)
	ListAsSelectOptionByFilter(ctx context.Context, f *domain.VideoCategoryListFilter) ([]*domain.VideoCategoryAsSelectOption, error)
	DeleteByID(ctx context.Context, id primitive.ObjectID) error
	PermanentlyDeleteByID(ctx context.Context, id primitive.ObjectID) error
}

type VideoCategoryControllerImpl struct {
	Config              *config.Conf
	Logger              *slog.Logger
	UUID                uuid.Provider
	S3                  s3_storage.S3Storager
	Emailer             mg.Emailer
	DbClient            *mongo.Client
	VideoCategoryStorer videocategory_s.VideoCategoryStorer
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
	org_storer videocategory_s.VideoCategoryStorer,
	ex_storer exercise_s.ExerciseStorer,
	usr_storer user_s.UserStorer,
) VideoCategoryController {
	s := &VideoCategoryControllerImpl{
		Config:              appCfg,
		Logger:              loggerp,
		UUID:                uuidp,
		S3:                  s3,
		DbClient:            client,
		Emailer:             emailer,
		VideoCategoryStorer: org_storer,
		UserStorer:          usr_storer,
		ExerciseStorer:      ex_storer,
	}
	s.Logger.Debug("videocategory controller initialization started...")
	s.Logger.Debug("videocategory controller initialized")
	return s
}
