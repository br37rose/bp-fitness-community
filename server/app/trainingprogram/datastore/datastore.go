package datastore

import (
	"context"
	"log"
	"log/slog"
	"time"

	w_s "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/workout/datastore"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/config"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	TrainingProgramStatusActive = iota + 1
	TrainingProgramStatusArchived
)

type TrainingProgram struct {
	ID                 primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	UserID             primitive.ObjectID `json:"user_id,omitempty" bson:"user_id,omitempty"`
	UserName           string             `json:"user_name,omitempty" bson:"user_name,omitempty"`
	OrganizationID     primitive.ObjectID `bson:"organization_id" json:"organization_id,omitempty"`
	Name               string             `json:"name,omitempty" bson:"name,omitempty"`
	Description        string             `json:"description,omitempty" bson:"description,omitempty"`
	Type               int8               `json:"type,omitempty" bson:"type,omitempty"`
	Phases             int64              `json:"phases,omitempty" bson:"phases,omitempty"`
	Weeks              int64              `json:"weeks,omitempty" bson:"weeks,omitempty"`
	DurationInWeeks    int64              `json:"duration_in_weeks,omitempty" bson:"duration_in_weeks,omitempty"`
	TrainingPhases     []*TrainingPhase   `json:"training_phases,omitempty" bson:"training_phases,omitempty"`
	StartTime          time.Time          `json:"start_time,omitempty" bson:"start_time,omitempty"`
	EndTime            time.Time          `json:"end_time,omitempty" bson:"end_time,omitempty"`
	CreatedAt          time.Time          `bson:"created_at" json:"created_at,omitempty"`
	CreatedByUserID    primitive.ObjectID `bson:"created_by_user_id" json:"created_by_user_id"`
	CreatedByUserName  string             `json:"created_by_user_name,omitempty" bson:"created_by_user_name,omitempty"`
	ModifiedByUserID   primitive.ObjectID `json:"modified_by_user_id,omitempty" bson:"modified_by_user_id,omitempty"`
	ModifiedByUserName string             `json:"modified_by_user_name,omitempty" bson:"modified_by_user_name,omitempty"`
	ModifiedAt         time.Time          `bson:"modified_at" json:"modified_at,omitempty"`
	Status             int8               `bson:"status" json:"status"`
}
type TrainingPhase struct {
	ID               primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name             string             `json:"name,omitempty" bson:"name,omitempty"`
	Description      string             `json:"description,omitempty" bson:"description,omitempty"`
	Phase            int64              `json:"phase,omitempty" bson:"phase,omitempty"`
	StartWeek        int64              `json:"start_week,omitempty" bson:"start_week,omitempty"`
	EndWeek          int64              `json:"end_week,omitempty" bson:"end_week,omitempty"`
	Type             int8               `json:"type,omitempty" bson:"type,omitempty"`
	TrainingRoutines []*TrainingRoutine `json:"training_routines,omitempty" bson:"training_routines,omitempty"`
	StartTime        time.Time          `json:"start_time,omitempty" bson:"start_time,omitempty"`
	EndTime          time.Time          `json:"end_time,omitempty" bson:"end_time,omitempty"`
}
type TrainingRoutine struct {
	ID           primitive.ObjectID `json:"id" bson:"_id"`
	WorkoutID    primitive.ObjectID `json:"workout_id" bson:"workout_id"`
	Workout      *w_s.Workout       `json:"workout,omitempty" bson:"workout,omitempty"`
	OrderNumber  int64              `json:"order_number" bson:"order_number"`
	Phase        int64              `json:"phase" bson:"phase"`
	TrainingDays []*TrainingDay     `json:"training_days,omitempty" bson:"training_days,omitempty"`
	Type         int8               `json:"type,omitempty" bson:"type,omitempty"`
}

