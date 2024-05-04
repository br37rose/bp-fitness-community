package datastore

import (
	"context"
	"log"
	"log/slog"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	c "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/config"
)

// Composite key = MetricID | Period | Start | End

type AggregatePoint struct {
	ID                 primitive.ObjectID `bson:"_id" json:"id,omitempty"`
	MetricID           primitive.ObjectID `bson:"metric_id" json:"metric_id,omitempty"`
	MetricDataTypeName string             `bson:"metric_data_type_name" json:"metric_data_type_name"`
	UserID             primitive.ObjectID `bson:"user_id" json:"user_id,omitempty"`
	Period             int8               `bson:"period" json:"period"`
	Start              time.Time          `bson:"start" json:"start"`
	End                time.Time          `bson:"end" json:"end"`
	Count              float64            `bson:"count" json:"count"`
	Average            float64            `bson:"average" json:"average"`
	Min                float64            `bson:"min" json:"min"`
	Max                float64            `bson:"max" json:"max"`
	Sum                float64            `bson:"sum" json:"sum"`
}

type AggregatePointListFilter struct {
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

type AggregatePointListResult struct {
	Results     []*AggregatePoint  `json:"results"`
	NextCursor  primitive.ObjectID `json:"next_cursor"`
	HasNextPage bool               `json:"has_next_page"`
}

// AggregatePointStorer Interface for organization.
type AggregatePointStorer interface {
	CheckIfExistsByCompositeKey(ctx context.Context, metricID primitive.ObjectID, timestamp time.Time) (bool, error)
	GetByID(ctx context.Context, metricID primitive.ObjectID) (*AggregatePoint, error)
	GetByCompositeKey(ctx context.Context, metricID primitive.ObjectID, period int8, start time.Time, end time.Time) (*AggregatePoint, error)
	Create(ctx context.Context, m *AggregatePoint) error
	ListByFilter(ctx context.Context, f *AggregatePointPaginationListFilter) (*AggregatePointPaginationListResult, error)
	UpdateByID(ctx context.Context, m *AggregatePoint) error
	//TODO: Add more...
}

type AggregatePointAsSelectOption struct {
	Value primitive.ObjectID `bson:"_id" json:"value"` // Extract from the database `_id` field and output through API as `value`.
	Label string             `bson:"name" json:"label"`
}

type AggregatePointStorerImpl struct {
	Logger     *slog.Logger
	DbClient   *mongo.Client
	Collection *mongo.Collection
}

func NewDatastore(appCfg *c.Conf, loggerp *slog.Logger, client *mongo.Client) AggregatePointStorer {
	// ctx := context.Background()
	uc := client.Database(appCfg.DB.Name).Collection("aggregate_points")

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
	_, err := uc.Indexes().CreateMany(context.TODO(), []mongo.IndexModel{
		{Keys: bson.D{{Key: "metric_id", Value: 1}}},
		{
			Keys: bson.D{ // Combined indeces for our `composite key`.
				{Key: "metric_id", Value: 1},
				{Key: "period", Value: 1},
				{Key: "start", Value: 1},
				{Key: "end", Value: 1},
			},
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	s := &AggregatePointStorerImpl{
		Logger:     loggerp,
		DbClient:   client,
		Collection: uc,
	}
	return s
}
