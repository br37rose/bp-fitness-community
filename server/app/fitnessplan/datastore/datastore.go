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

type FitnessPlan struct {
	ID                 primitive.ObjectID     `bson:"_id" json:"id"`
	CreatedAt          time.Time              `bson:"created_at,omitempty" json:"created_at,omitempty"`
	CreatedByUserName  string                 `bson:"created_by_user_name" json:"created_by_user_name"`
	CreatedByUserID    primitive.ObjectID     `bson:"created_by_user_id" json:"created_by_user_id"`
	ModifiedAt         time.Time              `bson:"modified_at,omitempty" json:"modified_at,omitempty"`
	ModifiedByUserName string                 `bson:"modified_by_user_name" json:"modified_by_user_name"`
	ModifiedByUserID   primitive.ObjectID     `bson:"modified_by_user_id" json:"modified_by_user_id"`
	EstimatedReadyDate time.Time              `bson:"estimated_ready_date,omitempty" json:"estimated_ready_date,omitempty"`
	Name               string                 `bson:"name" json:"name"`
	Status             int8                   `bson:"status" json:"status"`
	WasProcessed       bool                   `bson:"was_processed" json:"was_processed"`
	OrganizationID     primitive.ObjectID     `bson:"organization_id,omitempty" json:"organization_id,omitempty"`
	OrganizationName   string                 `bson:"organization_name" json:"organization_name"`
	ExerciseNames      []string               `bson:"exercise_names" json:"exercise_names"`
	Exercises          []*FitnessPlanExercise `bson:"exercises" json:"exercises"`
	Instructions       string                 `bson:"instructions" json:"instructions"`
	Error              string                 `bson:"error" json:"error"`
	PromptOne          string                 `bson:"prompt_one" json:"-"`
	PromptTwo          string                 `bson:"prompt_two" json:"-"`
	UserID             primitive.ObjectID     `bson:"user_id,omitempty" json:"user_id,omitempty"`
	UserName           string                 `bson:"user_name" json:"user_name"`
	UserLexicalName    string                 `bson:"user_lexical_name" json:"user_lexical_name"`
	WeeklyFitnessPlans []*WeeklyFitnessPlan   `bson:"weekly_fitness_plans" json:"weekly_fitness_plans"`
	ThreadID           string                 `bson:"thread_id"  json:"-"`
	RunnerID           string                 `bson:"runner_id"  json:"-"`
}

type FitnessPlanExercise struct {
	ID           primitive.ObjectID `bson:"_id" json:"id"`
	Name         string             `bson:"name" json:"name"`
	VideoURL     string             `bson:"video_url" json:"video_url"`
	ThumbnailURL string             `bson:"thumbnail_url" json:"thumbnail_url"`
	Description  string             `bson:"description" json:"description"`
	VideoType    int8               `bson:"video_type" json:"video_type"`
	// DEVELOPERS NOTE: Add more of your fields here...
}

type DailyPlan struct {
	Title        string                `bson:"title" json:"title"`
	Instructions string                `bson:"instructions" json:"instructions"`
	PlanDetails  []FitnessPlanExercise `bson:"plan_details" json:"plan_details"`
}

type WeeklyFitnessPlan struct {
	Title      string      `bson:"title" json:"title"`
	DailyPlans []DailyPlan `bson:"daily_plans" json:"daily_plans"`
}

type FitnessPlanListFilter struct {
	// Pagination related.
	Cursor    primitive.ObjectID
	PageSize  int64
	SortField string
	SortOrder int8 // 1=ascending | -1=descending

	// Filter related.
	OrganizationID primitive.ObjectID
	UserID         primitive.ObjectID
	StatusList     []int8
	SearchText     string
	ExerciseNames  []string
}

type FitnessPlanListResult struct {
	Results     []*FitnessPlan     `json:"results"`
	NextCursor  primitive.ObjectID `json:"next_cursor"`
	HasNextPage bool               `json:"has_next_page"`
}

type FitnessPlanAsSelectOption struct {
	Value primitive.ObjectID `bson:"_id" json:"value"` // Extract from the database `_id` field and output through API as `value`.
	Label string             `bson:"name" json:"label"`
}

// FitnessPlanStorer Interface for video category.
type FitnessPlanStorer interface {
	Create(ctx context.Context, m *FitnessPlan) error
	GetByID(ctx context.Context, id primitive.ObjectID) (*FitnessPlan, error)
	UpdateByID(ctx context.Context, m *FitnessPlan) error
	ListByFilter(ctx context.Context, m *FitnessPlanListFilter) (*FitnessPlanListResult, error)
	ListAsSelectOptionByFilter(ctx context.Context, f *FitnessPlanListFilter) ([]*FitnessPlanAsSelectOption, error)
	DeleteByID(ctx context.Context, id primitive.ObjectID) error
	// //TODO: Add more...
}

type FitnessPlanStorerImpl struct {
	Logger     *slog.Logger
	DbClient   *mongo.Client
	Collection *mongo.Collection
}

func NewDatastore(appCfg *c.Conf, loggerp *slog.Logger, client *mongo.Client) FitnessPlanStorer {
	// ctx := context.Background()
	uc := client.Database(appCfg.DB.Name).Collection("fitness_plans")

	// The following few lines of code will create the index for our app for this
	// colleciton.
	indexModel := mongo.IndexModel{
		Keys: bson.D{
			{Key: "name", Value: "text"},
		},
	}
	_, err := uc.Indexes().CreateOne(context.TODO(), indexModel)
	if err != nil {
		// It is important that we crash the app on startup to meet the
		// requirements of `google/wire` framework.
		log.Fatal(err)
	}

	s := &FitnessPlanStorerImpl{
		Logger:     loggerp,
		DbClient:   client,
		Collection: uc,
	}
	return s
}
