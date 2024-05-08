package controller

import (
	"context"
	"time"

	"log/slog"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	a_d "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/fitnessplan/datastore"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/user/datastore"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/config/constants"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/utils/httperror"
)

type FitnessPlanCreateRequestIDO struct {
	EstimatedReadyDate time.Time          `bson:"estimated_ready_date,omitempty" json:"estimated_ready_date,omitempty"`
	Name               string             `bson:"name" json:"name"`
	Status             int8               `bson:"status" json:"status"`
	TimePerDay         int8               `bson:"time_per_day" json:"time_per_day"`
	WasProcessed       bool               `bson:"was_processed" json:"was_processed"`
	UserId             primitive.ObjectID `bson:"user_id" json:"user_id"`
}

func ValidateCreateRequest(dirtyData *FitnessPlanCreateRequestIDO) error {
	e := make(map[string]string)

	if dirtyData.Name == "" {
		e["name"] = "missing value"
	}

	if len(e) != 0 {
		return httperror.NewForBadRequest(&e)
	}
	return nil
}

func (c *FitnessPlanControllerImpl) Create(ctx context.Context, req *FitnessPlanCreateRequestIDO) (*a_d.FitnessPlan, error) {
	// Extract from our session the following data.
	orgID := ctx.Value(constants.SessionUserOrganizationID).(primitive.ObjectID)
	orgName := ctx.Value(constants.SessionUserOrganizationName).(string)
	userID := ctx.Value(constants.SessionUserID).(primitive.ObjectID)
	userName := ctx.Value(constants.SessionUserName).(string)
	userLexicalName := ctx.Value(constants.SessionUserLexicalName).(string)
	userRole, _ := ctx.Value(constants.SessionUserRole).(int8)

	if err := ValidateCreateRequest(req); err != nil {
		return nil, err
	}
	if userRole == datastore.UserRoleAdmin {
		if req.UserId.IsZero() {
			return nil, httperror.NewForBadRequestWithSingleField("userid", "missing value")
		}
		userID = req.UserId
	}

	////
	//// Start the transaction.
	////

	session, err := c.DbClient.StartSession()
	if err != nil {
		c.Logger.Error("start session error",
			slog.Any("error", err))
		return nil, err
	}
	defer session.EndSession(ctx)

	// Define a transaction function with a series of operations
	transactionFunc := func(sessCtx mongo.SessionContext) (interface{}, error) {

		// Create our record in the database.
		res := &a_d.FitnessPlan{
			OrganizationID:     orgID,
			OrganizationName:   orgName,
			ID:                 primitive.NewObjectID(),
			CreatedAt:          time.Now(),
			CreatedByUserName:  userName,
			CreatedByUserID:    userID,
			ModifiedAt:         time.Now(),
			ModifiedByUserName: userName,
			ModifiedByUserID:   userID,
			Status:             a_d.StatusQueued, // Set to queued b/c we will wait on openai.
			Name:               req.Name,
			ExerciseNames:      []string{},
			Exercises:          make([]*a_d.FitnessPlanExercise, 0),
			UserID:             userID,
			UserName:           userName,
			UserLexicalName:    userLexicalName,
		}
		err := c.FitnessPlanStorer.Create(sessCtx, res)
		if err != nil {
			c.Logger.Error("fitnessplan create error", slog.Any("error", err))
			return nil, err
		}

		return res, nil
	}

	// Start a transaction
	result, err := session.WithTransaction(ctx, transactionFunc)
	if err != nil {
		c.Logger.Error("session failed error",
			slog.Any("error", err))
		return nil, err
	}

	// Execute in the background the call to OpenAI and generate the fitness
	// plan.
	if result != nil {
		go c.generateFitnessPlanInBackground(context.Background(), result.(*a_d.FitnessPlan))
	}

	return result.(*a_d.FitnessPlan), nil
}
