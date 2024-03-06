package controller

import (
	"context"
	"time"

	"log/slog"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	a_d "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/fitnessplan/datastore"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/config/constants"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/utils/httperror"
)

type FitnessPlanCreateRequestIDO struct {
	Birthday           time.Time `bson:"birthday,omitempty" json:"birthday,omitempty"`
	DaysPerWeek        int8      `bson:"days_per_week" json:"days_per_week"`
	EquipmentAccess    int8      `bson:"equipment_access" json:"equipment_access"`
	EstimatedReadyDate time.Time `bson:"estimated_ready_date,omitempty" json:"estimated_ready_date,omitempty"`
	Gender             int8      `bson:"gender" json:"gender"`
	GenderOther        string    `bson:"gender_other" json:"gender_other"`
	Goals              []int8    `bson:"goals" json:"goals"`
	HasWorkoutsAtHome  int8      `bson:"has_workouts_at_home" json:"has_workouts_at_home"`
	HeightFeet         float64   `bson:"height_feet" json:"height_feet"`
	HeightFeetInches   float64   `bson:"height_feet_inches" json:"height_feet_inches"`
	HeightInches       float64   `bson:"height_inches" json:"height_inches"`
	HomeGymEquipment   []int8    `bson:"home_gym_equipment" json:"home_gym_equipment"`
	IdealWeight        float64   `bson:"ideal_weight" json:"ideal_weight"`
	MaxWeeks           int8      `bson:"max_weeks" json:"max_weeks"`
	Name               string    `bson:"name" json:"name"`
	PhysicalActivity   int8      `bson:"physical_activity" json:"physical_activity"`
	Status             int8      `bson:"status" json:"status"`
	TimePerDay         int8      `bson:"time_per_day" json:"time_per_day"`
	WasProcessed       bool      `bson:"was_processed" json:"was_processed"`
	Weight             float64   `bson:"weight" json:"weight"`
	WorkoutIntensity   int8      `bson:"workout_intensity" json:"workout_intensity"`
	WorkoutPreferences []int8    `bson:"workout_preferences" json:"workout_preferences"`
}

func ValidateCreateRequest(dirtyData *FitnessPlanCreateRequestIDO) error {
	e := make(map[string]string)

	if dirtyData.EquipmentAccess > 3 || dirtyData.EquipmentAccess <= 0 {
		e["equipment_access"] = "missing value"
	}
	if dirtyData.HasWorkoutsAtHome > 2 || dirtyData.HasWorkoutsAtHome < 1 {
		e["has_workouts_at_home"] = "missing value"
	}
	if dirtyData.Birthday.IsZero() {
		e["birthday"] = "missing value"
	}
	if dirtyData.HeightFeet <= 0 {
		e["height_feet"] = "missing value"
	}
	if dirtyData.HeightInches < 0 { // Note: A value of zero is acceptable here.
		e["height_inches"] = "missing value"
	}
	if dirtyData.Weight == 0 {
		e["weight"] = "missing value"
	}
	if dirtyData.Gender == 0 {
		e["gender"] = "missing value"
	}
	if dirtyData.Gender == 3 { // OTHER
		if dirtyData.GenderOther == "" {
			e["gender_other"] = "missing value"
		}
	}
	if dirtyData.IdealWeight == 0 {
		e["ideal_weight"] = "missing value"
	}
	if dirtyData.PhysicalActivity == 0 {
		e["physical_activity"] = "missing value"
	}
	if dirtyData.WorkoutIntensity == 0 {
		e["workout_intensity"] = "missing value"
	}
	if dirtyData.DaysPerWeek == 0 {
		e["days_per_week"] = "missing value"
	}
	if dirtyData.TimePerDay == 0 {
		e["time_per_day"] = "missing value"
	}
	if dirtyData.MaxWeeks == 0 {
		e["max_weeks"] = "missing value"
	}
	if len(dirtyData.Goals) == 0 {
		e["goals"] = "missing value"
	}
	if len(dirtyData.WorkoutPreferences) == 0 {
		e["workout_preferences"] = "missing value"
	}
	if len(e) != 0 {
		return httperror.NewForBadRequest(&e)
	}
	return nil
}

func (c *FitnessPlanControllerImpl) Create(ctx context.Context, req *FitnessPlanCreateRequestIDO) (*a_d.FitnessPlan, error) {
	// Extract from our session the following data.
	orgID := ctx.Value(constants.SessionUserOrganizationID).(primitive.ObjectID)
	orgName := ctx.Value(constants.SessionUserOrganizationName).(string)
	userID := ctx.Value(constants.SessionUserID).(primitive.ObjectID)
	userName := ctx.Value(constants.SessionUserName).(string)
	userLexicalName := ctx.Value(constants.SessionUserLexicalName).(string)

	if err := ValidateCreateRequest(req); err != nil {
		return nil, err
	}

	////
	//// Start the transaction.
	////

	session, err := c.DbClient.StartSession()
	if err != nil {
		c.Logger.Error("start session error",
			slog.Any("error", err))
		return nil, err
	}
	defer session.EndSession(ctx)

	// Define a transaction function with a series of operations
	transactionFunc := func(sessCtx mongo.SessionContext) (interface{}, error) {

		// Create our record in the database.
		res := &a_d.FitnessPlan{
			OrganizationID:     orgID,
			OrganizationName:   orgName,
			ID:                 primitive.NewObjectID(),
			CreatedAt:          time.Now(),
			CreatedByUserName:  userName,
			CreatedByUserID:    userID,
			ModifiedAt:         time.Now(),
			ModifiedByUserName: userName,
			ModifiedByUserID:   userID,
			Status:             a_d.StatusQueued, // Set to queued b/c we will wait on openai.
			Name:               req.Name,
			DaysPerWeek:        req.DaysPerWeek,
			EquipmentAccess:    req.EquipmentAccess,
			EstimatedReadyDate: req.EstimatedReadyDate,
			Gender:             req.Gender,
			Birthday:           req.Birthday,
			Goals:              req.Goals,
			HasWorkoutsAtHome:  req.HasWorkoutsAtHome,
			HeightFeet:         req.HeightFeet,
			HeightFeetInches:   req.HeightFeetInches,
			HeightInches:       req.HeightInches,
			HomeGymEquipment:   req.HomeGymEquipment,
			IdealWeight:        req.IdealWeight,
			MaxWeeks:           req.MaxWeeks,
			PhysicalActivity:   req.PhysicalActivity,
			TimePerDay:         req.TimePerDay,
			Weight:             req.Weight,
			WorkoutIntensity:   req.WorkoutIntensity,
			WorkoutPreferences: req.WorkoutPreferences,
			ExerciseNames:      []string{},
			Exercises:          make([]*a_d.FitnessPlanExercise, 0),
			UserID:             userID,
			UserName:           userName,
			UserLexicalName:    userLexicalName,
		}
		err := c.FitnessPlanStorer.Create(sessCtx, res)
		if err != nil {
			c.Logger.Error("fitnessplan create error", slog.Any("error", err))
			return nil, err
		}

		return res, nil
	}

	// Start a transaction
	result, err := session.WithTransaction(ctx, transactionFunc)
	if err != nil {
		c.Logger.Error("session failed error",
			slog.Any("error", err))
		return nil, err
	}

	// Execute in the background the call to OpenAI and generate the fitness
	// plan.
	if result != nil {
		go c.generateFitnessPlanInBackground(context.Background(), result.(*a_d.FitnessPlan))
	}

	return result.(*a_d.FitnessPlan), nil
}
