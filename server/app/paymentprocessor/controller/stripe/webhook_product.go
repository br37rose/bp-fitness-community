package stripe

import (
	"context"
	"encoding/json"
	"time"

	el_d "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/eventlog/datastore"
	off_d "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/offer/datastore"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/utils/httperror"
	"github.com/stripe/stripe-go/v72"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log/slog"
)

// webhookForProductCreated function will handle Stripe's `product.created` webhook event to create an offer in our system.
func (c *StripePaymentProcessorControllerImpl) webhookForProductCreated(ctx context.Context, event stripe.Event, el *el_d.EventLog) error {
	c.Logger.Debug("webhookForProductCreated: starting...")

	// The following fields must be filled out:
	// - name
	// - description
	// - images (only one)
	// - metadata: { "OrganizationID": "648763d3f6fbead15f5bd4d2" },
	// - Recurring
	// - Monthly billing period
	// - price

	////
	//// Get the product.
	////

	var product stripe.Product
	// Successfully cast to []byte
	if err := json.Unmarshal(event.Data.Raw, &product); err != nil {
		c.Logger.Error("unmarshalling error", slog.Any("err", err))
		return err
	}
	el.Status = el_d.StatusOK

	var organizationID primitive.ObjectID = primitive.NilObjectID
	var organizationName string
	if organizationIDStr, ok := product.Metadata["OrganizationID"]; ok {
		oid, err := primitive.ObjectIDFromHex(organizationIDStr)
		if err != nil {
			c.Logger.Error("reading object id from hex error", slog.Any("err", err))
			return err
		}

		o, err := c.OrganizationStorer.GetByID(ctx, oid)
		if err != nil {
			c.Logger.Error("reading object id from hex error", slog.Any("err", err))
			return err
		}
		if o == nil {
			c.Logger.Error("org does not exist error", slog.String("organizationIDStr", organizationIDStr))
			return httperror.NewForBadRequestWithSingleField("metadata_organization_id", "organization does not exist")
		}
		organizationID = oid
		organizationName = o.Name
	}

	imageURL := product.Images[0]

	////
	//// Create the offer.
	////

	off := &off_d.Offer{
		OrganizationID:       organizationID,
		OrganizationName:     organizationName,
		ID:                   primitive.NewObjectID(),
		Name:                 product.Name,
		Description:          product.Description,
		Status:               off_d.StatusPending,
		PaymentProcessorName: "Stripe, Inc.",
		StripeProductID:      product.ID,
		StripeImageURL:       imageURL,
		CreatedAt:            time.Now(),
		ModifiedAt:           time.Now(),
		BusinessFunction:     off_d.BusinessFunctionUnspecified,
	}
	if err := c.OfferStorer.Create(ctx, off); err != nil {
		c.Logger.Error("create offer error", slog.Any("err", err))
		return err
	}
	c.Logger.Debug("webhookForProductCreated: successful processed")

	////
	//// Mark the logevent as processed.
	////

	if err := c.EventLogStorer.UpdateByID(ctx, el); err != nil {
		c.Logger.Error("create offer error", slog.Any("err", err))
		return err
	}

	c.Logger.Debug("webhookForProductCreated: finished")
	return nil
}

// webhookForProductUpdated function will handle Stripe's `product.updated` webhook event to update an offer in our system.
func (c *StripePaymentProcessorControllerImpl) webhookForProductUpdated(ctx context.Context, event stripe.Event, el *el_d.EventLog) error {
	c.Logger.Debug("webhookForProductUpdated: starting...")

	////
	//// Get the product.
	////

	var product stripe.Product
	// Successfully cast to []byte
	if err := json.Unmarshal(event.Data.Raw, &product); err != nil {
		c.Logger.Error("unmarshalling error", slog.Any("err", err))
		return err
	}
	el.Status = el_d.StatusOK

	off, err := c.OfferStorer.GetByStripeProductID(ctx, product.ID)
	if err != nil {
		c.Logger.Error("get offer error", slog.Any("plan.Product", product.ID))
		return err
	}

	imageURL := product.Images[0]

	////
	//// Update the offer.
	////

	off.Name = product.Name
	off.Description = product.Description
	off.StripeImageURL = imageURL
	off.ModifiedAt = time.Now()
	if err := c.OfferStorer.UpdateByID(ctx, off); err != nil {
		c.Logger.Error("create offer error", slog.Any("err", err))
		return err
	}
	c.Logger.Debug("webhookForProductUpdated: successful processed")

	////
	//// Mark the logevent as processed.
	////

	if err := c.EventLogStorer.UpdateByID(ctx, el); err != nil {
		c.Logger.Error("create offer error", slog.Any("err", err))
		return err
	}

	c.Logger.Debug("webhookForProductUpdated: finished")
	return nil
}
