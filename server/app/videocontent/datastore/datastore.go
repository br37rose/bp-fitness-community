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
	VideoContentStatusActive   = 1
	VideoContentStatusArchived = 2

	VideoContentVideoTypeS3      = 1
	VideoContentVideoTypeYouTube = 2
	VideoContentVideoTypeVimeo   = 3

	VideoContentThumbnailTypeS3          = 1
	VideoContentThumbnailTypeExternalURL = 2
	VideoContentThumbnailTypeLocal       = 3

	VideoContentTypeSystem = 1
	VideoContentTypeCustom = 2

	TimedLock1Day    = "24h"
	TimedLock2Days   = "48h"
	TimedLock1Week   = "168h"
	TimedLock2Weeks  = "336h"
	TimedLock3Weeks  = "504h"
	TimedLock4Weeks  = "672h"
	TimedLock1Month  = "720h"
	TimedLock2Months = "1440h"
	TimedLock3Months = "2160h"
	TimedLock4Months = "2880h"
	TimedLock5Months = "3600h"
	TimedLock6Months = "4320h"
	TimedLock1Year   = "8760h"
)

type VideoContent struct {
	ID          primitive.ObjectID `bson:"_id" json:"id"`
	No          int8               `bson:"no" json:"no"`
	Name        string             `bson:"name" json:"name"`
	Description string             `bson:"description" json:"description"`
	AuthorName  string             `bson:"author_name" json:"author_name"`
	AuthorURL   string             `bson:"author_url" json:"author_url"`
	Duration    string             `bson:"duration" json:"duration"`

	Type   int8 `bson:"type" json:"type"`
	Status int8 `bson:"status" json:"status"`

	VideoType               int8               `bson:"video_type" json:"video_type"`
	VideoAttachmentID       primitive.ObjectID `bson:"video_attachment_id" json:"video_attachment_id"`
	VideoAttachmentFilename string             `bson:"video_attachment_filename" json:"video_attachment_filename"`
	VideoObjectKey          string             `bson:"video_object_key" json:"video_object_key"`
	VideoObjectURL          string             `bson:"video_object_url" json:"video_object_url"`
	VideoObjectExpiry       time.Time          `bson:"video_object_expiry" json:"video_object_expiry"`
	VideoName               string             `bson:"video_name" json:"video_name"`
	VideoURL                string             `bson:"video_url" json:"video_url"`

	ThumbnailType               int8               `bson:"thumbnail_type" json:"thumbnail_type"`
	ThumbnailAttachmentID       primitive.ObjectID `bson:"thumbnail_attachment_id" json:"thumbnail_attachment_id"`
	ThumbnailAttachmentFilename string             `bson:"thumbnail_attachment_filename" json:"thumbnail_attachment_filename"`
	ThumbnailObjectKey          string             `bson:"thumbnail_object_key" json:"thumbnail_object_key"`
	ThumbnailObjectURL          string             `bson:"thumbnail_object_url" json:"thumbnail_object_url"`
	ThumbnailObjectExpiry       time.Time          `bson:"thumbnail_object_expiry" json:"thumbnail_object_expiry"`
	ThumbnailName               string             `bson:"thumbnail_name" json:"thumbnail_name"`
	ThumbnailURL                string             `bson:"thumbnail_url" json:"thumbnail_url"`
	VideoContentTags            []*VideoContentTag `bson:"videocontent_tags" json:"videocontent_tags"`

	CategoryID       primitive.ObjectID `bson:"category_id" json:"category_id,omitempty"`
	CategoryName     string             `bson:"category_name" json:"category_name"`
	CollectionID     primitive.ObjectID `bson:"collection_id" json:"collection_id,omitempty"`
	CollectionName   string             `bson:"collection_name" json:"collection_name"`
	OrganizationID   primitive.ObjectID `bson:"organization_id" json:"organization_id,omitempty"`
	OrganizationName string             `bson:"organization_name" json:"organization_name"`

	CreatedAt          time.Time          `bson:"created_at" json:"created_at,omitempty"`
	CreatedByUserID    primitive.ObjectID `bson:"created_by_user_id" json:"created_by_user_id"`
	CreatedByUserName  string             `bson:"created_by_user_name" json:"created_by_user_name"`
	ModifiedAt         time.Time          `bson:"modified_at" json:"modified_at,omitempty"`
	ModifiedByUserID   primitive.ObjectID `bson:"modified_by_user_id" json:"modified_by_user_id"`
	ModifiedByUserName string             `bson:"modified_by_user_name" json:"modified_by_user_name"`
	HasMonetization    bool               `bson:"has_monetization" json:"has_monetization"`
	OfferID            primitive.ObjectID `bson:"offer_id" json:"offer_id"`
	OfferName          string             `bson:"offer_name" json:"offer_name"`
	// OfferMembershipRank is unique identifier to specify this offer's value in the ranking system, higher is better.
	OfferMembershipRank int `bson:"offer_membership_rank" json:"offer_membership_rank"`
	// HasTimedLock controls whether this exercise has a user context specific lock
	// that will restrict access until the duration has elapsed for the lock to
	// unlock.
	HasTimedLock bool `bson:"has_timed_lock" json:"has_timed_lock"`
	// LockDuration is the duration from the user account creation date until
	// this exercise unlocks for the specific user account.
	TimedLock string `bson:"timed_lock" json:"timed_lock"`
	// TimedLockDurationValue is the actual golang value of the duration.
	TimedLockDuration time.Duration `bson:"timed_lock_duration" json:"timed_lock_duration"`
	// CurrentUserHasAccessGranted lets the user know if they can access.
	CurrentUserHasAccessGranted bool `bson:"-" json:"current_user_has_access_granted"`
	// CurrentUserLockedUntil is the date that in the future that this exercise becomes open to the authenticated user.
	CurrentUserLockedUntil time.Time `bson:"-" json:"current_user_locked_until"`
}

