package controller

import (
	"context"
	"log/slog"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/adapter/cache/mongodbcache"
	ap_s "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/aggregatepoint/datastore"
	dp_s "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/datapoint/datastore"
	fitbitapp_s "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/fitbitapp/datastore"
	fitbitdatum_s "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/fitbitdatum/datastore"
	organization_s "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/organization/datastore"
	user_s "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/user/datastore"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/config"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/provider/kmutex"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/provider/uuid"
)

// AggregatePointController Interface for organization business logic controller.
type AggregatePointController interface {
	ListByFilter(ctx context.Context, f *ap_s.AggregatePointPaginationListFilter) (*ap_s.AggregatePointPaginationListResult, error)
	AggregateThisHourForAllActiveFitBitApps(ctx context.Context) error
	AggregateLastHourForAllActiveFitBitApps(ctx context.Context) error
	AggregateTodayForAllActiveFitBitApps(ctx context.Context) error
	AggregateYesterdayForAllActiveFitBitApps(ctx context.Context) error
	AggregateThisISOWeekForAllActiveFitBitApps(ctx context.Context) error
	AggregateLastISOWeekForAllActiveFitBitApps(ctx context.Context) error
	AggregateThisMonthForAllActiveFitBitApps(ctx context.Context) error
	AggregateLastMonthForAllActiveFitBitApps(ctx context.Context) error
	AggregateThisYearForAllActiveFitBitApps(ctx context.Context) error
	AggregateLastYearForAllActiveFitBitApps(ctx context.Context) error
	GetSummary(ctx context.Context, userID primitive.ObjectID) (*AggregatePointSummaryResponse, error)
}

type AggregatePointControllerImpl struct {
	Config               *config.Conf
	Logger               *slog.Logger
	UUID                 uuid.Provider
	DbClient             *mongo.Client
	Cache                mongodbcache.Cacher
	CodeVerifierMap      map[primitive.ObjectID]string
	Kmutex               kmutex.Provider
	OrganizationStorer   organization_s.OrganizationStorer
	FitBitAppStorer      fitbitapp_s.FitBitAppStorer
	FitBitDatumStorer    fitbitdatum_s.FitBitDatumStorer
	UserStorer           user_s.UserStorer
	DataPointStorer      dp_s.DataPointStorer
	AggregatePointStorer ap_s.AggregatePointStorer
}

func NewController(
	appCfg *config.Conf,
	loggerp *slog.Logger,
	uuidp uuid.Provider,
	client *mongo.Client,
	cache mongodbcache.Cacher,
	kmutexp kmutex.Provider,
	org_storer organization_s.OrganizationStorer,
	fba_storer fitbitapp_s.FitBitAppStorer,
	usr_storer user_s.UserStorer,
	fbd_storer fitbitdatum_s.FitBitDatumStorer,
	dp_storer dp_s.DataPointStorer,
	ap_storer ap_s.AggregatePointStorer,
) AggregatePointController {
	s := &AggregatePointControllerImpl{
		Config:               appCfg,
		Logger:               loggerp,
		UUID:                 uuidp,
		DbClient:             client,
		Cache:                cache,
		CodeVerifierMap:      make(map[primitive.ObjectID]string, 0),
		Kmutex:               kmutexp,
		OrganizationStorer:   org_storer,
		FitBitAppStorer:      fba_storer,
		FitBitDatumStorer:    fbd_storer,
		UserStorer:           usr_storer,
		DataPointStorer:      dp_storer,
		AggregatePointStorer: ap_storer,
	}
	s.Logger.Debug("aggregate point controller initialization started...")
	s.Logger.Debug("aggregate point controller initialized")
	return s
}
