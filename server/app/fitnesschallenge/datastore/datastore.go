package datastore

import (
	"context"
	"log"
	"log/slog"
	"time"

	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/config"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	FitnessChallengeStatusActive = iota + 1
	FitnessChallengeStatusArchived
)

var Challenges = map[int]ChallengeRules{
	1: {Name: "Calorie Crunch Challenge", Description: " Burn calorie challenge", Type: 1},
	2: {Name: "Step-Up Challenge", Description: "Step-Up Challenge based on number of steps taken", Type: 2},
	3: {Name: "Rest & Rise Challenge (Sleep Focus)", Description: " Rest & Rise Challenge (Sleep Focus)", Type: 3},
	4: {Name: "Cardio King Challenge (VO2 Max Focus)", Description: " Cardio King Challenge (VO2 Max Focus)", Type: 4},
	5: {Name: "Swim Stride Challenge", Description: " Swim Stride Challenge", Type: 5},
	6: {Name: "Heartbeat Hero Challenge", Description: " Heartbeat Hero Challenge", Type: 6},
	7: {Name: "Oxygen Optimize Challenge (SpO2 Focus)", Description: " Oxygen Optimize Challenge (SpO2 Focus)", Type: 7},
}

type FitnessChallenge struct {
	ID                 primitive.ObjectID   `json:"id,omitempty" bson:"_id,omitempty"`
	UserIDs            []primitive.ObjectID `json:"user_ids,omitempty" bson:"user_ids"`
	UserNames          []string             `json:"user_names,omitempty" bson:"user_names"`
	OrganizationID     primitive.ObjectID   `bson:"organization_id" json:"organization_id,omitempty"`
	Name               string               `json:"name,omitempty" bson:"name,omitempty"`
	Description        string               `json:"description,omitempty" bson:"description,omitempty"`
	DurationInWeeks    int64                `json:"duration_in_weeks,omitempty" bson:"duration_in_weeks,omitempty"`
	Rules              []*ChallengeRules    `json:"rules,omitempty" bson:"rules"`
	StartTime          time.Time            `json:"start_time,omitempty" bson:"start_time,omitempty"`
	EndTime            time.Time            `json:"end_time,omitempty" bson:"end_time,omitempty"`
	CreatedAt          time.Time            `bson:"created_at" json:"created_at,omitempty"`
	CreatedByUserID    primitive.ObjectID   `bson:"created_by_user_id" json:"created_by_user_id"`
	CreatedByUserName  string               `json:"created_by_user_name,omitempty" bson:"created_by_user_name,omitempty"`
	ModifiedByUserID   primitive.ObjectID   `json:"modified_by_user_id,omitempty" bson:"modified_by_user_id,omitempty"`
	ModifiedByUserName string               `json:"modified_by_user_name,omitempty" bson:"modified_by_user_name,omitempty"`
	ModifiedAt         time.Time            `bson:"modified_at" json:"modified_at,omitempty"`
	Status             int8                 `bson:"status" json:"status"`
}
type ChallengeRules struct {
	Name        string `json:"name,omitempty" bson:"name,omitempty"`
	Description string `json:"description,omitempty" bson:"description,omitempty"`
	Type        int8   `json:"type,omitempty" bson:"type,omitempty"`
	Target      string `json:"target,omitempty" bson:"target,omitempty"`
}

type FitnessChallengeStorerImpl struct {
	Logger     *slog.Logger
	DbClient   *mongo.Client
	Collection *mongo.Collection
}

type FitnessChallengeListFilter struct {
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

type FitnessChallengeListResult struct {
	Results     []*FitnessChallenge `json:"results"`
	NextCursor  primitive.ObjectID  `json:"next_cursor"`
	HasNextPage bool                `json:"has_next_page"`
}

type FitnessChallengeStorer interface {
	Create(ctx context.Context, tp *FitnessChallenge) error
	GetByID(ctx context.Context, id primitive.ObjectID) (*FitnessChallenge, error)
	UpdateByID(ctx context.Context, tp *FitnessChallenge) error
	Delete(ctx context.Context, id primitive.ObjectID) error
	ListByFilter(ctx context.Context, f *FitnessChallengeListFilter) (*FitnessChallengeListResult, error)
}

func NewDatastore(appCfg *config.Conf, loggerp *slog.Logger, client *mongo.Client) FitnessChallengeStorer {

	uc := client.Database(appCfg.DB.Name).Collection("fitness_challenge")
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

	s := &FitnessChallengeStorerImpl{
		Logger:     loggerp,
		DbClient:   client,
		Collection: uc,
	}
	return s
}
