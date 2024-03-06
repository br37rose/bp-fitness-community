package controller

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"time"

	exercise_s "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/exercise/datastore"
	a_d "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/fitnessplan/datastore"
	domain "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/fitnessplan/datastore"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func (c *FitnessPlanControllerImpl) generateFitnessPlanInBackground(ctx context.Context, fp *domain.FitnessPlan) {
	c.Logger.Debug("getting fitness plan exercise recommondations from openai",
		slog.Any("modified_by_user_id", fp.ModifiedByUserID))

	c.Kmutex.Lockf("openai-api-for-fitness-plan-%s", fp.ID.Hex())
	defer c.Kmutex.Unlockf("openai-api-for-fitness-plan-%s", fp.ID.Hex())

	////
	//// Start the transaction.
	////

	session, err := c.DbClient.StartSession()
	if err != nil {
		c.Logger.Error("start session error",
			slog.Any("modified_by_user_id", fp.ModifiedByUserID),
			slog.Any("error", err))
		return
	}
	defer session.EndSession(ctx)

	// Define a transaction function with a series of operations
	transactionFunc := func(sessCtx mongo.SessionContext) (interface{}, error) {

		// if err := c.generateFitnessPlanRecommendedExercises(sessCtx, fp); err != nil {
		// 	c.Logger.Error("failed getting from openai",
		// 		slog.Any("modified_by_user_id", fp.ModifiedByUserID))
		// 	return nil, err
		// }

		if fp.Status != domain.StatusError {
			c.Logger.Debug("getting fitness plan instrucitons from openai",
				slog.Any("modified_by_user_id", fp.ModifiedByUserID))

			if err := c.generateFitnessPlanInstructions(sessCtx, fp); err != nil { // xxx
				c.Logger.Error("failed getting from openai",
					slog.Any("modified_by_user_id", fp.ModifiedByUserID))
				return nil, err
			}

			// Set the state.
			// fp.Status = domain.StatusActive
			fp.ModifiedAt = time.Now()
			if err := c.FitnessPlanStorer.UpdateByID(sessCtx, fp); err != nil {
				c.Logger.Error("response received",
					slog.Any("modified_by_user_id", fp.ModifiedByUserID))
				return nil, err
			}
		}
		c.Logger.Debug("finished getting from openai",
			slog.Any("modified_by_user_id", fp.ModifiedByUserID))

		return nil, nil
	}
	// Start a transaction

	if _, err := session.WithTransaction(ctx, transactionFunc); err != nil {
		c.Logger.Error("session failed error",
			slog.Any("error", err))
	}
}

