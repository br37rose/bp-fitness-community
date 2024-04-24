package datastore

import (
	"context"
	"log"
	"log/slog"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	c "github.com/bci-innovation-labs/bp8fitnesscommunity-cli/config"
)

const (
	TagStatusActive   = 1
	TagStatusArchived = 2
)

type Tag struct {
	Text                  string             `bson:"text" json:"text"`
	Description           string             `bson:"description" json:"description"`
	Status                int8               `bson:"status" json:"status"`
	ID                    primitive.ObjectID `bson:"_id" json:"id"`
	OrganizationID        primitive.ObjectID `bson:"organization_id" json:"organization_id,omitempty"`
	OrganizationName      string             `bson:"organization_name" json:"organization_name,omitempty"`
	PublicID              uint64             `bson:"public_id" json:"public_id"`
	CreatedAt             time.Time          `bson:"created_at" json:"created_at"`
	CreatedByUserID       primitive.ObjectID `bson:"created_by_user_id" json:"created_by_user_id,omitempty"`
	CreatedByUserName     string             `bson:"created_by_user_name" json:"created_by_user_name"`
	CreatedFromIPAddress  string             `bson:"created_from_ip_address" json:"created_from_ip_address"`
	ModifiedAt            time.Time          `bson:"modified_at" json:"modified_at"`
	ModifiedByUserID      primitive.ObjectID `bson:"modified_by_user_id" json:"modified_by_user_id,omitempty"`
	ModifiedByUserName    string             `bson:"modified_by_user_name" json:"modified_by_user_name"`
	ModifiedFromIPAddress string             `bson:"modified_from_ip_address" json:"modified_from_ip_address"`
}

type TagListFilter struct {
	// Pagination related.
	Cursor    primitive.ObjectID
	PageSize  int64
	SortField string
	SortOrder int8 // 1=ascending | -1=descending

	// Filter related.
	OrganizationID  primitive.ObjectID
	Status          int8
	ExcludeArchived bool
	SearchText      string
}

type TagListResult struct {
	Results     []*Tag             `json:"results"`
	NextCursor  primitive.ObjectID `json:"next_cursor"`
	HasNextPage bool               `json:"has_next_page"`
}

type TagAsSelectOption struct {
	Value primitive.ObjectID `bson:"_id" json:"value"` // Extract from the database `_id` field and output through API as `value`.
	Label string             `bson:"text" json:"label"`
}

// TagStorer Interface for user.
type TagStorer interface {
	Create(ctx context.Context, m *Tag) error
	GetByID(ctx context.Context, id primitive.ObjectID) (*Tag, error)
	GetByPublicID(ctx context.Context, oldID uint64) (*Tag, error)
	GetByEmail(ctx context.Context, email string) (*Tag, error)
	GetByVerificationCode(ctx context.Context, verificationCode string) (*Tag, error)
	GetLatestByOrganizationID(ctx context.Context, organizationID primitive.ObjectID) (*Tag, error)
	CheckIfExistsByEmail(ctx context.Context, email string) (bool, error)
	UpdateByID(ctx context.Context, m *Tag) error
	ListByFilter(ctx context.Context, f *TagPaginationListFilter) (*TagPaginationListResult, error)
	ListAsSelectOptionByFilter(ctx context.Context, f *TagListFilter) ([]*TagAsSelectOption, error)
	DeleteByID(ctx context.Context, id primitive.ObjectID) error
}

type TagStorerImpl struct {
	Logger     *slog.Logger
	DbClient   *mongo.Client
	Collection *mongo.Collection
}

func NewDatastore(appCfg *c.Conf, loggerp *slog.Logger, client *mongo.Client) TagStorer {
	// ctx := context.Background()
	uc := client.Database(appCfg.DB.Name).Collection("tags")

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
		{Keys: bson.D{{Key: "organization_id", Value: 1}}},
		{Keys: bson.D{{Key: "public_id", Value: -1}}},
		{Keys: bson.D{{Key: "status", Value: 1}}},
		{Keys: bson.D{{Key: "text", Value: 1}}},
		{Keys: bson.D{{Key: "created_at", Value: 1}}},
		{Keys: bson.D{
			{"text", "text"},
			{"description", "text"},
		}},
	})
	if err != nil {
		// It is important that we crash the app on startup to meet the
		// requirements of `google/wire` framework.
		log.Fatal(err)
	}

	s := &TagStorerImpl{
		Logger:     loggerp,
		DbClient:   client,
		Collection: uc,
	}
	return s
}
