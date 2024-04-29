package datastore

import (
	"context"
	"log"
	"log/slog"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	c "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/config"
)

type DataPoint struct {
	ID        primitive.ObjectID `bson:"_id" json:"id,omitempty"`
	MetricID  primitive.ObjectID `bson:"metric_id" json:"metric_id,omitempty"`
	Timestamp time.Time          `bson:"timestamp" json:"timestamp,omitempty"`
	Value     float64            `bson:"value" json:"value"`
	IsNull    bool               `bson:"is_null" json:"is_null"`
}

type DataPointListFilter struct {
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

type DataPointListResult struct {
	Results     []*DataPoint       `json:"results"`
	NextCursor  primitive.ObjectID `json:"next_cursor"`
	HasNextPage bool               `json:"has_next_page"`
}

// DataPointStorer Interface for organization.
type DataPointStorer interface {
	CheckIfExistsByCompositeKey(ctx context.Context, metricID primitive.ObjectID, timestamp time.Time) (bool, error)
	GetByCompositeKey(ctx context.Context, metricID primitive.ObjectID, timestamp time.Time) (*DataPoint, error)
	GetByLatestTimestampAndMetricID(ctx context.Context, metricID primitive.ObjectID) (*DataPoint, error)
	Create(ctx context.Context, m *DataPoint) error
	ListByFilter(ctx context.Context, f *DataPointPaginationListFilter) (*DataPointPaginationListResult, error)
	ListForToday(ctx context.Context, metricIDs []primitive.ObjectID) (*DataPointPaginationListResult, error)
	ListForYesterday(ctx context.Context, metricIDs []primitive.ObjectID) (*DataPointPaginationListResult, error)
	ListForThisISOWeek(ctx context.Context, metricIDs []primitive.ObjectID) (*DataPointPaginationListResult, error)
	ListForLastISOWeek(ctx context.Context, metricIDs []primitive.ObjectID) (*DataPointPaginationListResult, error)
	ListForThisMonth(ctx context.Context, metricIDs []primitive.ObjectID) (*DataPointPaginationListResult, error)
	ListForLastMonth(ctx context.Context, metricIDs []primitive.ObjectID) (*DataPointPaginationListResult, error)
	ListForThisYear(ctx context.Context, metricIDs []primitive.ObjectID) (*DataPointPaginationListResult, error)
	ListForLastYear(ctx context.Context, metricIDs []primitive.ObjectID) (*DataPointPaginationListResult, error)
	Aggregate(ctx context.Context, metricID primitive.ObjectID, startAt, endAt time.Time) (*DataPointAggregateResult, error)
	//TODO: Add more...
}

type DataPointAsSelectOption struct {
	Value primitive.ObjectID `bson:"_id" json:"value"` // Extract from the database `_id` field and output through API as `value`.
	Label string             `bson:"name" json:"label"`
}

type DataPointStorerImpl struct {
	Logger     *slog.Logger
	DbClient   *mongo.Client
	Collection *mongo.Collection
}

func NewDatastore(appCfg *c.Conf, loggerp *slog.Logger, client *mongo.Client) DataPointStorer {

	// Special thanks: https://www.mongodb.com/community/forums/t/is-there-an-easier-way-to-create-time-series-collections-in-go/193071/2

	// ctx := context.Background()
	err := client.Database(appCfg.DB.Name).CreateCollection(context.TODO(), "data_points", &options.CreateCollectionOptions{
		TimeSeriesOptions: &options.TimeSeriesOptions{
			TimeField:   "timestamp",
			MetaField:   aws.String("metric_id"),
			Granularity: aws.String("hours"),
		},
		ExpireAfterSeconds: nil, // Note: https://www.mongodb.com/docs/manual/core/timeseries/timeseries-automatic-removal/#disable-automatic-removal
	})

	if err != nil {
		if !strings.Contains(err.Error(), "already exists") {
			log.Fatal(err)
		}
	}

	uc := client.Database(appCfg.DB.Name).Collection("data_points")

	// // // For debugging purposes only.
	// // if _, err := uc.Indexes().DropAll(context.TODO()); err != nil {
	// // 	loggerp.Error("failed deleting all indexes",
	// // 		slog.Any("err", err))
	// //
	// // 	// It is important that we crash the app on startup to meet the
	// // 	// requirements of `google/wire` framework.
	// // 	log.Fatal(err)
	// // }
	//
	_, err = uc.Indexes().CreateMany(context.TODO(), []mongo.IndexModel{
		{Keys: bson.D{{Key: "user_id", Value: 1}}},
		{Keys: bson.D{{Key: "metric_id", Value: 1}}},
		{Keys: bson.D{{Key: "timestamp", Value: 1}}},
		{Keys: bson.D{{Key: "type", Value: 1}}},
		{
			Keys: bson.D{ // Combined indeces.
				{Key: "timestamp", Value: 1},
				{Key: "metric_id", Value: 1},
			},
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	s := &DataPointStorerImpl{
		Logger:     loggerp,
		DbClient:   client,
		Collection: uc,
	}
	return s
}
