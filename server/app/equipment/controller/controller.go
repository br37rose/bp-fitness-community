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
	domain "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/equipment/datastore"
	equipment_s "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/equipment/datastore"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/config"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/provider/uuid"
)

// EquipmentController Interface for equipment business logic controller.
type EquipmentController interface {
	Create(ctx context.Context, req *EquipmentCreateRequestIDO) (*domain.Equipment, error)
	GetByID(ctx context.Context, id primitive.ObjectID) (*domain.Equipment, error)
	UpdateByID(ctx context.Context, ns *EquipmentUpdateRequestIDO) (*domain.Equipment, error)
	ListByFilter(ctx context.Context, f *domain.EquipmentListFilter) (*domain.EquipmentListResult, error)
	ListAsSelectOptionByFilter(ctx context.Context, f *domain.EquipmentListFilter) ([]*domain.EquipmentAsSelectOption, error)
	DeleteByID(ctx context.Context, id primitive.ObjectID) error
	PermanentlyDeleteByID(ctx context.Context, id primitive.ObjectID) error
}

type EquipmentControllerImpl struct {
	Config              *config.Conf
	Logger              *slog.Logger
	UUID                uuid.Provider
	S3                  s3_storage.S3Storager
	Emailer             mg.Emailer
	DbClient            *mongo.Client
	EquipmentStorer equipment_s.EquipmentStorer
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
	org_storer equipment_s.EquipmentStorer,
	ex_storer exercise_s.ExerciseStorer,
	usr_storer user_s.UserStorer,
) EquipmentController {
	s := &EquipmentControllerImpl{
		Config:              appCfg,
		Logger:              loggerp,
		UUID:                uuidp,
		S3:                  s3,
		DbClient:            client,
		Emailer:             emailer,
		EquipmentStorer: org_storer,
		UserStorer:          usr_storer,
		ExerciseStorer:      ex_storer,
	}
	s.Logger.Debug("equipment controller initialization started...")
	s.Logger.Debug("equipment controller initialized")
	return s
}
