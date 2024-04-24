package datastore

import (
	"context"
	"log/slog"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	c "github.com/bci-innovation-labs/bp8fitnesscommunity-cli/config"
)

const (
	ExerciseStatusActive             = 1
	ExerciseStatusArchived           = 2
	ExerciseVideoTypeS3              = 1
	ExerciseVideoTypeYouTube         = 2
	ExerciseVideoTypeVimeo           = 3
	ExerciseThumbnailTypeS3          = 1
	ExerciseThumbnailTypeExternalURL = 2
	ExerciseThumbnailTypeLocal       = 3
	ExerciseTypeSystem               = 1
	ExerciseTypeCustom               = 2
)

type Exercise struct {
	OrganizationID              primitive.ObjectID `bson:"organization_id" json:"organization_id,omitempty"`
	OrganizationName            string             `bson:"organization_name" json:"organization_name"`
	ParentID                    primitive.ObjectID `bson:"parent_id" json:"parent_id"`
	ID                          primitive.ObjectID `bson:"_id" json:"id"`
	Type                        int8               `bson:"type" json:"type"`
	Name                        string             `bson:"name" json:"name"`
	AlternateName               string             `bson:"alternate_name" json:"alternate_name"`
	Gender                      string             `bson:"gender" json:"gender"`
	MovementType                int8               `bson:"movement_type" json:"movement_type"`
	Category                    int8               `bson:"category" json:"category"`
	Description                 string             `bson:"description" json:"description"`
	VideoType                   int8               `bson:"video_type" json:"video_type"`
	VideoAttachmentID           primitive.ObjectID `bson:"video_attachment_id" json:"video_attachment_id"`
	VideoAttachmentFilename     string             `bson:"video_attachment_filename" json:"video_attachment_filename"`
	VideoObjectKey              string             `bson:"video_object_key" json:"video_object_key"`
	VideoObjectURL              string             `bson:"video_object_url" json:"video_object_url"`
	VideoObjectExpiry           time.Time          `bson:"video_object_expiry" json:"video_object_expiry"`
	VideoName                   string             `bson:"video_name" json:"video_name"`
	VideoURL                    string             `bson:"video_url" json:"video_url"`
	ThumbnailType               int8               `bson:"thumbnail_type" json:"thumbnail_type"`
	ThumbnailAttachmentID       primitive.ObjectID `bson:"thumbnail_attachment_id" json:"thumbnail_attachment_id"`
	ThumbnailAttachmentFilename string             `bson:"thumbnail_attachment_filename" json:"thumbnail_attachment_filename"`
	ThumbnailObjectKey          string             `bson:"thumbnail_object_key" json:"thumbnail_object_key"`
	ThumbnailObjectURL          string             `bson:"thumbnail_object_url" json:"thumbnail_object_url"`
	ThumbnailObjectExpiry       time.Time          `bson:"thumbnail_object_expiry" json:"thumbnail_object_expiry"`
	ThumbnailName               string             `bson:"thumbnail_name" json:"thumbnail_name"`
	ThumbnailURL                string             `bson:"thumbnail_url" json:"thumbnail_url"`
	Tags                        []*ExerciseTag     `bson:"tags" json:"tags"`
	Status                      int8               `bson:"status" json:"status"`
	CreatedAt                   time.Time          `bson:"created_at" json:"created_at,omitempty"`
	CreatedByUserID             primitive.ObjectID `bson:"created_by_user_id" json:"created_by_user_id"`
	CreatedByUserName           string             `bson:"created_by_user_name" json:"created_by_user_name"`
	ModifiedAt                  time.Time          `bson:"modified_at" json:"modified_at,omitempty"`
	ModifiedByUserID            primitive.ObjectID `bson:"modified_by_user_id" json:"modified_by_user_id"`
	ModifiedByUserName          string             `bson:"modified_by_user_name" json:"modified_by_user_name"`
	HasMonetization             bool               `bson:"has_monetization" json:"has_monetization"`
	OfferID                     primitive.ObjectID `bson:"offer_id" json:"offer_id"`
	OfferName                   string             `bson:"offer_name" json:"offer_name"`
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
	CurrentUserLockedUntil time.Time            `bson:"-" json:"current_user_locked_until"`
	Equipments             []*ExerciseEquipment `bson:"equipments" json:"equipments"`
}

