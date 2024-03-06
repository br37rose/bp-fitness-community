package datastore

import (
	"context"
	"log/slog"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (impl TagStorerImpl) Create(ctx context.Context, u *Tag) error {
	// DEVELOPER NOTES:
	// According to mongodb documentaiton:
	//     Non-existent Databases and Collections
	//     If the necessary database and collection don't exist when you perform a write operation, the server implicitly creates them.
	//     Source: https://www.mongodb.com/docs/drivers/go/current/usage-examples/insertOne/

	if u.ID == primitive.NilObjectID {
		u.ID = primitive.NewObjectID()
		impl.Logger.Warn("database insert user not included id value, created id now.", slog.Any("id", u.ID))
	}

	// If `public_is` not explicitly set then we implicitly set it.
	if u.PublicID == 0 {
		publicID, err := impl.generatePublicID(ctx, u.OrganizationID)
		if err != nil {
			return err
		}
		u.PublicID = publicID
	}

	_, err := impl.Collection.InsertOne(ctx, u)

	// check for errors in the insertion
	if err != nil {
		impl.Logger.Error("database insert error", slog.Any("error", err))
	}

	return nil
}

func (impl TagStorerImpl) generatePublicID(ctx context.Context, organizationID primitive.ObjectID) (uint64, error) {
	var publicID uint64
	latest, err := impl.GetLatestByOrganizationID(ctx, organizationID)
	if err != nil {
		impl.Logger.Error("database get latest tag by organization id error",
			slog.Any("error", err),
			slog.Any("organization_id", organizationID))
		return 0, err
	}
	if latest == nil {
		impl.Logger.Debug("first tag creation detected, setting publicID to value of 1",
			slog.Any("organization_id", organizationID))
		publicID = 1
	} else {
		publicID = latest.PublicID + 1
		impl.Logger.Debug("system generated new tag publicID",
			slog.Int("organization_id", int(publicID)))
	}
	return publicID, nil
}
