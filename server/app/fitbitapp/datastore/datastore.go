package datastore

import (
	"context"
	"log"
	"time"

	"log/slog"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	c "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/config"
)

const (
	StatusActive   = 1
	StatusArchived = 2
	StatusError    = 3
	AuthTypeOAuth2 = 1
)

type FitBitApp struct {
	ID               primitive.ObjectID `bson:"_id" json:"id,omitempty"`
	UserFirstName    string             `bson:"user_first_name" json:"user_first_name"`
	UserLastName     string             `bson:"user_last_name" json:"user_last_name"`
	UserName         string             `bson:"user_name" json:"user_name"`
	UserLexicalName  string             `bson:"user_lexical_name" json:"user_lexical_name"`
	UserID           primitive.ObjectID `bson:"user_id" json:"user_id"`
	Status           int8               `bson:"status" json:"status"`
	CreatedAt        time.Time          `bson:"created_at" json:"created_at,omitempty"`
	ModifiedAt       time.Time          `bson:"modified_at" json:"modified_at,omitempty"`
	OrganizationID   primitive.ObjectID `bson:"organization_id" json:"organization_id"`
	OrganizationName string             `bson:"organization_name" json:"organization_name"`

	// FitBitUserID is the user id provided by fitbit web-services for our user account.
	FitBitUserID string `bson:"fitbit_user_id" json:"fitbit_user_id"`

	// AuthType tracks how the user authenticated their fitbit with our app.
	AuthType int8 `bson:"auth_type" json:"auth_type"`

	// Errors indicates what error was returned by FitBit web-services.
	Errors        string    `bson:"errors" json:"errors,omitempty"`
	Scope         string    `bson:"scope" json:"scope,omitempty"`
	TokenType     string    `bson:"token_type" json:"token_type,omitempty"`
	AccessToken   string    `bson:"access_token" json:"access_token,omitempty"`
	ExpiresIn     int64     `bson:"expires_in" json:"expires_in"`
	RefreshToken  string    `bson:"refresh_token" json:"refresh_token,omitempty"`
	ExpireTime    time.Time `bson:"expire_time" json:"expire_time,omitempty"`
	LastFetchedAt time.Time `bson:"last_fetched_at" json:"last_fetched_at,omitempty"`

	HeartRateMetricID  primitive.ObjectID `bson:"heart_rate_metric_id" json:"heart_rate_metric_id,omitempty"`
	StepsCountMetricID primitive.ObjectID `bson:"steps_count_metric_id" json:"steps_count_metric_id,omitempty"`

	IsTestMode         bool   `bson:"is_test_mode" json:"is_test_mode"`
	SimulatorAlgorithm string `bson:"simulator_algorithm,omitempty" json:"simulator_algorithm,omitempty"`
}

type FitBitAppListFilter struct {
	// Pagination related.
	Cursor    primitive.ObjectID
	PageSize  int64
	SortField string
	SortOrder int8 // 1=ascending | -1=descending

	// Filter related.
	OrganizationID primitive.ObjectID
	BranchID       primitive.ObjectID
	Status         int8
	SearchText     string
}

type FitBitAppListResult struct {
	Results     []*FitBitApp       `json:"results"`
	NextCursor  primitive.ObjectID `json:"next_cursor"`
	HasNextPage bool               `json:"has_next_page"`
}

// FitBitAppStorer Interface for organization.
type FitBitAppStorer interface {
	Create(ctx context.Context, m *FitBitApp) error
	GetByID(ctx context.Context, id primitive.ObjectID) (*FitBitApp, error)
	GetByUserID(ctx context.Context, userID primitive.ObjectID) (*FitBitApp, error)
	GetByName(ctx context.Context, name string) (*FitBitApp, error)
	GetByPaymentProcessorFitBitAppID(ctx context.Context, paymentProcessorFitBitAppID string) (*FitBitApp, error)
	UpdateByID(ctx context.Context, m *FitBitApp) error
	UpsertByUserID(ctx context.Context, fba *FitBitApp) error
	ListByFilter(ctx context.Context, m *FitBitAppListFilter) (*FitBitAppListResult, error)
	ListAsSelectOptionByFilter(ctx context.Context, f *FitBitAppListFilter) ([]*FitBitAppAsSelectOption, error)
	ListIDsByStatus(ctx context.Context, status int8) ([]primitive.ObjectID, error)
	ListDevicesByStatus(ctx context.Context, status int8) ([]*FitBitAppDevice, error)
	ListPhysicalIDsByStatus(ctx context.Context, status int8) ([]primitive.ObjectID, error)
	ListSimulatorIDsByStatus(ctx context.Context, status int8) ([]primitive.ObjectID, error)
	DeleteByID(ctx context.Context, id primitive.ObjectID) error
	CheckIfExistsByNameInOrgBranch(ctx context.Context, name string, orgID primitive.ObjectID, branchID primitive.ObjectID) (bool, error)
	// //TODO: Add more...
}

type FitBitAppAsSelectOption struct {
	Value primitive.ObjectID `bson:"_id" json:"value"` // Extract from the database `_id` field and output through API as `value`.
	Label string             `bson:"name" json:"label"`
}

type FitBitAppStorerImpl struct {
	Logger     *slog.Logger
	DbClient   *mongo.Client
	Collection *mongo.Collection
}

func NewDatastore(appCfg *c.Conf, loggerp *slog.Logger, client *mongo.Client) FitBitAppStorer {
	// ctx := context.Background()
	uc := client.Database(appCfg.DB.Name).Collection("fitbit_apps")

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
		{Keys: bson.D{{Key: "status", Value: 1}}},
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

	s := &FitBitAppStorerImpl{
		Logger:     loggerp,
		DbClient:   client,
		Collection: uc,
	}
	return s
}
