package datastore

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"log/slog"

	c "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/config"
)

const (
	StatusActive   = 1
	StatusArchived = 2
)

type Invoice struct {
	OrganizationID   primitive.ObjectID `bson:"organization_id" json:"organization_id"`
	OrganizationName string             `bson:"organization_name" json:"organization_name"`
	ID               primitive.ObjectID `bson:"_id" json:"id"`
	UserName         string             `bson:"user_name" json:"user_name"`
	UserLexicalName  string             `bson:"user_lexical_name" json:"user_lexical_name"`
	UserID           primitive.ObjectID `bson:"user_id" json:"user_id"`
	Status           int8               `bson:"status" json:"status"`
	CreatedAt        time.Time          `bson:"created_at,omitempty" json:"created_at,omitempty"`
	ModifiedAt       time.Time          `bson:"modified_at,omitempty" json:"modified_at,omitempty"`

	// The name of the payment processor we used to create this invoice.
	PaymentProcessorName string `bson:"payment_processor_name" json:"payment_processor_name"`

	// The unique id set by the payment processor for this particular invoice.
	PaymentProcessorInvoiceID string `bson:"payment_processor_invoice_id" json:"payment_processor_invoice_id"`

	// The Stripe specific representation of what an invoice should be.
	StripeInvoice *StripeInvoice `bson:"stripe_invoice,omitempty" json:"stripe_invoice,omitempty"`
}

type StripeInvoice struct {
	// The unique identification created by Stripe to present this particular invoice.
	ID string `bson:"id" json:"id"`
	// Time at which the object was created. Measured in seconds since the Unix epoch.
	Created int64 `json:"created"`
	// Whether payment was successfully collected for this invoice. An invoice can be paid (most commonly) with a charge or with credit from the customer's account balance.
	Paid bool `json:"paid"`
	// The URL for the hosted invoice page, which allows customers to view and pay an invoice. If the invoice has not been finalized yet, this will be null.
	HostedInvoiceURL string `bson:"hosted_invoice_url" json:"hosted_invoice_url"`
	// The link to download the PDF for the invoice. If the invoice has not been finalized yet, this will be null.
	InvoicePDF string `json:"invoice_pdf"`
}

type InvoiceListFilter struct {
	// Pagination related.
	Cursor    primitive.ObjectID
	PageSize  int64
	SortField string
	SortOrder int8 // 1=ascending | -1=descending

	// Filter related.
	OrganizationID  primitive.ObjectID
	BranchID        primitive.ObjectID
	ExcludeArchived bool
	SearchText      string
}

type InvoiceListResult struct {
	Results     []*Invoice         `json:"results"`
	NextCursor  primitive.ObjectID `json:"next_cursor"`
	HasNextPage bool               `json:"has_next_page"`
}

// InvoiceStorer Interface for organization.
type InvoiceStorer interface {
	Create(ctx context.Context, m *Invoice) error
	GetByID(ctx context.Context, id primitive.ObjectID) (*Invoice, error)
	GetByName(ctx context.Context, name string) (*Invoice, error)
	GetByPaymentProcessorInvoiceID(ctx context.Context, paymentProcessorInvoiceID string) (*Invoice, error)
	UpdateByID(ctx context.Context, m *Invoice) error
	ListByFilter(ctx context.Context, m *InvoiceListFilter) (*InvoiceListResult, error)
	ListAsSelectOptionByFilter(ctx context.Context, f *InvoiceListFilter) ([]*InvoiceAsSelectOption, error)
	DeleteByID(ctx context.Context, id primitive.ObjectID) error
	CheckIfExistsByNameInOrgBranch(ctx context.Context, name string, orgID primitive.ObjectID, branchID primitive.ObjectID) (bool, error)
	// //TODO: Add more...
}

type InvoiceAsSelectOption struct {
	Value primitive.ObjectID `bson:"_id" json:"value"` // Extract from the database `_id` field and output through API as `value`.
	Label string             `bson:"name" json:"label"`
}

type InvoiceStorerImpl struct {
	Logger     *slog.Logger
	DbClient   *mongo.Client
	Collection *mongo.Collection
}

func NewDatastore(appCfg *c.Conf, loggerp *slog.Logger, client *mongo.Client) InvoiceStorer {
	// ctx := context.Background()
	uc := client.Database(appCfg.DB.Name).Collection("invoices")

	// The following few lines of code will create the index for our app for
	// this colleciton.
	indexModel := mongo.IndexModel{
		Keys: bson.D{
			{"organization_name", "text"},
			{"branch_name", "text"},
			{"name", "text"},
			{"trainer_name", "text"},
		},
	}
	_, err := uc.Indexes().CreateOne(context.TODO(), indexModel)
	if err != nil {
		// It is important that we crash the app on startup to meet the
		// requirements of `google/wire` framework.
		log.Fatal(err)
	}

	s := &InvoiceStorerImpl{
		Logger:     loggerp,
		DbClient:   client,
		Collection: uc,
	}
	return s
}