func (c *FitnessPlanControllerImpl) generateFitnessPlanRecommendedExercises(sessCtx mongo.SessionContext, fp *domain.FitnessPlan) error {
	c.Logger.Debug("generating exercises recommendations prompt...",
		slog.Any("modified_by_user_id", fp.ModifiedByUserID))

	prompt, err := c.generateFitnessPlanRecommendedExercisesPrompt(sessCtx, fp)
	if err != nil {
		c.Logger.Error("failed generating prompt", slog.Any("error", err))
		return err
	}
	c.Logger.Debug("prompt ready",
		slog.Any("modified_by_user_id", fp.ModifiedByUserID))
	fp.PromptOne = prompt

	c.Logger.Debug("submitting prompt to openai",
		slog.Any("modified_by_user_id", fp.ModifiedByUserID))

	res, err := c.OpenAI.CreateChatCompletion(prompt)
	if err != nil {
		c.Logger.Error("failed executing `completion` api call to openai",
			slog.Any("error", err),
			slog.Any("modified_by_user_id", fp.ModifiedByUserID))
		return err
	}
	c.Logger.Debug("response received",
		slog.Any("modified_by_user_id", fp.ModifiedByUserID))

	// Attempt to unmarhsal.
	arr := []string{}
	if err := json.Unmarshal([]byte(res), &arr); err != nil {
		c.Logger.Error("failed unmarshaling",
			slog.String("response", res),
			slog.Any("modified_by_user_id", fp.ModifiedByUserID))

		fp.Error = fmt.Sprintf("%v", err)
		fp.Status = domain.StatusError

	} else {
		c.Logger.Debug("successfully unmarshalled",
			slog.Any("modified_by_user_id", fp.ModifiedByUserID))

		fp.Error = ""
		fp.ExerciseNames = arr
		fp.Exercises = make([]*a_d.FitnessPlanExercise, 0) // Reset our list because we will be repopulating it below.

		f := &exercise_s.ExerciseListFilter{
			Cursor:    primitive.NilObjectID,
			PageSize:  1_000_000,
			SortField: "_id",
			SortOrder: 1, // 1=ascending | -1=descending
			Names:     fp.ExerciseNames,
		}
		res, err := c.ExerciseStorer.ListByFilter(sessCtx, f)
		if err != nil {
			c.Logger.Error("list exercises error",
				slog.Any("error", err),
				slog.Any("modified_by_user_id", fp.ModifiedByUserID))
			return err
		}
		if len(res.Results) == 0 {
			c.Logger.Error("list exercises returned no results, failed mapping exericses",
				slog.Any("error", "no results returned"),
				slog.Any("exercise_names", fp.ExerciseNames),
				slog.Any("modified_by_user_id", fp.ModifiedByUserID))
			fp.Status = domain.StatusError
			fp.Error = "system returned no exercises from what openai provided us"
			return errors.New("system returned no exercises from what openai provided us")
		} else {
			for _, e := range res.Results {
				fe := &domain.FitnessPlanExercise{
					ID:   e.ID,
					Name: e.Name,
					// DEVELOPERS NOTE: Add more of your fields here...
				}
				fp.Exercises = append(fp.Exercises, fe)
			}
			c.Logger.Debug("successfully mapped exercises",
				slog.Any("exercises", res.Results),
				slog.Any("modified_by_user_id", fp.ModifiedByUserID))
		}
	}

	fp.ModifiedAt = time.Now()
	if err := c.FitnessPlanStorer.UpdateByID(sessCtx, fp); err != nil {
		c.Logger.Error("database update error",
			slog.Any("error", err),
			slog.Any("modified_by_user_id", fp.ModifiedByUserID))
		return err
	}
	return nil
}

func (c *FitnessPlanControllerImpl) generateFitnessPlanRecommendedExercisesPrompt(sessCtx mongo.SessionContext, fp *domain.FitnessPlan) (string, error) {
	//
	// In the next few lines of code, get all the exercises in our system and
	// get all the exercise names.
	//

	ee, err := c.ExerciseStorer.ListAllAsSelectOption(sessCtx, fp.OrganizationID)
	if err != nil {
		return "", err
	}
	en := ""
	for _, e := range ee {
		en += e.Label + ","
	}
	if en != "" { // Remove the last comma in the string.
		en = strings.TrimSuffix(en, ",")
	}

	//
	// Take our fitness plan and stringify it so we can plug into prompt.
	//

	// Step 1: Convert to a string the `HomeGymEquipment`field.
	var hge string = ""
	for _, s := range fp.HomeGymEquipment {
		hge += domain.HomeGymEquipmentMap[s] + ", "
	}
	if hge != "" { // Remove the last comma in the string.
		hge = strings.TrimSuffix(hge, ",")
	}

	// Step 2: Return yes or no if has access to`commercial-style gym`
	hatcg := "Yes" //TODO

	// Step 3: Convert to a string the `EquipmentAccess` field.
	wah := domain.HasWorkoutsAtHomeMap[fp.HasWorkoutsAtHome]

	// Step 4: Convert clients height.
	height := fmt.Sprintf("%v ft %v in", fp.HeightFeet, fp.HeightInches)

	// Step 5: Date of birth.
	bd := fmt.Sprintf("%v", fp.Birthday)

	// Step 6: Gender
	var gender string
	if fp.Gender == domain.GenderOther {
		gender = fp.GenderOther
	} else {
		if fp.Gender == domain.GenderMale {
			gender = "Male"
		}
		if fp.Gender == domain.GenderFemale {
			gender = "Female"
		}
	}

	// Step 7: IdealWeight
	iw := fmt.Sprintf("%v lbs", fp.IdealWeight)

	// Step 8: PhysicalActivity
	pa := domain.PhysicalActivityMap[fp.PhysicalActivity]

	// Step 9: WorkoutIntensity
	wi := domain.WorkoutIntensityMap[fp.WorkoutIntensity]

	// Step 10: DaysPerWeek
	dpw := domain.DaysPerWeekMap[fp.DaysPerWeek]

	// Step 11: TimePerDay
	tpd := domain.TimePerDayMap[fp.TimePerDay]

	// Step 12: MaxWeeks
	mw := domain.MaxWeekMap[fp.MaxWeeks]

	// Step 13: Goals
	var g string = ""
	for _, s := range fp.Goals {
		g += domain.GoalMap[s] + ", "
	}
	if g != "" { // Remove the last comma in the string.
		g = strings.TrimSuffix(g, ",")
	}

	// Step 14: WorkoutPreferences
	wp := ""
	for _, s := range fp.WorkoutPreferences {
		wp += domain.WorkoutPreferenceMap[s] + ", "
	}
	if wp != "" { // Remove the last comma in the string.
		wp = strings.TrimSuffix(wp, ",")
	}

	prompt := fmt.Sprintf(`
		As a personal trainer, please select from the following 'EXERCISES' as a JSON array of strings that you recommended based on the following 'CLIENT DETAILS'. Do not include any explanations, only provide a RFC8259 compliant JSON response following this format without deviation.

        EXERCISES:
		%s

        CLIENT DETAILS:
     	The equipment that client has access to: %s
		Does client have access to a commercial-style gym: %s
		Does client workout at home: %s
		Client Height: %s
		Client Date of Birth: %s
        Client Gender: %s
		Client ideal weight they wish to achieve (lbs): %s
		Client current level of daily physical activity: %s
		Client current level of intensity in exercise routine: %s
        Client wants wants to train the following number of days per week: %s
        Client's length of time per day that they can train: %s
        Client's goal of number of weeks that they would like the training plan to last: %s
        Client's fitness goals: %s
        Client's workout preference: %s
	`, en, hge, hatcg, wah,
		height, bd, gender, iw, pa, wi, dpw, tpd, mw, g, wp)
	return prompt, nil
}

