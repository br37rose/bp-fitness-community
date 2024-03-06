package controller

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"log/slog"

	a_d "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/equipment/datastore"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/config/constants"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/utils/httperror"
)

type EquipmentCreateRequestIDO struct {
	Name string
	No   int8
}

func ValidateCreateRequest(dirtyData *EquipmentCreateRequestIDO) error {
	e := make(map[string]string)

	if dirtyData.Name == "" {
		e["name"] = "missing value"
	}
	if dirtyData.No == 0 {
		e["no"] = "missing value"
	}
	if len(e) != 0 {
		return httperror.NewForBadRequest(&e)
	}
	return nil
}

func (c *EquipmentControllerImpl) Create(ctx context.Context, req *EquipmentCreateRequestIDO) (*a_d.Equipment, error) {
	// Extract from our session the following data.
	orgID := ctx.Value(constants.SessionUserOrganizationID).(primitive.ObjectID)
	orgName := ctx.Value(constants.SessionUserOrganizationName).(string)
	userID := ctx.Value(constants.SessionUserID).(primitive.ObjectID)
	userName := ctx.Value(constants.SessionUserName).(string)

	if err := ValidateCreateRequest(req); err != nil {
		return nil, err
	}

	// Create our meta record in the database.
	res := &a_d.Equipment{
		OrganizationID:     orgID,
		OrganizationName:   orgName,
		ID:                 primitive.NewObjectID(),
		CreatedAt:          time.Now(),
		CreatedByUserName:  userName,
		CreatedByUserID:    userID,
		ModifiedAt:         time.Now(),
		ModifiedByUserName: userName,
		ModifiedByUserID:   userID,
		Name:               req.Name,
		No:                 req.No,
		Status:             a_d.StatusActive,
	}
	err := c.EquipmentStorer.Create(ctx, res)
	if err != nil {
		c.Logger.Error("equipment create error", slog.Any("error", err))
		return nil, err
	}
	return res, nil
}
