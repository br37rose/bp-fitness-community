package controller

import (
	"context"
	"log/slog"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	a_d "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/fitnessplan/datastore"
	domain "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/fitnessplan/datastore"
	user_d "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/user/datastore"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/config/constants"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/utils/httperror"
)

type FitnessPlanUpdateRequestIDO struct {
	ID                 primitive.ObjectID `bson:"id,omitempty" json:"id,omitempty"`
	Birthday           time.Time          `bson:"birthday,omitempty" json:"birthday,omitempty"`
	DaysPerWeek        int8               `bson:"days_per_week" json:"days_per_week"`
	EquipmentAccess    int8               `bson:"equipment_access" json:"equipment_access"`
	EstimatedReadyDate time.Time          `bson:"estimated_ready_date,omitempty" json:"estimated_ready_date,omitempty"`
	Gender             int8               `bson:"gender" json:"gender"`
	GenderOther        string             `bson:"gender_other" json:"gender_other"`
	Goals              []int8             `bson:"goals" json:"goals"`
	HasWorkoutsAtHome  int8               `bson:"has_workouts_at_home" json:"has_workouts_at_home"`
	HeightFeet         float64            `bson:"height_feet" json:"height_feet"`
	HeightFeetInches   float64            `bson:"height_feet_inches" json:"height_feet_inches"`
	HeightInches       float64            `bson:"height_inches" json:"height_inches"`
	HomeGymEquipment   []int8             `bson:"home_gym_equipment" json:"home_gym_equipment"`
	IdealWeight        float64            `bson:"ideal_weight" json:"ideal_weight"`
	MaxWeeks           int8               `bson:"max_weeks" json:"max_weeks"`
	Name               string             `bson:"name" json:"name"`
	PhysicalActivity   int8               `bson:"physical_activity" json:"physical_activity"`
	Status             int8               `bson:"status" json:"status"`
	TimePerDay         int8               `bson:"time_per_day" json:"time_per_day"`
	WasProcessed       bool               `bson:"was_processed" json:"was_processed"`
	Weight             float64            `bson:"weight" json:"weight"`
	WorkoutIntensity   int8               `bson:"workout_intensity" json:"workout_intensity"`
	WorkoutPreferences []int8             `bson:"workout_preferences" json:"workout_preferences"`
}

