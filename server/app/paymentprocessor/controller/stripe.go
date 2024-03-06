package controller

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"log/slog"

	u_d "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/user/datastore"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/config/constants"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/utils/httperror"
)

func (c *PaymentProcessorControllerImpl) CancelStripeSubscription(ctx context.Context) error {
	// Extract from our session the following data.
	userID := ctx.Value(constants.SessionUserID).(primitive.ObjectID)

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

	// Defensive code: Prevent executing this function if different processor.
	if u.PaymentProcessorName != "Stripe, Inc." {
		c.Logger.Warn("not stripe payment processor assigned to user.")
		return errors.New("user is using payment processor which is not supported")
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

func (impl *PaymentProcessorControllerImpl) ListLatestStripeInvoices(ctx context.Context, userID primitive.ObjectID, cursor int64, limit int64) (*u_d.StripeInvoiceListResult, error) {
	// Extract from our session the following data.
	urole := ctx.Value(constants.SessionUserRole).(int8)
	uid := ctx.Value(constants.SessionUserID).(primitive.ObjectID)

	////
	//// Permission handling.
	////

	var pickedUserID primitive.ObjectID
	switch urole { // Security.
	case u_d.UserRoleRoot:
		return nil, httperror.NewForForbiddenWithSingleField("message", "you did not saasify product")
	case u_d.UserRoleTrainer:
		return nil, httperror.NewForForbiddenWithSingleField("message", "you do not have permission as a trainer") //TODO: Implement.
	case u_d.UserRoleMember:
		// Override whatever user ID was set to the logged in user because the
		// user is a member and thus a member only has access to his/her data.
		pickedUserID = uid
	case u_d.UserRoleAdmin:
		if userID.IsZero() {
			return nil, httperror.NewForBadRequestWithSingleField("user_id", "missing url parameter")
		}
		pickedUserID = userID
	}

	////
	//// Database lookups interactions.
	////

	// Lookup the user in our database, else return a `400 Bad Request` error.
	u, err := impl.UserStorer.GetByID(ctx, pickedUserID)
	if err != nil {
		impl.Logger.Error("database error", slog.Any("err", err))
		return nil, err
	}
	if u == nil {
		impl.Logger.Warn("user does not exist validation error")
		return nil, httperror.NewForBadRequestWithSingleField("id", "does not exist")
	}

	return impl.UserStorer.ListLatestStripeInvoices(ctx, userID, cursor, limit)
}
