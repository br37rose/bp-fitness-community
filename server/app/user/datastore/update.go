package datastore

import (
	"context"

	"log/slog"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (impl UserStorerImpl) UpdateByID(ctx context.Context, m *User) error {
	filter := bson.D{{"_id", m.ID}}

	update := bson.M{ // DEVELOPERS NOTE: https://stackoverflow.com/a/60946010
		"$set": m,
	}

	// execute the UpdateOne() function to update the first matching document
	_, err := impl.Collection.UpdateOne(ctx, filter, update)
	if err != nil {
		impl.Logger.Error("database update by user id error", slog.Any("error", err))
	}

	return nil
}

func (impl UserStorerImpl) UpdateLiteByID(ctx context.Context, m *UserLite) error {
	filter := bson.D{{"_id", m.ID}}

	update := bson.M{ // DEVELOPERS NOTE: https://stackoverflow.com/a/60946010
		"$set": m,
	}

	// execute the UpdateOne() function to update the first matching document
	_, err := impl.Collection.UpdateOne(ctx, filter, update)
	if err != nil {
		impl.Logger.Error("database update lite by user id error", slog.Any("error", err))
	}

	return nil
}

func (impl UserStorerImpl) UpdateStripeInvoiceByPaymentProcessorInvoiceID(ctx context.Context, newInvoice *StripeInvoice) error {
	// DEVELOPERS NOTE:
	// To learn more about querying inside nested fields, then please see:
	// https://www.mongodb.com/docs/manual/tutorial/query-embedded-documents/

	filter := bson.M{"stripe_invoices.invoice_id": newInvoice.InvoiceID}

	var result User
	err := impl.Collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// This error means your query did not match any documents.
			return nil
		}
		impl.Logger.Error("database get by verification code error", slog.Any("error", err))
		return err
	}

	// Now we need to find the specific StripeInvoice within the User's StripeInvoices slice.
	for _, invoice := range result.StripeInvoices {
		if invoice.InvoiceID == newInvoice.InvoiceID {
			invoice.Created = newInvoice.Created
			invoice.Paid = newInvoice.Paid
			invoice.HostedInvoiceURL = newInvoice.HostedInvoiceURL
			invoice.InvoicePDF = newInvoice.InvoicePDF

			return impl.UpdateByID(ctx, &result)
		}
	}

	// If the paymentProcessorInvoiceID is not found in the StripeInvoices, return nil.
	return nil
}
