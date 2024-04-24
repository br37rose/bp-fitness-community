package controller

import (
	"context"
	"log/slog"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/trainingprogram/datastore"
	tp_s "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/trainingprogram/datastore"
	u_s "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/user/datastore"
	wrk_c "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/workout/controller"
	wk_s "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/workout/datastore"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/config"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/provider/uuid"
)

type TrainingprogramController interface {
	Create(ctx context.Context, req *TrainingProgramCreateRequestIDO) (*tp_s.TrainingProgram, error)
	GetByID(ctx context.Context, id primitive.ObjectID) (*datastore.TrainingProgram, error)
	DeleteByID(ctx context.Context, id primitive.ObjectID) error
	ListByFilter(ctx context.Context, f *datastore.TrainingProgramListFilter) (*datastore.TrainingProgramListResult, error)
	UpdateTPPhase(ctx context.Context, TPid primitive.ObjectID, req PhaseUpdateRequestIDO) (*datastore.TrainingProgram, error)
}

type TrainingprogramControllerImpl struct {
	Config            *config.Conf
	Logger            *slog.Logger
	UUID              uuid.Provider
	DbClient          *mongo.Client
	TpStorer          tp_s.TrainingProgramStorer
	WorkoutStorer     wk_s.WorkoutStorer
	WorkoutController wrk_c.WorkoutController
	UserStorer        u_s.UserStorer
}

func NewController(
	appCfg *config.Conf,
	loggerp *slog.Logger,
	uuidp uuid.Provider,
	client *mongo.Client,
	tp_store tp_s.TrainingProgramStorer,
	exc_storer wk_s.WorkoutStorer,
	us_storer u_s.UserStorer,
	wrk_contr wrk_c.WorkoutController,
) TrainingprogramController {
	impl := &TrainingprogramControllerImpl{
		Config:            appCfg,
		Logger:            loggerp,
		UUID:              uuidp,
		DbClient:          client,
		TpStorer:          tp_store,
		WorkoutStorer:     exc_storer,
		UserStorer:        us_storer,
		WorkoutController: wrk_contr,
	}
	return impl
}
