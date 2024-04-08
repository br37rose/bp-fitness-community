package datastore

import (
	"context"
	"log"
	"log/slog"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	gcp_a "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/adapter/cloudprovider/google"
	c "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/config"
)

type GoogleFitDataPoint struct {
	ID              primitive.ObjectID `bson:"_id" json:"id,omitempty"`
	DataTypeName    string             `bson:"data_type_name" json:"data_type_name"`
	Status          int8               `bson:"status" json:"status"`
	UserID          primitive.ObjectID `bson:"user_id" json:"user_id"`
	UserName        string             `bson:"user_name" json:"user_name"`
	UserLexicalName string             `bson:"user_lexical_name" json:"user_lexical_name"`

	// GoogleFitAppID is the reference ID to the `Google Fit App` we have registered with our system.
	GoogleFitAppID primitive.ObjectID `bson:"googlefit_app_id" json:"googlefit_app_id"`

	// MetricID is the ID of the metric we assigned for this record.
	MetricID primitive.ObjectID `bson:"metric_id" json:"metric_id"`

	StartAt time.Time `bson:"start_at" json:"start_at"`
	EndAt   time.Time `bson:"end_at" json:"end_at"`

	ActivitySegment                  *gcp_a.ActivitySegmentStruct                  `bson:"activity_segment,omitempty" json:"activity_segment,omitempty"`
	BasalMetabolicRate               *gcp_a.BasalMetabolicRateStruct               `bson:"basal_metabolic_rate,omitempty" json:"basal_metabolic_rate,omitempty"`
	CaloriesBurned                   *gcp_a.CaloriesBurnedStruct                   `bson:"calories_burned,omitempty" json:"calories_burned,omitempty"`
	CyclingPedalingCadence           *gcp_a.CyclingPedalingCadenceStruct           `bson:"cycling_pedaling_cadence,omitempty" json:"cycling_pedaling_cadence,omitempty"`
	CyclingPedalingCumulative        *gcp_a.CyclingPedalingCumulativeStruct        `bson:"cycling_pedaling_cumulative,omitempty" json:"cycling_pedaling_cumulative,omitempty"`
	HeartPoints                      *gcp_a.HeartPointsStruct                      `bson:"heart_points,omitempty" json:"heart_points,omitempty"`
	MoveMinutes                      *gcp_a.MoveMinutesStruct                      `bson:"move_minutes,omitempty" json:"move_minutes,omitempty"`
	Power                            *gcp_a.PowerStruct                            `bson:"power,omitempty" json:"power,omitempty"`
	StepCountDelta                   *gcp_a.StepCountDeltaStruct                   `bson:"step_count_delta,omitempty" json:"step_count_delta,omitempty"`
	StepCountCadence                 *gcp_a.StepCountCadenceStruct                 `bson:"step_count_cadence,omitempty" json:"step_count_cadence,omitempty"`
	Workout                          *gcp_a.WorkoutStruct                          `bson:"workout,omitempty" json:"workout,omitempty"`
	CyclingWheelRevolutionRPM        *gcp_a.CyclingWheelRevolutionRPMStruct        `bson:"cycling_wheel_revolution_rpm,omitempty" json:"cycling_wheel_revolution_rpm,omitempty"`
	CyclingWheelRevolutionCumulative *gcp_a.CyclingWheelRevolutionCumulativeStruct `bson:"cycling_wheel_revolution_cumulative,omitempty" json:"cycling_wheel_revolution_cumulative,omitempty"`
	DistanceDelta                    *gcp_a.DistanceDeltaStruct                    `bson:"distance_delta,omitempty" json:"distance_delta,omitempty"`
	LocationSample                   *gcp_a.LocationSampleStruct                   `bson:"location_sample,omitempty" json:"location_sample,omitempty"`
	Speed                            *gcp_a.SpeedStruct                            `bson:"speed,omitempty" json:"speed,omitempty"`
	Hydration                        *gcp_a.HydrationStruct                        `bson:"hydration,omitempty" json:"hydration,omitempty"`
	Nutrition                        *gcp_a.NutritionStruct                        `bson:"nutrition,omitempty" json:"nutrition,omitempty"`
	BloodGlucose                     *gcp_a.BloodGlucoseStruct                     `bson:"blood_glucose,omitempty" json:"blood_glucose,omitempty"`
	BloodPressure                    *gcp_a.BloodPressureStruct                    `bson:"blood_pressure,omitempty" json:"blood_pressure,omitempty"`
	BodyFatPercentage                *gcp_a.BodyFatPercentageStruct                `bson:"body_fat_percentage,omitempty" json:"body_fat_percentage,omitempty"`
	BodyTemperature                  *gcp_a.BodyTemperatureStruct                  `bson:"body_temperature,omitempty" json:"body_temperature,omitempty"`
	HeartRateBPM                     *gcp_a.HeartRateBPMStruct                     `bson:"hearte_rate_bpm,omitempty" json:"hearte_rate_bpm,omitempty"`
	Height                           *gcp_a.HeightStruct                           `bson:"height,omitempty" json:"height,omitempty"`
	Sleep                            *gcp_a.SleepStruct                            `bson:"sleep,omitempty" json:"sleep,omitempty"`
	OxygenSaturation                 *gcp_a.OxygenSaturationStruct                 `bson:"oxygen_saturation,omitempty" json:"oxygen_saturation,omitempty"`
	Weight                           *gcp_a.WeightStruct                           `bson:"weight,omitempty" json:"weight,omitempty"`

	// Error is the error response content provided by `Google Fit` when making the API call.
	Error string `bson:"errors" json:"errors"`

	// CreatedAt represents the time this data point got inserted into our database.
	CreatedAt time.Time `bson:"created_at" json:"created_at,omitempty"`

	// ModifiedAt represents the time this data point was modified by our app.
	ModifiedAt time.Time `bson:"modified_at" json:"modified_at,omitempty"`

	// The organization this data point belongs to. Used for tenancy restrictions.
	OrganizationID primitive.ObjectID `bson:"organization_id" json:"organization_id"`
}

