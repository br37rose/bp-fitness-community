package stripe

import (
	"context"

	"github.com/stripe/stripe-go/v72"

	el_d "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/eventlog/datastore"
)

func (impl *StripePaymentProcessorControllerImpl) webhookForForCheckoutSessionCompleted(ctx context.Context, event stripe.Event, el *el_d.EventLog) error {

	return nil
}
