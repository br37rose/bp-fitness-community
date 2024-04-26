package controller

import (
	"context"
	"log/slog"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	fc_s "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/fitnesschallenge/datastore"
	rank_c "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/rankpoint/controller"
	rank "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/rankpoint/datastore"
	u_s "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/user/datastore"
	wrk_c "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/workout/controller"
	wk_s "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/workout/datastore"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/config"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/provider/uuid"
)

type FitnessChallengeController interface {
	Create(ctx context.Context, req *FitnessChallengeCreateRequestIDO) (*fc_s.FitnessChallenge, error)
	GetByID(ctx context.Context, id primitive.ObjectID) (*fc_s.FitnessChallenge, error)
	DeleteByID(ctx context.Context, id primitive.ObjectID) error
	ListByFilter(ctx context.Context, f *fc_s.FitnessChallengeListFilter) (*fc_s.FitnessChallengeListResult, error)
	UpdateByID(ctx context.Context, req *FitnessChallengeUpdateRequestIDO) (*fc_s.FitnessChallenge, error)
	ChangeParticipationStatus(ctx context.Context, id primitive.ObjectID) (*fc_s.FitnessChallenge, error)
	GetChallengeLeaderBoard(ctx context.Context, id primitive.ObjectID) (*rank.RankPointPaginationListResult, error)
}

type FitnessChallengeControllerImpl struct {
	Config            *config.Conf
	Logger            *slog.Logger
	UUID              uuid.Provider
	DbClient          *mongo.Client
	Storer            fc_s.FitnessChallengeStorer
	WorkoutStorer     wk_s.WorkoutStorer
	WorkoutController wrk_c.WorkoutController
	UserStorer        u_s.UserStorer
	RankPoint         rank_c.RankPointController
}

func NewController(
	appCfg *config.Conf,
	loggerp *slog.Logger,
	uuidp uuid.Provider,
	client *mongo.Client,
	tp_store fc_s.FitnessChallengeStorer,
	exc_storer wk_s.WorkoutStorer,
	us_storer u_s.UserStorer,
	wrk_contr wrk_c.WorkoutController,
	rank_c rank_c.RankPointController,
) FitnessChallengeController {
	impl := &FitnessChallengeControllerImpl{
		Config:            appCfg,
		Logger:            loggerp,
		UUID:              uuidp,
		DbClient:          client,
		Storer:            tp_store,
		WorkoutStorer:     exc_storer,
		UserStorer:        us_storer,
		WorkoutController: wrk_contr,
		RankPoint:         rank_c,
	}
	return impl
}
