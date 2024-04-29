package controller

import (
	"context"
	"log/slog"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/adapter/cache/mongodbcache"
	ap_s "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/aggregatepoint/datastore"
	dp_s "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/datapoint/datastore"
	gfa_ds "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/googlefitapp/datastore"
	gfdp_s "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/googlefitdatapoint/datastore"
	organization_s "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/organization/datastore"
	user_s "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/user/datastore"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/config"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/provider/kmutex"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/provider/uuid"
)

// AggregatePointController Interface for organization business logic controller.
type AggregatePointController interface {
	ListByFilter(ctx context.Context, f *ap_s.AggregatePointPaginationListFilter) (*ap_s.AggregatePointPaginationListResult, error)
	AggregateThisHourForAllActiveGoogleFitApps(ctx context.Context) error    // DEPRECATED
	AggregateLastHourForAllActiveGoogleFitApps(ctx context.Context) error    // DEPRECATED
	AggregateTodayForAllActiveGoogleFitApps(ctx context.Context) error       // DEPRECATED
	AggregateYesterdayForAllActiveGoogleFitApps(ctx context.Context) error   // DEPRECATED
	AggregateThisISOWeekForAllActiveGoogleFitApps(ctx context.Context) error // DEPRECATED
	AggregateLastISOWeekForAllActiveGoogleFitApps(ctx context.Context) error // DEPRECATED
	AggregateThisMonthForAllActiveGoogleFitApps(ctx context.Context) error   // DEPRECATED
	AggregateLastMonthForAllActiveGoogleFitApps(ctx context.Context) error   // DEPRECATED
	AggregateThisYearForAllActiveGoogleFitApps(ctx context.Context) error    // DEPRECATED
	AggregateLastYearForAllActiveGoogleFitApps(ctx context.Context) error    // DEPRECATED
	GetSummary(ctx context.Context, userID primitive.ObjectID) (*AggregatePointSummaryResponse, error)
	AggregateForAllActiveGoogleFitApps(ctx context.Context) error
	AggregateForGoogleFitAppID(ctx context.Context, gfaID primitive.ObjectID) error
}

type AggregatePointControllerImpl struct {
	Config                   *config.Conf
	Logger                   *slog.Logger
	UUID                     uuid.Provider
	DbClient                 *mongo.Client
	Cache                    mongodbcache.Cacher
	CodeVerifierMap          map[primitive.ObjectID]string
	Kmutex                   kmutex.Provider
	OrganizationStorer       organization_s.OrganizationStorer
	UserStorer               user_s.UserStorer
	GoogleFitAppStorer       gfa_ds.GoogleFitAppStorer
	GoogleFitDataPointStorer gfdp_s.GoogleFitDataPointStorer
	DataPointStorer          dp_s.DataPointStorer
	AggregatePointStorer     ap_s.AggregatePointStorer
}

func NewController(
	appCfg *config.Conf,
	loggerp *slog.Logger,
	uuidp uuid.Provider,
	client *mongo.Client,
	cache mongodbcache.Cacher,
	kmutexp kmutex.Provider,
	org_storer organization_s.OrganizationStorer,
	usr_storer user_s.UserStorer,
	gfa_storer gfa_ds.GoogleFitAppStorer,
	gfdp_storer gfdp_s.GoogleFitDataPointStorer,
	dp_storer dp_s.DataPointStorer,
	ap_storer ap_s.AggregatePointStorer,
) AggregatePointController {
	s := &AggregatePointControllerImpl{
		Config:                   appCfg,
		Logger:                   loggerp,
		UUID:                     uuidp,
		DbClient:                 client,
		Cache:                    cache,
		CodeVerifierMap:          make(map[primitive.ObjectID]string, 0),
		Kmutex:                   kmutexp,
		OrganizationStorer:       org_storer,
		UserStorer:               usr_storer,
		GoogleFitAppStorer:       gfa_storer,
		GoogleFitDataPointStorer: gfdp_storer,
		DataPointStorer:          dp_storer,
		AggregatePointStorer:     ap_storer,
	}
	s.Logger.Debug("aggregate point controller initialization started...")
	s.Logger.Debug("aggregate point controller initialized")
	return s
}
