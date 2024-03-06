package stripe

import (
	"context"
	"errors"
	"fmt"

	"log/slog"

	"go.mongodb.org/mongo-driver/bson/primitive"

	offer_s "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/offer/datastore"
	usr_d "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/user/datastore"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/config/constants"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/utils/httperror"
)

func (c *StripePaymentProcessorControllerImpl) CreateStripeCheckoutSessionURL(ctx context.Context, offerIDString string) (string, error) {
	// Extract from our session the following data.
	userID := ctx.Value(constants.SessionUserID).(primitive.ObjectID)

	// Lookup the user in our database, else return a `400 Bad Request` error.
	u, err := c.UserStorer.GetByID(ctx, userID)
	if err != nil {
		c.Logger.Error("database error", slog.Any("err", err))
		return "", err
	}
	if u == nil {
		c.Logger.Warn("user does not exist validation error")
		return "", errors.New("user id does not exist")
	}

	offerID, err := primitive.ObjectIDFromHex(offerIDString)
	if err != nil {
		return "", err
	}

	// Validate the offer id
	offer, err := c.OfferStorer.GetByID(ctx, offerID)
	if err != nil {
		c.Logger.Warn("offer d.n.e. error")
		return "", errors.New("offer does not exist")
	}

	// Defensive code: Prevent executing this function if different processor.
	if offer.StripePriceID == "" {
		c.Logger.Warn("this product is not ready")
		return "", errors.New("this product is not ready")
	}

	// DEVELOPERS NOTE: If for some reason the user has not created a
	// `CustomerID` with the payment merchant then we'll create it now.
	if u.PaymentProcessorCustomerID == "" {
		c.Logger.Debug("detected user missing stripe customer id", slog.Any("userId", u.ID))

		paymentProcessorCustomerID, err := c.PaymentProcessor.CreateCustomer(
			fmt.Sprintf("%s %s", u.FirstName, u.LastName),
			u.Email,
			"", // description...
			fmt.Sprintf("%s %s Shipping Address", u.FirstName, u.LastName),
			u.Phone,
			u.ShippingCity, u.ShippingCountry, u.ShippingAddressLine1, u.ShippingAddressLine2, u.ShippingPostalCode, u.ShippingRegion, // Shipping
			u.City, u.Country, u.AddressLine1, u.AddressLine2, u.PostalCode, u.Region, // Billing
		)
		if err != nil {
			c.Logger.Error("creating customer from payment processor error", slog.Any("error", err))
			return "", err
		}
		c.Logger.Debug("created stripe customer id",
			slog.Any("userId", u.ID),
			slog.Any("paymentProcessorCustomerID", paymentProcessorCustomerID))

		u.PaymentProcessorCustomerID = *paymentProcessorCustomerID
		if err := c.UserStorer.UpdateByID(ctx, u); err != nil {
			c.Logger.Error("unable to get checkout session from stripe error")
			return "", err
		}

		c.Logger.Debug("attached stripe customer id to user account",
			slog.Any("userId", u.ID),
			slog.Any("paymentProcessorCustomerID", paymentProcessorCustomerID))
	}

	hasShippingAddress := u.ShippingCity != "" || u.ShippingCountry != "" || u.ShippingAddressLine1 != ""

	c.Logger.Debug("creating stripe checkout session", slog.String("priceID", offer.StripePriceID), slog.Any("userId", u.ID))

	var redirectURL string
	switch offer.PayFrequency {
	case offer_s.PayFrequencyOneTime:
		redirectURL, err = c.PaymentProcessor.CreateOneTimeCheckoutSessionURL(
			c.Emailer.GetFrontendDomainName(),
			"/purchase/success",
			"/purchase/canceled",
			u.PaymentProcessorCustomerID,
			offer.StripePriceID,
			hasShippingAddress,
			"", // u.PaymentProcessorCouponID, // StripePaymentProcessorControllerImpl
		)
		if err != nil {
			return "", err
		}
	case offer_s.PayFrequencyMonthly, offer_s.PayFrequencyAnnual:

		// Defensive code: Prevent creating a new subscription if a previous
		// subscription already exists!
		if u.StripeSubscription != nil {
			if u.StripeSubscription.Status == usr_d.SubscriptionStatusActive {
				c.Logger.Warn("subscription already exists error")
				return "", httperror.NewForBadRequestWithSingleField("subscription", "you already have a subscription, please cancel existing subscription before enrolling again")
			}
		}
		redirectURL, err = c.PaymentProcessor.CreateSubscriptionCheckoutSessionURL(
			c.Emailer.GetFrontendDomainName(),
			"/purchase/success",
			"/purchase/canceled",
			u.PaymentProcessorCustomerID,
			offer.StripePriceID,
			hasShippingAddress,
			"", // u.PaymentProcessorCouponID, // u.PaymentProcessorCouponID
		)
		if err != nil {
			return "", err
		}
	}

	c.Logger.Debug("stripe checkout session ready",
		slog.String("redirectURL", redirectURL),
		slog.Bool("hasShippingAddress", hasShippingAddress))
	return redirectURL, nil
}

