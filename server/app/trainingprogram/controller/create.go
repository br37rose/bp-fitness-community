package controller

import (
	"context"
	"fmt"
	"time"

	tp_s "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/trainingprogram/datastore"
	u_s "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/user/datastore"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/config/constants"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/utils/httperror"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TrainingProgramCreateRequestIDO struct {
	Name           string             `json:"name,omitempty"`
	Description    string             `json:"description,omitempty"`
	Phases         int64              `json:"phases,omitempty,string"`
	Weeks          int64              `json:"weeks,omitempty,string"`
	OrganizationID primitive.ObjectID `json:"organization_id"`
	UserId         primitive.ObjectID `json:"user_id"`
}

func (impl *TrainingprogramControllerImpl) Create(ctx context.Context, req *TrainingProgramCreateRequestIDO) (*tp_s.TrainingProgram, error) {
	tp, err := impl.TrainingProgramFromrequest(ctx, req)
	if err != nil {
		return nil, err
	}

	err = impl.TpStorer.Create(ctx, tp)
	return tp, err
}
func (impl *TrainingprogramControllerImpl) TrainingProgramFromrequest(ctx context.Context, requestData *TrainingProgramCreateRequestIDO) (*tp_s.TrainingProgram, error) {

	orgID, _ := ctx.Value(constants.SessionUserOrganizationID).(primitive.ObjectID)
	userID, _ := ctx.Value(constants.SessionUserID).(primitive.ObjectID)
	userName, _ := ctx.Value(constants.SessionUserName).(string)
	userRole, _ := ctx.Value(constants.SessionUserRole).(int8)

	if userRole == u_s.UserRoleMember {
		requestData.UserId = userID
	} else if requestData.UserId.IsZero() {
		return nil, httperror.NewForBadRequestWithSingleField("userid", "not provided")
	}

	tp := tp_s.TrainingProgram{
		ID:                 primitive.NewObjectID(),
		UserID:             requestData.UserId,
		Name:               requestData.Name,
		Description:        requestData.Description,
		Phases:             requestData.Phases,
		Weeks:              requestData.Weeks,
		DurationInWeeks:    requestData.Phases * requestData.Weeks,
		CreatedAt:          time.Now().UTC(),
		ModifiedAt:         time.Now().UTC(),
		CreatedByUserID:    userID,
		CreatedByUserName:  userName,
		ModifiedByUserID:   userID,
		ModifiedByUserName: userName,
		Status:             tp_s.TrainingProgramStatusActive,
	}

	if requestData.OrganizationID.IsZero() {
		requestData.OrganizationID = orgID
	}

	var weekStart int64 = 1
	var weekEnd int64 = tp.Weeks
	var trainingPhases []*tp_s.TrainingPhase

	for phase := int64(1); phase <= tp.Phases; phase++ {
		tphase := &tp_s.TrainingPhase{
			ID:          primitive.NewObjectID(),
			Name:        fmt.Sprintf("Week %d - %d", weekStart, weekEnd),
			Description: tp.Description,
			Phase:       phase,
			StartWeek:   weekStart,
			EndWeek:     weekEnd,
		}
		trainingPhases = append(trainingPhases, tphase)
		// Get the weeks range of the next phase.
		weekStart = weekEnd + 1
		weekEnd = weekStart + tp.Weeks
	}
	tp.TrainingPhases = trainingPhases
	return &tp, nil
}
