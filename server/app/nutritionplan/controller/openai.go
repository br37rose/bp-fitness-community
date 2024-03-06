package controller

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/mongo"

	domain "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/nutritionplan/datastore"
)

func (c *NutritionPlanControllerImpl) generateNutritionPlanInBackground(ctx context.Context, fp *domain.NutritionPlan) {
	c.Logger.Debug("getting nutrition plan instrucitons from openai",
		slog.Any("modified_by_user_id", fp.ModifiedByUserID))

	c.Kmutex.Lockf("openai-api-for-nutrition-plan-%s", fp.ID.Hex())
	defer c.Kmutex.Unlockf("openai-api-for-nutrition-plan-%s", fp.ID.Hex())

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

		if err := c.generateNutritionPlanInstructions(sessCtx, fp); err != nil {
			c.Logger.Error("failed getting from openai",
				slog.Any("modified_by_user_id", fp.ModifiedByUserID))
			return nil, err
		}

		// Set the state.
		fp.Status = domain.StatusActive
		fp.ModifiedAt = time.Now()
		if err := c.NutritionPlanStorer.UpdateByID(sessCtx, fp); err != nil {
			c.Logger.Error("response received",
				slog.Any("modified_by_user_id", fp.ModifiedByUserID))
			return nil, err
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

func (c *NutritionPlanControllerImpl) generateNutritionPlanInstructions(sessCtx mongo.SessionContext, fp *domain.NutritionPlan) error {
	c.Logger.Debug("generating instructions prompt...",
		slog.Any("modified_by_user_id", fp.ModifiedByUserID))

	prompt, err := c.generateNutritionPlanInstructionsPrompt(sessCtx, fp)
	if err != nil {
		c.Logger.Error("failed generating prompt",
			slog.Any("error", err),
			slog.Any("modified_by_user_id", fp.ModifiedByUserID))
		return err
	}
	c.Logger.Debug("prompt ready",
		slog.Any("modified_by_user_id", fp.ModifiedByUserID))
	fp.Prompt = prompt

	c.Logger.Debug("submitting prompt to openai",
		slog.Any("modified_by_user_id", fp.ModifiedByUserID))

	res, err := c.OpenAI.CreateChatCompletion(prompt)
	if err != nil {
		c.Logger.Error("failed executing `completion` api call to openai",
			slog.Any("error", err),
			slog.Any("modified_by_user_id", fp.ModifiedByUserID))

		// Save error.
		fp.Status = domain.StatusError
		fp.Error = err.Error()
		fp.ModifiedAt = time.Now()
		if err := c.NutritionPlanStorer.UpdateByID(sessCtx, fp); err != nil {
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
		if err := c.NutritionPlanStorer.UpdateByID(sessCtx, fp); err != nil {
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
	if err := c.NutritionPlanStorer.UpdateByID(sessCtx, fp); err != nil {
		c.Logger.Error("datebase update error",
			slog.Any("error", err),
			slog.Any("modified_by_user_id", fp.ModifiedByUserID))
		return err
	}
	return nil
}

func (c *NutritionPlanControllerImpl) generateNutritionPlanInstructionsPrompt(sessCtx mongo.SessionContext, fp *domain.NutritionPlan) (string, error) {
	//
	// Take our nutrition plan and stringify it so we can plug into prompt.
	//

	// Step 4: Convert clients height.
	height := fmt.Sprintf("%v ft %v in", fp.HeightFeet, fp.HeightInches)

	// Step 5: Date of birth.
	bd := fmt.Sprintf("%s", fp.Birthday)

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

	// Step 9: Nutritional Goals
	var ng string = ""
	for _, s := range fp.Goals {
		ng += domain.GoalMap[s] + ", "
	}
	if ng != "" { // Remove the last comma in the string.
		ng = strings.TrimSuffix(ng, ",")
	}

	// Step 10: Allergies
	var allergies string = ""
	if fp.HasAllergies == domain.HasAllergiesYes {
		allergies = fp.Allergies
	} else {
		allergies = "No allergies"
	}

	// Step 11: How many meals do you typically eat in a day
	mpd := domain.MealsPerDayMap[fp.MealsPerDay]

	// Step 12: How often do you consume fast food or junk food?
	cff := domain.ConsumeFrequencyMap[fp.ConsumeJunkFood]

	// Step 13: How often do you consume fruits and vegetables?
	cfav := domain.ConsumeFrequencyMap[fp.ConsumeFruitsAndVegetables]

	// Step 14: Has Intermittent Fasting
	hif := domain.HasIntermittentFastingMap[fp.HasIntermittentFasting]

	// Step 15: max weeks
	mw := domain.MaxWeekMap[fp.MaxWeeks]

	prompt := fmt.Sprintf(`
		CONTEXT
		As a licensed dietician and nutritionist, you have a deep understanding of what the body needs to achieve specific goals and results, including avoiding foods that can cause allergies.

		With your expertise, you can help clients create personalized nutrition plans that are tailored to their individual needs and goals. You understand the importance of a balanced diet and can provide guidance on macronutrient ratios, micronutrient intake, and other factors that impact overall health and performance.

		Furthermore, you stay up-to-date with the latest research and developments in the field of nutrition, ensuring that you are always providing the most current and effective advice to your clients.

		Overall, your knowledge and expertise as a licensed dietician and nutritionist make you a valuable asset to anyone looking to optimize their nutrition and achieve their health and fitness goals.

		CLIENT DATA
		Date of Birth: %s
		Height: %s
        Gender: %s
		Ideal weight they wish to achieve (lbs): %s
		Current level of daily physical activity: %s
		Current level of intensity in exercise routine: %s
		Nutritional Goals: %s
		Allergies: %s
		How many meals do you typically eat in a day: %s
		How often do you consume fast food or junk food: %s
		How often do you consume fruits and vegetables: %s
		Include intermittent fasting in the nutrition plan: %s

		Based on above assessment answered you will also recommend the best nutrition plans for %s for me to follow remembering any allergies I said above and my goals stated above. I need to know ingredients of food, calories & macros of each meal and day, instructions on how to make meals, and grocery list to buy meals using the grocery store layout to be as efficient as possible.
	`, bd, height, gender, iw, pa, wi, ng, allergies, mpd, cff, cfav, hif, mw)
	return prompt, nil
}