func (c *FitnessPlanControllerImpl) generateFitnessPlanInstructions(sessCtx mongo.SessionContext, fp *domain.FitnessPlan) error {
	c.Logger.Debug("generating instructions prompt...",
		slog.Any("modified_by_user_id", fp.ModifiedByUserID))

	prompt, err := c.generateFitnessPlanInstructionsPrompt(sessCtx, fp)
	if err != nil {
		c.Logger.Error("failed generating prompt",
			slog.Any("error", err),
			slog.Any("modified_by_user_id", fp.ModifiedByUserID))
		return err
	}
	c.Logger.Debug("prompt ready",
		slog.Any("modified_by_user_id", fp.ModifiedByUserID))
	fp.PromptTwo = prompt

	c.Logger.Debug("submitting prompt to openai",
		slog.Any("modified_by_user_id", fp.ModifiedByUserID))

	res, err := c.OpenAI.CreateChatCompletion(prompt)
	if err != nil { // Error case.
		c.Logger.Error("failed executing `completion` api call to openai",
			slog.Any("error", err),
			slog.Any("modified_by_user_id", fp.ModifiedByUserID))

		// Save error.
		fp.Status = domain.StatusError
		fp.Error = err.Error()
		fp.ModifiedAt = time.Now()
		if err := c.FitnessPlanStorer.UpdateByID(sessCtx, fp); err != nil {
			c.Logger.Error("datebase update error",
				slog.Any("error", err),
				slog.Any("modified_by_user_id", fp.ModifiedByUserID))
			return err
		}

		return err
	}
	c.Logger.Debug("response received",
		slog.Any("modified_by_user_id", fp.ModifiedByUserID))

	// Empty case from OpenAI
	if res == "" {
		err := fmt.Errorf("no response from open ai for fp %v", fp.ID.Hex())

		// Save error.
		fp.Status = domain.StatusError
		fp.Error = err.Error()
		fp.ModifiedAt = time.Now()
		if err := c.FitnessPlanStorer.UpdateByID(sessCtx, fp); err != nil {
			c.Logger.Error("datebase update error",
				slog.Any("error", err),
				slog.Any("modified_by_user_id", fp.ModifiedByUserID))
			return err
		}
		return err
	}

	fp.Status = domain.StatusActive
	fp.Instructions = res
	fp.ModifiedAt = time.Now()
	if err := c.FitnessPlanStorer.UpdateByID(sessCtx, fp); err != nil {
		c.Logger.Error("datebase update error",
			slog.Any("error", err),
			slog.Any("modified_by_user_id", fp.ModifiedByUserID))
		return err
	}
	return nil
}

