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

func (impl *StripePaymentProcessorControllerImpl) webhookForInvoiceUpdated(ctx context.Context, event stripe.Event, el *el_d.EventLog) error {
	// TECHDEBT: Find a way to replace parts of this code into adapter.
	impl.Logger.Debug("processing stripe invoice updated webhook call", slog.String("webhook", string(event.Type)))

	// DEVELOPERS NOTE:
	// See: https://stripe.com/docs/api/events/types#event_types-invoice.updated

	var i stripe.Invoice
	err := json.Unmarshal(event.Data.Raw, &i)
	if err != nil {
		impl.Logger.Error("unmarshalling invoice from stripe error",
			slog.Any("err", err),
			slog.String("webhook", string(event.Type)))
		return err
	}

	// Lookup the user in our database, else return a `400 Bad Request` error.
	u, err := impl.UserStorer.GetByPaymentProcessorCustomerID(ctx, i.Customer.ID)
	if err != nil {
		impl.Logger.Error("database error",
			slog.Any("err", err),
			slog.String("invoiceID", i.ID),
			slog.String("webhook", string(event.Type)))
		return err
	}
	if u == nil {
		impl.Logger.Warn("user not found via customer id, attempting lookup by email...",
			slog.Any("invoice.Customer.ID", i.Customer.ID),
			slog.String("invoiceID", i.ID),
			slog.String("webhook", string(event.Type)))

		// Alternative case of looking up the user's email in case we cannot
		// find the use based on their `Customer ID`.
		u, err = impl.UserStorer.GetByEmail(ctx, i.Customer.Email)
		if err != nil {
			impl.Logger.Error("database error",
				slog.Any("err", err),
				slog.String("invoiceID", i.ID),
				slog.String("webhook", string(event.Type)))
			return err
		}

		// Finally if the user is not found then return error.
		if u == nil {
			impl.Logger.Error("user does not exist validation error",
				slog.String("invoiceID", i.ID),
				slog.String("customerID", i.Customer.ID),
				slog.String("customerEmail", i.Customer.Email),
				slog.String("webhook", string(event.Type)))
			return httperror.NewForBadRequestWithSingleField("user", fmt.Sprintf("does not exist for email of %s nor customer id %s", i.Customer.Email, i.Customer.ID))
		}

		impl.Logger.Warn("user found via customer id",
			slog.String("invoiceID", i.ID),
			slog.String("invoice.Customer.Email", i.Customer.Email),
			slog.String("webhook", string(event.Type)))
	}

	impl.Logger.Debug("found customer in our system",
		slog.Any("user_id", u.ID),
		slog.String("invoiceID", i.ID),
		slog.String("webhook", string(event.Type)))

	// Lookup the Stripe invoice that we have previously created.
	invoice, err := impl.UserStorer.GetStripeInvoiceByPaymentProcessorInvoiceID(ctx, i.ID)
	if err != nil {
		impl.Logger.Error("database error",
			slog.Any("err", err),
			slog.String("invoiceID", i.ID),
			slog.String("webhook", string(event.Type)))
		return err
	}
	if invoice == nil {
		impl.Logger.Warn("No invoice found, creating one now...",
			slog.String("webhook", "invoice.updated"))

		// Create our own interval stripe representation of the invoice and save
		// this new invoices it in the database.

		si := &usr_d.StripeInvoice{
			InvoiceID:            i.ID,
			Created:              i.Created,
			Paid:                 i.Paid,
			HostedInvoiceURL:     i.HostedInvoiceURL,
			InvoicePDF:           i.InvoicePDF,
			SubtotalExcludingTax: 0,
			Tax:                  0,
			Total:                0,
			Number:               i.Number,
			Currency:             string(i.Currency),
		}
		if i.SubtotalExcludingTax > 0 {
			si.SubtotalExcludingTax = fromStripeFormat(i.SubtotalExcludingTax)
		}
		if i.Tax > 0 {
			si.Tax = fromStripeFormat(i.Tax)
		}
		if i.Total > 0 {
			si.Total = fromStripeFormat(i.Total)
		}

		// The following code will "prepend into a slice". Special thanks to:
		// https://codingair.wordpress.com/2014/07/18/go-appendprepend-item-into-slice/
		u.StripeInvoices = append([]*usr_d.StripeInvoice{si}, u.StripeInvoices...)

		// Update the user record.
		if err := impl.UserStorer.UpdateByID(ctx, u); err != nil {
			impl.Logger.Error("update user error",
				slog.Any("err", err),
				slog.String("webhook", string(event.Type)))
			return err
		}
		impl.Logger.Debug("customer invoice created",
			slog.String("webhook", string(event.Type)),
			slog.Any("invoice_id", i.ID),
			slog.Any("customer_id", i.Customer.ID),
			slog.Bool("paid", i.Paid),
			slog.Int64("created", i.Created),
			slog.String("invoice_pdf", i.InvoicePDF),
			slog.String("hosted_invoice_url", i.HostedInvoiceURL))

		invoice = si
		impl.Logger.Warn("Created invoice",
			slog.String("webhook", "invoice.updated"))
	}

	invoice.Created = i.Created
	invoice.Paid = i.Paid
	invoice.HostedInvoiceURL = i.HostedInvoiceURL
	invoice.InvoicePDF = i.InvoicePDF
	invoice.Created = i.Created
	invoice.Paid = i.Paid
	invoice.HostedInvoiceURL = i.HostedInvoiceURL
	invoice.InvoicePDF = i.InvoicePDF
	if invoice.SubtotalExcludingTax > 0 {
		invoice.SubtotalExcludingTax = fromStripeFormat(i.SubtotalExcludingTax)
	} else {
		invoice.SubtotalExcludingTax = 0
	}
	if i.Tax > 0 {
		invoice.Tax = fromStripeFormat(i.Tax)
	} else {
		invoice.Tax = 0
	}
	if i.Total > 0 {
		invoice.Total = fromStripeFormat(i.Total)
	} else {
		invoice.Total = 0
	}
	invoice.Number = i.Number
	invoice.Currency = string(i.Currency)

	// Update the user record.
	if err := impl.UserStorer.UpdateStripeInvoiceByPaymentProcessorInvoiceID(ctx, invoice); err != nil {
		impl.Logger.Error("user update invoice error",
			slog.Any("err", err),
			slog.String("invoiceID", i.ID),
			slog.String("webhook", string(event.Type)))
		return err
	}
	impl.Logger.Debug("customer invoice updated",
		slog.String("webhook", string(event.Type)),
		slog.Any("invoice_id", i.ID),
		slog.Any("customer_id", i.Customer.ID),
		slog.Bool("paid", i.Paid),
		slog.Int64("created", i.Created),
		slog.String("invoice_pdf", i.InvoicePDF),
		slog.String("hosted_invoice_url", i.HostedInvoiceURL))

	// Send notification
	if invoice.Paid && i.InvoicePDF != "" && i.HostedInvoiceURL != "" {
		// // Use the user's provided time zone or default to UTC.
		// localLoc, err := time.LoadLocation("UTC")
		// if err != nil {
		// 	impl.Logger.Error("load location error",
		// 		slog.Any("error", err),
		// 		slog.Any("timezone", "UTC"),
		// 		slog.String("invoiceID", i.ID))
		// 	return nil
		// }
		//
		// if err := impl.TemplatedEmailer.SendMemberInvoicePaidEmailToMember(u.Email, u.FirstName, invoice.Number, invoice.Total, invoice.Currency, time.Now().In(localLoc), i.HostedInvoiceURL); err != nil { //TODO: FIX.
		// 	impl.Logger.Error("failed sending email error",
		// 		slog.Any("err", err),
		// 		slog.String("invoiceID", i.ID),
		// 		slog.String("webhook", string(event.Type)))
		// }

		// DEVELOPERS NOTE:
		// THE FOLLOWING CODE IS ADDED TO HANDLE A UNIQUE CASE IN WHICH USER
		// SUBSCRIPTION NEEDS TO BE ACTIVATED ONCE WE VERIFY AN INVOICE HAS
		// BEEN PAID. THIS IS DONE IF THERE WAS A COMPLICATION DURING CHECKOUT
		// AND WE HANDLE TURNING ON SUBSCRIPTIONS NOW THAT INVOICE IS PAID.
		if i.Subscription != nil {
			if i.Subscription.ID != "" {
				// Fetch the subscription record.
				subscription, err := impl.PaymentProcessor.GetSubscription(i.Subscription.ID)
				if err != nil {
					impl.Logger.Warn("get subscription from stripe error",
						slog.Any("err", err),
						slog.String("invoiceID", i.ID),
						slog.String("webhook", string(event.Type)))
					return nil
				}
				if subscription == nil {
					impl.Logger.Warn("subscription does not exist in stripe",
						slog.String("SubscriptionID", i.Subscription.ID),
						slog.String("invoiceID", i.ID),
						slog.String("webhook", string(event.Type)))
					return nil
				}
				if subscription.Status == usr_d.SubscriptionStatusActive {
					u.IsSubscriber = true
					u.SubscriptionStatus = string(subscription.Status)
					if u.StripeSubscription != nil {
						u.StripeSubscription.Status = string(subscription.Status)
					}
					impl.Logger.Debug("user is verified as a subscriber b/c invoice paid",
						slog.String("invoiceID", i.ID),
						slog.String("webhook", string(event.Type)))

					if err := impl.UserStorer.UpdateByID(ctx, u); err != nil {
						impl.Logger.Warn("database update by id error",
							slog.String("invoiceID", i.ID),
							slog.Any("err", err),
							slog.String("webhook", string(event.Type)))
						return nil
					}
					impl.Logger.Debug("user updated because invoice was verified as paid.",
						slog.String("invoiceID", i.ID),
						slog.String("webhook", string(event.Type)))
				}
			}
		}

		// DEVELOPERS NOTE:
		// HERE IS ANOTHER ATTEMPT AT HAVING A SUBSCRIPTION BECOME ENABLED
		// IF AN INVOICE IS DETECTED TO BE PAID AND THE REASON THE USER WAS
		// BILLED THIS INVOICE WAS BECAUSE A SUBSCRIPTION WAS CREATED.
		if i.BillingReason == stripe.InvoiceBillingReasonSubscriptionCreate {
			u.IsSubscriber = true
			u.SubscriptionStatus = usr_d.SubscriptionStatusActive
			impl.Logger.Debug("user is verified as a subscriber b/c invoice paid and billing reason is subscription was created",
				slog.String("invoiceID", i.ID),
				slog.String("webhook", string(event.Type)))

			if err := impl.UserStorer.UpdateByID(ctx, u); err != nil {
				impl.Logger.Warn("database update by id error",
					slog.String("invoiceID", i.ID),
					slog.Any("err", err),
					slog.String("webhook", string(event.Type)))
				return nil
			}
			impl.Logger.Debug("user updated because invoice was verified as paid and billing reason is subscription was created.",
				slog.String("invoiceID", i.ID),
				slog.String("webhook", string(event.Type)))
		}
	}

	return nil
}
