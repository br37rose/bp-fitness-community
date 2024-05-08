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
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/question/datastore"
	u_d "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/user/datastore"
	"github.com/sashabaranov/go-openai"
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

			// Set the state to in progress.
			// fp.Status = domain.S
			// fp.ModifiedAt = time.Now()
			// if err := c.FitnessPlanStorer.UpdateByID(sessCtx, fp); err != nil {
			// 	c.Logger.Error("response received",
			// 		slog.Any("modified_by_user_id", fp.ModifiedByUserID))
			// 	return nil, err
			// }
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

	excercises, err := c.listAvailableExcercises(sessCtx, fp)
	if err != nil {
		c.Logger.Error("failed to list excecises", slog.Any("error", err))
		return err
	}

	res, err := c.OpenAI.CreateFitnessPlan(sessCtx, prompt, excercises)
	if err != nil {
		c.Logger.Error("failed to `CreateFitnessPlan` api call to openai",
			slog.Any("error", err))

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
	}

	// res, err := c.OpenAI.CreateChatCompletion(prompt)
	// if err != nil { // Error case.
	// 	c.Logger.Error("failed executing `completion` api call to openai",
	// 		slog.Any("error", err),
	// 		slog.Any("modified_by_user_id", fp.ModifiedByUserID))

	// 	// Save error.
	// 	fp.Status = domain.StatusError
	// 	fp.Error = err.Error()
	// 	fp.ModifiedAt = time.Now()
	// 	if err := c.FitnessPlanStorer.UpdateByID(sessCtx, fp); err != nil {
	// 		c.Logger.Error("datebase update error",
	// 			slog.Any("error", err),
	// 			slog.Any("modified_by_user_id", fp.ModifiedByUserID))
	// 		return err
	// 	}

	// 	return err
	// }
	// c.Logger.Debug("response received",
	// 	slog.Any("modified_by_user_id", fp.ModifiedByUserID))

	// Empty case from OpenAI
	if res.ID == "" || res.ThreadID == "" {
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

	status := 0
	switch res.Status {
	case openai.RunStatusCancelling, openai.RunStatusExpired, openai.RunStatusFailed:
		status = domain.StatusError
	default:
		status = domain.StatusPending
	}

	fp.Status = int8(status)
	fp.RunnerID = res.ID
	fp.ThreadID = res.ThreadID
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

	// Step 1: Convert to a string the `HomeGymEquipment`field.
	userProfile, err := c.UserStorer.GetByID(sessCtx, fp.UserID)
	if err != nil {
		return "", err
	}
	AnsMap := make(map[string]u_d.Answer, 0)
	for _, v := range userProfile.OnboardingAnswers {
		AnsMap[v.QuestionID.Hex()] = *v
	}
	qstRes, err := c.QuestionController.ListByFilter(sessCtx,
		&datastore.QuestionListFilter{
			PageSize:  10000,
			SortField: "_id",
			SortOrder: 1,
			Status:    true},
	)
	prompt := `CLIENT DETAILS AS QUESTIONS AND ANSWERS:`
	for _, qstn := range qstRes.Results {
		if ans, ok := AnsMap[qstn.ID.Hex()]; ok {
			prompt = fmt.Sprintf("%s,\n %s (%s) : %s", prompt, qstn.Title, qstn.Subtitle, strings.Join(ans.Answers, ","))
		}
	}

	return prompt, nil
}

func (c *FitnessPlanControllerImpl) listAvailableExcercises(ctx context.Context, fp *domain.FitnessPlan) (string, error) {

	exercises, err := c.ExerciseStorer.ListByFilter(ctx, &exercise_s.ExerciseListFilter{
		Cursor:          primitive.NilObjectID,
		SortField:       "created_at",
		SortOrder:       -1,
		OrganizationID:  fp.OrganizationID,
		ExcludeArchived: true,
		PageSize:        500,
		// Gender:          a_d.GenderMap[fp.Gender],
	})
	if err != nil {
		c.Logger.Error("304")
		return "", err
	}

	if len(exercises.Results) == 0 {
		return "", nil
	}

	exerciseList := make([]string, 0, len(exercises.Results))
	for _, exercise := range exercises.Results {
		exerciseList = append(exerciseList, fmt.Sprintf("{ id=%s, name=%s, description=%s }", exercise.ID.Hex(), exercise.Name, exercise.Description))
	}

	return strings.Join(exerciseList, ","), nil

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
	userProfile, err := c.UserStorer.GetByID(sessCtx, fp.UserID)
	if err != nil {
		return "", err
	}
	AnsMap := make(map[string]u_d.Answer, 0)
	for _, v := range userProfile.OnboardingAnswers {
		AnsMap[v.QuestionID.Hex()] = *v
	}
	qstRes, err := c.QuestionController.ListByFilter(sessCtx,
		&datastore.QuestionListFilter{
			PageSize:  10000,
			SortField: "_id",
			SortOrder: 1,
			Status:    true},
	)
	clientProfile := ""
	for _, qstn := range qstRes.Results {
		if ans, ok := AnsMap[qstn.ID.Hex()]; ok {
			clientProfile = fmt.Sprintf("%s,\n %s (%s) : %s", clientProfile, qstn.Title, qstn.Subtitle, ans)
		}
	}

	prompt := fmt.Sprintf(`
		As a personal trainer, please select from the following 'EXERCISES' as a JSON array of strings that you recommended based on the following 'CLIENT DETAILS'. Do not include any explanations, only provide a RFC8259 compliant JSON response following this format without deviation.

        EXERCISES:
		%s.
		CLIENT DETAILS AS QUESTIONS AND ANSWERS:
		%s
	`, en, clientProfile)
	return prompt, nil
}
