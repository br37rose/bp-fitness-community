package controller

import (
	"context"
	"log/slog"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	s3_storage "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/adapter/storage/s3"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/adapter/templatedemailer"
	t_s "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/tag/datastore"
	tag_s "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/tag/datastore"
	user_s "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/user/datastore"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/config"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/provider/kmutex"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/provider/password"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/provider/uuid"
)

// TagController Interface for tag business logic controller.
type TagController interface {
	Create(ctx context.Context, requestData *TagCreateRequestIDO) (*tag_s.Tag, error)
	GetByID(ctx context.Context, id primitive.ObjectID) (*tag_s.Tag, error)
	UpdateByID(ctx context.Context, nu *TagUpdateRequestIDO) (*tag_s.Tag, error)
	ListByFilter(ctx context.Context, f *t_s.TagPaginationListFilter) (*t_s.TagPaginationListResult, error)
	ListAsSelectOptionByFilter(ctx context.Context, f *tag_s.TagListFilter) ([]*tag_s.TagAsSelectOption, error)
	ArchiveByID(ctx context.Context, id primitive.ObjectID) (*tag_s.Tag, error)
	DeleteByID(ctx context.Context, id primitive.ObjectID) error
}

type TagControllerImpl struct {
	Config           *config.Conf
	Logger           *slog.Logger
	UUID             uuid.Provider
	S3               s3_storage.S3Storager
	Password         password.Provider
	Kmutex           kmutex.Provider
	DbClient         *mongo.Client
	UserStorer       user_s.UserStorer
	TagStorer        t_s.TagStorer
	TemplatedEmailer templatedemailer.TemplatedEmailer
}

func NewController(
	appCfg *config.Conf,
	loggerp *slog.Logger,
	uuidp uuid.Provider,
	s3 s3_storage.S3Storager,
	passwordp password.Provider,
	kmux kmutex.Provider,
	client *mongo.Client,
	temailer templatedemailer.TemplatedEmailer,
	usr_storer user_s.UserStorer,
	t_storer t_s.TagStorer,
) TagController {
	s := &TagControllerImpl{
		Config:           appCfg,
		Logger:           loggerp,
		UUID:             uuidp,
		S3:               s3,
		Password:         passwordp,
		Kmutex:           kmux,
		TemplatedEmailer: temailer,
		DbClient:         client,
		UserStorer:       usr_storer,
		TagStorer:        t_storer,
	}
	s.Logger.Debug("tag controller initialization started...")
	s.Logger.Debug("tag controller initialized")
	return s
}
