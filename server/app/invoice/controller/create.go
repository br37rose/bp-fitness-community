package controller

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"log/slog"

	s_d "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/invoice/datastore"
	u_d "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/user/datastore"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/config/constants"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/utils/httperror"
)

func (c *InvoiceControllerImpl) Create(ctx context.Context, m *s_d.Invoice) (*s_d.Invoice, error) {
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

	// Add defaults.
	m.OrganizationID = oid
	m.OrganizationName = oname
	m.ID = primitive.NewObjectID()
	m.CreatedAt = time.Now()
	// m.CreatedByUserID = uid
	// m.CreatedByUserName = uname
	m.ModifiedAt = time.Now()
	// m.ModifiedByUserID = uid
	// m.ModifiedByUserName = uname
	m.Status = s_d.StatusActive

	// Save to our database.
	if err := c.InvoiceStorer.Create(ctx, m); err != nil {
		c.Logger.Error("database create error", slog.Any("error", err))
		return nil, err
	}

	return m, nil
}
