package stripe

import (
	"context"
	"log/slog"

	"github.com/stripe/stripe-go/v72"

	el_d "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/eventlog/datastore"
)

func (impl *StripePaymentProcessorControllerImpl) webhookForInvoicePaid(ctx context.Context, event stripe.Event, el *el_d.EventLog) error {
	// TECHDEBT: Find a way to replace parts of this code into adapter.
	impl.Logger.Debug("processing stripe invoice create webhook call", slog.String("webhook", string(event.Type)))

	// TECHDEBT: Find a way to replace parts of this code into adapter.
	impl.Logger.Debug("processing stripe invoice paid webhook call",
		slog.String("webhook", string(event.Type)))

	// DEVELOPERS NOTE:
	// See: https://stripe.com/docs/api/events/types#event_types-invoice.paid

	if err := impl.webhookForInvoiceUpdated(ctx, event, el); err != nil {
		return err
	}

	impl.Logger.Debug("processing stripe invoice paid webhook call",
		slog.String("webhook", string(event.Type)))

	impl.Logger.Debug("finished stripe invoice paid webhook call",
		slog.String("webhook", string(event.Type)))

	return nil
}
