package controller

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"log/slog"

	usr_d "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/user/datastore"
	vcol_d "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/videocollection/datastore"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/config/constants"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/utils/httperror"
)

type VideoCollectionUpdateRequestIDO struct {
	ID primitive.ObjectID `bson:"id" json:"id"`

	Name        string `bson:"name" json:"name"`
	Summary     string `bson:"summary" json:"summary"`
	Description string `bson:"description" json:"description"`

	ThumbnailType   int8   `bson:"thumbnail_type" json:"thumbnail_type"`
	ThumbnailUpload string `bson:"thumbnail_upload" json:"thumbnail_upload"`
	ThumbnailURL    string `bson:"thumbnail_url" json:"thumbnail_url"`

	Type       int8               `bson:"type" json:"type"`
	CategoryID primitive.ObjectID `bson:"category_id" json:"category_id,omitempty"`
	Status     int8               `bson:"status" json:"status"`
}

func ValidateUpdateRequest(dirtyData *VideoCollectionUpdateRequestIDO) error {
	e := make(map[string]string)

	if dirtyData.ID.IsZero() {
		e["id"] = "missing value"
	}

	if dirtyData.Name == "" {
		e["name"] = "missing value"
	}
	if dirtyData.Summary == "" {
		e["summary"] = "missing value"
	}
	if dirtyData.Description == "" {
		e["description"] = "missing value"
	}

	// Thumbnail upload
	if dirtyData.ThumbnailType == 0 {
		e["thumbnail_type"] = "missing value"
	} else {
		if dirtyData.ThumbnailType == vcol_d.VideoCollectionThumbnailTypeS3 && dirtyData.ThumbnailUpload == "" {
			e["thumbnail_upload"] = "missing value"
		}
		if dirtyData.ThumbnailType == vcol_d.VideoCollectionThumbnailTypeExternalURL && dirtyData.ThumbnailURL == "" {
			e["thumbnail_url"] = "missing value"
		}
	}

	if dirtyData.Type == 0 {
		e["type"] = "missing value"
	}
	if dirtyData.CategoryID.IsZero() {
		e["category_id"] = "missing value"
	}
	if dirtyData.Status == 0 {
		e["status"] = "missing value"
	}

	if len(e) != 0 {
		return httperror.NewForBadRequest(&e)
	}
	return nil
}

