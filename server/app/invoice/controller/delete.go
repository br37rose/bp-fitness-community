package controller

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"log/slog"

	s_d "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/invoice/datastore"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/utils/httperror"
)

func (impl *InvoiceControllerImpl) DeleteByID(ctx context.Context, id primitive.ObjectID) error {
	d, err := impl.GetByID(ctx, id)
	if err != nil {
		impl.Logger.Error("database get by id error", slog.Any("error", err))
		return err
	}
	if d == nil {
		impl.Logger.Error("database returns nothing from get by id")
		return httperror.NewForBadRequestWithSingleField("id", "workout program type does not exist")
	}
	d.Status = s_d.StatusArchived
	d.ModifiedAt = time.Now()

	// Save to the database the modified organization.
	if err := impl.InvoiceStorer.UpdateByID(ctx, d); err != nil {
		impl.Logger.Error("database update by id error", slog.Any("error", err))
		return err
	}

	return nil
}
