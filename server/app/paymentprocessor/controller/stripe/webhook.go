package stripe

import (
	"context"
	"log/slog"
	"time"

	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/webhook"
	"go.mongodb.org/mongo-driver/bson/primitive"

	el_d "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/eventlog/datastore"
)

func (c *StripePaymentProcessorControllerImpl) Webhook(ctx context.Context, header string, b []byte) error {
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

	// Log in our system the webhook event and process in our system this event.
	eventlog, err := c.logWebhookEvent(context.Background(), event)
	if err != nil {
		c.Logger.Error("failed logging stripe webhook event", slog.String("event_type", event.Type))
	}

	c.Logger.Debug("stripe webhook executing...", slog.String("event_type", event.Type))

	// DEVELOPERS NOTE: The full list can be found via https://stripe.com/docs/api/events/types
	switch eventlog.SecondaryType {
	case "customer.subscription.created":
		return c.webhookForSubscriptionCreated(ctx, event, eventlog) //TODO: QA
	case "customer.subscription.updated":
		return c.webhookForSubscriptionUpdated(ctx, event, eventlog) //TODO: QA
	case "customer.subscription.deleted":
		return c.webhookForSubscriptionDeleted(ctx, event, eventlog) //TODO: QA
	case "customer.updated":
		return c.webhookForCustomerUpdated(ctx, event, eventlog) //TODO: QA
	case "invoice.created":
		return c.webhookForInvoiceCreated(ctx, event, eventlog) //TODO: QA
	case "invoice.updated":
		return c.webhookForInvoiceUpdated(ctx, event, eventlog) //TODO: QA
	case "invoice.paid":
		return c.webhookForInvoicePaid(ctx, event, eventlog) //TODO: QA
	case "checkout.session.completed":
		return c.webhookForForCheckoutSessionCompleted(ctx, event, eventlog) //TODO: QA
	case "product.created":
		return c.webhookForProductCreated(ctx, event, eventlog)
	case "product.updated":
		return c.webhookForProductUpdated(ctx, event, eventlog)
	case "plan.created", "plan.updated":
		return c.webhookForPlanCreatedOrUpdated(ctx, event, eventlog)
	case "price.created", "price.updated":
		return c.webhookForPriceCreatedOrUpdated(ctx, event, eventlog)
	default:
		c.Logger.Warn("skip processing stripe event", slog.Any("eventType", event.Type))
		return nil
	}
}

func (c *StripePaymentProcessorControllerImpl) logWebhookEvent(ctx context.Context, event stripe.Event) (*el_d.EventLog, error) {
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
		Status:        el_d.StatusPending,
		ExternalID:    event.ID,
		ID:            primitive.NewObjectID(),
	}
	if err := c.EventLogStorer.Create(ctx, eventlog); err != nil {
		c.Logger.Error("marshalling create error", slog.Any("err", err))
		return nil, err
	}
	c.Logger.Debug("logged stripe webhook event")
	return eventlog, nil
}
