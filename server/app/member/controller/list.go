package controller

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"log/slog"

	user_s "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/user/datastore"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/config/constants"
)

func (c *MemberControllerImpl) ListByFilter(ctx context.Context, f *user_s.UserListFilter) (*user_s.UserListResult, error) {
	// // Extract from our session the following data.
	organizationID := ctx.Value(constants.SessionUserOrganizationID).(primitive.ObjectID)
	userRole := ctx.Value(constants.SessionUserRole).(int8)

	// Apply filtering based on ownership and role.
	if userRole == user_s.UserRoleAdmin {
		f.OrganizationID = organizationID
	}

	// We need to filter only by member users, not any other type of users.
	f.Role = user_s.UserRoleMember

	c.Logger.Debug("listing using filter options:",
		slog.Any("OrganizationID", f.OrganizationID),
		slog.Any("Cursor", f.Cursor),
		slog.Int64("PageSize", f.PageSize),
		slog.String("SortField", f.SortField),
		slog.Int("SortOrder", int(f.SortOrder)),
		slog.String("SearchText", f.SearchText),
		slog.Bool("ExcludeArchived", f.ExcludeArchived))

	m, err := c.UserStorer.ListByFilter(ctx, f)
	if err != nil {
		c.Logger.Error("database list by filter error", slog.Any("error", err))
		return nil, err
	}
	return m, err
}
func (c *MemberControllerImpl) ListAsSelectOptionByFilter(ctx context.Context, f *user_s.UserListFilter) ([]*user_s.UserAsSelectOption, error) {
	// Developers note: We want this unrestricted to account.

	c.Logger.Debug("listing using filter options:",
		slog.Any("OrganizationID", f.OrganizationID))

	// REQUIRED: Force filtering by members.
	f.Role = user_s.UserRoleMember

	// Filtering the database.
	m, err := c.UserStorer.ListAsSelectOptionByFilter(ctx, f)
	if err != nil {
		c.Logger.Error("database list by filter error", slog.Any("error", err))
		return nil, err
	}
	return m, err
}