// GoogleFitDataPointStorer Interface for organization.
type GoogleFitDataPointStorer interface {
	CheckIfExistsByCompositeKey(ctx context.Context, userID primitive.ObjectID, dataTypeName string, startAt time.Time, endAt time.Time) (bool, error)
	Create(ctx context.Context, m *GoogleFitDataPoint) error
	// GetByID(ctx context.Context, id primitive.ObjectID) (*GoogleFitDataPoint, error)
	// GetByUserID(ctx context.Context, userID primitive.ObjectID) (*GoogleFitDataPoint, error)
	// GetByName(ctx context.Context, name string) (*GoogleFitDataPoint, error)
	// GetByPaymentProcessorGoogleFitDataPointID(ctx context.Context, paymentProcessorGoogleFitDataPointID string) (*GoogleFitDataPoint, error)
	UpdateByID(ctx context.Context, m *GoogleFitDataPoint) error
	// UpsertByUserID(ctx context.Context, fba *GoogleFitDataPoint) error
	ListByFilter(ctx context.Context, m *GoogleFitDataPointPaginationListFilter) (*GoogleFitDataPointPaginationListResult, error)
	// ListAsSelectOptionByFilter(ctx context.Context, f *GoogleFitDataPointPaginationListFilter) ([]*GoogleFitDataPointAsSelectOption, error)
	ListByQueuedStatus(ctx context.Context) (*GoogleFitDataPointPaginationListResult, error)
	ListByQueuedStatusInDataTypeNames(ctx context.Context, dataTypeNames []string) (*GoogleFitDataPointPaginationListResult, error)
	ListByActiveStatusInDataTypeNames(ctx context.Context, dataTypeNames []string) (*GoogleFitDataPointPaginationListResult, error)
	// DeleteByID(ctx context.Context, id primitive.ObjectID) error
	// CheckIfExistsByNameInOrgBranch(ctx context.Context, name string, orgID primitive.ObjectID, branchID primitive.ObjectID) (bool, error)
	// // //TODO: Add more...
}

type GoogleFitDataPointAsSelectOption struct {
	Value primitive.ObjectID `bson:"_id" json:"value"` // Extract from the database `_id` field and output through API as `value`.
	Label string             `bson:"name" json:"label"`
}

type GoogleFitDataPointStorerImpl struct {
	Logger     *slog.Logger
	DbClient   *mongo.Client
	Collection *mongo.Collection
}

func NewDatastore(appCfg *c.Conf, loggerp *slog.Logger, client *mongo.Client) GoogleFitDataPointStorer {
	// ctx := context.Background()
	uc := client.Database(appCfg.DB.Name).Collection("googlefit_data_points")

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
		{Keys: bson.D{{Key: "googlefit_app_id", Value: 1}}},
		{Keys: bson.D{{Key: "start_at", Value: 1}}},
		{Keys: bson.D{{Key: "status", Value: 1}}},
		{Keys: bson.D{{Key: "type", Value: 1}}},
		{Keys: bson.D{
			{"user_id", "text"},
			{"data_type_name", "text"},
			{"start_at", "text"},
			{"end_at", "text"},
		}},
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

	s := &GoogleFitDataPointStorerImpl{
		Logger:     loggerp,
		DbClient:   client,
		Collection: uc,
	}
	return s
}
