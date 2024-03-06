package stripe

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/stripe/stripe-go/v72"

	el_d "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/eventlog/datastore"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/utils/httperror"
)

func (impl *StripePaymentProcessorControllerImpl) webhookForCustomerUpdated(ctx context.Context, event stripe.Event, el *el_d.EventLog) error {
	// TECHDEBT: Find a way to replace parts of this code into adapter.
	impl.Logger.Debug("processing stripe invoice create webhook call", slog.String("webhook", string(event.Type)))

	// Unmarshal our subscription record from stripe.

	var customer stripe.Customer
	err := json.Unmarshal(event.Data.Raw, &customer)
	if err != nil {
		impl.Logger.Error("unmarshalling customer from stripe error", slog.Any("err", err))
		return err
	}

	// Lookup the user in our database, else return a `400 Bad Request` error.
	u, err := impl.UserStorer.GetByPaymentProcessorCustomerID(ctx, customer.ID)
	if err != nil {
		impl.Logger.Error("database error", slog.Any("err", err))
		return err
	}
	if u == nil {
		impl.Logger.Error("user does not exist validation error", slog.Any("customerID", customer.ID))
		return httperror.NewForBadRequestWithSingleField("user", "does not exist")
	}

	impl.Logger.Debug("found customer in our system", slog.Any("user_id", u.ID))

	u.Email = customer.Email
	if customer.Shipping != nil {
		//TODO: Impl.
	}
	// if customer.Address != nil {
	// 	//TODO: Impl.
	// }

	return nil
}
