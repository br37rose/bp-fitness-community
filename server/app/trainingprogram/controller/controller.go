package controller

import (
	"context"
	"log/slog"

	"go.mongodb.org/mongo-driver/mongo"

	tp_s "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/trainingprogram/datastore"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/config"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/provider/uuid"
)

type TrainingprogramController interface {
	Create(ctx context.Context, req *TrainingProgramCreateRequestIDO) (*tp_s.TrainingProgram, error)
}

type TrainingprogramControllerImpl struct {
	Config   *config.Conf
	Logger   *slog.Logger
	UUID     uuid.Provider
	DbClient *mongo.Client
	TpStorer tp_s.TrainingProgramStorer
}

func NewController(
	appCfg *config.Conf,
	loggerp *slog.Logger,
	uuidp uuid.Provider,
	client *mongo.Client,
	tp_store tp_s.TrainingProgramStorer,
) TrainingprogramController {
	impl := &TrainingprogramControllerImpl{
		Config:   appCfg,
		Logger:   loggerp,
		UUID:     uuidp,
		DbClient: client,
		TpStorer: tp_store,
	}
	return impl
}
