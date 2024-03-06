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
	OrganizationStatusPending  = 1
	OrganizationStatusActive   = 2
	OrganizationStatusError    = 3
	OrganizationStatusInactive = 4
	GymType                    = 1
)

type Organization struct {
	ID              primitive.ObjectID `bson:"_id" json:"id"`
	CreatedAt       time.Time          `bson:"created_at,omitempty" json:"created_at,omitempty"`
	ModifiedAt      time.Time          `bson:"modified_at,omitempty" json:"modified_at,omitempty"`
	Type            int8               `bson:"type" json:"type"`
	Status          int8               `bson:"status" json:"status"`
	Name            string             `bson:"name" json:"name"`
	Description     string             `bson:"description" json:"description"`
	OpeningHours    []string           `bson:"opening_hours" json:"opening_hours"`
	WebsiteURL      []string           `bson:"website_url" json:"website_url"`
	Email           string             `bson:"email" json:"email"`
	Phone           string             `bson:"phone" json:"phone,omitempty"`
	Country         string             `bson:"country" json:"country,omitempty"`
	Region          string             `bson:"region" json:"region,omitempty"`
	City            string             `bson:"city" json:"city,omitempty"`
	PostalCode      string             `bson:"postal_code" json:"postal_code,omitempty"`
	AddressLine1    string             `bson:"address_line_1" json:"address_line_1,omitempty"`
	AddressLine2    string             `bson:"address_line_2" json:"address_line_2,omitempty"`
	CreatedByUserID primitive.ObjectID `bson:"created_by_user_id" json:"created_by_user_id"`
}

type OrganizationListFilter struct {
	PageSize  int64
	LastID    string
	SortField string
	UserID    primitive.ObjectID
	UserRole  int8

	// SortOrder string   `json:"sort_order"`
	// SortField string   `json:"sort_field"`
	// Offset    uint64   `json:"offset"`
	// Limit     uint64   `json:"limit"`
	// Statuss    []int8   `json:"statuss"`
	// UUIDs     []string `json:"uuids"`
}

type OrganizationListResult struct {
	Results []*Organization `json:"results"`
}

// OrganizationStorer Interface for organization.
type OrganizationStorer interface {
	Create(ctx context.Context, m *Organization) error
	GetByID(ctx context.Context, id primitive.ObjectID) (*Organization, error)
	GetByName(ctx context.Context, name string) (*Organization, error)
	UpdateByID(ctx context.Context, m *Organization) error
	ListByFilter(ctx context.Context, m *OrganizationListFilter) (*OrganizationListResult, error)
	DeleteByID(ctx context.Context, id primitive.ObjectID) error
	CheckIfExistsByName(ctx context.Context, name string) (bool, error)
	// //TODO: Add more...
}

type OrganizationStorerImpl struct {
	Logger     *slog.Logger
	DbClient   *mongo.Client
	Collection *mongo.Collection
}

func NewDatastore(appCfg *c.Conf, loggerp *slog.Logger, client *mongo.Client) OrganizationStorer {
	// ctx := context.Background()
	uc := client.Database(appCfg.DB.Name).Collection("organizations")

	// The following few lines of code will create the index for our app for this
	// colleciton.
	indexModel := mongo.IndexModel{
		Keys: bson.D{
			{"organization_name", "text"},
			{"name", "text"},
			{"description", "text"},
			{"email", "text"},
			{"phone", "text"},
			{"country", "text"},
			{"region", "text"},
			{"city", "text"},
			{"postal_code", "text"},
			{"address_line_1", "text"},
		},
	}
	_, err := uc.Indexes().CreateOne(context.TODO(), indexModel)
	if err != nil {
		// It is important that we crash the app on startup to meet the
		// requirements of `google/wire` framework.
		log.Fatal(err)
	}

	s := &OrganizationStorerImpl{
		Logger:     loggerp,
		DbClient:   client,
		Collection: uc,
	}
	return s
}
