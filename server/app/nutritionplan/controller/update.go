package controller

import (
	"context"
	"log/slog"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	a_d "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/nutritionplan/datastore"
	domain "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/nutritionplan/datastore"
	user_d "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/user/datastore"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/config/constants"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/utils/httperror"
)

type NutritionPlanUpdateRequestIDO struct {
	ID                         primitive.ObjectID `bson:"id,omitempty" json:"id,omitempty"`
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
}

func ValidateUpdateRequest(dirtyData *NutritionPlanUpdateRequestIDO) error {
	e := make(map[string]string)

	if dirtyData.ID.IsZero() {
		e["id"] = "missing value"
	}
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

func (c *NutritionPlanControllerImpl) UpdateByID(ctx context.Context, req *NutritionPlanUpdateRequestIDO) (*domain.NutritionPlan, error) {
	if err := ValidateUpdateRequest(req); err != nil {
		return nil, err
	}

	c.Kmutex.Lockf("update-for-nutrition-plan-%s", req.ID.Hex())
	defer c.Kmutex.Unlockf("update-for-nutrition-plan-%s", req.ID.Hex())

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

		// Fetch the original nutritionplan.
		os, err := c.NutritionPlanStorer.GetByID(ctx, req.ID)
		if err != nil {
			c.Logger.Error("database get by id error",
				slog.Any("error", err),
				slog.Any("nutritionplan_id", req.ID))
			return nil, err
		}
		if os == nil {
			c.Logger.Error("nutritionplan does not exist error",
				slog.Any("nutritionplan_id", req.ID))
			return nil, httperror.NewForBadRequestWithSingleField("message", "nutritionplan does not exist")
		}

		// Prevent any updates if we submitted to open-ai without open-ai finished.
		if os.Status == domain.StatusQueued {
			return nil, httperror.NewForBadRequestWithSingleField("message", "your nutrition plan is being queued by us and you will be able to resubmit when we finish")
		}

		// Extract from our session the following data.
		userID := ctx.Value(constants.SessionUserID).(primitive.ObjectID)
		userOrganizationID := ctx.Value(constants.SessionUserOrganizationID).(primitive.ObjectID)
		userRole := ctx.Value(constants.SessionUserRole).(int8)
		userName := ctx.Value(constants.SessionUserName).(string)

		// If user is not administrator nor belongs to the nutritionplan then error.
		if userRole != user_d.UserRoleRoot && os.OrganizationID != userOrganizationID {
			c.Logger.Error("authenticated user is not staff role nor belongs to the nutritionplan error",
				slog.Any("userRole", userRole),
				slog.Any("userOrganizationID", userOrganizationID))
			return nil, httperror.NewForForbiddenWithSingleField("message", "you do not belong to this nutritionplan")
		}

		// Modify our original nutritionplan.
		os.ModifiedAt = time.Now()
		os.Status = domain.StatusQueued // Do not change! This is required b/c of OpenAI submission.
		os.ModifiedByUserID = userID
		os.ModifiedByUserName = userName

		os.HasAllergies = req.HasAllergies
		os.Allergies = req.Allergies
		os.Birthday = req.Birthday
		os.MealsPerDay = req.MealsPerDay
		os.ConsumeJunkFood = req.ConsumeJunkFood
		os.ConsumeFruitsAndVegetables = req.ConsumeFruitsAndVegetables
		os.HasIntermittentFasting = req.HasIntermittentFasting
		os.HeightFeet = req.HeightFeet
		os.HeightFeetInches = req.HeightFeetInches
		os.HeightInches = req.HeightInches
		os.Weight = req.Weight
		os.Gender = req.Gender
		os.GenderOther = req.GenderOther
		os.IdealWeight = req.IdealWeight
		os.PhysicalActivity = req.PhysicalActivity
		os.WorkoutIntensity = req.WorkoutIntensity
		os.Goals = req.Goals
		os.MaxWeeks = req.MaxWeeks
		os.Name = req.Name

		// Save to the database the modified nutritionplan.
		if err := c.NutritionPlanStorer.UpdateByID(ctx, os); err != nil {
			c.Logger.Error("database update by id error", slog.Any("error", err))
			return nil, err
		}
		return os, nil

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

// func (c *NutritionPlanControllerImpl) updateNutritionPlanNameForAllUsers(ctx context.Context, ns *domain.NutritionPlan) error {
// 	c.Logger.Debug("Beginning to update nutritionplan name for all uses")
// 	f := &user_d.UserListFilter{NutritionPlanID: ns.ID}
// 	uu, err := c.UserStorer.ListByFilter(ctx, f)
// 	if err != nil {
// 		c.Logger.Error("database update by id error", slog.Any("error", err))
// 		return err
// 	}
// 	for _, u := range uu.Results {
// 		u.NutritionPlanName = ns.Name
// 		log.Println("--->", u)
// 		// if err := c.UserStorer.UpdateByID(ctx, u); err != nil {
// 		// 	c.Logger.Error("database update by id error", slog.Any("error", err))
// 		// 	return err
// 		// }
// 	}
// 	return nil
// }
//
// func (c *NutritionPlanControllerImpl) updateNutritionPlanNameForAllComicSubmissions(ctx context.Context, ns *domain.NutritionPlan) error {
// 	c.Logger.Debug("Beginning to update nutritionplan name for all submissions")
// 	f := &domain.ComicSubmissionListFilter{NutritionPlanID: ns.ID}
// 	uu, err := c.ComicSubmissionStorer.ListByFilter(ctx, f)
// 	if err != nil {
// 		c.Logger.Error("database update by id error", slog.Any("error", err))
// 		return err
// 	}
// 	for _, u := range uu.Results {
// 		u.NutritionPlanName = ns.Name
// 		log.Println("--->", u)
// 		// if err := c.ComicSubmissionStorer.UpdateByID(ctx, u); err != nil {
// 		// 	c.Logger.Error("database update by id error", slog.Any("error", err))
// 		// 	return err
// 		// }
// 	}
// 	return nil
// }
