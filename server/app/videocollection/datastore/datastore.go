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
	StatusActive   = 1
	StatusArchived = 2

	VideoCollectionThumbnailTypeS3          = 1
	VideoCollectionThumbnailTypeExternalURL = 2
	VideoCollectionThumbnailTypeLocal       = 3

	VideoCollectionTypeManyVideos  = 1
	VideoCollectionTypeSingleVideo = 2
)

type VideoCollection struct {
	ID          primitive.ObjectID `bson:"_id" json:"id"`
	Name        string             `bson:"name" json:"name"`
	Summary     string             `bson:"summary" json:"summary"`
	Description string             `bson:"description" json:"description"`

	ThumbnailType               int8               `bson:"thumbnail_type" json:"thumbnail_type"`
	ThumbnailAttachmentID       primitive.ObjectID `bson:"thumbnail_attachment_id" json:"thumbnail_attachment_id"`
	ThumbnailAttachmentFilename string             `bson:"thumbnail_attachment_filename" json:"thumbnail_attachment_filename"`
	ThumbnailObjectKey          string             `bson:"thumbnail_object_key" json:"thumbnail_object_key"`
	ThumbnailObjectURL          string             `bson:"thumbnail_object_url" json:"thumbnail_object_url"`
	ThumbnailObjectExpiry       time.Time          `bson:"thumbnail_object_expiry" json:"thumbnail_object_expiry"`
	ThumbnailName               string             `bson:"thumbnail_name" json:"thumbnail_name"`
	ThumbnailURL                string             `bson:"thumbnail_url" json:"thumbnail_url"`

	Type   int8 `bson:"type" json:"type"`
	Status int8 `bson:"status" json:"status"`
	Count  int8 `bson:"count" json:"count"`

	CreatedAt          time.Time          `bson:"created_at" json:"created_at,omitempty"`
	CreatedByUserID    primitive.ObjectID `bson:"created_by_user_id" json:"created_by_user_id"`
	CreatedByUserName  string             `bson:"created_by_user_name" json:"created_by_user_name"`
	ModifiedAt         time.Time          `bson:"modified_at" json:"modified_at,omitempty"`
	ModifiedByUserID   primitive.ObjectID `bson:"modified_by_user_id" json:"modified_by_user_id"`
	ModifiedByUserName string             `bson:"modified_by_user_name" json:"modified_by_user_name"`

	CategoryID       primitive.ObjectID `bson:"category_id" json:"category_id,omitempty"`
	CategoryName     string             `bson:"category_name" json:"category_name"`
	OrganizationID   primitive.ObjectID `bson:"organization_id" json:"organization_id,omitempty"`
	OrganizationName string             `bson:"organization_name" json:"organization_name"`
}

type VideoCollectionListFilter struct {
	// Pagination related.
	Cursor    primitive.ObjectID
	PageSize  int64
	SortField string
	SortOrder int8 // 1=ascending | -1=descending

	// Filter related.
	OrganizationID primitive.ObjectID `bson:"organization_id" json:"organization_id,omitempty"`
	SearchText     string             `json:"search_text"`
	Type           int8               `bson:"type" json:"type"`
	Status         int8               `bson:"status" json:"status"`
}

type VideoCollectionListResult struct {
	Results     []*VideoCollection `json:"results"`
	NextCursor  primitive.ObjectID `json:"next_cursor"`
	HasNextPage bool               `json:"has_next_page"`
}

type VideoCollectionAsSelectOption struct {
	Value primitive.ObjectID `bson:"_id" json:"value"` // Extract from the database `_id` field and output through API as `value`.
	Label string             `bson:"name" json:"label"`
}

// VideoCollectionStorer Interface for user.
type VideoCollectionStorer interface {
	CheckIfExistsByID(ctx context.Context, id primitive.ObjectID) (bool, error)
	Create(ctx context.Context, m *VideoCollection) error
	GetByID(ctx context.Context, id primitive.ObjectID) (*VideoCollection, error)
	UpdateByID(ctx context.Context, m *VideoCollection) error
	UpsertByID(ctx context.Context, e *VideoCollection) error
	ListByFilter(ctx context.Context, f *VideoCollectionListFilter) (*VideoCollectionListResult, error)
	ListAsSelectOptionByFilter(ctx context.Context, f *VideoCollectionListFilter) ([]*VideoCollectionAsSelectOption, error)
	DeleteByID(ctx context.Context, id primitive.ObjectID) error
	// //TODO: Add more...
}

type VideoCollectionStorerImpl struct {
	Logger     *slog.Logger
	DbClient   *mongo.Client
	Collection *mongo.Collection
}

func NewDatastore(appCfg *c.Conf, loggerp *slog.Logger, client *mongo.Client) VideoCollectionStorer {
	// ctx := context.Background()
	uc := client.Database(appCfg.DB.Name).Collection("video_collections")

	// The following few lines of code will create the index for our app for this
	// colleciton.
	indexModel := mongo.IndexModel{
		Keys: bson.D{
			{"name", "text"},
			{"summary", "text"},
			{"description", "text"},
			{"thumbnail_attachment_filename", "text"},
			{"thumbnail_name", "text"},
			{"category_name", "text"},
		},
	}
	_, err := uc.Indexes().CreateOne(context.TODO(), indexModel)
	if err != nil {
		// It is important that we crash the app on startup to meet the
		// requirements of `google/wire` framework.
		log.Fatal(err)
	}

	s := &VideoCollectionStorerImpl{
		Logger:     loggerp,
		DbClient:   client,
		Collection: uc,
	}
	return s
}
