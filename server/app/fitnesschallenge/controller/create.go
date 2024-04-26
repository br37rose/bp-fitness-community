package controller

import (
	"context"
	"time"

	fc "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/fitnesschallenge/datastore"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/config/constants"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type FitnessChallengeCreateRequestIDO struct {
	Name            string               `json:"name,omitempty"`
	Description     string               `json:"description,omitempty"`
	DurationInWeeks int64                `json:"duration,omitempty"`
	OrganizationID  primitive.ObjectID   `json:"organization_id"`
	Rules           []int                `json:"rules"`
	StartOn         time.Time            `json:"start_on"`
	Users           []primitive.ObjectID `json:"users"`
}

func (impl *FitnessChallengeControllerImpl) Create(ctx context.Context, req *FitnessChallengeCreateRequestIDO) (*fc.FitnessChallenge, error) {
	tp, err := impl.TrainingProgramFromrequest(ctx, req)
	if err != nil {
		return nil, err
	}

	err = impl.Storer.Create(ctx, tp)
	return tp, err
}
func (impl *FitnessChallengeControllerImpl) TrainingProgramFromrequest(
	ctx context.Context,
	requestData *FitnessChallengeCreateRequestIDO) (*fc.FitnessChallenge, error) {

	userID, _ := ctx.Value(constants.SessionUserID).(primitive.ObjectID)
	userName, _ := ctx.Value(constants.SessionUserName).(string)

	rules := []*fc.ChallengeRules{}
	for _, v := range requestData.Rules {
		if rule, ok := fc.Challenges[v]; ok {
			rules = append(rules, &rule)
		}
	}
	userids := make([]primitive.ObjectID, 0)
	usernames := make([]string, 0)
	for _, v := range requestData.Users {
		user, err := impl.UserStorer.GetByID(ctx, v)
		if err != nil {
			return nil, err
		}
		userids = append(userids, user.ID)
		usernames = append(usernames, user.Name)
	}

	fc := fc.FitnessChallenge{
		ID:                primitive.NewObjectID(),
		Name:              requestData.Name,
		Description:       requestData.Description,
		OrganizationID:    requestData.OrganizationID,
		DurationInWeeks:   requestData.DurationInWeeks,
		StartTime:         requestData.StartOn,
		Rules:             rules,
		CreatedAt:         time.Now().UTC(),
		CreatedByUserID:   userID,
		CreatedByUserName: userName,
		Status:            fc.FitnessChallengeStatusActive,
		UserIDs:           userids,
		UserNames:         usernames,
	}

	return &fc, nil
}
