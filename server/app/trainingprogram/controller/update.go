package controller

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/trainingprogram/datastore"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/utils/httperror"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PhaseUpdateRequestIDO struct {
	Phases []*PhaseRequestIDO `bson:"phases" json:"phases"`
}

type PhaseRequestIDO struct {
	PhaseID  primitive.ObjectID        `bson:"phase_id" json:"phase_id"`
	Phase    int64                     `bson:"phase" json:"phase"`
	Routines []*PhaseRoutinerequestIDO `bson:"routines" json:"routines"`
}

type PhaseRoutinerequestIDO struct {
	WorkoutId   primitive.ObjectID `bson:"workout_id" json:"workout_id"`
	RoutineDay  int8               `bson:"routine_day" json:"routine_day"`
	RoutineWeek int8               `bson:"routine_week" json:"routine_week"`
}

func (c *TrainingprogramControllerImpl) UpdateTPPhase(ctx context.Context, TPid primitive.ObjectID, req PhaseUpdateRequestIDO) (
	*datastore.TrainingProgram,
	error) {

	TP, err := c.TpStorer.GetByID(ctx, TPid)
	if err != nil {
		return nil, err
	}
	if TP == nil {
		return nil, httperror.NewForSingleField(http.StatusBadRequest, "Program id", "invalid program id")
	}
	oldPhaseMap := make(map[string]datastore.TrainingPhase, 0)
	for _, v := range TP.TrainingPhases {
		oldPhaseMap[v.ID.Hex()] = *v
	}
	ns := []*datastore.TrainingPhase{}
	for _, v := range req.Phases {
		tp := datastore.TrainingPhase{}
		routines := []*datastore.TrainingRoutine{}
		tp.ID = v.PhaseID
		tp.Phase = v.Phase
		old, ok := oldPhaseMap[v.PhaseID.Hex()]
		if ok {
			tp.Description = old.Description
			tp.Name = old.Name
			tp.EndTime = old.EndTime
			tp.EndWeek = old.EndWeek
			tp.StartTime = old.StartTime
			tp.Type = old.Type
		}
		for i, r := range v.Routines {
			wk, err := c.WorkoutStorer.GetByID(ctx, r.WorkoutId)
			if err != nil {
				return nil, err
			}
			routines = append(routines, &datastore.TrainingRoutine{
				ID:          primitive.NewObjectID(),
				WorkoutID:   wk.ID,
				Workout:     wk,
				OrderNumber: int64(i),
				Phase:       v.Phase,
				TrainingDays: []*datastore.TrainingDay{
					{Day: int64(r.RoutineDay), Week: int64(r.RoutineWeek)},
				},
			})
		}
		tp.TrainingRoutines = routines
		ns = append(ns, &tp)
	}

	err = c.TpStorer.UpdatePhase(ctx, TPid, ns)
	if err != nil {
		c.Logger.Error("database get by id error", slog.Any("error", err))
		return nil, err
	}
	return &datastore.TrainingProgram{}, err
}
