package controller

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"log/slog"

	domain "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/invoice/datastore"
	u_d "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/user/datastore"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/config/constants"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/utils/httperror"
)

func (c *InvoiceControllerImpl) UpdateByID(ctx context.Context, ns *domain.Invoice) (*domain.Invoice, error) {
	// Extract from our session the following data.
	urole := ctx.Value(constants.SessionUserRole).(int8)
	// uid := ctx.Value(constants.SessionUserID).(primitive.ObjectID)
	// uname := ctx.Value(constants.SessionUserName).(string)
	oid := ctx.Value(constants.SessionUserOrganizationID).(primitive.ObjectID)
	oname := ctx.Value(constants.SessionUserOrganizationName).(string)

	switch urole { // Security.
	case u_d.UserRoleRoot:
		return nil, httperror.NewForForbiddenWithSingleField("message", "you did not saasify invoice")
	case u_d.UserRoleTrainer:
		return nil, httperror.NewForForbiddenWithSingleField("message", "you do not have permission")
	case u_d.UserRoleMember:
		return nil, httperror.NewForForbiddenWithSingleField("message", "you do not have permission")
	}

	// Fetch the original organization.
	os, err := c.InvoiceStorer.GetByID(ctx, ns.ID)
	if err != nil {
		c.Logger.Error("database get by id error", slog.Any("error", err))
		return nil, err
	}
	if os == nil {
		return nil, httperror.NewForBadRequestWithSingleField("id", "workout program type does not exist")
	}

	// Modify our original organization.
	os.OrganizationID = oid
	os.OrganizationName = oname
	os.ModifiedAt = time.Now()
	os.Status = ns.Status
	// os.Name = ns.Name

	// Save to the database the modified organization.
	if err := c.InvoiceStorer.UpdateByID(ctx, os); err != nil {
		c.Logger.Error("database update by id error", slog.Any("error", err))
		return nil, err
	}

	return os, nil
}
