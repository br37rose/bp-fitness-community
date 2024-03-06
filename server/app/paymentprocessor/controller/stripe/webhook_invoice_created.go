package stripe

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/stripe/stripe-go/v72"

	el_d "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/eventlog/datastore"
	usr_d "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/user/datastore"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/utils/httperror"
)

func (impl *StripePaymentProcessorControllerImpl) webhookForInvoiceCreated(ctx context.Context, event stripe.Event, el *el_d.EventLog) error {
	// TECHDEBT: Find a way to replace parts of this code into adapter.
	impl.Logger.Debug("processing stripe invoice create webhook call", slog.String("webhook", string(event.Type)))

	// DEVELOPERS NOTE:
	// See: https://stripe.com/docs/api/events/types#event_types-invoice.created
	// Occurs whenever a new invoice is created.

	var invoice stripe.Invoice
	err := json.Unmarshal(event.Data.Raw, &invoice)
	if err != nil {
		impl.Logger.Error("unmarshalling invoice from stripe error",
			slog.Any("err", err),
			slog.String("webhook", string(event.Type)))
		return err
	}

	// Lookup the user in our database, else return a `400 Bad Request` error.
	u, err := impl.UserStorer.GetByPaymentProcessorCustomerID(ctx, invoice.Customer.ID)
	if err != nil {
		impl.Logger.Error("database error",
			slog.Any("err", err),
			slog.String("invoiceID", invoice.ID),
			slog.String("webhook", string(event.Type)))
		return err
	}
	if u == nil {
		impl.Logger.Warn("user not found via customer id, attempting lookup by email...",
			slog.String("invoiceID", invoice.ID),
			slog.Any("invoice.Customer.ID", invoice.Customer.ID),
			slog.String("webhook", string(event.Type)))

		// Alternative case of looking up the user's email in case we cannot
		// find the use based on their `Customer ID`.
		u, err = impl.UserStorer.GetByEmail(ctx, invoice.Customer.Email)
		if err != nil {
			impl.Logger.Error("database error",
				slog.String("invoiceID", invoice.ID),
				slog.Any("err", err),
				slog.String("webhook", string(event.Type)))
			return err
		}

		// Finally if the user is not found then return error.
		if u == nil {
			impl.Logger.Error("user does not exist validation error",
				slog.String("invoiceID", invoice.ID),
				slog.String("customerID", invoice.Customer.ID),
				slog.String("customerEmail", invoice.Customer.Email),
				slog.String("webhook", string(event.Type)))
			return httperror.NewForBadRequestWithSingleField("user", fmt.Sprintf("does not exist for email of %s nor customer id %s", invoice.Customer.Email, invoice.Customer.ID))
		}

		impl.Logger.Warn("user found via customer id",
			slog.String("invoiceID", invoice.ID),
			slog.String("invoice.Customer.Email", invoice.Customer.Email),
			slog.String("webhook", string(event.Type)))
	}

	impl.Logger.Debug("found customer in our system",
		slog.String("invoiceID", invoice.ID),
		slog.Any("user_id", u.ID),
		slog.String("webhook", string(event.Type)))

	// Lookup the Stripe invoice that we have previously created.
	found, err := impl.UserStorer.GetStripeInvoiceByPaymentProcessorInvoiceID(ctx, invoice.ID)
	if err != nil {
		impl.Logger.Error("database error",
			slog.String("invoiceID", invoice.ID),
			slog.Any("err", err),
			slog.String("webhook", string(event.Type)))
		return err
	}
	if found != nil {
		impl.Logger.Warn("Invoice already exists, skipping this webhook call...",
			slog.String("invoice.ID", invoice.ID),
			slog.String("webhook", string(event.Type)))
		return nil
	}

	// Create our own interval stripe representation of the invoice and save
	// this new invoices it in the database.

	si := &usr_d.StripeInvoice{
		InvoiceID:            invoice.ID,
		Created:              invoice.Created,
		Paid:                 invoice.Paid,
		HostedInvoiceURL:     invoice.HostedInvoiceURL,
		InvoicePDF:           invoice.InvoicePDF,
		SubtotalExcludingTax: 0,
		Tax:                  0,
		Total:                0,
		Number:               invoice.Number,
		Currency:             string(invoice.Currency),
	}
	if invoice.SubtotalExcludingTax > 0 {
		si.SubtotalExcludingTax = fromStripeFormat(invoice.SubtotalExcludingTax)
	}
	if invoice.Tax > 0 {
		si.Tax = fromStripeFormat(invoice.Tax)
	}
	if invoice.Total > 0 {
		si.Total = fromStripeFormat(invoice.Total)
	}

	// The following code will "prepend into a slice". Special thanks to:
	// https://codingair.wordpress.com/2014/07/18/go-appendprepend-item-into-slice/
	u.StripeInvoices = append([]*usr_d.StripeInvoice{si}, u.StripeInvoices...)

	// Update the user record.
	if err := impl.UserStorer.UpdateByID(ctx, u); err != nil {
		impl.Logger.Error("update user error",
			slog.String("invoiceID", invoice.ID),
			slog.Any("err", err),
			slog.String("webhook", string(event.Type)))
		return err
	}
	impl.Logger.Debug("customer invoice created",
		slog.String("webhook", string(event.Type)),
		slog.Any("invoice_id", invoice.ID),
		slog.Any("customer_id", invoice.Customer.ID),
		slog.Bool("paid", invoice.Paid),
		slog.Int64("created", invoice.Created),
		slog.String("invoice_pdf", invoice.InvoicePDF),
		slog.String("webhook", "invoice.created"),
		slog.String("hosted_invoice_url", invoice.HostedInvoiceURL))

	return nil
}
