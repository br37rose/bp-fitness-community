package controller

import (
	"context"
	"log/slog"

	"go.mongodb.org/mongo-driver/bson/primitive"

	s3_storage "github.com/bci-innovation-labs/bp8fitnesscommunity-cli/adapter/storage/s3"
	equip_s "github.com/bci-innovation-labs/bp8fitnesscommunity-cli/app/equipement/datastore"
	o_s "github.com/bci-innovation-labs/bp8fitnesscommunity-cli/app/exercise/datastore"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-cli/config"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-cli/provider/uuid"
)

// ExerciseController Interface for exercise business logic controller.
type ExerciseController interface {
	Create(ctx context.Context, m *o_s.Exercise) (*o_s.Exercise, error)
	GetByID(ctx context.Context, id primitive.ObjectID) (*o_s.Exercise, error)
	UpdateByID(ctx context.Context, m *o_s.Exercise) (*o_s.Exercise, error)
	ListByFilter(ctx context.Context, f *o_s.ExerciseListFilter) (*o_s.ExerciseListResult, error)
	DeleteByID(ctx context.Context, id primitive.ObjectID) error
}

type ExerciseControllerImpl struct {
	Config          *config.Conf
	Logger          *slog.Logger
	UUID            uuid.Provider
	S3              s3_storage.S3Storager
	EquipmentStorer equip_s.EquipementStorer
	ExerciseStorer  o_s.ExerciseStorer
}

func NewController(
	appCfg *config.Conf,
	loggerp *slog.Logger,
	uuidp uuid.Provider,
	s3 s3_storage.S3Storager,
	eq_storer equip_s.EquipementStorer,
	org_storer o_s.ExerciseStorer,
) ExerciseController {
	s := &ExerciseControllerImpl{
		Config:          appCfg,
		Logger:          loggerp,
		UUID:            uuidp,
		S3:              s3,
		EquipmentStorer: eq_storer,
		ExerciseStorer:  org_storer,
	}
	s.Logger.Debug("exercise controller initialization started...")
	s.Logger.Debug("exercise controller initialized")
	return s
}
