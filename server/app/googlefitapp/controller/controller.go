package controller

import (
	"context"
	"log/slog"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/adapter/cache/mongodbcache"
	gcp_a "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/adapter/cloudprovider/google"
	ap_s "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/aggregatepoint/datastore"
	dp_s "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/datapoint/datastore"
	googlefitapp_s "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/googlefitapp/datastore"
	gfdp_s "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/googlefitdatapoint/datastore"
	organization_s "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/organization/datastore"
	user_s "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/user/datastore"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/config"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/provider/kmutex"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/provider/uuid"
)

// GoogleFitAppController Interface for organization business logic controller.
type GoogleFitAppController interface {
	GetGoogleLoginURL(ctx context.Context) (*GoogleLoginURLResponse, error)
	GoogleCallback(ctx context.Context, state, code string) (*GoogleCallbackResponse, error)
	RefreshTokensFromGoogle() error
	ProcessAllQueuedData() error
}

type GoogleFitAppControllerImpl struct {
	Config                   *config.Conf
	Logger                   *slog.Logger
	UUID                     uuid.Provider
	DbClient                 *mongo.Client
	Cache                    mongodbcache.Cacher
	CodeVerifierMap          map[primitive.ObjectID]string
	Kmutex                   kmutex.Provider
	GCP                      gcp_a.GoogleCloudPlatformAdapter
	OrganizationStorer       organization_s.OrganizationStorer
	GoogleFitDataPointStorer gfdp_s.GoogleFitDataPointStorer
	GoogleFitAppStorer       googlefitapp_s.GoogleFitAppStorer
	UserStorer               user_s.UserStorer
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
	gcpa gcp_a.GoogleCloudPlatformAdapter,
	org_storer organization_s.OrganizationStorer,
	gfdp_storer gfdp_s.GoogleFitDataPointStorer,
	gfa_storer googlefitapp_s.GoogleFitAppStorer,
	usr_storer user_s.UserStorer,
	dp_storer dp_s.DataPointStorer,
	ap_storer ap_s.AggregatePointStorer,
) GoogleFitAppController {
	s := &GoogleFitAppControllerImpl{
		Config:                   appCfg,
		Logger:                   loggerp,
		UUID:                     uuidp,
		DbClient:                 client,
		Cache:                    cache,
		CodeVerifierMap:          make(map[primitive.ObjectID]string, 0),
		Kmutex:                   kmutexp,
		GCP:                      gcpa,
		OrganizationStorer:       org_storer,
		GoogleFitDataPointStorer: gfdp_storer,
		GoogleFitAppStorer:       gfa_storer,
		UserStorer:               usr_storer,
		DataPointStorer:          dp_storer,
		AggregatePointStorer:     ap_storer,
	}
	s.Logger.Debug("googlefit app controller initialization started...")
	s.Logger.Debug("googlefit app controller initialized")
	return s
}
