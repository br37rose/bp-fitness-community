package controller

import (
	"context"
	"log/slog"

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
					c.Logger.Error("excercise does not exist", slog.Any("error", err))
					return nil, err
				}
				dailyPlan.PlanDetails[i].Name = e.Name
				dailyPlan.PlanDetails[i].Description = e.Description
				dailyPlan.PlanDetails[i].VideoURL = e.VideoObjectURL
				dailyPlan.PlanDetails[i].ThumbnailURL = e.ThumbnailObjectURL

			}
		}
	}

	//This is to fix the issue with the expiration of url stored in the fitness Plan
	for index, excercise := range m.Exercises {
		e, err := c.ExcerciseContr.GetByID(ctx, excercise.ID)
		if err != nil {
			c.Logger.Error("excercise does not exist", slog.Any("error", err))
			return nil, err
		}
		m.Exercises[index].VideoURL = e.VideoObjectURL
		m.Exercises[index].ThumbnailURL = e.ThumbnailObjectURL
		m.Exercises[index].Description = e.Description
	}

	return m, err
}
