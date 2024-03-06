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

type NutritionPlan struct {
	ID                         primitive.ObjectID `bson:"_id" json:"id"`
	CreatedAt                  time.Time          `bson:"created_at,omitempty" json:"created_at,omitempty"`
	CreatedByUserName          string             `bson:"created_by_user_name" json:"created_by_user_name"`
	CreatedByUserID            primitive.ObjectID `bson:"created_by_user_id" json:"created_by_user_id"`
	ModifiedAt                 time.Time          `bson:"modified_at,omitempty" json:"modified_at,omitempty"`
	ModifiedByUserName         string             `bson:"modified_by_user_name" json:"modified_by_user_name"`
	ModifiedByUserID           primitive.ObjectID `bson:"modified_by_user_id" json:"modified_by_user_id"`
	HasAllergies               int8               `bson:"has_allergies" json:"has_allergies"`
	Allergies                  string             `bson:"allergies" json:"allergies"`
	Birthday                   time.Time          `bson:"birthday,omitempty" json:"birthday,omitempty"`
	MealsPerDay                int8               `bson:"meals_per_day" json:"meals_per_day"`
	ConsumeJunkFood            int8               `bson:"consume_junk_food" json:"consume_junk_food"`
	ConsumeFruitsAndVegetables int8               `bson:"consume_fruits_and_vegetables" json:"consume_fruits_and_vegetables"`
	HasIntermittentFasting     int8               `bson:"has_intermittent_fasting" json:"has_intermittent_fasting"`
	HeightFeet                 float64            `bson:"height_feet" json:"height_feet"`
	HeightFeetInches           float64            `bson:"height_feet_inches" json:"height_feet_inches"`
	HeightInches               float64            `bson:"height_inches" json:"height_inches"`
	Weight                     float64            `bson:"weight" json:"weight"`
	Gender                     int8               `bson:"gender" json:"gender"`
	GenderOther                string             `bson:"gender_other" json:"gender_other"`
	IdealWeight                float64            `bson:"ideal_weight" json:"ideal_weight"`
	PhysicalActivity           int8               `bson:"physical_activity" json:"physical_activity"`
	Status                     int8               `bson:"status" json:"status"`
	WorkoutIntensity           int8               `bson:"workout_intensity" json:"workout_intensity"`
	Goals                      []int8             `bson:"goals" json:"goals"`
	MaxWeeks                   int8               `bson:"max_weeks" json:"max_weeks"`
	Name                       string             `bson:"name" json:"name"`
	EstimatedReadyDate         time.Time          `bson:"estimated_ready_date,omitempty" json:"estimated_ready_date,omitempty"`
	WasProcessed               bool               `bson:"was_processed" json:"was_processed"`
	OrganizationID             primitive.ObjectID `bson:"organization_id,omitempty" json:"organization_id,omitempty"`
	OrganizationName           string             `bson:"organization_name" json:"organization_name"`
	Instructions               string             `bson:"instructions" json:"instructions"`
	Error                      string             `bson:"error" json:"error"`
	Prompt                     string             `bson:"prompt" json:"-"`
	UserID                     primitive.ObjectID `bson:"user_id,omitempty" json:"user_id,omitempty"`
	UserName                   string             `bson:"user_name" json:"user_name"`
	UserLexicalName            string             `bson:"user_lexical_name" json:"user_lexical_name"`
}

type NutritionPlanExercise struct {
	ID   primitive.ObjectID `bson:"_id" json:"id"`
	Name string             `bson:"name" json:"name"`
	// DEVELOPERS NOTE: Add more of your fields here...
}

type NutritionPlanListFilter struct {
	// Pagination related.
	Cursor    primitive.ObjectID
	PageSize  int64
	SortField string
	SortOrder int8 // 1=ascending | -1=descending

	// Filter related.
	OrganizationID primitive.ObjectID
	UserID         primitive.ObjectID
	Status         int8
	SearchText     string
	ExerciseNames  []string
}

type NutritionPlanListResult struct {
	Results     []*NutritionPlan   `json:"results"`
	NextCursor  primitive.ObjectID `json:"next_cursor"`
	HasNextPage bool               `json:"has_next_page"`
}

type NutritionPlanAsSelectOption struct {
	Value primitive.ObjectID `bson:"_id" json:"value"` // Extract from the database `_id` field and output through API as `value`.
	Label string             `bson:"name" json:"label"`
}

// NutritionPlanStorer Interface for video category.
type NutritionPlanStorer interface {
	Create(ctx context.Context, m *NutritionPlan) error
	GetByID(ctx context.Context, id primitive.ObjectID) (*NutritionPlan, error)
	UpdateByID(ctx context.Context, m *NutritionPlan) error
	ListByFilter(ctx context.Context, m *NutritionPlanListFilter) (*NutritionPlanListResult, error)
	ListAsSelectOptionByFilter(ctx context.Context, f *NutritionPlanListFilter) ([]*NutritionPlanAsSelectOption, error)
	DeleteByID(ctx context.Context, id primitive.ObjectID) error
	// //TODO: Add more...
}

type NutritionPlanStorerImpl struct {
	Logger     *slog.Logger
	DbClient   *mongo.Client
	Collection *mongo.Collection
}

func NewDatastore(appCfg *c.Conf, loggerp *slog.Logger, client *mongo.Client) NutritionPlanStorer {
	// ctx := context.Background()
	uc := client.Database(appCfg.DB.Name).Collection("nutrition_plans")

	// The following few lines of code will create the index for our app for this
	// colleciton.
	indexModel := mongo.IndexModel{
		Keys: bson.D{
			{"name", "text"},
		},
	}
	_, err := uc.Indexes().CreateOne(context.TODO(), indexModel)
	if err != nil {
		// It is important that we crash the app on startup to meet the
		// requirements of `google/wire` framework.
		log.Fatal(err)
	}

	s := &NutritionPlanStorerImpl{
		Logger:     loggerp,
		DbClient:   client,
		Collection: uc,
	}
	return s
}
