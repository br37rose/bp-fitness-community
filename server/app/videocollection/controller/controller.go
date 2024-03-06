package controller

import (
	"context"
	"log/slog"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	s3_storage "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/adapter/storage/s3"
	a_s "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/attachment/datastore"
	vcat_d "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/videocategory/datastore"
	vcol_d "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/videocollection/datastore"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/config"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/provider/kmutex"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/provider/uuid"
)

// VideoCollectionController Interface for videocollection business logic controller.
type VideoCollectionController interface {
	Create(ctx context.Context, req *VideoCollectionCreateRequestIDO) (*vcol_d.VideoCollection, error)
	GetByID(ctx context.Context, id primitive.ObjectID) (*vcol_d.VideoCollection, error)
	UpdateByID(ctx context.Context, req *VideoCollectionUpdateRequestIDO) (*vcol_d.VideoCollection, error)
	ListByFilter(ctx context.Context, f *vcol_d.VideoCollectionListFilter) (*vcol_d.VideoCollectionListResult, error)
	ListAsSelectOptionByFilter(ctx context.Context, f *vcol_d.VideoCollectionListFilter) ([]*vcol_d.VideoCollectionAsSelectOption, error)
	DeleteByID(ctx context.Context, id primitive.ObjectID) error
}

type VideoCollectionControllerImpl struct {
	Config                *config.Conf
	Logger                *slog.Logger
	UUID                  uuid.Provider
	S3                    s3_storage.S3Storager
	Kmutex                kmutex.Provider
	DbClient              *mongo.Client
	AttachmentStorer      a_s.AttachmentStorer
	VideoCategoryStorer   vcat_d.VideoCategoryStorer
	VideoCollectionStorer vcol_d.VideoCollectionStorer
}

func NewController(
	appCfg *config.Conf,
	loggerp *slog.Logger,
	uuidp uuid.Provider,
	kmutexp kmutex.Provider,
	s3 s3_storage.S3Storager,
	client *mongo.Client,
	a_storer a_s.AttachmentStorer,
	vcat_storer vcat_d.VideoCategoryStorer,
	vcol_storer vcol_d.VideoCollectionStorer,
) VideoCollectionController {
	s := &VideoCollectionControllerImpl{
		Config:                appCfg,
		Logger:                loggerp,
		UUID:                  uuidp,
		Kmutex:                kmutexp,
		S3:                    s3,
		DbClient:              client,
		AttachmentStorer:      a_storer,
		VideoCategoryStorer:   vcat_storer,
		VideoCollectionStorer: vcol_storer,
	}
	s.Logger.Debug("videocollection controller initialization started...")
	s.Logger.Debug("videocollection controller initialized")
	return s
}
