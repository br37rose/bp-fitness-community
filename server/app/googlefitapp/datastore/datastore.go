package datastore

import (
	"context"
	"log"
	"time"

	"log/slog"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/oauth2"

	c "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/config"
)

const (
	StatusActive   = 1
	StatusArchived = 2
	StatusError    = 3
	AuthTypeOAuth2 = 1
)

type GoogleFitApp struct {
	ID               primitive.ObjectID `bson:"_id" json:"id,omitempty"`
	UserFirstName    string             `bson:"user_first_name" json:"user_first_name"`
	UserLastName     string             `bson:"user_last_name" json:"user_last_name"`
	UserName         string             `bson:"user_name" json:"user_name"`
	UserLexicalName  string             `bson:"user_lexical_name" json:"user_lexical_name"`
	UserID           primitive.ObjectID `bson:"user_id" json:"user_id"`
	Status           int8               `bson:"status" json:"status"`
	CreatedAt        time.Time          `bson:"created_at" json:"created_at,omitempty"`
	ModifiedAt       time.Time          `bson:"modified_at" json:"modified_at,omitempty"`
	OrganizationID   primitive.ObjectID `bson:"organization_id" json:"organization_id"`
	OrganizationName string             `bson:"organization_name" json:"organization_name"`

	// GoogleFitUserID is the user id provided by googlefit web-services for our user account.
	GoogleFitUserID string `bson:"googlefit_user_id" json:"googlefit_user_id"`

	// AuthType tracks how the user authenticated their googlefit with our app.
	AuthType int8 `bson:"auth_type" json:"auth_type"`

	// Errors indicates what error was returned by GoogleFit web-services.
	Errors string `bson:"errors" json:"errors,omitempty"`

	// The token returned through the successful oAuth exchange for the user.
	Token *oauth2.Token `bson:"token" json:"token,omitempty"`

	// The last time we made a fetch to Google API.
	LastFetchedAt time.Time `bson:"last_fetched_at" json:"last_fetched_at,omitempty"`

	ActivitySegmentMetricID                  primitive.ObjectID `bson:"activity_segment_metric_id" json:"activity_segment_metric_id,omitempty"`
	BasalMetabolicRateMetricID               primitive.ObjectID `bson:"basal_metabolic_rate_metric_id" json:"basal_metabolic_rate_metric_id,omitempty"`
	CaloriesBurnedMetricID                   primitive.ObjectID `bson:"calories_burned_metric_id" json:"calories_burned_metric_id,omitempty"`
	CyclingPedalingCadenceMetricID           primitive.ObjectID `bson:"cycling_pedaling_cadence_metric_id" json:"cycling_pedaling_cadence_metric_id,omitempty"`
	CyclingPedalingCumulativeMetricID        primitive.ObjectID `bson:"cycling_pedaling_cumulative_metric_id" json:"cycling_pedaling_cumulative_metric_id,omitempty"`
	HeartPointsMetricID                      primitive.ObjectID `bson:"heart_points_id" json:"heart_points_id,omitempty"`
	MoveMinutesMetricID                      primitive.ObjectID `bson:"move_minutes_metric_id" json:"move_minutes_metric_id,omitempty"`
	PowerMetricID                            primitive.ObjectID `bson:"power_metric_id" json:"power_metric_id,omitempty"`
	StepCountDeltaMetricID                   primitive.ObjectID `bson:"step_count_delta_metric_id" json:"step_count_delta_metric_id,omitempty"`
	StepCountCadenceMetricID                 primitive.ObjectID `bson:"step_count_cadence_metric_id" json:"step_count_cadence_metric_id,omitempty"`
	WorkoutMetricID                          primitive.ObjectID `bson:"workout_metric_id" json:"workout_metric_id,omitempty"` // is deprecated.
	CyclingWheelRevolutionRPMMetricID        primitive.ObjectID `bson:"cycling_wheel_revolution_rpm_metric_id" json:"cycling_wheel_revolution_rpm_metric_id,omitempty"`
	CyclingWheelRevolutionCumulativeMetricID primitive.ObjectID `bson:"cycling_wheel_revolution_cumulative_metric_id" json:"cycling_wheel_revolution_cumulative_metric_id,omitempty"`
	DistanceDeltaMetricID                    primitive.ObjectID `bson:"distance_delta_metric_id" json:"distance_delta_metric_id,omitempty"`
	LocationSampleMetricID                   primitive.ObjectID `bson:"location_sample_metric_id" json:"location_sample_metric_id,omitempty"`
	SpeedMetricID                            primitive.ObjectID `bson:"speed_metric_id" json:"speed_metric_id,omitempty"`
	HydrationMetricID                        primitive.ObjectID `bson:"hydration_metric_id" json:"hydration_metric_id,omitempty"`
	NutritionMetricID                        primitive.ObjectID `bson:"nutrition_metric_id" json:"nutrition_metric_id,omitempty"`
	BloodGlucoseMetricID                     primitive.ObjectID `bson:"blood_glucose_metric_id" json:"blood_glucose_metric_id,omitempty"`
	BloodPressureMetricID                    primitive.ObjectID `bson:"blood_pressure_metric_id" json:"blood_pressure_metric_id,omitempty"`
	BodyFatPercentageMetricID                primitive.ObjectID `bson:"body_fat_percentage_metric_id" json:"body_fat_percentage_metric_id,omitempty"`
	BodyTemperatureMetricID                  primitive.ObjectID `bson:"body_temperature_percentage_metric_id" json:"body_temperature_percentage_metric_id,omitempty"`
	HeartRateBPMMetricID                     primitive.ObjectID `bson:"heart_rate_bpm_metric_id" json:"heart_rate_bpm_metric_id,omitempty"`
	HeightMetricID                           primitive.ObjectID `bson:"height_metric_id" json:"height_metric_id,omitempty"`
	OxygenSaturationMetricID                 primitive.ObjectID `bson:"oxygen_saturation_metric_id" json:"oxygen_saturation_metric_id,omitempty"`
	SleepMetricID                            primitive.ObjectID `bson:"sleep_metric_id" json:"sleep_metric_id,omitempty"`
	WeightMetricID                           primitive.ObjectID `bson:"weight_metric_id" json:"weight_metric_id,omitempty"`

	IsTestMode         bool   `bson:"is_test_mode" json:"is_test_mode"`
	SimulatorAlgorithm string `bson:"simulator_algorithm,omitempty" json:"simulator_algorithm,omitempty"`

	// RequiresGoogleLoginAgain indicates whether the user must log in again
	// into Google for whatever reason to re-authorize our app.
	RequiresGoogleLoginAgain bool `bson:"requires_google_login_again" json:"requires_google_login_again"`
}

