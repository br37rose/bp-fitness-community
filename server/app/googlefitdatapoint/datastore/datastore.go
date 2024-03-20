package datastore

import (
	"context"
	"log"
	"log/slog"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	gcp_a "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/adapter/cloudprovider/google"
	c "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/config"
)

type GoogleFitDataPoint struct {
	ID              primitive.ObjectID `bson:"_id" json:"id,omitempty"`
	DataTypeName    string             `bson:"data_type_name" json:"data_type_name"`
	Status          int8               `bson:"status" json:"status"`
	UserID          primitive.ObjectID `bson:"user_id" json:"user_id"`
	UserName        string             `bson:"user_name" json:"user_name"`
	UserLexicalName string             `bson:"user_lexical_name" json:"user_lexical_name"`

	// GoogleFitAppID is the reference ID to the `Google Fit App` we have registered with our system.
	GoogleFitAppID primitive.ObjectID `bson:"googlefit_app_id" json:"googlefit_app_id"`

	// MetricID is the ID of the metric we assigned for this record.
	MetricID primitive.ObjectID `bson:"metric_id" json:"metric_id"`

	StartAt time.Time `bson:"start_at" json:"start_at"`
	EndAt   time.Time `bson:"end_at" json:"end_at"`

	// Hydration data provided by `Google Fit`.
	StepCountDelta *gcp_a.StepCountDeltaStruct `bson:"step_count_delta,omitempty" json:"step_count_delta,omitempty"`
	Hydration        *gcp_a.HydrationStruct        `bson:"hydration,omitempty" json:"hydration,omitempty"`
	HeartRateBPM     *gcp_a.HeartRateBPMStruct     `bson:"hearte_rate_bpm,omitempty" json:"hearte_rate_bpm,omitempty"`

	// Error is the error response content provided by `Google Fit` when making the API call.
	Error string `bson:"errors" json:"errors"`

	// CreatedAt represents the time this data point got inserted into our database.
	CreatedAt time.Time `bson:"created_at" json:"created_at,omitempty"`

	// ModifiedAt represents the time this data point was modified by our app.
	ModifiedAt time.Time `bson:"modified_at" json:"modified_at,omitempty"`

	// The organization this data point belongs to. Used for tenancy restrictions.
	OrganizationID primitive.ObjectID `bson:"organization_id" json:"organization_id"`
}

type GoogleFitDataPointListFilter struct {
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

type GoogleFitDataPointListResult struct {
	Results     []*GoogleFitDataPoint `json:"results"`
	NextCursor  primitive.ObjectID    `json:"next_cursor"`
	HasNextPage bool                  `json:"has_next_page"`
}

// GoogleFitDataPointStorer Interface for organization.
type GoogleFitDataPointStorer interface {
	CheckIfExistsByCompositeKey(ctx context.Context, userID primitive.ObjectID, dataTypeName string, startAt time.Time, endAt time.Time) (bool, error)
	Create(ctx context.Context, m *GoogleFitDataPoint) error
	// GetByID(ctx context.Context, id primitive.ObjectID) (*GoogleFitDataPoint, error)
	// GetByUserID(ctx context.Context, userID primitive.ObjectID) (*GoogleFitDataPoint, error)
	// GetByName(ctx context.Context, name string) (*GoogleFitDataPoint, error)
	// GetByPaymentProcessorGoogleFitDataPointID(ctx context.Context, paymentProcessorGoogleFitDataPointID string) (*GoogleFitDataPoint, error)
	UpdateByID(ctx context.Context, m *GoogleFitDataPoint) error
	// UpsertByUserID(ctx context.Context, fba *GoogleFitDataPoint) error
	// ListByFilter(ctx context.Context, m *GoogleFitDataPointListFilter) (*GoogleFitDataPointListResult, error)
	// ListAsSelectOptionByFilter(ctx context.Context, f *GoogleFitDataPointListFilter) ([]*GoogleFitDataPointAsSelectOption, error)
	ListByQueuedStatus(ctx context.Context) (*GoogleFitDataPointListResult, error)
	// DeleteByID(ctx context.Context, id primitive.ObjectID) error
	// CheckIfExistsByNameInOrgBranch(ctx context.Context, name string, orgID primitive.ObjectID, branchID primitive.ObjectID) (bool, error)
	// // //TODO: Add more...
}

type GoogleFitDataPointAsSelectOption struct {
	Value primitive.ObjectID `bson:"_id" json:"value"` // Extract from the database `_id` field and output through API as `value`.
	Label string             `bson:"name" json:"label"`
}

type GoogleFitDataPointStorerImpl struct {
	Logger     *slog.Logger
	DbClient   *mongo.Client
	Collection *mongo.Collection
}

func NewDatastore(appCfg *c.Conf, loggerp *slog.Logger, client *mongo.Client) GoogleFitDataPointStorer {
	// ctx := context.Background()
	uc := client.Database(appCfg.DB.Name).Collection("googlefit_data_points")

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
		{Keys: bson.D{{Key: "googlefit_app_id", Value: 1}}},
		{Keys: bson.D{{Key: "start_at", Value: 1}}},
		{Keys: bson.D{{Key: "status", Value: 1}}},
		{Keys: bson.D{{Key: "type", Value: 1}}},
		{Keys: bson.D{
			{"user_id", "text"},
			{"data_type_name", "text"},
			{"start_at", "text"},
			{"end_at", "text"},
		}},
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

	s := &GoogleFitDataPointStorerImpl{
		Logger:     loggerp,
		DbClient:   client,
		Collection: uc,
	}
	return s
}
