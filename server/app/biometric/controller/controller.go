package controller

import (
	"context"
	"log/slog"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/adapter/cache/mongodbcache"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/adapter/distributedmutex"
	s3_storage "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/adapter/storage/s3"
	ap_s "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/aggregatepoint/datastore"
	dp_s "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/datapoint/datastore"
	organization_s "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/organization/datastore"
	rp_s "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/rankpoint/datastore"
	user_s "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/user/datastore"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/config"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/provider/uuid"
)

// BiometricController Interface for organization business logic controller.
type BiometricController interface {
	Leaderboard(ctx context.Context, req *LeaderboardRequest) (*rp_s.RankPointPaginationListResult, error)
	GetSummary(ctx context.Context, userID primitive.ObjectID) (*AggregatePointSummaryResponse, error)
	HistoricData(ctx context.Context, req *HistoricDataRequest) (*rp_s.RankPointPaginationListResult, error)
}

type BiometricControllerImpl struct {
	Config               *config.Conf
	Logger               *slog.Logger
	UUID                 uuid.Provider
	S3                   s3_storage.S3Storager
	DbClient             *mongo.Client
	Cache                mongodbcache.Cacher
	CodeVerifierMap      map[primitive.ObjectID]string
	DistributedMutex     distributedmutex.Adapter
	OrganizationStorer   organization_s.OrganizationStorer
	UserStorer           user_s.UserStorer
	DataPointStorer      dp_s.DataPointStorer
	AggregatePointStorer ap_s.AggregatePointStorer
	RankPointStorer      rp_s.RankPointStorer
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
	dp_storer dp_s.DataPointStorer,
	ap_storer ap_s.AggregatePointStorer,
	rp_storer rp_s.RankPointStorer,
) BiometricController {
	s := &BiometricControllerImpl{
		Config:               appCfg,
		Logger:               loggerp,
		UUID:                 uuidp,
		DbClient:             client,
		Cache:                cache,
		CodeVerifierMap:      make(map[primitive.ObjectID]string, 0),
		DistributedMutex:     dlocker,
		S3:                   s3,
		OrganizationStorer:   org_storer,
		UserStorer:           usr_storer,
		DataPointStorer:      dp_storer,
		AggregatePointStorer: ap_storer,
		RankPointStorer:      rp_storer,
	}
	s.Logger.Debug("biometric controller initialization started...")
	s.Logger.Debug("biometric controller initialized")
	return s
}
