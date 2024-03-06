package stripe

import (
	"context"
	"errors"
	"time"

	"log/slog"

	"go.mongodb.org/mongo-driver/bson/primitive"

	u_d "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/user/datastore"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/config/constants"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/utils/httperror"
)

func (c *StripePaymentProcessorControllerImpl) CancelStripeSubscription(ctx context.Context, userID primitive.ObjectID) error {
	// Extract from our session the following data.
	sessionUserID, _ := ctx.Value(constants.SessionUserID).(primitive.ObjectID)

	if userID.IsZero() {
		userID = sessionUserID
	}

	// Lookup the user in our database, else return a `400 Bad Request` error.
	u, err := c.UserStorer.GetByID(ctx, userID)
	if err != nil {
		c.Logger.Error("database error", slog.Any("err", err))
		return err
	}
	if u == nil {
		c.Logger.Warn("user does not exist validation error")
		return errors.New("user id does not exist")
	}

	// Defensive code: Prevent cancelling a subscription for various issues.
	if u.StripeSubscription == nil {
		c.Logger.Warn("subscription does not exists error")
		return httperror.NewForBadRequestWithSingleField("subscription", "you do not have a subscription, please create a subscription before cancelling again")
	}

	c.Logger.Debug("beginning the stripe subscription cancelling process",
		slog.String("SubscriptionID", u.StripeSubscription.SubscriptionID),
		slog.String("PriceID", u.StripeSubscription.PriceID),
		slog.Any("UserID", u.ID))

	subscription, err := c.PaymentProcessor.CancelSubscription(u.StripeSubscription.SubscriptionID)
	if err != nil {
		c.Logger.Error("stripe cancel subscription error", slog.Any("err", err))
		return err
	}
	if subscription == nil {
		c.Logger.Warn("subscription does not exists error")
		return httperror.NewForBadRequestWithSingleField("subscription", "was not returned when cancelled")
	}

	// Update our subscription to be canceled.
	u.StripeSubscription = &u_d.StripeSubscription{
		Status: u_d.SubscriptionStatusCanceled,
	}
	u.IsSubscriber = false

	// Update the user record.
	if err := c.UserStorer.UpdateByID(ctx, u); err != nil {
		c.Logger.Error("update user error", slog.Any("err", err))
		return err
	}

	//
	// Send notification email.
	//

	// Lookup the offer...
	offer, err := c.OfferStorer.GetByStripePriceID(ctx, subscription.Plan.ID) // aka `price_id`
	if err != nil {
		c.Logger.Error("database error", slog.Any("err", err), slog.String("webhook", "customer.subscription.updated"))
		return err
	}
	if offer == nil {
		c.Logger.Error("offer does not exist error", slog.String("webhook", "customer.subscription.updated"))
		return errors.New("offer does not exist in our system")
	}
	c.Logger.Debug("fetched offer", slog.Any("name", offer.Name), slog.Any("id", offer.ID), slog.String("webhook", "customer.subscription.updated"))

	// Send cancelation email.
	if err := c.TemplatedEmailer.SendMemberCancelledSubscriptionEmailToMember(u.Email, u.FirstName, offer.Name, time.Now()); err != nil {
		c.Logger.Error("email error", slog.Any("err", err))
		return err
	}

	return nil
}
