package controller

import (
	"context"
	"log/slog"
	"sync"

	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/trainingprogram/datastore"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (c *TrainingprogramControllerImpl) GetByID(ctx context.Context, id primitive.ObjectID) (*datastore.TrainingProgram, error) {
	tp, err := c.TpStorer.GetByID(ctx, id)
	if err != nil {
		c.Logger.Error("database get by id error", slog.Any("error", err))
		return nil, err
	}
	var wg sync.WaitGroup
	var errorChan = make(chan error, len(tp.TrainingPhases))
	for _, p := range tp.TrainingPhases {
		for i, v := range p.TrainingRoutines {
			wg.Add(1)
			go c.fetchWorkout(ctx, *v, p, i, errorChan, &wg)
		}
	}
	// this will wait till all the goroutines are  complete and will close the channel
	go func() {
		wg.Wait()
		close(errorChan)
	}()
	// parallelly range over error chan till its closed by goroutine above
	for err := range errorChan {
		if err != nil {
			return tp, err
		}
	}
	return tp, err
}

func (c *TrainingprogramControllerImpl) fetchWorkout(
	ctx context.Context,
	routine datastore.TrainingRoutine,
	phase *datastore.TrainingPhase,
	index int,
	resultChannel chan<- error,
	wg *sync.WaitGroup,
) {
	defer wg.Done()
	workout, err := c.WorkoutController.GetByID(ctx, routine.WorkoutID)
	if err != nil {
		resultChannel <- err
		return
	}
	phase.TrainingRoutines[index].Workout = workout
	resultChannel <- nil
}
