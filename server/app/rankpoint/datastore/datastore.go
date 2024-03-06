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

// Composite key = UserID | MetricID | Period | Start | End

// RankPoint is user's global rank in the application for the specific user
// period and datetime range.
type RankPoint struct {
	ID primitive.ObjectID `bson:"_id" json:"id,omitempty"`

	// Place represents the spot the user is in the global ranking, for example first place would be `1` and second place would be 2.
	Place uint64 `bson:"place" json:"place"`

	MetricType int8               `bson:"metric_type" json:"metric_type"`
	MetricID   primitive.ObjectID `bson:"metric_id" json:"metric_id,omitempty"`

	Function int8      `bson:"function" json:"function"`
	Period   int8      `bson:"period" json:"period"`
	Start    time.Time `bson:"start" json:"start"`
	End      time.Time `bson:"end" json:"end"`

	// Value represents the metric value that was selected for, can be either average or summation but it all depends on the metric.
	Value float64 `bson:"value" json:"value"`

	UserFirstName string             `bson:"user_first_name" json:"user_first_name"`
	UserLastName  string             `bson:"user_last_name" json:"user_last_name"`
	UserID        primitive.ObjectID `bson:"user_id" json:"user_id"`

	UserAvatarObjectExpiry time.Time `bson:"user_avatar_object_expiry" json:"user_avatar_object_expiry"`
	UserAvatarObjectURL    string    `bson:"user_avatar_object_url" json:"user_avatar_object_url"`
	UserAvatarObjectKey    string    `bson:"user_avatar_object_key" json:"user_avatar_object_key"`
	UserAvatarFileType     string    `bson:"user_avatar_file_type" json:"user_avatar_file_type"`
	UserAvatarFileName     string    `bson:"user_avatar_file_name" json:"user_avatar_file_name"`

	OrganizationID   primitive.ObjectID `bson:"organization_id" json:"organization_id,omitempty"`
	OrganizationName string             `bson:"organization_name" json:"organization_name"`
}

// RankPointStorer Interface for organization.
type RankPointStorer interface {
	CheckIfExistsByCompositeKey(ctx context.Context, metricID primitive.ObjectID, timestamp time.Time) (bool, error)
	GetByID(ctx context.Context, metricID primitive.ObjectID) (*RankPoint, error)
	GetByCompositeKey(ctx context.Context, metricID primitive.ObjectID, function int8, period int8, start time.Time, end time.Time) (*RankPoint, error)
	GetByCompositeKeyForToday(ctx context.Context, metricID primitive.ObjectID, function int8) (*RankPoint, error)
	GetByCompositeKeyForYesterday(ctx context.Context, metricID primitive.ObjectID, function int8) (*RankPoint, error)
	GetByCompositeKeyForThisISOWeek(ctx context.Context, metricID primitive.ObjectID, function int8) (*RankPoint, error)
	GetByCompositeKeyForLastISOWeek(ctx context.Context, metricID primitive.ObjectID, function int8) (*RankPoint, error)
	GetByCompositeKeyForThisMonth(ctx context.Context, metricID primitive.ObjectID, function int8) (*RankPoint, error)
	GetByCompositeKeyForLastMonth(ctx context.Context, metricID primitive.ObjectID, function int8) (*RankPoint, error)
	GetByCompositeKeyForThisYear(ctx context.Context, metricID primitive.ObjectID, function int8) (*RankPoint, error)
	GetByCompositeKeyForLastYear(ctx context.Context, metricID primitive.ObjectID, function int8) (*RankPoint, error)
	Create(ctx context.Context, m *RankPoint) error
	ListByFilter(ctx context.Context, f *RankPointPaginationListFilter) (*RankPointPaginationListResult, error)
	ListWithinPlace(ctx context.Context, metricTypes []int8, function int8, period int8, start, end uint64) (*RankPointPaginationListResult, error)
	ListWithinPlaceAndToday(ctx context.Context, metricTypes []int8, function int8, period int8, start, end uint64) (*RankPointPaginationListResult, error)
	ListWithinPlaceAndISOWeek(ctx context.Context, metricTypes []int8, function int8, period int8, start, end uint64) (*RankPointPaginationListResult, error)
	ListWithinPlaceAndMonth(ctx context.Context, metricTypes []int8, function int8, period int8, start, end uint64) (*RankPointPaginationListResult, error)
	ListWithinPlaceAndYear(ctx context.Context, metricTypes []int8, function int8, period int8, start, end uint64) (*RankPointPaginationListResult, error)
	UpdateByID(ctx context.Context, m *RankPoint) error
	//TODO: Add more...
}

type RankPointAsSelectOption struct {
	Value primitive.ObjectID `bson:"_id" json:"value"` // Extract from the database `_id` field and output through API as `value`.
	Label string             `bson:"name" json:"label"`
}

type RankPointStorerImpl struct {
	Logger     *slog.Logger
	DbClient   *mongo.Client
	Collection *mongo.Collection
}

func NewDatastore(appCfg *c.Conf, loggerp *slog.Logger, client *mongo.Client) RankPointStorer {
	// ctx := context.Background()
	uc := client.Database(appCfg.DB.Name).Collection("rank_points")

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
				{Key: "function", Value: 1},
				{Key: "period", Value: 1},
				{Key: "start", Value: 1},
				{Key: "end", Value: 1},
			},
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	s := &RankPointStorerImpl{
		Logger:     loggerp,
		DbClient:   client,
		Collection: uc,
	}
	return s
}
