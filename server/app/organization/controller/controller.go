package controller

import (
	"context"
	"log"
	"log/slog"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	s3_storage "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/adapter/storage/s3"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/adapter/templatedemailer"
	o_s "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/organization/datastore"
	user_s "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/user/datastore"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/config"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/provider/password"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/provider/uuid"
)

// OrganizationController Interface for organization business logic controller.
type OrganizationController interface {
	Create(ctx context.Context, m *o_s.Organization) (*o_s.Organization, error)
	GetByID(ctx context.Context, id primitive.ObjectID) (*o_s.Organization, error)
	UpdateByID(ctx context.Context, m *o_s.Organization) (*o_s.Organization, error)
	ListByFilter(ctx context.Context, f *o_s.OrganizationListFilter) (*o_s.OrganizationListResult, error)
	DeleteByID(ctx context.Context, id primitive.ObjectID) error
}

type OrganizationControllerImpl struct {
	Config             *config.Conf
	Logger             *slog.Logger
	Password           password.Provider
	UUID               uuid.Provider
	S3                 s3_storage.S3Storager
	DbClient           *mongo.Client
	OrganizationStorer o_s.OrganizationStorer
	UserStorer         user_s.UserStorer
	TemplatedEmailer   templatedemailer.TemplatedEmailer
}

func NewController(
	appCfg *config.Conf,
	loggerp *slog.Logger,
	uuidp uuid.Provider,
	passwordp password.Provider,
	s3 s3_storage.S3Storager,
	client *mongo.Client,
	te templatedemailer.TemplatedEmailer,
	org_storer o_s.OrganizationStorer,
	usr_storer user_s.UserStorer,
) OrganizationController {
	s := &OrganizationControllerImpl{
		Config:             appCfg,
		Logger:             loggerp,
		UUID:               uuidp,
		S3:                 s3,
		DbClient:           client,
		Password:           passwordp,
		TemplatedEmailer:   te,
		OrganizationStorer: org_storer,
		UserStorer:         usr_storer,
	}
	s.Logger.Debug("organization controller initialization started...")

	// Execute the code which will check to see if we have an initial accounts
	// if not then we'll need to create it.
	if err := s.createInitialRootAdmin(context.Background()); err != nil {
		log.Fatal(err) // We terminate app here b/c dependency injection not allowed to fail, so fail here at startup of our app.
	}
	org, err := s.createInitialOrg(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	if org != nil {
		if err := s.createInitialOrgAdmin(context.Background(), org); err != nil {
			log.Fatal(err) // We terminate app here b/c dependency injection not allowed to fail, so fail here at startup of our app.
		}
	}
	s.Logger.Debug("organization controller initialized")
	return s
}
