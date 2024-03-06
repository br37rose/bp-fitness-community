package controller

import (
	"context"
	"log"
	"log/slog"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	e_s "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/exercise/datastore"
	domain "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/offer/datastore"
	organization_s "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/organization/datastore"
	user_s "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/user/datastore"
	vcon_s "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/videocontent/datastore"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/config"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/provider/uuid"
)

// Offerontroller Interface for organization business logic controller.
type Offerontroller interface {
	// Create(ctx context.Context, m *domain.Offer) (*domain.Offer, error)
	GetByID(ctx context.Context, id primitive.ObjectID) (*domain.Offer, error)
	UpdateByID(ctx context.Context, m *domain.Offer) (*domain.Offer, error)
	ListByFilter(ctx context.Context, f *domain.OfferListFilter) (*domain.OfferListResult, error)
	ListAsSelectOptionByFilter(ctx context.Context, f *domain.OfferListFilter) ([]*domain.OfferAsSelectOption, error)
	// DeleteByID(ctx context.Context, id primitive.ObjectID) error
}

type OfferControllerImpl struct {
	Config             *config.Conf
	Logger             *slog.Logger
	UUID               uuid.Provider
	DbClient           *mongo.Client
	OrganizationStorer organization_s.OrganizationStorer
	OfferStorer        domain.OfferStorer
	UserStorer         user_s.UserStorer
	VideoContentStorer vcon_s.VideoContentStorer
	ExerciseStorer     e_s.ExerciseStorer
}

func NewController(
	appCfg *config.Conf,
	loggerp *slog.Logger,
	uuidp uuid.Provider,
	client *mongo.Client,
	org_storer organization_s.OrganizationStorer,
	sub_storer domain.OfferStorer,
	usr_storer user_s.UserStorer,
	e_storer e_s.ExerciseStorer,
	con_storer vcon_s.VideoContentStorer,
) Offerontroller {
	loggerp.Debug("offer controller initialization started...")
	s := &OfferControllerImpl{
		Config:             appCfg,
		Logger:             loggerp,
		UUID:               uuidp,
		DbClient:           client,
		OrganizationStorer: org_storer,
		OfferStorer:        sub_storer,
		UserStorer:         usr_storer,
		VideoContentStorer: con_storer,
		ExerciseStorer:     e_storer,
	}
	if appCfg.AppServer.IsDeveloperMode {
		if err := s.createDefaults(context.Background()); err != nil {
			log.Fatal(err)
		}
	}
	s.Logger.Debug("offer controller initialized")
	return s
}
