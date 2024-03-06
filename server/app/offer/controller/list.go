package controller

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"log/slog"

	domain "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/offer/datastore"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/config/constants"
)

func (c *OfferControllerImpl) ListByFilter(ctx context.Context, f *domain.OfferListFilter) (*domain.OfferListResult, error) {
	// // Extract from our session the following data.
	// userID := ctx.Value(constants.SessionUserID).(primitive.ObjectID)
	userOID := ctx.Value(constants.SessionUserOrganizationID).(primitive.ObjectID)
	userID := ctx.Value(constants.SessionUserID).(primitive.ObjectID)

	// Apply filtering based on tenancy.
	f.OrganizationID = userOID

	c.Logger.Debug("listing using filter options:",
		slog.Any("OrganizationID", f.OrganizationID),
		slog.Any("BranchID", f.BranchID),
		slog.Any("Cursor", f.Cursor),
		slog.Int64("PageSize", f.PageSize),
		slog.String("SortField", f.SortField),
		slog.Int("SortOrder", int(f.SortOrder)),
		slog.String("SearchText", f.SearchText),
		slog.Any("Status", f.Status))

	m, err := c.OfferStorer.ListByFilter(ctx, f)
	if err != nil {
		c.Logger.Error("database list by filter error", slog.Any("error", err))
		return nil, err
	}

	// DEVELOPERS NOTE:
	// We need to restrict offerings if the user previously purchased it.
	u, err := c.UserStorer.GetByID(ctx, userID)
	if err != nil {
		c.Logger.Error("database get error", slog.Any("error", err))
		return nil, err
	}
	if u == nil {
		return nil, errors.New("user does not exist")
	}
	// for _, offer := range m.Results {
	// 	// Automatically assume the user has access.
	// 	offer.CurrentUserHasAccessGranted = true

	// 	// Set the variables we will use.
	// 	currentPurchaseCount := 0
	// 	purchaseLimit := offer.PurchaseLimit

	// 	// DEVELOPERS NOTE: If the purchase limit is greater then zero then there
	// 	// is a limit, if the purchase limit is zero then there is no limit.
	// 	if purchaseLimit > 0 {
	// 		for _, purchase := range u.Purchases {
	// 			if purchase.OfferID == offer.ID {
	// 				currentPurchaseCount += 1
	// 			}
	// 		}
	// 		offer.CurrentUserHasAccessGranted = currentPurchaseCount <= purchaseLimit
	// 	}
	// }

	return m, err
}

func (c *OfferControllerImpl) ListAsSelectOptionByFilter(ctx context.Context, f *domain.OfferListFilter) ([]*domain.OfferAsSelectOption, error) {
	// // Extract from our session the following data.
	// userID := ctx.Value(constants.SessionUserID).(primitive.ObjectID)
	userOID := ctx.Value(constants.SessionUserOrganizationID).(primitive.ObjectID)

	// Apply filtering based on tenancy.
	f.OrganizationID = userOID

	c.Logger.Debug("listing using filter options:",
		slog.Any("OrganizationID", f.OrganizationID),
		slog.Any("BranchID", f.BranchID),
		slog.Any("Cursor", f.Cursor),
		slog.Int64("PageSize", f.PageSize),
		slog.String("SortField", f.SortField),
		slog.Int("SortOrder", int(f.SortOrder)),
		slog.String("SearchText", f.SearchText),
		slog.Any("Status", f.Status))

	c.Logger.Debug("listing using filter options:",
		slog.Any("OrganizationID", f.OrganizationID))

	// Filtering the database.
	m, err := c.OfferStorer.ListAsSelectOptionByFilter(ctx, f)
	if err != nil {
		c.Logger.Error("database list by filter error", slog.Any("error", err))
		return nil, err
	}
	return m, err
}