type CompleteStripeCheckoutSessionResponse struct {
	Name          string  `json:"name"`
	Description   string  `bson:"description" json:"description"`
	Price         float64 `bson:"price" json:"price"`
	PriceCurrency string  `bson:"price_currency" json:"price_currency"`
	PayFrequency  int8    `bson:"pay_frequency" json:"pay_frequency"`
	SessionID     string  `json:"session_id"`
	PaymentStatus string  `json:"payment_status"`
	Status        string  `json:"status"`
}

// CompleteStripeCheckoutSession API endpoint handles return the checkout success details.
func (c *StripePaymentProcessorControllerImpl) CompleteStripeCheckoutSession(ctx context.Context, sessionID string) (*CompleteStripeCheckoutSessionResponse, error) {
	// Extract from our session the following data.
	userID := ctx.Value(constants.SessionUserID).(primitive.ObjectID)

	// Lookup the user in our database, else return a `400 Bad Request` error.
	u, err := c.UserStorer.GetByID(ctx, userID)
	if err != nil {
		c.Logger.Error("database error", slog.Any("err", err))
		return nil, err
	}
	if u == nil {
		c.Logger.Error("user does not exist validation error")
		return nil, httperror.NewForBadRequestWithSingleField("user", "does not exist")
	}

	// Lookup the `Checkout Session` on `Stripe` for the particular unique id.
	checkoutSession, err := c.PaymentProcessor.GetCheckoutSession(sessionID)
	if err != nil {
		c.Logger.Error("unable to get checkout session from stripe error")
		return nil, err
	}
	if checkoutSession == nil {
		c.Logger.Error("checkout session does not exist error")
		return nil, httperror.NewForBadRequestWithSingleField("session_id", "checkout session does not exist")
	}
	if checkoutSession.Customer == nil {
		c.Logger.Error("customer does not exiserror")
		return nil, httperror.NewForBadRequestWithSingleField("customer", "does not exist")
	}

	// For debugging purposes only.
	c.Logger.Debug("fetched checkout session from stripe",
		slog.Any("BP8 UserID", userID),
		slog.String("Stripe SessionID", sessionID),
		slog.String("Stripe CustomerID", checkoutSession.Customer.ID))

	// Defensive code: Confirm that the `customer id` in the looked up session
	// matches the authenticated users payment processor customer id to
	// confirm they are the same users!
	if u.PaymentProcessorCustomerID != checkoutSession.Customer.ID {
		c.Logger.Error("unauthorized user fetching session error")
		return nil, httperror.NewForForbiddenWithSingleField("session_id", "you do not belong to this session")
	}

	switch checkoutSession.PaymentStatus {
	case "paid":
		c.Logger.Debug("fetched subscription is paid")
	case "unpaid":
		c.Logger.Error("unauthorized user fetching session error")
		return nil, httperror.NewForForbiddenWithSingleField("message", "you do not pay yet for the subscription")
	case "no_payment_required":
		c.Logger.Debug("fetched subscription payment is not required")
	}

	switch checkoutSession.Status {
	case "open":
		c.Logger.Warn("unauthorized user fetching session error")
		return nil, httperror.NewForForbiddenWithSingleField("message", "your checkout session is still open")
	case "expired":
		c.Logger.Error("unauthorized user fetching session error")
		return nil, httperror.NewForForbiddenWithSingleField("message", "your checkout session expired")
	case "complete":
		c.Logger.Debug("fetched subscription is completed")
	}

	// According to Stripe, Inc. documentation, we need to look at the line items
	// in the response, so we will make API call to them to get the details.

	lineItems, err := c.PaymentProcessor.GetCheckoutSessionLineItems(sessionID)
	if err != nil {
		c.Logger.Error("payment processor error", slog.Any("err", err), slog.String("webhook", "checkout.session.completed"))
		return nil, err
	}
	for _, lineItem := range lineItems {
		// Lookup the offer that our web-app fullfilled.
		offer, err := c.OfferStorer.GetByStripePriceID(ctx, lineItem.Price.ID)
		if err != nil {
			c.Logger.Error("database error", slog.Any("err", err), slog.String("webhook", "checkout.session.completed"))
			return nil, err
		}
		if offer == nil {
			c.Logger.Error("offer does not exist error", slog.String("webhook", "checkout.session.completed"))
			return nil, errors.New("offer does not exist in our system")
		}

		res := &CompleteStripeCheckoutSessionResponse{
			Name:          offer.Name,
			Description:   offer.Description,
			Price:         offer.Price,
			PriceCurrency: offer.PriceCurrency,
			PayFrequency:  offer.PayFrequency,
			SessionID:     sessionID,
			PaymentStatus: string(checkoutSession.PaymentStatus),
			Status:        string(checkoutSession.Status),
		}

		return res, nil
	}
	return nil, nil
}
