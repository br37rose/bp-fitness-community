package controller

import (
	"context"
	"log/slog"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	tag_s "github.com/bci-innovation-labs/bp8fitnesscommunity-cli/app/tag/datastore"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-cli/config/constants"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-cli/utils/httperror"
)

type TagCreateRequestIDO struct {
	Text        string `bson:"text" json:"text"`
	Description string `bson:"description" json:"description"`
}

func (impl *TagControllerImpl) userFromCreateRequest(requestData *TagCreateRequestIDO) (*tag_s.Tag, error) {

	return &tag_s.Tag{
		Text:        requestData.Text,
		Description: requestData.Description,
	}, nil
}

func ValidateCreateRequest(dirtyData *TagCreateRequestIDO) error {
	e := make(map[string]string)

	if dirtyData.Text == "" {
		e["text"] = "missing value"
	}
	if dirtyData.Description == "" {
		e["description"] = "missing value"
	}

	if len(e) != 0 {
		return httperror.NewForBadRequest(&e)
	}
	return nil
}

func (impl *TagControllerImpl) Create(ctx context.Context, requestData *TagCreateRequestIDO) (*tag_s.Tag, error) {

	//
	// Get variables from our user authenticated session.
	//

	tid, _ := ctx.Value(constants.SessionUserOrganizationID).(primitive.ObjectID)
	// role, _ := ctx.Value(constants.SessionUserRole).(int8)
	userID, _ := ctx.Value(constants.SessionUserID).(primitive.ObjectID)
	userName, _ := ctx.Value(constants.SessionUserName).(string)
	ipAddress, _ := ctx.Value(constants.SessionIPAddress).(string)

	// DEVELOPERS NOTE:
	// Every submission needs to have a unique `public id` (PID)
	// generated. The following needs to happen to generate the unique PID:
	// 1. Make the `Create` function be `atomic` and thus lock this function.
	// 2. Count total records in system (for particular organization).
	// 3. Generate PID.
	// 4. Apply the PID to the record.
	// 5. Unlock this `Create` function to be usable again by other calls after
	//    the function successfully submits the record into our system.
	impl.Kmutex.Lockf("create-tag-by-organization-%s", tid.Hex())
	defer impl.Kmutex.Unlockf("create-tag-by-organization-%s", tid.Hex())

	//
	// Perform our validation and return validation error on any issues detected.
	//

	if err := ValidateCreateRequest(requestData); err != nil {
		return nil, err
	}

	////
	//// Start the transaction.
	////

	session, err := impl.DbClient.StartSession()
	if err != nil {
		impl.Logger.Error("start session error",
			slog.Any("error", err))
		return nil, err
	}
	defer session.EndSession(ctx)

	// Define a transaction function with a series of operations
	transactionFunc := func(sessCtx mongo.SessionContext) (interface{}, error) {

		m, err := impl.userFromCreateRequest(requestData)
		if err != nil {
			return nil, err
		}

		// Add defaults.
		m.OrganizationID = tid
		m.ID = primitive.NewObjectID()
		m.CreatedAt = time.Now()
		m.CreatedByUserID = userID
		m.CreatedByUserName = userName
		m.CreatedFromIPAddress = ipAddress
		m.ModifiedAt = time.Now()
		m.ModifiedByUserID = userID
		m.ModifiedByUserName = userName
		m.ModifiedFromIPAddress = ipAddress
		m.Text = requestData.Text
		m.Description = requestData.Description
		m.Status = tag_s.TagStatusActive

		// Save to our database.
		if err := impl.TagStorer.Create(sessCtx, m); err != nil {
			impl.Logger.Error("database create error", slog.Any("error", err))
			return nil, err
		}

		////
		//// Exit our transaction successfully.
		////

		return m, nil
	}

	// Start a transaction
	result, err := session.WithTransaction(ctx, transactionFunc)
	if err != nil {
		impl.Logger.Error("session failed error",
			slog.Any("error", err))
		return nil, err
	}

	return result.(*tag_s.Tag), nil
}
