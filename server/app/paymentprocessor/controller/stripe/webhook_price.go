package stripe

import (
	"context"
	"encoding/json"

	"log/slog"

	el_d "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/eventlog/datastore"
	"github.com/stripe/stripe-go/v72"
)

// webhookForPriceCreatedOrUpdated function will handle Stripe's `price.created` webhook event to create an offer in our system.
func (c *StripePaymentProcessorControllerImpl) webhookForPriceCreatedOrUpdated(ctx context.Context, event stripe.Event, el *el_d.EventLog) error {
	c.Logger.Debug("webhookForPriceCreatedOrUpdated: starting...")

	////
	//// Get the price and product.
	////

	var price stripe.Price

	// Successfully cast to []byte
	if err := json.Unmarshal(event.Data.Raw, &price); err != nil {
		c.Logger.Error("unmarshalling error", slog.Any("err", err))
		return err
	}
	el.Status = el_d.StatusOK

	off, err := c.OfferStorer.GetByStripeProductID(ctx, price.Product.ID)
	if err != nil {
		c.Logger.Error("get offer error", slog.Any("price.Product", price.Product))
		return err
	}

	////
	//// Update the offer.
	////

	off.StripePriceID = price.ID

	if err := c.OfferStorer.UpdateByID(ctx, off); err != nil {
		c.Logger.Error("create offer error", slog.Any("err", err))
		return err
	}
	c.Logger.Debug("webhookForPriceCreatedOrUpdated: successful processed")

	////
	//// Mark the logevent as processed.
	////

	if err := c.EventLogStorer.UpdateByID(ctx, el); err != nil {
		c.Logger.Error("create offer error", slog.Any("err", err))
		return err
	}

	c.Logger.Debug("webhookForPriceCreatedOrUpdated: finished")
	return nil
}
