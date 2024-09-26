package datastore

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"log/slog"
)

func (impl UserStorerImpl) GetByID(ctx context.Context, id primitive.ObjectID) (*User, error) {
	filter := bson.D{{"_id", id}}

	var result User
	err := impl.Collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// This error means your query did not match any documents.
			return nil, nil
		}
		impl.Logger.Error("database get by user id error", slog.Any("error", err))
		return nil, err
	}
	return &result, nil
}

func (impl UserStorerImpl) GetByEmail(ctx context.Context, email string) (*User, error) {
	filter := bson.D{{"email", email}}

	var result User
	err := impl.Collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// This error means your query did not match any documents.
			return nil, nil
		}
		impl.Logger.Error("database get by email error", slog.Any("error", err))
		return nil, err
	}
	return &result, nil
}

func (impl UserStorerImpl) GetByVerificationCode(ctx context.Context, verificationCode string) (*User, error) {
	filter := bson.D{{"email_verification_code", verificationCode}}

	var result User
	err := impl.Collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// This error means your query did not match any documents.
			return nil, nil
		}
		impl.Logger.Error("database get by verification code error", slog.Any("error", err))
		return nil, err
	}
	return &result, nil
}

func (impl UserStorerImpl) GetByPaymentProcessorCustomerID(ctx context.Context, paymentProcessorCustomerID string) (*User, error) {
	filter := bson.D{{"payment_processor_customer_id", paymentProcessorCustomerID}}

	var result User
	err := impl.Collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// This error means your query did not match any documents.
			return nil, nil
		}
		impl.Logger.Error("database get by verification code error", slog.Any("error", err))
		return nil, err
	}
	return &result, nil
}

func (impl UserStorerImpl) GetStripeInvoiceByPaymentProcessorInvoiceID(ctx context.Context, paymentProcessorInvoiceID string) (*StripeInvoice, error) {
	// DEVELOPERS NOTE:
	// To learn more about querying inside nested fields, then please see:
	// https://www.mongodb.com/docs/manual/tutorial/query-embedded-documents/

	filter := bson.M{"stripe_invoices.invoice_id": paymentProcessorInvoiceID}

	var result User
	err := impl.Collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// This error means your query did not match any documents.
			return nil, nil
		}
		impl.Logger.Error("database get by verification code error", slog.Any("error", err))
		return nil, err
	}

	// Now we need to find the specific StripeInvoice within the User's StripeInvoices slice.
	for _, invoice := range result.StripeInvoices {
		if invoice.InvoiceID == paymentProcessorInvoiceID {
			return invoice, nil
		}
	}

	// If the paymentProcessorInvoiceID is not found in the StripeInvoices, return nil.
	return nil, nil
}

func (impl UserStorerImpl) GetLiteByID(ctx context.Context, id primitive.ObjectID) (*UserLite, error) {
	filter := bson.D{{"_id", id}}

	var result UserLite
	err := impl.Collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// This error means your query did not match any documents.
			return nil, nil
		}
		impl.Logger.Error("database get user lite by id error", slog.Any("error", err))
		return nil, err
	}
	return &result, nil
}

func (impl UserStorerImpl) GetLiteByEmail(ctx context.Context, email string) (*UserLite, error) {
	filter := bson.D{{"email", email}}

	var result UserLite
	err := impl.Collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// This error means your query did not match any documents.
			return nil, nil
		}
		impl.Logger.Error("database get user lite by email error", slog.Any("error", err))
		return nil, err
	}
	return &result, nil
}

// Auto-generated comment for change 20