type VideoContentTag struct {
	ID             primitive.ObjectID `bson:"id" json:"id,omitempty"`
	OrganizationID primitive.ObjectID `bson:"organization_id" json:"organization_id,omitempty"`
	VideoContentID primitive.ObjectID `bson:"videocontent_id" json:"videocontent_id,omitempty"`
	TagID          uint64             `bson:"tag_id" json:"tag_id,omitempty"`
	TagText        string             `bson:"tag_text" json:"tag_text,omitempty"`
}

type VideoContentListFilter struct {
	// Pagination related.
	Cursor    primitive.ObjectID
	PageSize  int64
	SortField string
	SortOrder int8 // 1=ascending | -1=descending

	// Filter related.
	OfferID        primitive.ObjectID `bson:"offer_id" json:"offer_id"`
	CollectionID   primitive.ObjectID `bson:"collection_id" json:"collection_id,omitempty"`
	OrganizationID primitive.ObjectID `bson:"organization_id" json:"organization_id,omitempty"`
	SearchText     string             `json:"search_text"`
	Gender         string             `bson:"gender" json:"gender"`
	MovementType   int8               `bson:"movement_type" json:"movement_type"`
	CategoryID     primitive.ObjectID `bson:"category_id" json:"category_id,omitempty"`
	Status         int8               `bson:"status" json:"status"`
	VideoType      int8               `bson:"video_type" json:"video_type"`
}

type VideoContentListResult struct {
	Results     []*VideoContent    `json:"results"`
	NextCursor  primitive.ObjectID `json:"next_cursor"`
	HasNextPage bool               `json:"has_next_page"`
}

type VideoContentAsSelectOption struct {
	Value primitive.ObjectID `bson:"_id" json:"value"` // Extract from the database `_id` field and output through API as `value`.
	Label string             `bson:"name" json:"label"`
}

// VideoContentStorer Interface for user.
type VideoContentStorer interface {
	CheckIfExistsByID(ctx context.Context, id primitive.ObjectID) (bool, error)
	Create(ctx context.Context, m *VideoContent) error
	GetByID(ctx context.Context, id primitive.ObjectID) (*VideoContent, error)
	UpdateByID(ctx context.Context, m *VideoContent) error
	UpsertByID(ctx context.Context, e *VideoContent) error
	ListByFilter(ctx context.Context, f *VideoContentListFilter) (*VideoContentListResult, error)
	ListAsSelectOptionByFilter(ctx context.Context, f *VideoContentListFilter) ([]*VideoContentAsSelectOption, error)
	DeleteByID(ctx context.Context, id primitive.ObjectID) error
	// //TODO: Add more...
}

type VideoContentStorerImpl struct {
	Logger     *slog.Logger
	DbClient   *mongo.Client
	Collection *mongo.Collection
}

func NewDatastore(appCfg *c.Conf, loggerp *slog.Logger, client *mongo.Client) VideoContentStorer {
	// ctx := context.Background()
	uc := client.Database(appCfg.DB.Name).Collection("video_contents")

	// The following few lines of code will create the index for our app for this
	// colleciton.
	indexModel := mongo.IndexModel{
		Keys: bson.D{
			{"name", "text"},
			{"description", "text"},
			{"author_name", "text"},
			{"author_url", "text"},
			{"video_attachment_filename", "text"},
			{"thumbnail_attachment_filename", "text"},
			{"category_name", "text"},
			{"collection_name", "text"},
		},
	}
	_, err := uc.Indexes().CreateOne(context.TODO(), indexModel)
	if err != nil {
		// It is important that we crash the app on startup to meet the
		// requirements of `google/wire` framework.
		log.Fatal(err)
	}

	s := &VideoContentStorerImpl{
		Logger:     loggerp,
		DbClient:   client,
		Collection: uc,
	}
	return s
}