type ExerciseTag struct {
	ID             primitive.ObjectID `bson:"id" json:"id,omitempty"` // References `_id` from inside the `tags` table.
	OrganizationID primitive.ObjectID `bson:"organization_id" json:"organization_id,omitempty"`
	Text           string             `bson:"text" json:"text,omitempty"`
	Status         int8               `bson:"status" json:"status"`
}

type ExerciseEquipment struct {
	ID     primitive.ObjectID `bson:"_id" json:"id"` // Note: ID refers to the ID within `Equipment` table.
	Name   string             `bson:"name" json:"name"`
	No     int8               `bson:"no" json:"no"`
	Status int8               `bson:"status" json:"status"`
}

type ExerciseListFilter struct {
	// Pagination related.
	Cursor    primitive.ObjectID
	PageSize  int64
	SortField string
	SortOrder int8 // 1=ascending | -1=descending

	// Filter related.
	OfferID         primitive.ObjectID
	OrganizationID  primitive.ObjectID
	ExcludeArchived bool
	SearchText      string
	Gender          string
	MovementType    int8
	Category        int8
	Status          int8
	VideoType       int8
	Names           []string
}

type ExerciseListResult struct {
	Results     []*Exercise        `json:"results"`
	NextCursor  primitive.ObjectID `json:"next_cursor"`
	HasNextPage bool               `json:"has_next_page"`
}

type ExerciseAsSelectOption struct {
	Value primitive.ObjectID `bson:"_id" json:"value"` // Extract from the database `_id` field and output through API as `value`.
	Label string             `bson:"name" json:"label"`
}

// ExerciseStorer Interface for user.
type ExerciseStorer interface {
	CheckIfExistsByID(ctx context.Context, id primitive.ObjectID) (bool, error)
	Create(ctx context.Context, m *Exercise) error
	GetByID(ctx context.Context, id primitive.ObjectID) (*Exercise, error)
	GetByVideoURL(ctx context.Context, videoURL string) (*Exercise, error)
	UpdateByID(ctx context.Context, m *Exercise) error
	UpsertByID(ctx context.Context, e *Exercise) error
	ListByFilter(ctx context.Context, f *ExerciseListFilter) (*ExerciseListResult, error)
	ListAsSelectOptionByFilter(ctx context.Context, f *ExerciseListFilter) ([]*ExerciseAsSelectOption, error)
	DeleteByID(ctx context.Context, id primitive.ObjectID) error
	// //TODO: Add more...
}

type ExerciseStorerImpl struct {
	Logger     *slog.Logger
	DbClient   *mongo.Client
	Collection *mongo.Collection
}

func NewDatastore(appCfg *c.Conf, loggerp *slog.Logger, client *mongo.Client) ExerciseStorer {
	// ctx := context.Background()
	uc := client.Database(appCfg.DB.Name).Collection("exercises")

	// // The following few lines of code will create the index for our app for this
	// // colleciton.
	// indexModel := mongo.IndexModel{
	// 	Keys: bson.D{
	// 		{"name", "text"},
	// 		{"description", "text"},
	// 	},
	// }
	// _, err := uc.Indexes().CreateOne(context.TODO(), indexModel)
	// if err != nil {
	// 	// It is important that we crash the app on startup to meet the
	// 	// requirements of `google/wire` framework.
	// 	log.Fatal(err)
	// }

	s := &ExerciseStorerImpl{
		Logger:     loggerp,
		DbClient:   client,
		Collection: uc,
	}
	return s
}