func (c *VideoCollectionControllerImpl) UpdateByID(ctx context.Context, req *VideoCollectionUpdateRequestIDO) (*vcol_d.VideoCollection, error) {
	if err := ValidateUpdateRequest(req); err != nil {
		c.Logger.Error("validation failure", slog.Any("error", err))
		return nil, err
	}

	// Lock this database record when we are updating it so in case we don't cause any data inconsistency when system refreshes the presigned URL in the detail and list functions.
	c.Kmutex.Lockf("videocollection-%v", req.ID.Hex()) // Step 1
	defer func() {
		c.Kmutex.Unlockf("videocollection-%v", req.ID.Hex()) // Step 2
	}()

	// Extract from our session the following data.
	orgID, _ := ctx.Value(constants.SessionUserOrganizationID).(primitive.ObjectID)
	userID, _ := ctx.Value(constants.SessionUserID).(primitive.ObjectID)
	userName, _ := ctx.Value(constants.SessionUserName).(string)
	userRole, _ := ctx.Value(constants.SessionUserRole).(int8)

	// Fetch the original videocollection.
	os, err := c.VideoCollectionStorer.GetByID(ctx, req.ID)
	if err != nil {
		c.Logger.Error("database get by id error", slog.Any("error", err))
		return nil, err
	}
	if os == nil {
		return nil, httperror.NewForBadRequestWithSingleField("id", "videocollection does not exist")
	}

	// Defensive code: Do not edit system.
	if userRole != usr_d.UserRoleAdmin && userRole != usr_d.UserRoleRoot {
		c.Logger.Error("system videocollection edit error")
		return nil, httperror.NewForForbiddenWithSingleField("message", "you are not an administrator")
	}

	// Modify our original videocollection.
	os.ModifiedByUserName = userName
	os.ModifiedByUserID = userID
	os.ModifiedAt = time.Now()
	os.OrganizationID = orgID
	os.Status = req.Status

	// Base information.
	os.Type = req.Type
	os.Name = req.Name
	os.Summary = req.Summary
	os.Description = req.Description

	// Associate information.
	vc, err := c.VideoCategoryStorer.GetByID(ctx, req.CategoryID)
	if err != nil {
		c.Logger.Error("get video category id from error", slog.Any("error", err))
		return nil, err
	}
	if vc == nil {
		c.Logger.Error("video category does not exist", slog.Any("category_id", req.CategoryID))
		return nil, httperror.NewForBadRequestWithSingleField("category_id", "does not exist")
	}
	os.CategoryID = vc.ID
	os.CategoryName = vc.Name

	// Mendia content.
	os.ThumbnailType = req.ThumbnailType
	os.ThumbnailURL = req.ThumbnailURL

	//
	// Process thumbnail file uploads here...
	//

	if os.ThumbnailType == vcol_d.VideoCollectionThumbnailTypeS3 {

		//
		// STEP 1: Fetch the video attachment or error if D.N.E.
		//

		aid, err := primitive.ObjectIDFromHex(req.ThumbnailUpload)
		if err != nil {
			c.Logger.Error("object id from hex error", slog.Any("error", err))
			return nil, err
		}
		attachment, err := c.AttachmentStorer.GetByID(ctx, aid)
		if err != nil {
			c.Logger.Error("attachment get by id error", slog.Any("error", err))
			return nil, err
		}

		//
		// STEP 2: Copy attachment into our video collection IF we have a new attachment.
		//

		if os.ThumbnailAttachmentID != aid {
			if attachment == nil {
				c.Logger.Error("thumbnail attachment does not exist", slog.Any("aid", aid))
				return nil, httperror.NewForBadRequestWithSingleField("thumbnail_upload", "does not exist")
			}

			// Generate the key of our videocollection video upload.
			thumbnailObjectKey := fmt.Sprintf("org/%v/video-collection/%v/%v", orgID.Hex(), os.ID.Hex(), attachment.Filename)

			// Cut the object from the attachment temporary location into our new videocollection location.
			if cutErr := c.S3.Cut(ctx, attachment.ObjectKey, thumbnailObjectKey); cutErr != nil {
				c.Logger.Error("thumbnail s3 cut error", slog.Any("error", err))
				return nil, err
			}

			// Generate a presigned URL for today.
			expiryDur := time.Hour * 12
			thumbnailObjectURL, presignErr := c.S3.GetPresignedURL(ctx, thumbnailObjectKey, expiryDur)
			if presignErr != nil {
				c.Logger.Error("thumbnail s3 presign url error", slog.Any("presignErr", presignErr))
				return nil, err
			}

			// Update the videocollection.
			os.ThumbnailObjectKey = thumbnailObjectKey
			os.ThumbnailObjectURL = thumbnailObjectURL
			os.ThumbnailObjectExpiry = time.Now().Add(expiryDur)
			os.ThumbnailAttachmentID = attachment.ID
			os.ThumbnailAttachmentFilename = attachment.Filename

			//
			// STEP 3: Update the attachment.
			//

			attachment.ObjectKey = thumbnailObjectKey
			attachment.ObjectURL = thumbnailObjectURL
			if updateErr := c.AttachmentStorer.UpdateByID(ctx, attachment); updateErr != nil {
				c.Logger.Error("video attachment update error", slog.Any("updateErr", updateErr))
				return nil, updateErr
			}

		}
	}

	//
	// Save to the database the modified videocollection.
	//

	if err := c.VideoCollectionStorer.UpdateByID(ctx, os); err != nil {
		c.Logger.Error("database update by id error", slog.Any("error", err))
		return nil, err
	}

	return os, nil
}