func (c *FitnessPlanControllerImpl) generateFitnessPlanInstructionsPrompt(sessCtx mongo.SessionContext, fp *domain.FitnessPlan) (string, error) {
	//
	// The following lines of code will iterate through all the exercise names
	// in this fitness plan and put it into a string to plug into our prompt.
	//

	en := ""
	for _, e := range fp.ExerciseNames {
		en += e + ","
	}
	if en != "" { // Remove the last comma in the string.
		en = strings.TrimSuffix(en, ",")
	}

	//
	// Take our fitness plan and stringify it so we can plug into prompt.
	//

	// Step 1: Convert to a string the `HomeGymEquipment`field.
	var hge string = ""
	for _, s := range fp.HomeGymEquipment {
		hge += domain.HomeGymEquipmentMap[s] + ", "
	}
	if hge != "" { // Remove the last comma in the string.
		hge = strings.TrimSuffix(hge, ",")
	}

	// Step 2: Return yes or no if has access to`commercial-style gym`
	hatcg := "Yes" //TODO

	// Step 3: Convert to a string the `EquipmentAccess` field.
	wah := domain.HasWorkoutsAtHomeMap[fp.HasWorkoutsAtHome]

	// Step 4: Convert clients height.
	height := fmt.Sprintf("%v ft %v in", fp.HeightFeet, fp.HeightInches)

	// Step 5: Date of birth.
	bd := fmt.Sprintf("%v", fp.Birthday)

	// Step 6: Gender
	var gender string
	if fp.Gender == domain.GenderOther {
		gender = fp.GenderOther
	} else {
		if fp.Gender == domain.GenderMale {
			gender = "Male"
		}
		if fp.Gender == domain.GenderFemale {
			gender = "Female"
		}
	}

	// Step 7: IdealWeight
	iw := fmt.Sprintf("%v lbs", fp.IdealWeight)

	// Step 8: PhysicalActivity
	pa := domain.PhysicalActivityMap[fp.PhysicalActivity]

	// Step 9: WorkoutIntensity
	wi := domain.WorkoutIntensityMap[fp.WorkoutIntensity]

	// Step 10: DaysPerWeek
	dpw := domain.DaysPerWeekMap[fp.DaysPerWeek]

	// Step 11: TimePerDay
	tpd := domain.TimePerDayMap[fp.TimePerDay]

	// Step 12: MaxWeeks
	mw := domain.MaxWeekMap[fp.MaxWeeks]

	// Step 13: Goals
	var g string = ""
	for _, s := range fp.Goals {
		g += domain.GoalMap[s] + ", "
	}
	if g != "" { // Remove the last comma in the string.
		g = strings.TrimSuffix(g, ",")
	}

	// Step 14: WorkoutPreferences
	wp := ""
	for _, s := range fp.WorkoutPreferences {
		wp += domain.WorkoutPreferenceMap[s] + ", "
	}
	if wp != "" { // Remove the last comma in the string.
		wp = strings.TrimSuffix(wp, ",")
	}

	prompt := fmt.Sprintf(`
		Act as a personal trainer and create a customized workout plan based on the following EXERCISES and DETAILS. In your response, please consider the client's fitness goals, current fitness level, available equipment, time commitment, and preferences in order to provide a safe, effective, and enjoyable workout plan that is tailored to their specific needs and circumstances. All output must be in English.

		EXERCISES:
		%s

        CLIENT DETAILS:
     	The equipment that client has access to: %s
		Does client have access to a commercial-style gym: %s
		Does client workout at home: %s
		Client Height: %s
		Client Date of Birth: %s
        Client Gender: %s
		Client ideal weight they wish to achieve (lbs): %s
		Client current level of daily physical activity: %s
		Client current level of intensity in exercise routine: %s
        Client wants wants to train the following number of days per week: %s
        Client's length of time per day that they can train: %s
        Client's goal of number of weeks that they would like the training plan to last: %s
        Client's fitness goals: %s
        Client's workout preference: %s
	`, en, hge, hatcg, wah, height, bd, gender, iw, pa, wi, dpw, tpd, mw, g, wp)
	return prompt, nil
}
