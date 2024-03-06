package controller

import (
	"context"
	"log/slog"
	"time"

	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/webhook"
	"go.mongodb.org/mongo-driver/bson/primitive"

	el_d "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/eventlog/datastore"
)

func (c *PaymentProcessorControllerImpl) StripeWebhook(ctx context.Context, header string, b []byte) error {
	// TECHDEBT: Find a way to replace parts of this code into adapter.

	// DEVELOPERS NOTE:
	// HOW DO WE MANUALLY RESEND AN EVENT THAT WAS PREVIOUSLY CALLED?
	// STEP 1:
	// Perform the action you want in the app: Ex: purchase.
	//
	// STEP 2:
	// Go to: https://dashboard.stripe.com/test/events
	//
	// STEP 3:
	// Lookup the event ID, for example: `evt_1NfoHAC1dNpgYbqFJfSeeCNK`.
	//
	// STEP 4:
	// In a new terminal window, run the code to resend (note: https://stripe.com/docs/cli/events/resend):
	// $ stripe events resend evt_1NfoHAC1dNpgYbqFJfSeeCNK

	event, err := webhook.ConstructEvent(b, header, c.PaymentProcessor.GetWebhookSecretKey())
	if err != nil {
		c.Logger.Error("update user error", slog.Any("err", err))
		return err
	}

	go func() {
		if err := c.logStripeWebhookEvent(context.Background(), event); err != nil {
			c.Logger.Error("failed logging stripe webhook event", slog.String("event_type", event.Type))
		}
	}()

	c.Logger.Debug("stripe webhook executing...", slog.String("event_type", event.Type))

	// DEVELOPERS NOTE: The full list can be found via https://stripe.com/docs/api/events/types
	switch event.Type {
	case "customer.subscription.created":
		// return c.stripeWebhookForSubscriptionCreated(ctx, event)
		return nil
	case "customer.subscription.updated":
		// return c.stripeWebhookForSubscriptionUpdated(ctx, event)
		return nil
	case "customer.subscription.deleted":
		// 	Sent when a customer’s subscription ends.
		// return c.stripeWebhookForSubscriptionDeleted(ctx, event)
		return nil
	case "invoice.created":
		// return c.stripeWebhookForInvoiceCreated(ctx, event)
		return nil
	case "invoice.updated":
		// Sent each billing interval when a payment succeeds.
		// return c.stripeWebhookForInvoiceUpdated(ctx, event)
		return nil
	case "invoice.paid":
		// Continue to provision the subscription as payments continue to be made.
		// Store the status in your database and check when a user accesses your service.
		// This approach helps you avoid hitting rate limits.
		// return c.stripeWebhookForInvoicePaid(ctx, event)
		return nil
	case "invoice.payment_failed":
		// The payment failed or the customer does not have a valid payment method.
		// The subscription becomes past_due. Notify your customer and send them to the
		// customer portal to update their payment information.
		c.Logger.Error("please implement - invoice.payment_failed")
		return nil //TODO: IMPLEMENT.
	case "checkout.session.completed":
		// Payment is successful and the subscription is created or a product baught.
		// You should provision the subscription and save the customer ID to your database.
		// return c.stripeWebhookForCheckoutSessionCompleted(ctx, event)
		return nil
	default:
		c.Logger.Warn("skip processing stripe event", slog.Any("eventType", event.Type))
		return nil
	}
}

func (c *PaymentProcessorControllerImpl) logStripeWebhookEvent(ctx context.Context, event stripe.Event) error {
	c.Logger.Debug("logging stripe webhook event...")

	// DEVELOPERS NOTE:
	// We will take advantage of Golang's new "generics" functionality and store
	// the event data into our database. The mongodb library will understand
	// how to handle the generic.

	eventlog := &el_d.EventLog{
		PrimaryType:   el_d.PrimaryTypeStripeWebhookEvent,
		SecondaryType: event.Type,
		CreatedAt:     time.Now(),
		Content:       event.Data.Object, // Store the event payload, not metadata.
		Status:        el_d.StatusOK,
		ExternalID:    event.ID,
		ID:            primitive.NewObjectID(),
	}
	if err := c.EventLogStorer.Create(ctx, eventlog); err != nil {
		c.Logger.Error("marshalling create error", slog.Any("err", err))
		return err
	}
	c.Logger.Debug("logged stripe webhook event")
	return nil
}
