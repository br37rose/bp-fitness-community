package controller

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"log/slog"

	user_s "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/user/datastore"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/utils/httperror"
)

func (c *MemberControllerImpl) GetByID(ctx context.Context, id primitive.ObjectID) (*user_s.User, error) {
	// Retrieve from our database the record for the specific id.
	m, err := c.UserStorer.GetByID(ctx, id)
	if err != nil {
		c.Logger.Error("database get by id error", slog.Any("error", err))
		return nil, err
	}
	if m == nil {
		return nil, httperror.NewForBadRequestWithSingleField("id", "member does not exist")
	}
	return m, err
}