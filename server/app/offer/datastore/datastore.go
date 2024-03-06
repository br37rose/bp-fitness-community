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
	// StatusPending indicates the user has not posted this offer and remains hidden from all the users until active status is used.
	StatusPending = 1
	// StatusActive controls that offer is active.
	StatusActive = 2
	// StatusArchived controls that offer is deleted.
	StatusArchived = 3
	// BusinessFunctionUnspecified indicates no business logic was set with this offer.
	BusinessFunctionUnspecified = 0
	// BusinessFunctionProvideMembershipAccessToContentAccess indicates the offer is based on this applications ability to grant access to content based on this membership rank.
	BusinessFunctionProvideMembershipAccessToContentAccess = 1
	// PayFrequencyOneTime indicates user only pays once
	PayFrequencyOneTime = 1
	PayFrequencyDay     = 2
	PayFrequencyWeek    = 3
	// PayFrequencyMonthly indicates user pays monthly.
	PayFrequencyMonthly = 4
	// PayFrequencyAnnual indicates user pays annual.
	PayFrequencyAnnual = 5
	// OfferTypeService indicates user gets access to Gym or staff
	OfferTypeService = 1
	// OfferTypeProduct indicates user gets physical product
	OfferTypeProduct = 2
)

type Offer struct {
	Name                 string             `bson:"name" json:"name"`
	Description          string             `bson:"description" json:"description"`
	Price                float64            `bson:"price" json:"price"`
	PriceCurrency        string             `bson:"price_currency" json:"price_currency"`
	PayFrequency         int8               `bson:"pay_frequency" json:"pay_frequency"`
	Status               int8               `bson:"status" json:"status"`
	Type                 int8               `bson:"type" json:"type"`
	CreatedAt            time.Time          `bson:"created_at,omitempty" json:"created_at,omitempty"`
	ModifiedAt           time.Time          `bson:"modified_at,omitempty" json:"modified_at,omitempty"`
	OrganizationID       primitive.ObjectID `bson:"organization_id" json:"organization_id"`
	OrganizationName     string             `bson:"organization_name" json:"organization_name"`
	ID                   primitive.ObjectID `bson:"_id" json:"id"`
	PaymentProcessorName string             `bson:"payment_processor_name" json:"payment_processor_name"`
	StripeProductID      string             `bson:"stripe_product_id" json:"stripe_product_id"`
	StripePriceID        string             `bson:"stripe_price_id" json:"stripe_price_id"`
	StripeImageURL       string             `bson:"stripe_image_url" json:"stripe_image_url"`
	IsSubscription       bool               `bson:"is_subscription" json:"is_subscription"`

	// Controls how the user is able to book in our system. Special thanks to http://www.heppnetz.de/ontologies/goodrelations/v1#BusinessFunction.
	BusinessFunction int8 `bson:"business_function" json:"business_function"`
	// MembershipRank is unique identifier to specify this offer's value in the ranking system, higher is better.
	MembershipRank int `bson:"membership_rank" json:"membership_rank"`
}

type OfferListFilter struct {
	// Pagination related.
	Cursor    primitive.ObjectID
	PageSize  int64
	SortField string
	SortOrder int8 // 1=ascending | -1=descending

	// Filter related.
	OrganizationID primitive.ObjectID
	BranchID       primitive.ObjectID
	SearchText     string
	Status         int8
}

type OfferListResult struct {
	Results     []*Offer           `json:"results"`
	NextCursor  primitive.ObjectID `json:"next_cursor"`
	HasNextPage bool               `json:"has_next_page"`
}

// OfferStorer Interface for organization.
type OfferStorer interface {
	Create(ctx context.Context, m *Offer) error
	GetByID(ctx context.Context, id primitive.ObjectID) (*Offer, error)
	GetByStripeProductID(ctx context.Context, stripeProductID string) (*Offer, error)
	GetByStripePriceID(ctx context.Context, stripePriceID string) (*Offer, error)
	GetByName(ctx context.Context, name string) (*Offer, error)
	UpdateByID(ctx context.Context, m *Offer) error
	// Upsert will create the `Offer` record if the `ID` field doesn't exist, or update the `Offer` record if the `ID` exists.
	Upsert(ctx context.Context, offer *Offer) error
	ListByFilter(ctx context.Context, m *OfferListFilter) (*OfferListResult, error)
	ListAsSelectOptionByFilter(ctx context.Context, f *OfferListFilter) ([]*OfferAsSelectOption, error)
	DeleteByID(ctx context.Context, id primitive.ObjectID) error
	CheckIfExistsByNameInOrgBranch(ctx context.Context, name string, orgID primitive.ObjectID, branchID primitive.ObjectID) (bool, error)
	CheckIfExistsByID(ctx context.Context, id primitive.ObjectID) (bool, error)
	// //TODO: Add more...
}

type OfferAsSelectOption struct {
	Value primitive.ObjectID `bson:"_id" json:"value"` // Extract from the database `_id` field and output through API as `value`.
	Label string             `bson:"name" json:"label"`
}

type OfferStorerImpl struct {
	Logger     *slog.Logger
	DbClient   *mongo.Client
	Collection *mongo.Collection
}

func NewDatastore(appCfg *c.Conf, loggerp *slog.Logger, client *mongo.Client) OfferStorer {
	// ctx := context.Background()
	uc := client.Database(appCfg.DB.Name).Collection("offers")

	// The following few lines of code will create the index for our app for
	// this colleciton.
	indexModel := mongo.IndexModel{
		Keys: bson.D{
			{"organization_name", "text"},
			{"branch_name", "text"},
			{"name", "text"},
			{"description", "text"},
		},
	}
	_, err := uc.Indexes().CreateOne(context.TODO(), indexModel)
	if err != nil {
		// It is important that we crash the app on startup to meet the
		// requirements of `google/wire` framework.
		log.Fatal(err)
	}

	s := &OfferStorerImpl{
		Logger:     loggerp,
		DbClient:   client,
		Collection: uc,
	}
	return s
}
