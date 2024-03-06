package controller

import (
	"context"

	o_s "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/organization/datastore"
	user_s "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/user/datastore"
	"log/slog"
)

func (c *OrganizationControllerImpl) updateRelated(ctx context.Context, org *o_s.Organization) error {

	if err := c.updateRelatedUsers(ctx, org); err != nil {
		return err
	}

	return nil
}

func (c *OrganizationControllerImpl) updateRelatedUsers(ctx context.Context, org *o_s.Organization) error {
	f := &user_s.UserListFilter{
		OrganizationID: org.ID,
		SortField:      "_id",
		SortOrder:      1, //1=ascending
		PageSize:       1_000_000_000,
	}
	bb, err := c.UserStorer.ListByFilter(ctx, f)
	if err != nil {
		c.Logger.Error("database list error", slog.Any("error", err))
		return err
	}
	for _, b := range bb.Results {
		b.OrganizationName = org.Name
		err := c.UserStorer.UpdateByID(ctx, b)
		if err != nil {
			c.Logger.Error("database update error", slog.Any("error", err))
			return err
		}
	}
	return nil
}
