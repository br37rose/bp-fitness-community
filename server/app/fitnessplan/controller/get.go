package controller

import (
	"context"
	"log/slog"

	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/exercise/datastore"
	domain "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/fitnessplan/datastore"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (c *FitnessPlanControllerImpl) GetByID(ctx context.Context, id primitive.ObjectID) (*domain.FitnessPlan, error) {
	// Retrieve from our database the record for the specific id.
	m, err := c.FitnessPlanStorer.GetByID(ctx, id)
	if err != nil {
		c.Logger.Error("database get by id error", slog.Any("error", err))
		return nil, err
	}

	for _, wPlan := range m.WeeklyFitnessPlans {
		for _, dailyPlan := range wPlan.DailyPlans {
			for i, dailyExercise := range dailyPlan.PlanDetails {
				e, err := c.ExcerciseContr.GetByID(ctx, dailyExercise.ID)
				if err != nil {
					ex, err := c.ExcerciseContr.ListByFilter(ctx, &datastore.ExerciseListFilter{
						ExcludeArchived: true,
						SearchText:      dailyExercise.Name,
						// Gender:          domain.GenderMap[m.Gender], TODO - get gender from user profile
						SortField: "created_at",
						SortOrder: -1,
						PageSize:  1,
					})

					if err != nil {
						c.Logger.Error("excercise does not exist", slog.Any("error", err))
						continue
					}
					if len(ex.Results) > 0 {
						e = ex.Results[0]
					}
				}

				if e != nil {
					dailyPlan.PlanDetails[i].VideoURL = e.VideoObjectURL
					if e.VideoObjectURL == "" {
						dailyPlan.PlanDetails[i].VideoURL = e.VideoURL
					}
					dailyPlan.PlanDetails[i].Name = e.Name
					dailyPlan.PlanDetails[i].Description = e.Description
					dailyPlan.PlanDetails[i].ThumbnailURL = e.ThumbnailObjectURL
					dailyPlan.PlanDetails[i].VideoType = e.VideoType
				}

			}
		}
	}

	//This is to fix the issue with the expiration of url stored in the fitness Plan
	for index, excercise := range m.Exercises {
		e, err := c.ExcerciseContr.GetByID(ctx, excercise.ID)
		if err != nil {
			ex, err := c.ExcerciseContr.ListByFilter(ctx, &datastore.ExerciseListFilter{
				ExcludeArchived: true,
				SearchText:      excercise.Name,
				// Gender:          domain.GenderMap[m.Gender],TODO - get gender from user profile
				SortField: "created_at",
				SortOrder: -1,
				PageSize:  1,
			})

			if err != nil {
				c.Logger.Error("excercise does not exist", slog.Any("error", err))
				continue
			}
			if len(ex.Results) > 0 {
				e = ex.Results[0]
			}
		}

		if e != nil {
			m.Exercises[index].VideoURL = e.VideoObjectURL
			if e.VideoObjectURL == "" {
				m.Exercises[index].VideoURL = e.VideoURL
			}
			m.Exercises[index].ThumbnailURL = e.ThumbnailObjectURL
			m.Exercises[index].Description = e.Description
			m.Exercises[index].VideoType = e.VideoType
		}

	}

	return m, err
}
