package stripe

import (
	"context"
	"log/slog"

	"go.mongodb.org/mongo-driver/bson/primitive"

	u_d "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/user/datastore"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/config/constants"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/utils/httperror"
)

func (impl *StripePaymentProcessorControllerImpl) ListLatestStripeInvoices(ctx context.Context, userID primitive.ObjectID, cursor int64, limit int64) (*u_d.StripeInvoiceListResult, error) {
	// Extract from our session the following data.
	urole := ctx.Value(constants.SessionUserRole).(int8)
	uid := ctx.Value(constants.SessionUserID).(primitive.ObjectID)

	////
	//// Permission handling.
	////

	var pickedUserID primitive.ObjectID
	switch urole { // Security.
	case u_d.UserRoleRoot:
		return nil, httperror.NewForForbiddenWithSingleField("message", "you did not saasify product")
	case u_d.UserRoleTrainer:
		return nil, httperror.NewForForbiddenWithSingleField("message", "you do not have permission as a trainer") //TODO: Implement.
	case u_d.UserRoleMember:
		// Override whatever user ID was set to the logged in user because the
		// user is a member and thus a member only has access to his/her data.
		pickedUserID = uid
	case u_d.UserRoleAdmin:
		if userID.IsZero() {
			return nil, httperror.NewForBadRequestWithSingleField("user_id", "missing url parameter")
		}
		pickedUserID = userID
	}

	////
	//// Database lookups interactions.
	////

	// Lookup the user in our database, else return a `400 Bad Request` error.
	u, err := impl.UserStorer.GetByID(ctx, pickedUserID)
	if err != nil {
		impl.Logger.Error("database error", slog.Any("err", err))
		return nil, err
	}
	if u == nil {
		impl.Logger.Warn("user does not exist validation error")
		return nil, httperror.NewForBadRequestWithSingleField("id", "does not exist")
	}

	return impl.UserStorer.ListLatestStripeInvoices(ctx, userID, cursor, limit)
}
