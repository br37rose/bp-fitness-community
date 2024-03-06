package controller

import (
	"context"
	"fmt"
	"mime/multipart"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"log/slog"

	a_d "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/attachment/datastore"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/config/constants"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/utils/httperror"
)

type AttachmentCreateRequestIDO struct {
	Name          string
	Description   string
	OwnershipID   primitive.ObjectID
	OwnershipType int8
	FileName      string
	FileType      string
	File          multipart.File
}

func ValidateCreateRequest(dirtyData *AttachmentCreateRequestIDO) error {
	e := make(map[string]string)

	// if dirtyData.Name == "" {
	// 	e["name"] = "missing value"
	// }
	// if dirtyData.Description == "" {
	// 	e["description"] = "missing value"
	// }
	// if dirtyData.OwnershipID.IsZero() {
	// 	e["ownership_id"] = "missing value"
	// }
	// if dirtyData.OwnershipType == 0 {
	// 	e["ownership_type"] = "missing value"
	// }
	if dirtyData.FileName == "" {
		e["file"] = "missing value"
	}
	if len(e) != 0 {
		return httperror.NewForBadRequest(&e)
	}
	return nil
}

func (c *AttachmentControllerImpl) Create(ctx context.Context, req *AttachmentCreateRequestIDO) (*a_d.Attachment, error) {
	// Extract from our session the following data.
	orgID := ctx.Value(constants.SessionUserOrganizationID).(primitive.ObjectID)
	orgName := ctx.Value(constants.SessionUserOrganizationName).(string)
	userID := ctx.Value(constants.SessionUserID).(primitive.ObjectID)
	userName := ctx.Value(constants.SessionUserName).(string)

	if err := ValidateCreateRequest(req); err != nil {
		return nil, err
	}

	// The following code will choose the directory we will upload based on the image type.
	var directory string
	switch req.OwnershipType {
	case a_d.OwnershipTypeTemporary:
		directory = "temps"
	case a_d.OwnershipTypeUser:
		directory = "users"
	case a_d.OwnershipTypeExerciseVideo:
		directory = "exercises/videos"
	case a_d.OwnershipTypeExerciseThumbnail:
		directory = "exercises/thumbnails"
	case a_d.OwnershipTypeOrganization:
		directory = "organizations"
	default:
		// If not specified then automatically assign to the temporary folder.
		directory = "temps"
		req.OwnershipType = a_d.OwnershipTypeTemporary
	}

	// Generate the key of our upload.
	objectKey := fmt.Sprintf("org/%v/%v/%v/%v", orgID.Hex(), directory, req.OwnershipID.Hex(), req.FileName)

	// For debugging purposes only.
	c.Logger.Debug("pre-upload meta",
		slog.String("FileName", req.FileName),
		slog.String("FileType", req.FileType),
		slog.String("Directory", directory),
		slog.String("ObjectKey", objectKey),
		slog.String("Name", req.Name),
		slog.String("Desc", req.Description),
	)

	go func(file multipart.File, objkey string) {
		c.Logger.Debug("beginning private s3 file upload...")
		if err := c.S3.UploadContentFromMulipart(context.Background(), objkey, file); err != nil {
			c.Logger.Error("private s3 file upload error", slog.Any("error", err))
			// Do not return an error, simply continue this function as there might
			// be a case were the file was removed on the s3 bucket by ourselves
			// or some other reason.
		}
		c.Logger.Debug("Finished private s3 file upload")
	}(req.File, objectKey)

	// Create our meta record in the database.
	res := &a_d.Attachment{
		OrganizationID:     orgID,
		OrganizationName:   orgName,
		ID:                 primitive.NewObjectID(),
		CreatedAt:          time.Now(),
		CreatedByUserName:  userName,
		CreatedByUserID:    userID,
		ModifiedAt:         time.Now(),
		ModifiedByUserName: userName,
		ModifiedByUserID:   userID,
		Name:               req.Name,
		Description:        req.Description,
		Filename:           req.FileName,
		ObjectKey:          objectKey,
		ObjectURL:          "",
		OwnershipID:        req.OwnershipID,
		OwnershipType:      req.OwnershipType,
		Status:             a_d.StatusActive,
	}
	err := c.AttachmentStorer.Create(ctx, res)
	if err != nil {
		c.Logger.Error("attachment create error", slog.Any("error", err))
		return nil, err
	}
	return res, nil
}
