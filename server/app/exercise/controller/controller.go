package controller

import (
	"context"
	"log/slog"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	s3_storage "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/adapter/storage/s3"
	a_s "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/attachment/datastore"
	equip_s "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/equipment/datastore"
	e_s "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/exercise/datastore"
	o_s "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/offer/datastore"
	u_s "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/user/datastore"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/config"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/provider/kmutex"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/provider/uuid"
)

// ExerciseController Interface for exercise business logic controller.
type ExerciseController interface {
	Create(ctx context.Context, req *ExerciseCreateRequestIDO) (*e_s.Exercise, error)
	GetByID(ctx context.Context, id primitive.ObjectID) (*e_s.Exercise, error)
	UpdateByID(ctx context.Context, req *ExerciseUpdateRequestIDO) (*e_s.Exercise, error)
	ListByFilter(ctx context.Context, f *e_s.ExerciseListFilter) (*e_s.ExerciseListResult, error)
	DeleteByID(ctx context.Context, id primitive.ObjectID) error
}

type ExerciseControllerImpl struct {
	Config           *config.Conf
	Logger           *slog.Logger
	UUID             uuid.Provider
	S3               s3_storage.S3Storager
	Kmutex           kmutex.Provider
	AttachmentStorer a_s.AttachmentStorer
	DbClient         *mongo.Client
	EquipmentStorer  equip_s.EquipmentStorer
	ExerciseStorer   e_s.ExerciseStorer
	OfferStorer      o_s.OfferStorer
	UserStorer       u_s.UserStorer
}

func NewController(
	appCfg *config.Conf,
	loggerp *slog.Logger,
	uuidp uuid.Provider,
	kmutexp kmutex.Provider,
	s3 s3_storage.S3Storager,
	client *mongo.Client,
	a_storer a_s.AttachmentStorer,
	eq_storer equip_s.EquipmentStorer,
	e_storer e_s.ExerciseStorer,
	o_storer o_s.OfferStorer,
	u_storer u_s.UserStorer,
) ExerciseController {
	loggerp.Debug("exercise controller initialization started...")
	s := &ExerciseControllerImpl{
		Config:           appCfg,
		Logger:           loggerp,
		UUID:             uuidp,
		Kmutex:           kmutexp,
		S3:               s3,
		DbClient:         client,
		AttachmentStorer: a_storer,
		EquipmentStorer:  eq_storer,
		ExerciseStorer:   e_storer,
		OfferStorer:      o_storer,
		UserStorer:       u_storer,
	}
	s.Logger.Debug("exercise controller initialized")
	return s
}
