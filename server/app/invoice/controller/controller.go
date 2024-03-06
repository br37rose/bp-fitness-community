package controller

import (
	"context"
	"log/slog"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	domain "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/invoice/datastore"
	organization_s "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/organization/datastore"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/config"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/provider/uuid"
)

// InvoiceController Interface for organization business logic controller.
type InvoiceController interface {
	Create(ctx context.Context, m *domain.Invoice) (*domain.Invoice, error)
	GetByID(ctx context.Context, id primitive.ObjectID) (*domain.Invoice, error)
	UpdateByID(ctx context.Context, m *domain.Invoice) (*domain.Invoice, error)
	ListByFilter(ctx context.Context, f *domain.InvoiceListFilter) (*domain.InvoiceListResult, error)
	ListAsSelectOptionByFilter(ctx context.Context, f *domain.InvoiceListFilter) ([]*domain.InvoiceAsSelectOption, error)
	DeleteByID(ctx context.Context, id primitive.ObjectID) error
}

type InvoiceControllerImpl struct {
	Config             *config.Conf
	Logger             *slog.Logger
	UUID               uuid.Provider
	DbClient           *mongo.Client
	OrganizationStorer organization_s.OrganizationStorer
	InvoiceStorer      domain.InvoiceStorer
}

func NewController(
	appCfg *config.Conf,
	loggerp *slog.Logger,
	uuidp uuid.Provider,
	client *mongo.Client,
	org_storer organization_s.OrganizationStorer,
	sub_storer domain.InvoiceStorer,
) InvoiceController {
	s := &InvoiceControllerImpl{
		Config:             appCfg,
		Logger:             loggerp,
		UUID:               uuidp,
		DbClient:           client,
		OrganizationStorer: org_storer,
		InvoiceStorer:      sub_storer,
	}
	s.Logger.Debug("organization controller initialization started...")
	s.Logger.Debug("organization controller initialized")
	return s
}
