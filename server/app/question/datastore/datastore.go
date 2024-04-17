package datastore

import (
	"context"
	"log"
	"log/slog"
	"time"

	c "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/config"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	QuestionStatusActive = iota + 1
	QuestionStatusArchive
)

type Question struct {
	ID            primitive.ObjectID `bson:"_id" json:"id"`
	Question      string             `bson:"question" json:"question"`
	IsMultiSelect bool               `bson:"is_multiselect" json:"is_multiselect"`
	Content       []QuestionContent  `bson:"content" json:"content"`
	Status        bool               `bson:"status" json:"status"`

	CreatedAt          time.Time          `bson:"created_at" json:"created_at,omitempty"`
	CreatedByUserID    primitive.ObjectID `bson:"created_by_user_id" json:"created_by_user_id"`
	CreatedByUserName  string             `bson:"created_by_user_name" json:"created_by_user_name"`
	ModifiedAt         time.Time          `bson:"modified_at" json:"modified_at,omitempty"`
	ModifiedByUserID   primitive.ObjectID `bson:"modified_by_user_id" json:"modified_by_user_id"`
	ModifiedByUserName string             `bson:"modified_by_user_name" json:"modified_by_user_name"`
}

type QuestionContent struct {
	Title    string   `json:"title"`
	Subtitle string   `json:"subtitle,omitempty"`
	Options  []string `json:"options,omitempty"`
}

type QuestionListFilter struct {
	// Pagination related.
	Cursor    primitive.ObjectID
	PageSize  int64
	SortField string
	SortOrder int8 // 1=ascending | -1=descending

	IsMultiSelect bool // Filter by whether the question allows multi-select
	Status        []int8
}

type QuestionListResult struct {
	Results     []*Question        `json:"results"`
	NextCursor  primitive.ObjectID `json:"next_cursor"`
	HasNextPage bool               `json:"has_next_page"`
}

type QuestionStorer interface {
	Create(ctx context.Context, question *Question) error
	GetByID(ctx context.Context, id primitive.ObjectID) (*Question, error)
	UpdateByID(ctx context.Context, question *Question) error
	DeleteByID(ctx context.Context, id primitive.ObjectID) error
	ListByFilter(ctx context.Context, f *QuestionListFilter) (*QuestionListResult, error)
}

type QuestionStorerImpl struct {
	Logger     *slog.Logger
	DbClient   *mongo.Client
	Collection *mongo.Collection
}

func NewDatastore(appCfg *c.Conf, loggerp *slog.Logger, client *mongo.Client) QuestionStorer {
	// ctx := context.Background()
	uc := client.Database(appCfg.DB.Name).Collection("questions")

	// The following few lines of code will create the index for our app for this
	// colleciton.
	indexModel := mongo.IndexModel{
		Keys: bson.D{
			{Key: "content.title", Value: "text"},
		},
	}
	_, err := uc.Indexes().CreateOne(context.TODO(), indexModel)
	if err != nil {
		log.Fatal(err)
	}

	s := &QuestionStorerImpl{
		Logger:     loggerp,
		DbClient:   client,
		Collection: uc,
	}

	return s
}
