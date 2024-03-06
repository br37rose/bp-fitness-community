package controller

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"log/slog"

	u_s "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/user/datastore"
	domain "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/invoice/datastore"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/config/constants"
)

func (c *InvoiceControllerImpl) ListByFilter(ctx context.Context, f *domain.InvoiceListFilter) (*domain.InvoiceListResult, error) {
	// // Extract from our session the following data.
	// userID := ctx.Value(constants.SessionUserID).(primitive.ObjectID)
	userOID := ctx.Value(constants.SessionUserOrganizationID).(primitive.ObjectID)
	userRole := ctx.Value(constants.SessionUserRole).(int8)
	//
	// Apply filtering based on ownership and role.
	if userRole != u_s.UserRoleRoot {
		f.OrganizationID = userOID
	}

	c.Logger.Debug("listing using filter options:",
		slog.Any("OrganizationID", f.OrganizationID),
		slog.Any("BranchID", f.BranchID),
		slog.Any("Cursor", f.Cursor),
		slog.Int64("PageSize", f.PageSize),
		slog.String("SortField", f.SortField),
		slog.Int("SortOrder", int(f.SortOrder)),
		slog.String("SearchText", f.SearchText),
		slog.Bool("ExcludeArchived", f.ExcludeArchived))

	m, err := c.InvoiceStorer.ListByFilter(ctx, f)
	if err != nil {
		c.Logger.Error("database list by filter error", slog.Any("error", err))
		return nil, err
	}
	return m, err
}

func (c *InvoiceControllerImpl) ListAsSelectOptionByFilter(ctx context.Context, f *domain.InvoiceListFilter) ([]*domain.InvoiceAsSelectOption, error) {
	// Developers note: We want this unrestricted to account.

	c.Logger.Debug("listing using filter options:",
		slog.Any("OrganizationID", f.OrganizationID))

	// Filtering the database.
	m, err := c.InvoiceStorer.ListAsSelectOptionByFilter(ctx, f)
	if err != nil {
		c.Logger.Error("database list by filter error", slog.Any("error", err))
		return nil, err
	}
	return m, err
}
