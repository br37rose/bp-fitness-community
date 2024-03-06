package controller

import (
	"context"
	"log/slog"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	s3_storage "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/adapter/storage/s3"
	a_s "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/attachment/datastore"
	o_s "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/offer/datastore"
	u_s "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/user/datastore"
	vcat_s "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/videocategory/datastore"
	vcol_s "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/videocollection/datastore"
	vcon_s "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/videocontent/datastore"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/config"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/provider/kmutex"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/provider/uuid"
)

// VideoContentController Interface for videocontent business logic controller.
type VideoContentController interface {
	Create(ctx context.Context, req *VideoContentCreateRequestIDO) (*vcon_s.VideoContent, error)
	GetByID(ctx context.Context, id primitive.ObjectID) (*vcon_s.VideoContent, error)
	UpdateByID(ctx context.Context, req *VideoContentUpdateRequestIDO) (*vcon_s.VideoContent, error)
	ListByFilter(ctx context.Context, f *vcon_s.VideoContentListFilter) (*vcon_s.VideoContentListResult, error)
	DeleteByID(ctx context.Context, id primitive.ObjectID) error
}

type VideoContentControllerImpl struct {
	Config                *config.Conf
	Logger                *slog.Logger
	UUID                  uuid.Provider
	S3                    s3_storage.S3Storager
	Kmutex                kmutex.Provider
	DbClient              *mongo.Client
	AttachmentStorer      a_s.AttachmentStorer
	VideoCategoryStorer   vcat_s.VideoCategoryStorer
	VideoCollectionStorer vcol_s.VideoCollectionStorer
	VideoContentStorer    vcon_s.VideoContentStorer
	OfferStorer           o_s.OfferStorer
	UserStorer            u_s.UserStorer
}

func NewController(
	appCfg *config.Conf,
	loggerp *slog.Logger,
	uuidp uuid.Provider,
	kmutexp kmutex.Provider,
	s3 s3_storage.S3Storager,
	a_storer a_s.AttachmentStorer,
	client *mongo.Client,
	cat_storer vcat_s.VideoCategoryStorer,
	col_storer vcol_s.VideoCollectionStorer,
	con_storer vcon_s.VideoContentStorer,
	o_storer o_s.OfferStorer,
	u_storer u_s.UserStorer,
) VideoContentController {
	loggerp.Debug("videocontent controller initialization started...")
	s := &VideoContentControllerImpl{
		Config:                appCfg,
		Logger:                loggerp,
		UUID:                  uuidp,
		Kmutex:                kmutexp,
		S3:                    s3,
		DbClient:              client,
		AttachmentStorer:      a_storer,
		VideoCategoryStorer:   cat_storer,
		VideoCollectionStorer: col_storer,
		VideoContentStorer:    con_storer,
		OfferStorer:           o_storer,
		UserStorer:            u_storer,
	}
	s.Logger.Debug("videocontent controller initialized")
	return s
}
