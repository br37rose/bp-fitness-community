package controller

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"log/slog"

	user_d "github.com/bci-innovation-labs/bp8fitnesscommunity-cli/app/user/datastore"
	domain "github.com/bci-innovation-labs/bp8fitnesscommunity-cli/app/videocategory/datastore"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-cli/config/constants"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-cli/utils/httperror"
)

type VideoCategoryUpdateRequestIDO struct {
	ID     primitive.ObjectID
	Name   string
	No     int8
	Status int8
}

func ValidateUpdateRequest(dirtyData *VideoCategoryUpdateRequestIDO) error {
	e := make(map[string]string)

	if dirtyData.ID.IsZero() {
		e["id"] = "missing value"
	}
	if dirtyData.Name == "" {
		e["name"] = "missing value"
	}
	if dirtyData.No == 0 {
		e["no"] = "missing value"
	}
	if dirtyData.Status == 0 {
		e["status"] = "missing value"
	}
	if len(e) != 0 {
		return httperror.NewForBadRequest(&e)
	}
	return nil
}

func (c *VideoCategoryControllerImpl) UpdateByID(ctx context.Context, req *VideoCategoryUpdateRequestIDO) (*domain.VideoCategory, error) {
	if err := ValidateUpdateRequest(req); err != nil {
		return nil, err
	}

	// Fetch the original videocategory.
	os, err := c.VideoCategoryStorer.GetByID(ctx, req.ID)
	if err != nil {
		c.Logger.Error("database get by id error",
			slog.Any("error", err),
			slog.Any("videocategory_id", req.ID))
		return nil, err
	}
	if os == nil {
		c.Logger.Error("videocategory does not exist error",
			slog.Any("videocategory_id", req.ID))
		return nil, httperror.NewForBadRequestWithSingleField("message", "videocategory does not exist")
	}

	// Extract from our session the following data.
	userID := ctx.Value(constants.SessionUserID).(primitive.ObjectID)
	userOrganizationID := ctx.Value(constants.SessionUserOrganizationID).(primitive.ObjectID)
	userRole := ctx.Value(constants.SessionUserRole).(int8)
	userName := ctx.Value(constants.SessionUserName).(string)

	// If user is not administrator nor belongs to the videocategory then error.
	if userRole != user_d.UserRoleRoot && os.OrganizationID != userOrganizationID {
		c.Logger.Error("authenticated user is not staff role nor belongs to the videocategory error",
			slog.Any("userRole", userRole),
			slog.Any("userOrganizationID", userOrganizationID))
		return nil, httperror.NewForForbiddenWithSingleField("message", "you do not belong to this videocategory")
	}

	// Modify our original videocategory.
	os.ModifiedAt = time.Now()
	os.ModifiedByUserID = userID
	os.ModifiedByUserName = userName
	os.Name = req.Name
	os.No = req.No
	os.Status = req.Status

	// Save to the database the modified videocategory.
	if err := c.VideoCategoryStorer.UpdateByID(ctx, os); err != nil {
		c.Logger.Error("database update by id error", slog.Any("error", err))
		return nil, err
	}

	return os, nil
}

// func (c *VideoCategoryControllerImpl) updateVideoCategoryNameForAllUsers(ctx context.Context, ns *domain.VideoCategory) error {
// 	c.Logger.Debug("Beginning to update videocategory name for all uses")
// 	f := &user_d.UserListFilter{VideoCategoryID: ns.ID}
// 	uu, err := c.UserStorer.ListByFilter(ctx, f)
// 	if err != nil {
// 		c.Logger.Error("database update by id error", slog.Any("error", err))
// 		return err
// 	}
// 	for _, u := range uu.Results {
// 		u.VideoCategoryName = ns.Name
// 		log.Println("--->", u)
// 		// if err := c.UserStorer.UpdateByID(ctx, u); err != nil {
// 		// 	c.Logger.Error("database update by id error", slog.Any("error", err))
// 		// 	return err
// 		// }
// 	}
// 	return nil
// }
//
// func (c *VideoCategoryControllerImpl) updateVideoCategoryNameForAllComicSubmissions(ctx context.Context, ns *domain.VideoCategory) error {
// 	c.Logger.Debug("Beginning to update videocategory name for all submissions")
// 	f := &domain.ComicSubmissionListFilter{VideoCategoryID: ns.ID}
// 	uu, err := c.ComicSubmissionStorer.ListByFilter(ctx, f)
// 	if err != nil {
// 		c.Logger.Error("database update by id error", slog.Any("error", err))
// 		return err
// 	}
// 	for _, u := range uu.Results {
// 		u.VideoCategoryName = ns.Name
// 		log.Println("--->", u)
// 		// if err := c.ComicSubmissionStorer.UpdateByID(ctx, u); err != nil {
// 		// 	c.Logger.Error("database update by id error", slog.Any("error", err))
// 		// 	return err
// 		// }
// 	}
// 	return nil
// }

// Auto-generated comment for change 25
