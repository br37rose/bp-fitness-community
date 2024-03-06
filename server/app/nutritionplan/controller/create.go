package controller

import (
	"context"
	"time"

	"log/slog"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	a_d "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/nutritionplan/datastore"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/config/constants"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/utils/httperror"
)

type NutritionPlanCreateRequestIDO struct {
	HasAllergies               int8      `bson:"has_allergies" json:"has_allergies"`
	Allergies                  string    `bson:"allergies" json:"allergies"`
	Birthday                   time.Time `bson:"birthday,omitempty" json:"birthday,omitempty"`
	MealsPerDay                int8      `bson:"meals_per_day" json:"meals_per_day"`
	ConsumeJunkFood            int8      `bson:"consume_junk_food" json:"consume_junk_food"`
	ConsumeFruitsAndVegetables int8      `bson:"consume_fruits_and_vegetables" json:"consume_fruits_and_vegetables"`
	HasIntermittentFasting     int8      `bson:"has_intermittent_fasting" json:"has_intermittent_fasting"`
	HeightFeet                 float64   `bson:"height_feet" json:"height_feet"`
	HeightFeetInches           float64   `bson:"height_feet_inches" json:"height_feet_inches"`
	HeightInches               float64   `bson:"height_inches" json:"height_inches"`
	Weight                     float64   `bson:"weight" json:"weight"`
	Gender                     int8      `bson:"gender" json:"gender"`
	GenderOther                string    `bson:"gender_other" json:"gender_other"`
	IdealWeight                float64   `bson:"ideal_weight" json:"ideal_weight"`
	PhysicalActivity           int8      `bson:"physical_activity" json:"physical_activity"`
	Status                     int8      `bson:"status" json:"status"`
	WorkoutIntensity           int8      `bson:"workout_intensity" json:"workout_intensity"`
	Goals                      []int8    `bson:"goals" json:"goals"`
	MaxWeeks                   int8      `bson:"max_weeks" json:"max_weeks"`
	Name                       string    `bson:"name" json:"name"`
}

func ValidateCreateRequest(dirtyData *NutritionPlanCreateRequestIDO) error {
	e := make(map[string]string)

	if dirtyData.HasAllergies == 0 {
		e["has_allergies"] = "missing value"
	} else {
		if dirtyData.HasAllergies == 1 {
			if dirtyData.Allergies == "" {
				e["allergies"] = "missing value"
			}
		}
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
	if dirtyData.MealsPerDay == 0 {
		e["meals_per_day"] = "missing value"
	}
	if dirtyData.ConsumeJunkFood == 0 {
		e["consume_junk_food"] = "missing value"
	}
	if dirtyData.ConsumeFruitsAndVegetables == 0 {
		e["consume_fruits_and_vegetables"] = "missing value"
	}
	if dirtyData.HasIntermittentFasting == 0 {
		e["has_intermittent_fasting"] = "missing value"
	}
	if dirtyData.MaxWeeks == 0 {
		e["max_weeks"] = "missing value"
	}
	if len(dirtyData.Goals) == 0 {
		e["goals"] = "missing value"
	}
	if len(e) != 0 {
		return httperror.NewForBadRequest(&e)
	}
	return nil
}

func (c *NutritionPlanControllerImpl) Create(ctx context.Context, req *NutritionPlanCreateRequestIDO) (*a_d.NutritionPlan, error) {
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
		res := &a_d.NutritionPlan{
			OrganizationID:             orgID,
			OrganizationName:           orgName,
			ID:                         primitive.NewObjectID(),
			CreatedAt:                  time.Now(),
			CreatedByUserName:          userName,
			CreatedByUserID:            userID,
			ModifiedAt:                 time.Now(),
			ModifiedByUserName:         userName,
			ModifiedByUserID:           userID,
			Status:                     a_d.StatusQueued, // Do not change! This is required b/c of OpenAI submission
			Name:                       req.Name,
			HasAllergies:               req.HasAllergies,
			Allergies:                  req.Allergies,
			Birthday:                   req.Birthday,
			MealsPerDay:                req.MealsPerDay,
			ConsumeJunkFood:            req.ConsumeJunkFood,
			ConsumeFruitsAndVegetables: req.ConsumeFruitsAndVegetables,
			HasIntermittentFasting:     req.HasIntermittentFasting,
			HeightFeet:                 req.HeightFeet,
			HeightFeetInches:           req.HeightFeetInches,
			HeightInches:               req.HeightInches,
			Weight:                     req.Weight,
			Gender:                     req.Gender,
			GenderOther:                req.GenderOther,
			IdealWeight:                req.IdealWeight,
			PhysicalActivity:           req.PhysicalActivity,
			WorkoutIntensity:           req.WorkoutIntensity,
			Goals:                      req.Goals,
			MaxWeeks:                   req.MaxWeeks,
			UserID:                     userID,
			UserName:                   userName,
			UserLexicalName:            userLexicalName,
		}
		err := c.NutritionPlanStorer.Create(ctx, res)
		if err != nil {
			c.Logger.Error("nutritionplan create error", slog.Any("error", err))
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
		go c.generateNutritionPlanInBackground(context.Background(), result.(*a_d.NutritionPlan))
	}

	return result.(*a_d.NutritionPlan), nil
}
