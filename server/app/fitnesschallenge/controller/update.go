package controller

import (
	"context"
	"errors"
	"time"

	"log/slog"

	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/fitnesschallenge/datastore"
	fc "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/fitnesschallenge/datastore"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/config/constants"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/utils"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/utils/httperror"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type FitnessChallengeUpdateRequestIDO struct {
	ID              primitive.ObjectID   `json:"id"`
	Name            string               `json:"name,omitempty"`
	Description     string               `json:"description,omitempty"`
	DurationInWeeks int64                `json:"duration,omitempty"`
	OrganizationID  primitive.ObjectID   `json:"organization_id,omitempty"`
	Rules           []int                `json:"rules,omitempty"`
	StartOn         time.Time            `json:"start_on,omitempty"`
	Users           []primitive.ObjectID `json:"users,omitempty"`
}

func ValidateUpdateRequest(dirtyData *FitnessChallengeUpdateRequestIDO) error {
	e := make(map[string]string)
	if dirtyData.ID.IsZero() {
		e["id"] = "missing value"
	}

	if len(e) != 0 {
		return httperror.NewForBadRequest(&e)
	}
	return nil
}

func (c *FitnessChallengeControllerImpl) UpdateByID(
	ctx context.Context,
	req *FitnessChallengeUpdateRequestIDO) (*datastore.FitnessChallenge, error) {

	userID, _ := ctx.Value(constants.SessionUserID).(primitive.ObjectID)
	userName, _ := ctx.Value(constants.SessionUserName).(string)

	if err := ValidateUpdateRequest(req); err != nil {
		c.Logger.Error("validation failure", slog.Any("error", err))
		return nil, err
	}
	os, err := c.Storer.GetByID(ctx, req.ID)
	if err != nil {
		return nil, err
	}
	if os == nil {
		return nil, errors.New("no data found")
	}
	ns := datastore.FitnessChallenge{
		ModifiedByUserID:   userID,
		ModifiedByUserName: userName,
		ModifiedAt:         time.Now().UTC(),
		ID:                 req.ID,
		OrganizationID:     req.OrganizationID,
		Name:               req.Name,
		Description:        req.Description,
		StartTime:          req.StartOn,
		DurationInWeeks:    req.DurationInWeeks,
		CreatedAt:          os.CreatedAt,
		CreatedByUserID:    os.CreatedByUserID,
		CreatedByUserName:  os.CreatedByUserName,
		EndTime:            os.EndTime,
		Status:             datastore.FitnessChallengeStatusActive,
	}
	rules := []*datastore.ChallengeRules{}
	for _, v := range req.Rules {
		if rule, ok := fc.Challenges[v]; ok {
			rules = append(rules, &rule)
		}
	}
	userids := make([]primitive.ObjectID, 0)
	usernames := make([]string, 0)
	for _, v := range req.Users {
		if v.IsZero() {
			continue
		}
		user, err := c.UserStorer.GetByID(ctx, v)
		if err != nil {
			return nil, err
		}
		userids = append(userids, user.ID)
		usernames = append(usernames, user.Name)
	}
	ns.Rules = rules
	ns.UserIDs = userids
	ns.UserNames = usernames
	if err := c.Storer.UpdateByID(ctx, &ns); err != nil {
		c.Logger.Error("database update by id error", slog.Any("error", err))
		return nil, err
	}

	return &ns, nil
}

func (c *FitnessChallengeControllerImpl) ChangeParticipationStatus(
	ctx context.Context, id primitive.ObjectID) (*datastore.FitnessChallenge, error) {

	userID, ok := ctx.Value(constants.SessionUserID).(primitive.ObjectID)
	if !ok {
		return nil, httperror.NewForBadRequestWithSingleField("user", "couldnot find the user")
	}
	userName, ok := ctx.Value(constants.SessionUserName).(string)
	if !ok {
		return nil, httperror.NewForBadRequestWithSingleField("user", "couldnot find the user")
	}

	os, err := c.Storer.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	registerd := false
	for _, v := range os.UserIDs {
		if v == userID {
			registerd = true
			break
		}
	}
	if !registerd {
		os.UserIDs = append(os.UserIDs, userID)
		os.UserNames = append(os.UserNames, userName)
	} else {
		os.UserIDs = utils.RemoveObjectIDFromArray(os.UserIDs, userID)
		os.UserNames = utils.RemoveElementFromArray(os.UserNames, userName)
	}

	err = c.Storer.UpdateByID(ctx, os)

	return os, err
}
