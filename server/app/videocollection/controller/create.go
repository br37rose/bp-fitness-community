package controller

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"log/slog"

	vcol_d "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/videocollection/datastore"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/config/constants"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/utils/httperror"
)

type VideoCollectionCreateRequestIDO struct {
	Name        string `bson:"name" json:"name"`
	Summary     string `bson:"summary" json:"summary"`
	Description string `bson:"description" json:"description"`

	ThumbnailType   int8   `bson:"thumbnail_type" json:"thumbnail_type"`
	ThumbnailUpload string `bson:"thumbnail_upload" json:"thumbnail_upload"`
	ThumbnailURL    string `bson:"thumbnail_url" json:"thumbnail_url"`

	Type       int8               `bson:"type" json:"type"`
	CategoryID primitive.ObjectID `bson:"category_id" json:"category_id,omitempty"`
}

func ValidateCreateRequest(dirtyData *VideoCollectionCreateRequestIDO) error {
	e := make(map[string]string)

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

		if dirtyData.Type == 0 {
			e["type"] = "missing value"
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
		if dirtyData.Type == 0 {
			e["type"] = "missing value"
		}
		if dirtyData.CategoryID.IsZero() {
			e["category_id"] = "missing value"
		}
	}

	if len(e) != 0 {
		return httperror.NewForBadRequest(&e)
	}
	return nil
}

func (c *VideoCollectionControllerImpl) Create(ctx context.Context, req *VideoCollectionCreateRequestIDO) (*vcol_d.VideoCollection, error) {
	if err := ValidateCreateRequest(req); err != nil {
		c.Logger.Error("validation failure", slog.Any("error", err))
		return nil, err
	}

	// Extract from our session the following data.
	orgID, _ := ctx.Value(constants.SessionUserOrganizationID).(primitive.ObjectID)
	userID, _ := ctx.Value(constants.SessionUserID).(primitive.ObjectID)
	userName, _ := ctx.Value(constants.SessionUserName).(string)

	e := &vcol_d.VideoCollection{}

	// Meta information.
	e.OrganizationID = orgID
	e.ID = primitive.NewObjectID()
	e.CreatedByUserName = userName
	e.CreatedByUserID = userID
	e.CreatedAt = time.Now()
	e.ModifiedByUserName = userName
	e.ModifiedByUserID = userID
	e.ModifiedAt = time.Now()
	e.Status = vcol_d.StatusActive

	// Base information.
	e.Type = req.Type
	e.Name = req.Name
	e.Summary = req.Summary
	e.Description = req.Description

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
	e.CategoryID = vc.ID
	e.CategoryName = vc.Name

	// Mendia content.
	e.ThumbnailType = req.ThumbnailType
	e.ThumbnailURL = req.ThumbnailURL

	//
	// Process thumbnail file uploads here...
	//

	if e.ThumbnailType == vcol_d.VideoCollectionThumbnailTypeS3 {

		//
		// STEP 1: Fetch the video attachment or error if D.N.E.
		//

		aid, err := primitive.ObjectIDFromHex(req.ThumbnailUpload)
		if err != nil {
			c.Logger.Error("thumbnail object id from hex error", slog.Any("error", err))
			return nil, err
		}
		attachment, err := c.AttachmentStorer.GetByID(ctx, aid)
		if err != nil {
			c.Logger.Error("thumbnail attachment get by id error", slog.Any("error", err))
			return nil, err
		}
		if attachment == nil {
			c.Logger.Error("thumbnail attachment does not exist", slog.Any("aid", aid))
			return nil, httperror.NewForBadRequestWithSingleField("thumbnail_upload", "does not exist")
		}

		//
		// STEP 2: Copy attachment into our videocollection.
		//

		// Generate the key of our videocollection video upload.
		thumbnailObjectKey := fmt.Sprintf("org/%v/video-collection/%v/%v", orgID.Hex(), e.ID.Hex(), attachment.Filename)

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
		e.ThumbnailObjectKey = thumbnailObjectKey
		e.ThumbnailObjectURL = thumbnailObjectURL
		e.ThumbnailObjectExpiry = time.Now().Add(expiryDur)
		e.ThumbnailAttachmentID = attachment.ID
		e.ThumbnailAttachmentFilename = attachment.Filename

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

	//
	// Save to our database.
	//

	if err := c.VideoCollectionStorer.Create(ctx, e); err != nil {
		c.Logger.Error("database create error", slog.Any("error", err))
		return nil, err
	}

	return e, nil
}