type GoogleFitAppListFilter struct {
	// Pagination related.
	Cursor    primitive.ObjectID
	PageSize  int64
	SortField string
	SortOrder int8 // 1=ascending | -1=descending

	// Filter related.
	OrganizationID primitive.ObjectID
	BranchID       primitive.ObjectID
	Status         int8
	SearchText     string
}

type GoogleFitAppListResult struct {
	Results     []*GoogleFitApp    `json:"results"`
	NextCursor  primitive.ObjectID `json:"next_cursor"`
	HasNextPage bool               `json:"has_next_page"`
}

// GoogleFitAppStorer Interface for organization.
type GoogleFitAppStorer interface {
	Create(ctx context.Context, m *GoogleFitApp) error
	GetByID(ctx context.Context, id primitive.ObjectID) (*GoogleFitApp, error)
	GetByUserID(ctx context.Context, userID primitive.ObjectID) (*GoogleFitApp, error)
	GetByName(ctx context.Context, name string) (*GoogleFitApp, error)
	GetByPaymentProcessorGoogleFitAppID(ctx context.Context, paymentProcessorGoogleFitAppID string) (*GoogleFitApp, error)
	UpdateByID(ctx context.Context, m *GoogleFitApp) error
	UpsertByUserID(ctx context.Context, fba *GoogleFitApp) error
	ListByFilter(ctx context.Context, m *GoogleFitAppListFilter) (*GoogleFitAppListResult, error)
	ListAsSelectOptionByFilter(ctx context.Context, f *GoogleFitAppListFilter) ([]*GoogleFitAppAsSelectOption, error)
	ListIDsByStatus(ctx context.Context, status int8) ([]primitive.ObjectID, error)
	ListDevicesByStatus(ctx context.Context, status int8) ([]*GoogleFitAppDevice, error)
	ListPhysicalIDsByStatus(ctx context.Context, status int8) ([]primitive.ObjectID, error)
	ListSimulatorIDsByStatus(ctx context.Context, status int8) ([]primitive.ObjectID, error)
	DeleteByID(ctx context.Context, id primitive.ObjectID) error
	CheckIfExistsByNameInOrgBranch(ctx context.Context, name string, orgID primitive.ObjectID, branchID primitive.ObjectID) (bool, error)
	// //TODO: Add more...
}

type GoogleFitAppAsSelectOption struct {
	Value primitive.ObjectID `bson:"_id" json:"value"` // Extract from the database `_id` field and output through API as `value`.
	Label string             `bson:"name" json:"label"`
}

type GoogleFitAppStorerImpl struct {
	Logger     *slog.Logger
	DbClient   *mongo.Client
	Collection *mongo.Collection
}

func NewDatastore(appCfg *c.Conf, loggerp *slog.Logger, client *mongo.Client) GoogleFitAppStorer {
	// ctx := context.Background()
	uc := client.Database(appCfg.DB.Name).Collection("google_fit_apps")

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
		{Keys: bson.D{{Key: "user_id", Value: 1}}},
		{Keys: bson.D{{Key: "googlefit_user_id", Value: 1}}},
		{Keys: bson.D{{Key: "status", Value: 1}}},
		// {Keys: bson.D{
		// 	{"name", "text"},
		// 	{"lexical_name", "text"},
		// 	{"email", "text"},
		// 	{"phone", "text"},
		// 	{"country", "text"},
		// 	{"region", "text"},
		// 	{"city", "text"},
		// 	{"postal_code", "text"},
		// 	{"address_line1", "text"},
		// 	{"description", "text"},
		// }},
	})
	if err != nil {
		// It is important that we crash the app on startup to meet the
		// requirements of `google/wire` framework.
		log.Fatal(err)
	}

	s := &GoogleFitAppStorerImpl{
		Logger:     loggerp,
		DbClient:   client,
		Collection: uc,
	}
	return s
}
