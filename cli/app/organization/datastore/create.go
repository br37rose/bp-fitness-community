package datastore

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"log/slog"
)

func (impl OrganizationStorerImpl) Create(ctx context.Context, u *Organization) error {
	// DEVELOPER NOTES:
	// According to mongodb documentaiton:
	//     Non-existent Databases and Collections
	//     If the necessary database and collection don't exist when you perform a write operation, the server implicitly creates them.
	//     Source: https://www.mongodb.com/docs/drivers/go/current/usage-examples/insertOne/

	if u.ID == primitive.NilObjectID {
		u.ID = primitive.NewObjectID()
		impl.Logger.Warn("database insert organization not included id value, created id now.", slog.Any("id", u.ID))
	}

	result, err := impl.Collection.InsertOne(ctx, u)

	// check for errors in the insertion
	if err != nil {
		impl.Logger.Error("database insert error", slog.Any("error", err))
	}

	// display the id of the newly inserted object
	fmt.Println(result.InsertedID)

	return nil
}