func ValidateUpdateRequest(dirtyData *FitnessPlanUpdateRequestIDO) error {
	e := make(map[string]string)

	if dirtyData.ID.IsZero() {
		e["id"] = "missing value"
	}
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

func (c *FitnessPlanControllerImpl) UpdateByID(ctx context.Context, req *FitnessPlanUpdateRequestIDO) (*domain.FitnessPlan, error) {
	if err := ValidateUpdateRequest(req); err != nil {
		return nil, err
	}

	c.Kmutex.Lockf("update-for-fitness-plan-%s", req.ID.Hex())
	defer c.Kmutex.Unlockf("update-for-fitness-plan-%s", req.ID.Hex())

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
		// Fetch the original fitnessplan.
		os, err := c.FitnessPlanStorer.GetByID(sessCtx, req.ID)
		if err != nil {
			c.Logger.Error("database get by id error",
				slog.Any("error", err),
				slog.Any("fitnessplan_id", req.ID))
			return nil, err
		}
		if os == nil {
			c.Logger.Error("fitnessplan does not exist error",
				slog.Any("fitnessplan_id", req.ID))
			return nil, httperror.NewForBadRequestWithSingleField("message", "fitnessplan does not exist")
		}

		// Prevent any updates if we submitted to open-ai without open-ai finished.
		if os.Status == domain.StatusQueued {
			return nil, httperror.NewForBadRequestWithSingleField("message", "your fitness plan is being queued by us and you will be able to resubmit when we finish")
		}

		// Extract from our session the following data.
		userID := sessCtx.Value(constants.SessionUserID).(primitive.ObjectID)
		userOrganizationID := sessCtx.Value(constants.SessionUserOrganizationID).(primitive.ObjectID)
		userRole := sessCtx.Value(constants.SessionUserRole).(int8)
		userName := sessCtx.Value(constants.SessionUserName).(string)

		// If user is not administrator nor belongs to the fitnessplan then error.
		if userRole != user_d.UserRoleRoot && os.OrganizationID != userOrganizationID {
			c.Logger.Error("authenticated user is not staff role nor belongs to the fitnessplan error",
				slog.Any("userRole", userRole),
				slog.Any("userOrganizationID", userOrganizationID))
			return nil, httperror.NewForForbiddenWithSingleField("message", "you do not belong to this fitnessplan")
		}

		// Modify our original fitnessplan.
		os.ModifiedAt = time.Now()
		os.ModifiedByUserID = userID
		os.ModifiedByUserName = userName
		os.Birthday = req.Birthday
		os.DaysPerWeek = req.DaysPerWeek
		os.EquipmentAccess = req.EquipmentAccess
		os.EstimatedReadyDate = req.EstimatedReadyDate
		os.Gender = req.Gender
		os.Goals = req.Goals
		os.HasWorkoutsAtHome = req.HasWorkoutsAtHome
		os.HeightFeet = req.HeightFeet
		os.HeightFeetInches = req.HeightFeetInches
		os.HeightInches = req.HeightInches
		os.HomeGymEquipment = req.HomeGymEquipment
		os.IdealWeight = req.IdealWeight
		os.MaxWeeks = req.MaxWeeks
		os.Name = req.Name
		os.PhysicalActivity = req.PhysicalActivity
		os.TimePerDay = req.TimePerDay
		os.WasProcessed = req.WasProcessed
		os.Weight = req.Weight
		os.WorkoutIntensity = req.WorkoutIntensity
		os.WorkoutPreferences = req.WorkoutPreferences
		os.Status = domain.StatusQueued

		// Save to the database the modified fitnessplan.
		if err := c.FitnessPlanStorer.UpdateByID(sessCtx, os); err != nil {
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
		go c.generateFitnessPlanInBackground(context.Background(), result.(*a_d.FitnessPlan))
	}

	return result.(*a_d.FitnessPlan), nil
}

// func (c *FitnessPlanControllerImpl) updateFitnessPlanNameForAllUsers(ctx context.Context, ns *domain.FitnessPlan) error {
// 	c.Logger.Debug("Beginning to update fitnessplan name for all uses")
// 	f := &user_d.UserListFilter{FitnessPlanID: ns.ID}
// 	uu, err := c.UserStorer.ListByFilter(ctx, f)
// 	if err != nil {
// 		c.Logger.Error("database update by id error", slog.Any("error", err))
// 		return err
// 	}
// 	for _, u := range uu.Results {
// 		u.FitnessPlanName = ns.Name
// 		log.Println("--->", u)
// 		// if err := c.UserStorer.UpdateByID(ctx, u); err != nil {
// 		// 	c.Logger.Error("database update by id error", slog.Any("error", err))
// 		// 	return err
// 		// }
// 	}
// 	return nil
// }
//
// func (c *FitnessPlanControllerImpl) updateFitnessPlanNameForAllComicSubmissions(ctx context.Context, ns *domain.FitnessPlan) error {
// 	c.Logger.Debug("Beginning to update fitnessplan name for all submissions")
// 	f := &domain.ComicSubmissionListFilter{FitnessPlanID: ns.ID}
// 	uu, err := c.ComicSubmissionStorer.ListByFilter(ctx, f)
// 	if err != nil {
// 		c.Logger.Error("database update by id error", slog.Any("error", err))
// 		return err
// 	}
// 	for _, u := range uu.Results {
// 		u.FitnessPlanName = ns.Name
// 		log.Println("--->", u)
// 		// if err := c.ComicSubmissionStorer.UpdateByID(ctx, u); err != nil {
// 		// 	c.Logger.Error("database update by id error", slog.Any("error", err))
// 		// 	return err
// 		// }
// 	}
// 	return nil
// }
