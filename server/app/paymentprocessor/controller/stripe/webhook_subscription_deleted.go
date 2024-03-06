package stripe

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/stripe/stripe-go/v72"

	el_d "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/eventlog/datastore"
	usr_d "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/user/datastore"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/utils/httperror"
)

// webhookForSubscriptionDeleted function will handle Stripe's `customer.subscription.deleted` webhook event to create an offer in our system.
func (impl *StripePaymentProcessorControllerImpl) webhookForSubscriptionDeleted(ctx context.Context, event stripe.Event, el *el_d.EventLog) error {
	// TECHDEBT: Find a way to replace parts of this code into adapter.
	impl.Logger.Debug("webhookForSubscriptionDeleted: starting...")

	// Unmarshal our subscription record from stripe.
	var subscription stripe.Subscription
	err := json.Unmarshal(event.Data.Raw, &subscription)
	if err != nil {
		impl.Logger.Error("unmarshalling subscription from stripe error",
			slog.Any("err", err),
			slog.String("webhook", string(event.Type)))
		return err
	}

	// Lookup the user in our database, else return a `400 Bad Request` error.
	u, err := impl.UserStorer.GetByPaymentProcessorCustomerID(ctx, subscription.Customer.ID)
	if err != nil {
		impl.Logger.Error("database error",
			slog.Any("err", err),
			slog.String("subscriptionID", subscription.ID),
			slog.String("webhook", string(event.Type)))
		return err
	}
	if u == nil {
		impl.Logger.Warn("user not found via customer id, attempting lookup by email...",
			slog.String("subscriptionID", subscription.ID),
			slog.Any("subscription.Customer.ID", subscription.Customer.ID),
			slog.String("webhook", string(event.Type)))

		// Alternative case of looking up the user's email in case we cannot
		// find the use based on their `Customer ID`.
		u, err = impl.UserStorer.GetByEmail(ctx, subscription.Customer.Email)
		if err != nil {
			impl.Logger.Error("database error", slog.Any("err", err),
				slog.String("subscriptionID", subscription.ID),
				slog.String("webhook", string(event.Type)))
			return err
		}

		// Finally if the user is not found then return error.
		if u == nil {
			impl.Logger.Error("user does not exist validation error",
				slog.String("subscriptionID", subscription.ID),
				slog.String("customerID", subscription.Customer.ID),
				slog.String("customerEmail", subscription.Customer.Email),
				slog.String("webhook", string(event.Type)))
			return httperror.NewForBadRequestWithSingleField("user", fmt.Sprintf("does not exist for email of %s nor customer id %s", subscription.Customer.Email, subscription.Customer.ID))
		}
	}

	impl.Logger.Debug("found customer in our system", slog.Any("user_id", u.ID),
		slog.String("webhook", string(event.Type)))

	impl.Logger.Debug("processing delete stripe subscription...",
		slog.String("Subscription ID", subscription.ID),
		slog.Any("Subscription Plan", subscription.Plan),
		slog.Any("Subscription Items", subscription.Items),
		slog.String("webhook", string(event.Type)),
	)

	// Create our new user subscription.
	u.StripeSubscription = &usr_d.StripeSubscription{
		PriceID:        subscription.Plan.ID, // aka `price_id`
		SubscriptionID: subscription.ID,
		Interval:       string(subscription.Plan.Interval),
		Status:         string(subscription.Status),
	}
	u.IsSubscriber = subscription.Status == usr_d.SubscriptionStatusActive
	u.SubscriptionStatus = "canceled"         // Value set by Stripe.
	u.SubscriptionStartedAt = *new(time.Time) // Set to zero value of time.

	// Update the user record.
	if err := impl.UserStorer.UpdateByID(ctx, u); err != nil {
		impl.Logger.Error("update user error",
			slog.Any("err", err),
			slog.String("subscriptionID", subscription.ID),
			slog.String("webhook", string(event.Type)))
		return err
	}
	impl.Logger.Debug("customer subscription deleted",
		slog.String("webhook", string(event.Type)),
		slog.Any("stripe_subscription_id", subscription.ID),
		slog.Any("customer_id", subscription.Customer.ID))

	impl.Logger.Debug("webhookForSubscriptionDeleted: finished")
	return nil
}
