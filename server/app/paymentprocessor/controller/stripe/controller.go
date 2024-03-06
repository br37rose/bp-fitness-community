package stripe

import (
	"context"
	"log/slog"

	"go.mongodb.org/mongo-driver/bson/primitive"

	mg "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/adapter/emailer/mailgun"
	pm "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/adapter/paymentprocessor/stripe"
	s3_storage "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/adapter/storage/s3"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/adapter/templatedemailer"
	eventlog_s "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/eventlog/datastore"
	i_s "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/invoice/datastore"
	offer_s "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/offer/datastore"
	org_s "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/organization/datastore"
	user_s "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/user/datastore"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/config"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/provider/password"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/provider/uuid"
)

type StripePaymentProcessorController interface {
	CreateStripeCheckoutSessionURL(ctx context.Context, offerIDString string) (string, error)
	CompleteStripeCheckoutSession(ctx context.Context, sessionID string) (*CompleteStripeCheckoutSessionResponse, error)
	CancelStripeSubscription(ctx context.Context, userID primitive.ObjectID) error
	ListLatestStripeInvoices(ctx context.Context, userID primitive.ObjectID, cursor int64, limit int64) (*user_s.StripeInvoiceListResult, error)
	// SendStripeCheckoutSessionEmail(ctx context.Context, requestData *SendSubscriptionRequestEmailRequestIDO) error
	// GrantFreeCredits(ctx context.Context, requestData *GrantFreeCreditRequestIDO) error
	Webhook(ctx context.Context, header string, b []byte) error
}

type StripePaymentProcessorControllerImpl struct {
	Config             *config.Conf
	Logger             *slog.Logger
	UUID               uuid.Provider
	S3                 s3_storage.S3Storager
	Password           password.Provider
	Emailer            mg.Emailer
	TemplatedEmailer   templatedemailer.TemplatedEmailer
	PaymentProcessor   pm.PaymentProcessor
	OrganizationStorer org_s.OrganizationStorer
	UserStorer         user_s.UserStorer
	InvoiceStorer      i_s.InvoiceStorer
	OfferStorer        offer_s.OfferStorer
	EventLogStorer     eventlog_s.EventLogStorer
}

func NewController(
	appCfg *config.Conf,
	loggerp *slog.Logger,
	uuidp uuid.Provider,
	s3 s3_storage.S3Storager,
	passwordp password.Provider,
	emailer mg.Emailer,
	te templatedemailer.TemplatedEmailer,
	paymentProcessor pm.PaymentProcessor,
	org_storer org_s.OrganizationStorer,
	sub_storer user_s.UserStorer,
	is i_s.InvoiceStorer,
	offs offer_s.OfferStorer,
	evel eventlog_s.EventLogStorer,
) StripePaymentProcessorController {
	loggerp.Debug("member controller initialization started...")
	s := &StripePaymentProcessorControllerImpl{
		Config:             appCfg,
		Logger:             loggerp,
		UUID:               uuidp,
		S3:                 s3,
		Password:           passwordp,
		Emailer:            emailer,
		TemplatedEmailer:   te,
		PaymentProcessor:   paymentProcessor,
		OrganizationStorer: org_storer,
		UserStorer:         sub_storer,
		InvoiceStorer:      is,
		OfferStorer:        offs,
		EventLogStorer:     evel,
	}
	s.Logger.Debug("member controller initialized")
	return s
}
