package controller

import (
	"context"
	"time"

	domain "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/organization/datastore"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/utils/httperror"
	"log/slog"
)

func (c *OrganizationControllerImpl) UpdateByID(ctx context.Context, ns *domain.Organization) (*domain.Organization, error) {
	// Fetch the original organization.
	os, err := c.OrganizationStorer.GetByID(ctx, ns.ID)
	if err != nil {
		c.Logger.Error("database get by id error", slog.Any("error", err))
		return nil, err
	}
	if os == nil {
		return nil, httperror.NewForBadRequestWithSingleField("id", "organization does not exist")
	}

	// Modify our original organization.
	os.ModifiedAt = time.Now()
	os.Type = ns.Type
	os.Status = ns.Status
	os.Name = ns.Name

	// Save to the database the modified organization.
	if err := c.OrganizationStorer.UpdateByID(ctx, os); err != nil {
		c.Logger.Error("database update by id error", slog.Any("error", err))
		return nil, err
	}

	// Process all the records that related to this organization with the new
	// changes made by the update.
	if err := c.updateRelated(ctx, os); err != nil {
		return nil, err
	}

	return os, nil
}
