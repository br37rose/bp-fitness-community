package controller

import (
	"context"
	"log/slog"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	a_d "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/fitnessplan/datastore"
	domain "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/fitnessplan/datastore"
	user_d "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/user/datastore"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/config/constants"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/utils/httperror"
)

type FitnessPlanUpdateRequestIDO struct {
	ID                 primitive.ObjectID `bson:"id,omitempty" json:"id,omitempty"`
	EstimatedReadyDate time.Time          `bson:"estimated_ready_date,omitempty" json:"estimated_ready_date,omitempty"`
	MaxWeeks           int8               `bson:"max_weeks" json:"max_weeks"`
	Name               string             `bson:"name" json:"name"`
	Status             int8               `bson:"status" json:"status"`
	TimePerDay         int8               `bson:"time_per_day" json:"time_per_day"`
	WasProcessed       bool               `bson:"was_processed" json:"was_processed"`
	UserId             primitive.ObjectID `bson:"user_id" json:"user_id"`
}

func ValidateUpdateRequest(dirtyData *FitnessPlanUpdateRequestIDO) error {
	e := make(map[string]string)

	if dirtyData.ID.IsZero() {
		e["id"] = "missing value"
	}

	if dirtyData.Name == "" {
		e["name"] = "missing value"
	}

	if len(e) != 0 {
		return httperror.NewForBadRequest(&e)
	}
	return nil
}

func (c *FitnessPlanControllerImpl) UpdateByID(ctx context.Context, req *FitnessPlanUpdateRequestIDO) (*domain.FitnessPlan, error) {
	if err := ValidateUpdateRequest(req); err != nil {
		return nil, err
	}

	c.Kmutex.Lockf("update-for-fitness-plan-%s", req.ID.Hex())
	defer c.Kmutex.Unlockf("update-for-fitness-plan-%s", req.ID.Hex())

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
		// Fetch the original fitnessplan.
		os, err := c.FitnessPlanStorer.GetByID(sessCtx, req.ID)
		if err != nil {
			c.Logger.Error("database get by id error",
				slog.Any("error", err),
				slog.Any("fitnessplan_id", req.ID))
			return nil, err
		}
		if os == nil {
			c.Logger.Error("fitnessplan does not exist error",
				slog.Any("fitnessplan_id", req.ID))
			return nil, httperror.NewForBadRequestWithSingleField("message", "fitnessplan does not exist")
		}

		// Prevent any updates if we submitted to open-ai without open-ai finished.
		if os.Status == domain.StatusQueued {
			return nil, httperror.NewForBadRequestWithSingleField("message", "your fitness plan is being queued by us and you will be able to resubmit when we finish")
		}

		// Extract from our session the following data.
		userID := sessCtx.Value(constants.SessionUserID).(primitive.ObjectID)
		userOrganizationID := sessCtx.Value(constants.SessionUserOrganizationID).(primitive.ObjectID)
		userRole := sessCtx.Value(constants.SessionUserRole).(int8)
		userName := sessCtx.Value(constants.SessionUserName).(string)

		// If user is not administrator nor belongs to the fitnessplan then error.
		if userRole != user_d.UserRoleRoot && os.OrganizationID != userOrganizationID {
			c.Logger.Error("authenticated user is not staff role nor belongs to the fitnessplan error",
				slog.Any("userRole", userRole),
				slog.Any("userOrganizationID", userOrganizationID))
			return nil, httperror.NewForForbiddenWithSingleField("message", "you do not belong to this fitnessplan")
		}

		if userRole == user_d.UserRoleAdmin {
			if req.UserId.IsZero() {
				return nil, httperror.NewForBadRequestWithSingleField("userid", "missing value")
			}
			userID = req.UserId
		}

		// Modify our original fitnessplan.
		os.ModifiedAt = time.Now()
		os.ModifiedByUserID = userID
		os.ModifiedByUserName = userName
		os.EstimatedReadyDate = req.EstimatedReadyDate
		os.Name = req.Name
		os.WasProcessed = req.WasProcessed
		os.Status = domain.StatusQueued
		os.UserID = userID

		// Save to the database the modified fitnessplan.
		if err := c.FitnessPlanStorer.UpdateByID(sessCtx, os); err != nil {
			c.Logger.Error("database update by id error", slog.Any("error", err))
			return nil, err
		}
		return os, nil
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

// func (c *FitnessPlanControllerImpl) updateFitnessPlanNameForAllUsers(ctx context.Context, ns *domain.FitnessPlan) error {
// 	c.Logger.Debug("Beginning to update fitnessplan name for all uses")
// 	f := &user_d.UserListFilter{FitnessPlanID: ns.ID}
// 	uu, err := c.UserStorer.ListByFilter(ctx, f)
// 	if err != nil {
// 		c.Logger.Error("database update by id error", slog.Any("error", err))
// 		return err
// 	}
// 	for _, u := range uu.Results {
// 		u.FitnessPlanName = ns.Name
// 		log.Println("--->", u)
// 		// if err := c.UserStorer.UpdateByID(ctx, u); err != nil {
// 		// 	c.Logger.Error("database update by id error", slog.Any("error", err))
// 		// 	return err
// 		// }
// 	}
// 	return nil
// }
//
// func (c *FitnessPlanControllerImpl) updateFitnessPlanNameForAllComicSubmissions(ctx context.Context, ns *domain.FitnessPlan) error {
// 	c.Logger.Debug("Beginning to update fitnessplan name for all submissions")
// 	f := &domain.ComicSubmissionListFilter{FitnessPlanID: ns.ID}
// 	uu, err := c.ComicSubmissionStorer.ListByFilter(ctx, f)
// 	if err != nil {
// 		c.Logger.Error("database update by id error", slog.Any("error", err))
// 		return err
// 	}
// 	for _, u := range uu.Results {
// 		u.FitnessPlanName = ns.Name
// 		log.Println("--->", u)
// 		// if err := c.ComicSubmissionStorer.UpdateByID(ctx, u); err != nil {
// 		// 	c.Logger.Error("database update by id error", slog.Any("error", err))
// 		// 	return err
// 		// }
// 	}
// 	return nil
// }
