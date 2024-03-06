package stripe

import (
	"context"
	"encoding/json"

	el_d "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/eventlog/datastore"
	off_d "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/offer/datastore"
	"github.com/stripe/stripe-go/v72"
	"log/slog"
)

// webhookForPlanCreatedOrUpdated function will handle Stripe's `plan.created` webhook event to create an offer in our system.
func (c *StripePaymentProcessorControllerImpl) webhookForPlanCreatedOrUpdated(ctx context.Context, event stripe.Event, el *el_d.EventLog) error {
	c.Logger.Debug("webhookForPlanCreatedOrUpdated: starting...")

	////
	//// Get the plan and product.
	////

	var plan stripe.Plan

	// Successfully cast to []byte
	if err := json.Unmarshal(event.Data.Raw, &plan); err != nil {
		c.Logger.Error("unmarshalling error", slog.Any("err", err))
		return err
	}
	el.Status = el_d.StatusOK

	off, err := c.OfferStorer.GetByStripeProductID(ctx, plan.Product.ID)
	if err != nil {
		c.Logger.Error("get offer error", slog.Any("plan.Product", plan.Product))
		return err
	}

	////
	//// Update the offer.
	////

	off.Price = fromStripeFormat(plan.Amount)
	off.PriceCurrency = string(plan.Currency)
	switch plan.Interval {
	case stripe.PlanIntervalDay:
		off.PayFrequency = off_d.PayFrequencyDay
		break
	case stripe.PlanIntervalMonth:
		off.PayFrequency = off_d.PayFrequencyMonthly
		break
	case stripe.PlanIntervalWeek:
		off.PayFrequency = off_d.PayFrequencyWeek
		break
	case stripe.PlanIntervalYear:
		off.PayFrequency = off_d.PayFrequencyAnnual
		break
	default:
		off.PayFrequency = off_d.PayFrequencyOneTime
		break
	}

	if plan.Interval != "" {
		off.IsSubscription = true
	}

	if err := c.OfferStorer.UpdateByID(ctx, off); err != nil {
		c.Logger.Error("create offer error", slog.Any("err", err))
		return err
	}
	c.Logger.Debug("webhookForPlanCreatedOrUpdated: successful processed")

	////
	//// Mark the logevent as processed.
	////

	if err := c.EventLogStorer.UpdateByID(ctx, el); err != nil {
		c.Logger.Error("create offer error", slog.Any("err", err))
		return err
	}

	c.Logger.Debug("webhookForPlanCreatedOrUpdated: finished")
	return nil
}
