package controller

import (
	"context"
	"log/slog"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	mg "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/adapter/emailer/mailgun"
	s3_storage "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/adapter/storage/s3"
	attachment_s "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/attachment/datastore"
	domain "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/attachment/datastore"
	exercise_s "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/exercise/datastore"
	user_s "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/user/datastore"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/config"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/provider/uuid"
)

// AttachmentController Interface for attachment business logic controller.
type AttachmentController interface {
	Create(ctx context.Context, req *AttachmentCreateRequestIDO) (*domain.Attachment, error)
	GetByID(ctx context.Context, id primitive.ObjectID) (*domain.Attachment, error)
	UpdateByID(ctx context.Context, ns *AttachmentUpdateRequestIDO) (*domain.Attachment, error)
	ListByFilter(ctx context.Context, f *domain.AttachmentListFilter) (*domain.AttachmentListResult, error)
	ListAsSelectOptionByFilter(ctx context.Context, f *domain.AttachmentListFilter) ([]*domain.AttachmentAsSelectOption, error)
	DeleteByID(ctx context.Context, id primitive.ObjectID) error
	PermanentlyDeleteByID(ctx context.Context, id primitive.ObjectID) error
}

type AttachmentControllerImpl struct {
	Config           *config.Conf
	Logger           *slog.Logger
	UUID             uuid.Provider
	S3               s3_storage.S3Storager
	Emailer          mg.Emailer
	DbClient         *mongo.Client
	AttachmentStorer attachment_s.AttachmentStorer
	UserStorer       user_s.UserStorer
	ExerciseStorer   exercise_s.ExerciseStorer
}

func NewController(
	appCfg *config.Conf,
	loggerp *slog.Logger,
	uuidp uuid.Provider,
	s3 s3_storage.S3Storager,
	client *mongo.Client,
	emailer mg.Emailer,
	org_storer attachment_s.AttachmentStorer,
	ex_storer exercise_s.ExerciseStorer,
	usr_storer user_s.UserStorer,
) AttachmentController {
	s := &AttachmentControllerImpl{
		Config:           appCfg,
		Logger:           loggerp,
		UUID:             uuidp,
		S3:               s3,
		Emailer:          emailer,
		DbClient:         client,
		AttachmentStorer: org_storer,
		UserStorer:       usr_storer,
		ExerciseStorer:   ex_storer,
	}
	s.Logger.Debug("attachment controller initialization started...")
	s.Logger.Debug("attachment controller initialized")
	return s
}