// type Workout struct {
// 	ID                        primitive.ObjectID `json:"id" bson:"_id"`
// 	Name                      string             `json:"name" bson:"name"`
// 	Type                      int8               `json:"type" bson:"type"`
// 	Status                    int8               `json:"status" bson:"status"`
// 	WorkoutExercises          []*WorkoutExercise `json:"workout_exercises,omitempty" bson:"workout_exercises,omitempty"`
// 	WorkoutExerciseTimeInMins int64              `json:"workout_exercise_time_in_mins" bson:"workout_exercise_time_in_mins"`
// 	CreatedAt                 time.Time          `bson:"created_at" json:"created_at,omitempty"`
// 	CreatedByUserID           primitive.ObjectID `bson:"created_by_user_id" json:"created_by_user_id"`
// 	CreatedByUserName         string             `json:"created_by_user_name,omitempty" bson:"created_by_user_name,omitempty"`
// 	ModifiedByUserID          primitive.ObjectID `json:"modified_by_user_id,omitempty" bson:"modified_by_user_id,omitempty"`
// 	ModifiedByUserName        string             `json:"modified_by_user_name,omitempty" bson:"modified_by_user_name,omitempty"`
// 	ModifiedAt                time.Time          `bson:"modified_at" json:"modified_at,omitempty"`
// }
// type WorkoutExercise struct {
// 	ID               primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
// 	ExerciseID       primitive.ObjectID `json:"exercise_id,omitempty" bson:"exercise_id,omitempty"`
// 	ExerciseName     string             `json:"exercise_name,omitempty" bson:"exercise_name,omitempty"`
// 	IsRest           bool               `json:"is_rest" bson:"is_rest"`
// 	OrderNumber      int64              `json:"order_number,omitempty" bson:"order_number,omitempty"`
// 	Sets             int64              `json:"sets,omitempty" bson:"sets,omitempty"`
// 	Type             int8               `bson:"type" json:"type"`
// 	TargetTimeInSecs int64              `json:"target_time_in_secs,omitempty" bson:"target_time_in_secs,omitempty"`
// 	RestPeriodInSecs int64              `json:"rest_period_in_secs,omitempty" bson:"rest_period_in_secs,omitempty"`
// 	CreatedAt        time.Time          `json:"created_at,omitempty" bson:"created_at,omitempty"`
// 	ModifiedAt       time.Time          `json:"modified_at,omitempty" bson:"modified_at,omitempty"`
// }

type TrainingDay struct {
	Week int64      `json:"week,omitempty" bson:"week,omitempty"`
	Day  int64      `json:"day,omitempty" bson:"day,omitempty"`
	Time *time.Time `json:"time,omitempty" bson:"time,omitempty"`
	Type int8       `json:"type,omitempty" bson:"type,omitempty"`
}

type TrainingProgramStorerImpl struct {
	Logger     *slog.Logger
	DbClient   *mongo.Client
	Collection *mongo.Collection
}

type TrainingProgramListFilter struct {
	Cursor          primitive.ObjectID
	PageSize        int64
	SortField       string
	SortOrder       int8
	UserID          primitive.ObjectID
	OrganizationID  primitive.ObjectID
	Name            string
	Phases          int64
	Weeks           int64
	DurationInWeeks int64
	StartTime       time.Time
	EndTime         time.Time
	StatusList      []int8
	SearchText      string
}

type TrainingProgramListResult struct {
	Results     []*TrainingProgram `json:"results"`
	NextCursor  primitive.ObjectID `json:"next_cursor"`
	HasNextPage bool               `json:"has_next_page"`
}

type TrainingProgramStorer interface {
	Create(ctx context.Context, tp *TrainingProgram) error
	GetByID(ctx context.Context, id primitive.ObjectID) (*TrainingProgram, error)
	UpdateByID(ctx context.Context, tp *TrainingProgram) error
	UpdatePhase(ctx context.Context, tpId primitive.ObjectID, phase []*TrainingPhase) error
	Delete(ctx context.Context, id primitive.ObjectID) error
	ListByFilter(ctx context.Context, f *TrainingProgramListFilter) (*TrainingProgramListResult, error)
}

func NewDatastore(appCfg *config.Conf, loggerp *slog.Logger, client *mongo.Client) TrainingProgramStorer {

	uc := client.Database(appCfg.DB.Name).Collection("training_program")
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

	s := &TrainingProgramStorerImpl{
		Logger:     loggerp,
		DbClient:   client,
		Collection: uc,
	}
	return s
}
