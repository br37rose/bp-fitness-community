package controller

import (
	"context"
	"errors"
	"log/slog"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	offer_s "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/offer/datastore"
	usr_d "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/user/datastore"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/config/constants"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/utils/httperror"
)

func (impl *PaymentProcessorControllerImpl) CreateStripeCheckoutSessionURL(ctx context.Context, offerIDString string) (string, error) {
	// Extract from our session the following data.
	userID := ctx.Value(constants.SessionUserID).(primitive.ObjectID)

	////
	//// Start the transaction.
	////

	session, err := impl.DbClient.StartSession()
	if err != nil {
		impl.Logger.Error("start session error",
			slog.Any("error", err))
		return "", err
	}
	defer session.EndSession(ctx)

	// Define a transaction function with a series of operations
	transactionFunc := func(sessCtx mongo.SessionContext) (interface{}, error) {

		// Lookup the user in our database, else return a `400 Bad Request` error.
		u, err := impl.UserStorer.GetByID(sessCtx, userID)
		if err != nil {
			impl.Logger.Error("database error", slog.Any("err", err))
			return "", err
		}
		if u == nil {
			impl.Logger.Warn("user does not exist validation error")
			return "", errors.New("user id does not exist")
		}

		offerID, err := primitive.ObjectIDFromHex(offerIDString)
		if err != nil {
			return "", err
		}

		// Validate the offer id
		offer, err := impl.OfferStorer.GetByID(sessCtx, offerID)
		if err != nil {
			impl.Logger.Warn("offer d.n.e. error")
			return "", errors.New("offer does not exist")
		}

		// Defensive code: Prevent executing this function if different processor.
		if u.PaymentProcessorName != "Stripe, Inc." {
			impl.Logger.Warn("not stripe payment processor assigned to user.")
			return "", errors.New("user is using payment processor which is not supported")
		}

		// Defensive code: Prevent executing this function if different processor.
		if offer.StripePriceID == "" {
			impl.Logger.Warn("this product is not ready")
			return "", errors.New("this product is not ready")
		}

		impl.Logger.Debug("creating stripe checkout session", slog.String("priceID", offer.StripePriceID), slog.Any("userId", u.ID))

		var redirectURL string
		switch offer.PayFrequency {
		case offer_s.PayFrequencyOneTime:
			redirectURL, err = impl.PaymentProcessor.CreateOneTimeCheckoutSessionURL(impl.Emailer.GetFrontendDomainName(), "/purchase/success", "/purchase/canceled", u.PaymentProcessorCustomerID, offer.StripePriceID, false, "")
			if err != nil {
				return "", err
			}
		case offer_s.PayFrequencyMonthly, offer_s.PayFrequencyAnnual:

			// Defensive code: Prevent creating a new subscription if a previous
			// subscription already exists!
			if u.StripeSubscription != nil {
				if u.StripeSubscription.Status == usr_d.SubscriptionStatusActive {
					impl.Logger.Warn("subscription already exists error")
					return "", httperror.NewForBadRequestWithSingleField("subscription", "you already have a subscription, please cancel existing subscription before enrolling again")
				}
			}
			redirectURL, err = impl.PaymentProcessor.CreateSubscriptionCheckoutSessionURL(impl.Emailer.GetFrontendDomainName(), "/purchase/success", "/purchase/canceled", u.PaymentProcessorCustomerID, offer.StripePriceID, false, "")
			if err != nil {
				return "", err
			}
		}

		////
		//// Exit our transaction successfully.
		////

		impl.Logger.Debug("stripe checkout session ready", slog.String("redirectURL", redirectURL))
		return redirectURL, nil
	}

	// Start a transaction
	redirectURL, err := session.WithTransaction(ctx, transactionFunc)
	if err != nil {
		impl.Logger.Error("session failed error",
			slog.Any("error", err))
		return "", err
	}

	return redirectURL.(string), nil
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
func (impl *PaymentProcessorControllerImpl) CompleteStripeCheckoutSession(ctx context.Context, sessionID string) (*CompleteStripeCheckoutSessionResponse, error) {

	// Extract from our session the following data.
	userID := ctx.Value(constants.SessionUserID).(primitive.ObjectID)

	////
	//// Start the transaction.
	////

	session, err := impl.DbClient.StartSession()
	if err != nil {
		impl.Logger.Error("start session error",
			slog.Any("error", err))
		return nil, err
	}
	defer session.EndSession(ctx)

	// Define a transaction function with a series of operations
	transactionFunc := func(sessCtx mongo.SessionContext) (interface{}, error) {

		// Lookup the user in our database, else return a `400 Bad Request` error.
		u, err := impl.UserStorer.GetByID(sessCtx, userID)
		if err != nil {
			impl.Logger.Error("database error", slog.Any("err", err))
			return nil, err
		}
		if u == nil {
			impl.Logger.Error("user does not exist validation error")
			return nil, httperror.NewForBadRequestWithSingleField("user", "does not exist")
		}

		// Defensive code: Prevent executing this function if different processor.
		if u.PaymentProcessorName != "Stripe, Inc." {
			impl.Logger.Error("not stripe payment processor assigned to user")
			return nil, errors.New("user is using payment processor which is not stripe")
		}

		// Lookup the `Checkout Session` on `Stripe` for the particular unique id.
		checkoutSession, err := impl.PaymentProcessor.GetCheckoutSession(sessionID)
		if err != nil {
			impl.Logger.Error("unable to get checkout session from stripe error")
			return nil, err
		}
		if checkoutSession == nil {
			impl.Logger.Error("checkout session does not exist error")
			return nil, httperror.NewForBadRequestWithSingleField("session_id", "checkout session does not exist")
		}
		if checkoutSession.Customer == nil {
			impl.Logger.Error("customer does not exiserror")
			return nil, httperror.NewForBadRequestWithSingleField("customer", "does not exist")
		}

		// For debugging purposes only.
		impl.Logger.Debug("fetched checkout session from stripe",
			slog.Any("BP8 UserID", userID),
			slog.String("Stripe SessionID", sessionID),
			slog.String("Stripe CustomerID", checkoutSession.Customer.ID))

		// Defensive code: Confirm that the `customer id` in the looked up session
		// matches the authenticated users payment processor customer id to
		// confirm they are the same users!
		if u.PaymentProcessorCustomerID != checkoutSession.Customer.ID {
			impl.Logger.Error("unauthorized user fetching session error")
			return nil, httperror.NewForForbiddenWithSingleField("session_id", "you do not belong to this session")
		}

		switch checkoutSession.PaymentStatus {
		case "paid":
			impl.Logger.Debug("fetched subscription is paid")
		case "unpaid":
			impl.Logger.Error("unauthorized user fetching session error")
			return nil, httperror.NewForForbiddenWithSingleField("message", "you do not pay yet for the subscription")
		case "no_payment_required":
			impl.Logger.Debug("fetched subscription payment is not required")
		}

		switch checkoutSession.Status {
		case "open":
			impl.Logger.Warn("unauthorized user fetching session error")
			return nil, httperror.NewForForbiddenWithSingleField("message", "your checkout session is still open")
		case "expired":
			impl.Logger.Error("unauthorized user fetching session error")
			return nil, httperror.NewForForbiddenWithSingleField("message", "your checkout session expired")
		case "complete":
			impl.Logger.Debug("fetched subscription is completed")
		}

		// According to Stripe, Inc. documentation, we need to look at the line items
		// in the response, so we will make API call to them to get the details.

		lineItems, err := impl.PaymentProcessor.GetCheckoutSessionLineItems(sessionID)
		if err != nil {
			impl.Logger.Error("payment processor error", slog.Any("err", err), slog.String("webhook", "checkout.session.completed"))
			return nil, err
		}
		for _, lineItem := range lineItems {
			// Lookup the offer that our web-app fullfilled.
			offer, err := impl.OfferStorer.GetByStripePriceID(sessCtx, lineItem.Price.ID)
			if err != nil {
				impl.Logger.Error("database error", slog.Any("err", err), slog.String("webhook", "checkout.session.completed"))
				return nil, err
			}
			if offer == nil {
				impl.Logger.Error("offer does not exist error", slog.String("webhook", "checkout.session.completed"))
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

		////
		//// Exit our transaction successfully.
		////

		return nil, nil
	}

	// Start a transaction
	res, err := session.WithTransaction(ctx, transactionFunc)
	if err != nil {
		impl.Logger.Error("session failed error",
			slog.Any("error", err))
		return nil, err
	}

	return res.(*CompleteStripeCheckoutSessionResponse), nil
}
