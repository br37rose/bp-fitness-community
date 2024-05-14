package controller

import (
	"context"
	"log/slog"

	"go.mongodb.org/mongo-driver/mongo"

	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/adapter/cache/mongodbcache"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/adapter/distributedmutex"
	s3_storage "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/adapter/storage/s3"
	ap_s "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/aggregatepoint/datastore"
	dp_s "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/datapoint/datastore"
	gfa_ds "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/googlefitapp/datastore"
	gfdp_s "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/googlefitdatapoint/datastore"
	organization_s "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/organization/datastore"
	rp_s "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/rankpoint/datastore"
	user_s "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/user/datastore"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/config"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/provider/uuid"
)

// RankPointController Interface for organization business logic controller.
type RankPointController interface {
	ListByFilter(ctx context.Context, f *rp_s.RankPointPaginationListFilter) (*rp_s.RankPointPaginationListResult, error)
	GenerateGlobalRankingForActiveGoogleFitApps(ctx context.Context) error
}

type RankPointControllerImpl struct {
	Config                   *config.Conf
	Logger                   *slog.Logger
	UUID                     uuid.Provider
	S3                       s3_storage.S3Storager
	DbClient                 *mongo.Client
	Cache                    mongodbcache.Cacher
	DistributedMutex         distributedmutex.Adapter
	OrganizationStorer       organization_s.OrganizationStorer
	UserStorer               user_s.UserStorer
	GoogleFitAppStorer       gfa_ds.GoogleFitAppStorer
	GoogleFitDataPointStorer gfdp_s.GoogleFitDataPointStorer
	DataPointStorer          dp_s.DataPointStorer
	AggregatePointStorer     ap_s.AggregatePointStorer
	RankPointStorer          rp_s.RankPointStorer
}

func NewController(
	appCfg *config.Conf,
	loggerp *slog.Logger,
	uuidp uuid.Provider,
	client *mongo.Client,
	cache mongodbcache.Cacher,
	dlocker distributedmutex.Adapter,
	s3 s3_storage.S3Storager,
	org_storer organization_s.OrganizationStorer,
	usr_storer user_s.UserStorer,
	gfa_storer gfa_ds.GoogleFitAppStorer,
	gfdp_storer gfdp_s.GoogleFitDataPointStorer,
	dp_storer dp_s.DataPointStorer,
	ap_storer ap_s.AggregatePointStorer,
	rp_storer rp_s.RankPointStorer,
) RankPointController {
	s := &RankPointControllerImpl{
		Config:                   appCfg,
		Logger:                   loggerp,
		UUID:                     uuidp,
		DbClient:                 client,
		Cache:                    cache,
		DistributedMutex:         dlocker,
		S3:                       s3,
		OrganizationStorer:       org_storer,
		UserStorer:               usr_storer,
		GoogleFitAppStorer:       gfa_storer,
		GoogleFitDataPointStorer: gfdp_storer,
		DataPointStorer:          dp_storer,
		AggregatePointStorer:     ap_storer,
		RankPointStorer:          rp_storer,
	}
	s.Logger.Debug("rank point controller initialization started...")
	// if err := rp_storer.DeleteAll(context.Background()); err != nil {
	// 	panic(err)
	// }
	s.Logger.Debug("rank point controller initialized")
	return s
}
