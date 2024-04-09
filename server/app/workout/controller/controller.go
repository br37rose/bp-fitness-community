package controller

import (
	"context"
	"log/slog"

	w_d "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/workout/datastore"
	w_s "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/workout/datastore"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/config"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/provider/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type WorkoutController interface {
	Create(ctx context.Context, req *WorkoutCreateRequestIDO) (*w_d.Workout, error)
	GetByID(ctx context.Context, id primitive.ObjectID) (*w_d.Workout, error)
	DeleteByID(ctx context.Context, id primitive.ObjectID) error
	UpdateByID(ctx context.Context, req *WorkoutUpdateRequest) (*w_d.Workout, error)
	ListByFilter(ctx context.Context, f *w_d.WorkoutListFilter) (*w_d.WorkoutistResult, error)
}

type WorkoutControllerImpl struct {
	Config        *config.Conf
	Logger        *slog.Logger
	UUID          uuid.Provider
	DbClient      *mongo.Client
	WorkoutStorer w_s.WorkoutStorer
}

func NewController(
	appCfg *config.Conf,
	loggerp *slog.Logger,
	uuidp uuid.Provider,
	client *mongo.Client,
	workout_store w_s.WorkoutStorer,
) WorkoutController {
	impl := &WorkoutControllerImpl{
		Config:        appCfg,
		Logger:        loggerp,
		UUID:          uuidp,
		DbClient:      client,
		WorkoutStorer: workout_store,
	}
	return impl
}
