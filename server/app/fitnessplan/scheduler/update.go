package scheduler

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	fp_d "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/fitnessplan/datastore"
	"github.com/sashabaranov/go-openai"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (impl *fitnessPlanSchedulerImpl) RunEveryMinuteUpdateFitnessPlans() error {
	impl.Logger.Debug("scheduled: update fitness plans", slog.String("interval", "every minute"))
	err := impl.EventScheduler.ScheduleEveryMinuteFunc(func() {
		impl.updateFitnessPlans()
	})
	if err != nil {
		impl.Logger.Error("error with pinging scheduler", slog.Any("err", err))
	}
	return nil
}

// exercisePlan represents the structured fitness plan returned by the OpenAI API.
type exercisePlan struct {
	WeeklyFitnessPlans []*fp_d.WeeklyFitnessPlan   `json:"weekly_fitness_plans"`
	MainExercises      []*fp_d.FitnessPlanExercise `json:"main_exercises"`
	Instructions       string                      `json:"instructions"`
}

// updateFitnessPlans retrieves queued fitness plans, processes them using the OpenAI API, and updates their status and details.
func (impl *fitnessPlanSchedulerImpl) updateFitnessPlans() {
	ctx := context.Background()

	fitnessPlans, err := impl.FitnessPlanStorer.ListByFilter(ctx, &fp_d.FitnessPlanListFilter{
		Cursor:     primitive.NilObjectID,
		PageSize:   100000,
		SortField:  "created_at",
		SortOrder:  1,
		StatusList: []int8{fp_d.StatusPending, fp_d.StatusInProgress},
	})
	if err != nil {
		impl.Logger.Error("Error listing fitness plans", slog.Any("error", err))
		return
	}

	for _, fp := range fitnessPlans.Results {
		response, err := impl.OpenAIConnector.GetRunner(ctx, fp.ThreadID, fp.RunnerID)
		if err != nil {
			impl.Logger.Error("Error retrieving runner", slog.Any("error", err))
			continue
		}

		var (
			plan   *exercisePlan
			status int8
		)
		switch response.Status {
		case openai.RunStatusRequiresAction:
			for _, tool := range response.RequiredAction.SubmitToolOutputs.ToolCalls {
				if tool.Function.Name == "generateFitnessPlan" {
					if err := json.Unmarshal([]byte(tool.Function.Arguments), &plan); err != nil {
						impl.Logger.Error("Error unmarshaling data", slog.Any("error", err))
						status = fp_d.StatusError
						fp.Error = fmt.Sprintf("Error while unmarshaling data=%s", err.Error())
						break
					}

					status = fp_d.StatusInProgress

					_, err := impl.OpenAIConnector.SubmitToolOutputs(ctx, fp.ThreadID, fp.RunnerID, openai.SubmitToolOutputsRequest{
						ToolOutputs: []openai.ToolOutput{
							{
								ToolCallID: tool.ID,
								Output:     tool.Function.Arguments,
							},
						},
					})
					if err != nil {
						impl.Logger.Error("Error submitting tool output", slog.Any("error", err))
						continue
					}

				}
			}
		case openai.RunStatusCompleted:

			status = fp_d.StatusActive
		case openai.RunStatusCancelling, openai.RunStatusExpired, openai.RunStatusFailed:
			status = fp_d.StatusError

		}

		//for debugging
		if response.LastError != nil {
			fp.Error = response.LastError.Message
		}

		if plan != nil {
			fp.Exercises = plan.MainExercises
			fp.Instructions = plan.Instructions
			fp.WeeklyFitnessPlans = plan.WeeklyFitnessPlans
		}

		if status != 0 {
			fp.Status = status
			fp.ModifiedAt = time.Now()
			if err := impl.FitnessPlanStorer.UpdateByID(ctx, fp); err != nil {
				impl.Logger.Error("Error updating the fitness plan", slog.Any("error", err))
				continue
			}
		}

	}
}
