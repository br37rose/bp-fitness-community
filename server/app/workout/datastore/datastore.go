package datastore

import (
	"context"
	"log"
	"log/slog"
	"time"

	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/exercise/datastore"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/config"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	WorkoutStatusActive    = 1
	WorkoutStatusArchived  = 2
	WorkoutVisibileToAll   = 1
	WorkoutPersonalVisible = 2
)

type Workout struct {
	ID                        primitive.ObjectID `json:"id" bson:"_id"`
	UserId                    primitive.ObjectID `json:"user_id" bson:"user_id"`
	UserName                  string             `json:"user_name" bson:"user_name"`
	Name                      string             `json:"name" bson:"name"`
	Description               string             `json:"description"`
	Type                      int8               `json:"type" bson:"type"`
	Status                    int8               `json:"status" bson:"status"`
	Visibility                int8               `json:"visibility"`
	WorkoutExercises          []*WorkoutExercise `json:"workout_exercises,omitempty" bson:"workout_exercises,omitempty"`
	WorkoutExerciseTimeInMins int64              `json:"workout_exercise_time_in_mins" bson:"workout_exercise_time_in_mins"`
	CreatedAt                 time.Time          `bson:"created_at" json:"created_at,omitempty"`
	CreatedByUserID           primitive.ObjectID `bson:"created_by_user_id" json:"created_by_user_id"`
	CreatedByUserName         string             `json:"created_by_user_name,omitempty" bson:"created_by_user_name,omitempty"`
	ModifiedByUserID          primitive.ObjectID `json:"modified_by_user_id,omitempty" bson:"modified_by_user_id,omitempty"`
	ModifiedByUserName        string             `json:"modified_by_user_name,omitempty" bson:"modified_by_user_name,omitempty"`
	ModifiedAt                time.Time          `bson:"modified_at" json:"modified_at,omitempty"`
}
type WorkoutExercise struct {
	ID               primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	ExerciseID       primitive.ObjectID `json:"exercise_id,omitempty" bson:"exercise_id,omitempty"`
	ExerciseName     string             `json:"exercise_name,omitempty" bson:"exercise_name,omitempty"`
	IsRest           bool               `json:"is_rest" bson:"is_rest"`
	OrderNumber      int64              `json:"order_number,omitempty" bson:"order_number,omitempty"`
	Sets             int64              `json:"sets,omitempty" bson:"sets,omitempty"`
	Type             int8               `bson:"type" json:"type"`
	TargetTimeInSecs int64              `json:"target_time_in_secs,omitempty" bson:"target_time_in_secs,omitempty"`
	RestPeriodInSecs int64              `json:"rest_period_in_secs,omitempty" bson:"rest_period_in_secs,omitempty"`
	TargetText       string             `json:"target_text,omitempty" bson:"target_text,omitempty"`
	CreatedAt        time.Time          `json:"created_at,omitempty" bson:"created_at,omitempty"`
	ModifiedAt       time.Time          `json:"modified_at,omitempty" bson:"modified_at,omitempty"`
	Excercise        datastore.Exercise `json:"excercise,omitempty" bson:"excercise,omitempty"`
}

type WorkoutListFilter struct {
	Cursor    primitive.ObjectID
	PageSize  int64
	SortField string
	SortOrder int8 // 1=ascending | -1=descending

	StatusList      []int8
	Visibility      int8
	CreatedByUserID primitive.ObjectID
	Types           []int8
	ExcludeArchived bool
	SearchText      string
	UserId          primitive.ObjectID
	GetExcercise    bool
}

type WorkouStorerImpl struct {
	Logger     *slog.Logger
	DbClient   *mongo.Client
	Collection *mongo.Collection
}

type WorkoutistResult struct {
	Results     []*Workout         `json:"results"`
	NextCursor  primitive.ObjectID `json:"next_cursor"`
	HasNextPage bool               `json:"has_next_page"`
}

type WorkoutStorer interface {
	Create(ctx context.Context, w *Workout) error
	GetByID(ctx context.Context, id primitive.ObjectID) (*Workout, error)
	UpdateByID(ctx context.Context, m *Workout) error
	DeleteByID(ctx context.Context, id primitive.ObjectID) error
	ListByFilter(ctx context.Context, f *WorkoutListFilter) (*WorkoutistResult, error)
}

func NewDatastore(appCfg *config.Conf, loggerp *slog.Logger, client *mongo.Client) WorkoutStorer {

	uc := client.Database(appCfg.DB.Name).Collection("workouts")
	indexModel := mongo.IndexModel{
		Keys: bson.D{
			{Key: "name", Value: "text"},
			{Key: "description", Value: "text"},
		},
	}
	_, err := uc.Indexes().CreateOne(context.TODO(), indexModel)
	if err != nil {
		// It is important that we crash the app on startup to meet the
		// requirements of `google/wire` framework.
		log.Fatal(err)
	}

	s := &WorkouStorerImpl{
		Logger:     loggerp,
		DbClient:   client,
		Collection: uc,
	}
	return s
}
