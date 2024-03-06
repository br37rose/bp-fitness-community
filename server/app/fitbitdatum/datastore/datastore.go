package datastore

import (
	"context"
	"log"
	"time"

	"log/slog"

	fb_api "github.com/Thomas2500/go-fitbit"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	c "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/config"
)

type FitBitFormattedDatum struct {
	ActivitiesInterdayLog *fb_api.ActivitiesInterdayLog `bson:"activities_interday_log,omitempty" json:"activities_interday_log,omitempty"`
	HeartIntraday         *fb_api.HeartIntraday         `bson:"heart_intraday,omitempty" json:"user_heart_intraday,omitempty"`
	HeartDay              *fb_api.HeartDay              `bson:"heart_day,omitempty" json:"heart_day,omitempty"`
}

type FitBitDatum struct {
	ID              primitive.ObjectID `bson:"_id" json:"id,omitempty"`
	Type            int64              `bson:"type" json:"type"`
	Status          int8               `bson:"status" json:"status"`
	UserID          primitive.ObjectID `bson:"user_id" json:"user_id"`
	UserName        string             `bson:"user_name" json:"user_name"`
	UserLexicalName string             `bson:"user_lexical_name" json:"user_lexical_name"`

	// FitBitUserID is the user id provided by fitbit web-services for our user account.
	FitBitUserID string `bson:"fitbit_user_id" json:"fitbit_user_id"`

	// FitBitAppID is the authorized fitbit we have in our system.
	FitBitAppID primitive.ObjectID `bson:"fitbit_app_id" json:"fitbit_app_id"`

	// FitBitAppID is the authorized fitbit we have in our system.
	MetricID primitive.ObjectID `bson:"metric_id" json:"metric_id"`

	Path      string               `bson:"path" json:"path"`
	StartAt   time.Time            `bson:"start_at" json:"start_at"`
	EndAt     time.Time            `bson:"end_at" json:"end_at"`
	Raw       string               `bson:"raw" json:"raw"`
	Formatted FitBitFormattedDatum `bson:"formatted" json:"formatted"`
	Errors    string               `bson:"errors" json:"errors"`
	// SecondaryTypeID int64     `json:"secondary_type_id"` // ???
	CreatedAt      time.Time          `bson:"created_at" json:"created_at,omitempty"`
	ModifiedAt     time.Time          `bson:"modified_at" json:"modified_at,omitempty"`
	OrganizationID primitive.ObjectID `bson:"organization_id" json:"organization_id"`
}

type FitBitDatumListFilter struct {
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
	Status          int8
}

type FitBitDatumListResult struct {
	Results     []*FitBitDatum     `json:"results"`
	NextCursor  primitive.ObjectID `json:"next_cursor"`
	HasNextPage bool               `json:"has_next_page"`
}

// FitBitDatumStorer Interface for organization.
type FitBitDatumStorer interface {
	Create(ctx context.Context, m *FitBitDatum) error
	// GetByID(ctx context.Context, id primitive.ObjectID) (*FitBitDatum, error)
	// GetByUserID(ctx context.Context, userID primitive.ObjectID) (*FitBitDatum, error)
	// GetByName(ctx context.Context, name string) (*FitBitDatum, error)
	// GetByPaymentProcessorFitBitDatumID(ctx context.Context, paymentProcessorFitBitDatumID string) (*FitBitDatum, error)
	UpdateByID(ctx context.Context, m *FitBitDatum) error
	// UpsertByUserID(ctx context.Context, fba *FitBitDatum) error
	// ListByFilter(ctx context.Context, m *FitBitDatumListFilter) (*FitBitDatumListResult, error)
	// ListAsSelectOptionByFilter(ctx context.Context, f *FitBitDatumListFilter) ([]*FitBitDatumAsSelectOption, error)
	ListByQueuedStatus(ctx context.Context) (*FitBitDatumListResult, error)
	// DeleteByID(ctx context.Context, id primitive.ObjectID) error
	// CheckIfExistsByNameInOrgBranch(ctx context.Context, name string, orgID primitive.ObjectID, branchID primitive.ObjectID) (bool, error)
	// // //TODO: Add more...
}

type FitBitDatumAsSelectOption struct {
	Value primitive.ObjectID `bson:"_id" json:"value"` // Extract from the database `_id` field and output through API as `value`.
	Label string             `bson:"name" json:"label"`
}

type FitBitDatumStorerImpl struct {
	Logger     *slog.Logger
	DbClient   *mongo.Client
	Collection *mongo.Collection
}

func NewDatastore(appCfg *c.Conf, loggerp *slog.Logger, client *mongo.Client) FitBitDatumStorer {
	// ctx := context.Background()
	uc := client.Database(appCfg.DB.Name).Collection("fitbit_data")

	// // For debugging purposes only.
	// if _, err := uc.Indexes().DropAll(context.TODO()); err != nil {
	// 	loggerp.Error("failed deleting all indexes",
	// 		slog.Any("err", err))
	//
	// 	// It is important that we crash the app on startup to meet the
	// 	// requirements of `google/wire` framework.
	// 	log.Fatal(err)
	// }

	_, err := uc.Indexes().CreateMany(context.TODO(), []mongo.IndexModel{
		{Keys: bson.D{{Key: "user_id", Value: 1}}},
		{Keys: bson.D{{Key: "fitbit_user_id", Value: 1}}},
		{Keys: bson.D{{Key: "fitbit_app_id", Value: 1}}},
		{Keys: bson.D{{Key: "start_at", Value: 1}}},
		{Keys: bson.D{{Key: "status", Value: 1}}},
		{Keys: bson.D{{Key: "type", Value: 1}}},
		// {Keys: bson.D{
		// 	{"name", "text"},
		// 	{"lexical_name", "text"},
		// 	{"email", "text"},
		// 	{"phone", "text"},
		// 	{"country", "text"},
		// 	{"region", "text"},
		// 	{"city", "text"},
		// 	{"postal_code", "text"},
		// 	{"address_line1", "text"},
		// 	{"description", "text"},
		// }},
	})
	if err != nil {
		// It is important that we crash the app on startup to meet the
		// requirements of `google/wire` framework.
		log.Fatal(err)
	}

	s := &FitBitDatumStorerImpl{
		Logger:     loggerp,
		DbClient:   client,
		Collection: uc,
	}
	return s
}
