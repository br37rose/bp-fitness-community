package controller

import (
	"context"

	"log/slog"

	"go.mongodb.org/mongo-driver/bson/primitive"

	t_s "github.com/bci-innovation-labs/bp8fitnesscommunity-cli/app/tag/datastore"
	tag_s "github.com/bci-innovation-labs/bp8fitnesscommunity-cli/app/tag/datastore"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-cli/config/constants"
)

func (c *TagControllerImpl) ListByFilter(ctx context.Context, f *t_s.TagPaginationListFilter) (*t_s.TagPaginationListResult, error) {
	// // Extract from our session the following data.
	organizationID := ctx.Value(constants.SessionUserOrganizationID).(primitive.ObjectID)

	// Apply filtering based on ownership and role.
	f.OrganizationID = organizationID // Manditory

	c.Logger.Debug("listing using filter options:",
		slog.Any("Cursor", f.Cursor),
		slog.Int64("PageSize", f.PageSize),
		slog.String("SortField", f.SortField),
		// slog.Int("SortOrder", int(f.SortOrder)),
		// slog.Any("OrganizationID", f.OrganizationID),
		// slog.Any("Type", f.Type),
		// slog.Any("Status", f.Status),
		// slog.Bool("ExcludeArchived", f.ExcludeArchived),
		// slog.String("SearchText", f.SearchText),
		// slog.Any("FirstName", f.FirstName),
		// slog.Any("LastName", f.LastName),
		// slog.Any("Email", f.Email),
		// slog.Any("Phone", f.Phone),
		// slog.Time("CreatedAtGTE", f.CreatedAtGTE)
	)

	m, err := c.TagStorer.ListByFilter(ctx, f)
	if err != nil {
		c.Logger.Error("database list by filter error", slog.Any("error", err))
		return nil, err
	}
	return m, err
}

// func (c *TagControllerImpl) LiteListByFilter(ctx context.Context, f *t_s.TagListFilter) (*t_s.TagLiteListResult, error) {
// 	// // Extract from our session the following data.
// 	organizationID := ctx.Value(constants.SessionUserOrganizationID).(primitive.ObjectID)
//
// 	// Apply filtering based on ownership and role.
// 	f.OrganizationID = organizationID // Manditory
//
// 	c.Logger.Debug("listing using filter options:",
// 		slog.Any("Cursor", f.Cursor),
// 		slog.Int64("PageSize", f.PageSize),
// 		slog.String("SortField", f.SortField),
// 		slog.Int("SortOrder", int(f.SortOrder)),
// 		slog.Any("OrganizationID", f.OrganizationID),
// 	)
//
// 	m, err := c.TagStorer.LiteListByFilter(ctx, f)
// 	if err != nil {
// 		c.Logger.Error("database list by filter error", slog.Any("error", err))
// 		return nil, err
// 	}
// 	return m, err
// }

func (c *TagControllerImpl) ListAsSelectOptionByFilter(ctx context.Context, f *tag_s.TagListFilter) ([]*tag_s.TagAsSelectOption, error) {
	// // Extract from our session the following data.
	organizationID := ctx.Value(constants.SessionUserOrganizationID).(primitive.ObjectID)

	// Apply filtering based on ownership and role.
	f.OrganizationID = organizationID // Manditory

	c.Logger.Debug("listing using filter options:",
		slog.Any("Cursor", f.Cursor),
		slog.Int64("PageSize", f.PageSize),
		slog.String("SortField", f.SortField),
		slog.Int("SortOrder", int(f.SortOrder)),
		slog.Any("OrganizationID", f.OrganizationID),
	)

	// Filtering the database.
	m, err := c.TagStorer.ListAsSelectOptionByFilter(ctx, f)
	if err != nil {
		c.Logger.Error("database list by filter error", slog.Any("error", err))
		return nil, err
	}
	return m, err
}
