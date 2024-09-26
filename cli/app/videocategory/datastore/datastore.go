package datastore

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"log/slog"

	c "github.com/bci-innovation-labs/bp8fitnesscommunity-cli/config"
)

const (
	StatusActive   = 1
	StatusArchived = 2
	StatusError    = 3
	// OwnershipTypeTemporary indicates file has been uploaded and saved in our system but not assigned ownership to anything. As a result, if this videocategory is not assigned within 24 hours then the crontab will delete this videocategory record and the uploaded file.
	OwnershipTypeTemporary         = 1
	OwnershipTypeExerciseVideo     = 2
	OwnershipTypeExerciseThumbnail = 3
	OwnershipTypeUser              = 4
	OwnershipTypeOrganization      = 5
	ContentTypeFile                = 6
	ContentTypeImage               = 7
)

type VideoCategory struct {
	OrganizationID     primitive.ObjectID `bson:"organization_id,omitempty" json:"organization_id,omitempty"`
	OrganizationName   string             `bson:"organization_name" json:"organization_name"`
	ID                 primitive.ObjectID `bson:"_id" json:"id"`
	CreatedAt          time.Time          `bson:"created_at,omitempty" json:"created_at,omitempty"`
	CreatedByUserName  string             `bson:"created_by_user_name" json:"created_by_user_name"`
	CreatedByUserID    primitive.ObjectID `bson:"created_by_user_id" json:"created_by_user_id"`
	ModifiedAt         time.Time          `bson:"modified_at,omitempty" json:"modified_at,omitempty"`
	ModifiedByUserName string             `bson:"modified_by_user_name" json:"modified_by_user_name"`
	ModifiedByUserID   primitive.ObjectID `bson:"modified_by_user_id" json:"modified_by_user_id"`
	Name               string             `bson:"name" json:"name"`
	No                 int8               `bson:"no" json:"no"`
	Status             int8               `bson:"status" json:"status"`
}

type VideoCategoryListFilter struct {
	// Pagination related.
	Cursor    primitive.ObjectID
	PageSize  int64
	SortField string
	SortOrder int8 // 1=ascending | -1=descending

	// Filter related.
	OrganizationID primitive.ObjectID
	Status         int8
	SearchText     string
}

type VideoCategoryListResult struct {
	Results     []*VideoCategory   `json:"results"`
	NextCursor  primitive.ObjectID `json:"next_cursor"`
	HasNextPage bool               `json:"has_next_page"`
}

type VideoCategoryAsSelectOption struct {
	Value primitive.ObjectID `bson:"_id" json:"value"` // Extract from the database `_id` field and output through API as `value`.
	Label string             `bson:"name" json:"label"`
}

// VideoCategoryStorer Interface for video category.
type VideoCategoryStorer interface {
	Create(ctx context.Context, m *VideoCategory) error
	GetByID(ctx context.Context, id primitive.ObjectID) (*VideoCategory, error)
	UpdateByID(ctx context.Context, m *VideoCategory) error
	ListByFilter(ctx context.Context, m *VideoCategoryListFilter) (*VideoCategoryListResult, error)
	ListAsSelectOptionByFilter(ctx context.Context, f *VideoCategoryListFilter) ([]*VideoCategoryAsSelectOption, error)
	DeleteByID(ctx context.Context, id primitive.ObjectID) error
	// //TODO: Add more...
}

type VideoCategoryStorerImpl struct {
	Logger     *slog.Logger
	DbClient   *mongo.Client
	Collection *mongo.Collection
}

func NewDatastore(appCfg *c.Conf, loggerp *slog.Logger, client *mongo.Client) VideoCategoryStorer {
	// ctx := context.Background()
	uc := client.Database(appCfg.DB.Name).Collection("video_categories")

	// // The following few lines of code will create the index for our app for this
	// // colleciton.
	// indexModel := mongo.IndexModel{
	// 	Keys: bson.D{
	// 		{"organization_name", "text"},
	// 		{"name", "text"},
	// 	},
	// }
	// _, err := uc.Indexes().CreateOne(context.TODO(), indexModel)
	// if err != nil {
	// 	// It is important that we crash the app on startup to meet the
	// 	// requirements of `google/wire` framework.
	// 	log.Fatal(err)
	// }

	s := &VideoCategoryStorerImpl{
		Logger:     loggerp,
		DbClient:   client,
		Collection: uc,
	}
	return s
}

// Auto-generated comment for change 32
