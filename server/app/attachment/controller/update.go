package controller

import (
	"context"
	"fmt"
	"mime/multipart"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"log/slog"

	a_d "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/attachment/datastore"
	domain "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/attachment/datastore"
	user_d "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/user/datastore"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/config/constants"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/utils/httperror"
)

type AttachmentUpdateRequestIDO struct {
	ID            primitive.ObjectID
	Name          string
	Description   string
	OwnershipID   primitive.ObjectID
	OwnershipType int8
	FileName      string
	FileType      string
	File          multipart.File
}

func ValidateUpdateRequest(dirtyData *AttachmentUpdateRequestIDO) error {
	e := make(map[string]string)

	if dirtyData.ID.IsZero() {
		e["id"] = "missing value"
	}
	if dirtyData.Name == "" {
		e["name"] = "missing value"
	}
	if dirtyData.Description == "" {
		e["description"] = "missing value"
	}
	if dirtyData.OwnershipID.IsZero() {
		e["ownership_id"] = "missing value"
	}
	if dirtyData.OwnershipType == 0 {
		e["ownership_type"] = "missing value"
	}
	if len(e) != 0 {
		return httperror.NewForBadRequest(&e)
	}
	return nil
}

func (c *AttachmentControllerImpl) UpdateByID(ctx context.Context, req *AttachmentUpdateRequestIDO) (*domain.Attachment, error) {
	if err := ValidateUpdateRequest(req); err != nil {
		return nil, err
	}

	// Fetch the original attachment.
	os, err := c.AttachmentStorer.GetByID(ctx, req.ID)
	if err != nil {
		c.Logger.Error("database get by id error",
			slog.Any("error", err),
			slog.Any("attachment_id", req.ID))
		return nil, err
	}
	if os == nil {
		c.Logger.Error("attachment does not exist error",
			slog.Any("attachment_id", req.ID))
		return nil, httperror.NewForBadRequestWithSingleField("message", "attachment does not exist")
	}

	// Extract from our session the following data.
	userID := ctx.Value(constants.SessionUserID).(primitive.ObjectID)
	userOrganizationID := ctx.Value(constants.SessionUserOrganizationID).(primitive.ObjectID)
	userRole := ctx.Value(constants.SessionUserRole).(int8)
	userName := ctx.Value(constants.SessionUserName).(string)

	// If user is not administrator nor belongs to the attachment then error.
	if userRole != user_d.UserRoleRoot && os.OrganizationID != userOrganizationID {
		c.Logger.Error("authenticated user is not staff role nor belongs to the attachment error",
			slog.Any("userRole", userRole),
			slog.Any("userOrganizationID", userOrganizationID))
		return nil, httperror.NewForForbiddenWithSingleField("message", "you do not belong to this attachment")
	}

	// Update the file if the user uploaded a new file.
	if req.File != nil {
		// Proceed to delete the physical files from AWS s3.
		if err := c.S3.DeleteByKeys(ctx, []string{os.ObjectKey}); err != nil {
			c.Logger.Warn("s3 delete by keys error", slog.Any("error", err))
			// Do not return an error, simply continue this function as there might
			// be a case were the file was removed on the s3 bucket by ourselves
			// or some other reason.
		}

		// The following code will choose the directory we will upload based on the image type.
		var directory string
		switch req.OwnershipType {
		case a_d.OwnershipTypeUser:
			directory = "user"
		case a_d.OwnershipTypeExerciseVideo:
			directory = "submission"
		case a_d.OwnershipTypeOrganization:
			directory = "organization"
		default:
			c.Logger.Error("unsupported ownership type format", slog.Any("ownership_type", req.OwnershipType))
			return nil, fmt.Errorf("unsuported iownership type  of %v, please pick another type", req.OwnershipType)
		}

		// Generate the key of our upload.
		objectKey := fmt.Sprintf("org/%v/%v/%v/%v", userOrganizationID.Hex(), directory, req.OwnershipID.Hex(), req.FileName)

		go func(file multipart.File, objkey string) {
			c.Logger.Debug("beginning private s3 image upload...")
			if err := c.S3.UploadContentFromMulipart(context.Background(), objkey, file); err != nil {
				c.Logger.Error("private s3 upload error", slog.Any("error", err))
				// Do not return an error, simply continue this function as there might
				// be a case were the file was removed on the s3 bucket by ourselves
				// or some other reason.
			}
			c.Logger.Debug("Finished private s3 image upload")
		}(req.File, objectKey)

		// Update file.
		os.ObjectKey = objectKey
		os.Filename = req.FileName
	}

	// Modify our original attachment.
	os.ModifiedAt = time.Now()
	os.ModifiedByUserID = userID
	os.ModifiedByUserName = userName
	os.Name = req.Name
	os.Description = req.Description
	os.OwnershipID = req.OwnershipID
	os.OwnershipType = req.OwnershipType

	// Save to the database the modified attachment.
	if err := c.AttachmentStorer.UpdateByID(ctx, os); err != nil {
		c.Logger.Error("database update by id error", slog.Any("error", err))
		return nil, err
	}

	// go func(org *domain.Attachment) {
	// 	c.updateAttachmentNameForAllUsers(ctx, org)
	// }(os)
	// go func(org *domain.Attachment) {
	// 	c.updateAttachmentNameForAllComicSubmissions(ctx, org)
	// }(os)

	return os, nil
}

// func (c *AttachmentControllerImpl) updateAttachmentNameForAllUsers(ctx context.Context, ns *domain.Attachment) error {
// 	c.Logger.Debug("Beginning to update attachment name for all uses")
// 	f := &user_d.UserListFilter{AttachmentID: ns.ID}
// 	uu, err := c.UserStorer.ListByFilter(ctx, f)
// 	if err != nil {
// 		c.Logger.Error("database update by id error", slog.Any("error", err))
// 		return err
// 	}
// 	for _, u := range uu.Results {
// 		u.AttachmentName = ns.Name
// 		log.Println("--->", u)
// 		// if err := c.UserStorer.UpdateByID(ctx, u); err != nil {
// 		// 	c.Logger.Error("database update by id error", slog.Any("error", err))
// 		// 	return err
// 		// }
// 	}
// 	return nil
// }
//
// func (c *AttachmentControllerImpl) updateAttachmentNameForAllComicSubmissions(ctx context.Context, ns *domain.Attachment) error {
// 	c.Logger.Debug("Beginning to update attachment name for all submissions")
// 	f := &domain.ComicSubmissionListFilter{AttachmentID: ns.ID}
// 	uu, err := c.ComicSubmissionStorer.ListByFilter(ctx, f)
// 	if err != nil {
// 		c.Logger.Error("database update by id error", slog.Any("error", err))
// 		return err
// 	}
// 	for _, u := range uu.Results {
// 		u.AttachmentName = ns.Name
// 		log.Println("--->", u)
// 		// if err := c.ComicSubmissionStorer.UpdateByID(ctx, u); err != nil {
// 		// 	c.Logger.Error("database update by id error", slog.Any("error", err))
// 		// 	return err
// 		// }
// 	}
// 	return nil
// }
