package stripe

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/stripe/stripe-go/v72"

	el_d "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/eventlog/datastore"
	usr_d "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/user/datastore"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/utils/httperror"
)

// webhookForSubscriptionCreated function will handle Stripe's `customer.subscription.updated` webhook event to create an offer in our system.
func (impl *StripePaymentProcessorControllerImpl) webhookForSubscriptionUpdated(ctx context.Context, event stripe.Event, el *el_d.EventLog) error {
	// TECHDEBT: Find a way to replace parts of this code into adapter.
	impl.Logger.Debug("webhookForSubscriptionUpdated: starting...")

	// Unmarshal our subscription record from stripe.
	var subscription stripe.Subscription
	err := json.Unmarshal(event.Data.Raw, &subscription)
	if err != nil {
		impl.Logger.Error("unmarshalling subscription from stripe error",
			slog.Any("err", err), slog.String("webhook", "customer.subscription.updated"),
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
			slog.Any("subscription.Customer.ID", subscription.Customer.ID),
			slog.String("subscriptionID", subscription.ID),
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

	impl.Logger.Debug("found customer in our system", slog.Any("user_id", u.ID))

	impl.Logger.Debug("processing update stripe subscription...",
		slog.String("webhook", string(event.Type)),
		slog.String("Subscription ID", subscription.ID),
		slog.Any("Subscription Plan", subscription.Plan),
		slog.Any("Subscription Items", subscription.Items),
	)

	// Lookup the offer...
	offer, err := impl.OfferStorer.GetByStripePriceID(ctx, subscription.Plan.ID) // aka `price_id`
	if err != nil {
		impl.Logger.Error("database error", slog.Any("err", err),
			slog.String("subscriptionID", subscription.ID),
			slog.String("webhook", string(event.Type)))
		return err
	}
	if offer == nil {
		impl.Logger.Error("offer does not exist error",
			slog.String("subscriptionID", subscription.ID),
			slog.String("webhook", string(event.Type)))
		return errors.New("offer does not exist in our system")
	}
	impl.Logger.Debug("fetched offer",
		slog.String("subscriptionID", subscription.ID),
		slog.Any("name", offer.Name),
		slog.Any("id", offer.ID),
		slog.String("webhook", string(event.Type)))

	// Use the user's provided time zone or default to UTC.
	location, _ := time.LoadLocation("UTC")

	// // Create the user purchase record.
	// purchase := &usr_d.UserPurchase{
	// 	ID:                    primitive.NewObjectID(),
	// 	OrganizationID:        u.OrganizationID,
	// 	CreatedAt:             time.Now().In(location),
	// 	ModifiedAt:            time.Now().In(location),
	// 	OfferID:               offer.ID,
	// 	OfferName:             offer.Name,
	// 	OfferDescription:      offer.Description,
	// 	OfferPrice:            offer.Price,
	// 	OfferPriceCurrency:    offer.PriceCurrency,
	// 	OfferPayFrequency:     offer.PayFrequency,
	// 	OfferBusinessFunction: offer.BusinessFunction,
	// 	OfferType:             offer.Type,
	// }

	// Create our new user subscription.
	// u.StripeSubscription = &usr_d.StripeSubscription{
	// 	PriceID:        subscription.Plan.ID, // aka `price_id`
	// 	SubscriptionID: subscription.ID,
	// 	Interval:       string(subscription.Plan.Interval),
	// 	Status:         string(subscription.Status),
	// 	// OfferPurchase:  purchase,
	// }
	u.IsSubscriber = subscription.Status == usr_d.SubscriptionStatusActive
	u.SubscriptionOfferID = offer.ID
	u.SubscriptionOfferName = offer.Name
	u.SubscriptionStatus = string(subscription.Status)
	u.OfferID = offer.ID
	u.OfferMembershipRank = offer.MembershipRank
	u.OfferName = offer.Name
	u.ModifiedAt.In(location)

	// Update the user record.
	if err := impl.UserStorer.UpdateByID(ctx, u); err != nil {
		impl.Logger.Error("update user error",
			slog.Any("err", err),
			slog.String("subscriptionID", subscription.ID),
			slog.String("webhook", string(event.Type)))
		return err
	}
	impl.Logger.Debug("customer subscription updated",
		slog.String("webhook", string(event.Type)),
		slog.Any("stripe_subscription_id", subscription.ID),
		slog.Any("customer_id", subscription.Customer.ID))

	impl.Logger.Debug("webhookForSubscriptionUpdated: finished")
	return nil
}
