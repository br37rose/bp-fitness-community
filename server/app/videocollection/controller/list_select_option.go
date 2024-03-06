package controller

import (
	"context"

	vcol_d "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/videocollection/datastore"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/config/constants"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log/slog"
)

func (c *VideoCollectionControllerImpl) ListAsSelectOptionByFilter(ctx context.Context, f *vcol_d.VideoCollectionListFilter) ([]*vcol_d.VideoCollectionAsSelectOption, error) {
	// Extract from our session the following data.
	userID, _ := ctx.Value(constants.SessionUserID).(primitive.ObjectID)
	userOrgID, _ := ctx.Value(constants.SessionUserOrganizationID).(primitive.ObjectID)

	// Force tenancy on list.
	c.Logger.Debug("fetching videocategorys now...", slog.Any("userID", userID))
	f.OrganizationID = userOrgID

	// List.
	m, err := c.VideoCollectionStorer.ListAsSelectOptionByFilter(ctx, f)
	if err != nil {
		c.Logger.Error("database list by filter error", slog.Any("error", err))
		return nil, err
	}
	return m, err
}
