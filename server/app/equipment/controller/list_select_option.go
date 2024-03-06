package controller

import (
	"context"

	domain "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/equipment/datastore"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/config/constants"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log/slog"
)

func (c *EquipmentControllerImpl) ListAsSelectOptionByFilter(ctx context.Context, f *domain.EquipmentListFilter) ([]*domain.EquipmentAsSelectOption, error) {
	// Extract from our session the following data.
	userID, _ := ctx.Value(constants.SessionUserID).(primitive.ObjectID)
	userOrgID, _ := ctx.Value(constants.SessionUserOrganizationID).(primitive.ObjectID)

	// Force tenancy on list.
	c.Logger.Debug("fetching equipments now...", slog.Any("userID", userID))
	f.OrganizationID = userOrgID

	// List.
	m, err := c.EquipmentStorer.ListAsSelectOptionByFilter(ctx, f)
	if err != nil {
		c.Logger.Error("database list by filter error", slog.Any("error", err))
		return nil, err
	}
	c.Logger.Debug("fetched equipments", slog.Any("m", m))
	return m, err
}
