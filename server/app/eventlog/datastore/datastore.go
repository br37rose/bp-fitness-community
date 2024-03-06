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
	StatusPending                 = 1
	StatusOK                      = 2
	StatusError                   = 3
	StatusArchived                = 4
	PrimaryTypeStripeWebhookEvent = 1
)

type EventLog struct {
	PrimaryType   int8   `bson:"primary_type" json:"primary_type"`
	SecondaryType string `bson:"secondary_type" json:"secondary_type"`
	// Content field will store any datatype because we want the ability to
	// store complex data-structures from remote services (ex: Stripe, Inc.) and
	// be able to store simply stuff like plain ol' text.
	Content    any       `bson:"content" json:"content"`
	Error      error     `bson:"error" json:"error"`
	CreatedAt  time.Time `bson:"created_at,omitempty" json:"created_at,omitempty"`
	ModifiedAt time.Time `bson:"modified_at,omitempty" json:"modified_at,omitempty"`
	Status     int8      `bson:"status" json:"status"`
	// ExternalID represent the unique identifier provided by a remote system.
	ExternalID       string             `bson:"external_id" json:"external_id"`
	OrganizationID   primitive.ObjectID `bson:"organization_id" json:"organization_id"`
	OrganizationName string             `bson:"organization_name" json:"organization_name"`
	ID               primitive.ObjectID `bson:"_id" json:"id"`
}

type EventLogListFilter struct {
	// Pagination related.
	Cursor    primitive.ObjectID
	PageSize  int64
	SortField string
	SortOrder int8 // 1=ascending | -1=descending

	// Filter related.
	SearchText     string
	OrganizationID primitive.ObjectID
	BranchID       primitive.ObjectID
	PrimaryType    int8   `bson:"primary_type" json:"primary_type"`
	SecondaryType  string `bson:"secondary_type" json:"secondary_type"`
	Status         int8   `bson:"status" json:"status"`
}

type EventLogListResult struct {
	Results     []*EventLog        `json:"results"`
	NextCursor  primitive.ObjectID `json:"next_cursor"`
	HasNextPage bool               `json:"has_next_page"`
}

// EventLogStorer Interface for organization.
type EventLogStorer interface {
	Create(ctx context.Context, m *EventLog) error
	GetByID(ctx context.Context, id primitive.ObjectID) (*EventLog, error)
	GetByName(ctx context.Context, name string) (*EventLog, error)
	GetByPaymentProcessorEventLogID(ctx context.Context, paymentProcessorEventLogID string) (*EventLog, error)
	UpdateByID(ctx context.Context, m *EventLog) error
	ListByFilter(ctx context.Context, m *EventLogListFilter) (*EventLogListResult, error)
	ListAsSelectOptionByFilter(ctx context.Context, f *EventLogListFilter) ([]*EventLogAsSelectOption, error)
	DeleteByID(ctx context.Context, id primitive.ObjectID) error
	CheckIfExistsByNameInOrgBranch(ctx context.Context, name string, orgID primitive.ObjectID, branchID primitive.ObjectID) (bool, error)
	// //TODO: Add more...
}

type EventLogAsSelectOption struct {
	Value primitive.ObjectID `bson:"_id" json:"value"` // Extract from the database `_id` field and output through API as `value`.
	Label string             `bson:"name" json:"label"`
}

type EventLogStorerImpl struct {
	Logger     *slog.Logger
	DbClient   *mongo.Client
	Collection *mongo.Collection
}

func NewDatastore(appCfg *c.Conf, loggerp *slog.Logger, client *mongo.Client) EventLogStorer {
	// ctx := context.Background()
	uc := client.Database(appCfg.DB.Name).Collection("eventlogs")

	// The following few lines of code will create the index for our app for
	// this colleciton.
	indexModel := mongo.IndexModel{
		Keys: bson.D{
			{"content", "text"},
		},
	}
	_, err := uc.Indexes().CreateOne(context.TODO(), indexModel)
	if err != nil {
		// It is important that we crash the app on startup to meet the
		// requirements of `google/wire` framework.
		log.Fatal(err)
	}

	s := &EventLogStorerImpl{
		Logger:     loggerp,
		DbClient:   client,
		Collection: uc,
	}
	return s
}
