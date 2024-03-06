package controller

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"log/slog"

	vcon_s "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/videocontent/datastore"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/config/constants"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/utils/httperror"
)

type VideoContentCreateRequestIDO struct {
	Name        string `bson:"name" json:"name"`
	No          int8   `bson:"no" json:"no"`
	Description string `bson:"description" json:"description"`
	AuthorName  string `bson:"author_name" json:"author_name"`
	AuthorURL   string `bson:"author_url" json:"author_url"`
	Duration    string `bson:"duration" json:"duration"`

	Type int8 `bson:"type" json:"type"`

	VideoType   int8   `bson:"video_type" json:"video_type"`
	VideoUpload string `bson:"video_upload" json:"video_upload"`
	VideoURL    string `bson:"video_url" json:"video_url"`

	ThumbnailType   int8   `bson:"thumbnail_type" json:"thumbnail_type"`
	ThumbnailUpload string `bson:"thumbnail_upload" json:"thumbnail_upload"`
	ThumbnailURL    string `bson:"thumbnail_url" json:"thumbnail_url"`

	CategoryID   primitive.ObjectID `bson:"category_id" json:"category_id,omitempty"`
	CollectionID primitive.ObjectID `bson:"collection_id" json:"collection_id,omitempty"`

	HasMonetization bool               `bson:"has_monetization" json:"has_monetization"`
	OfferID         primitive.ObjectID `bson:"offer_id" json:"offer_id"`
	HasTimedLock    bool               `bson:"has_timed_lock" json:"has_timed_lock"`
	TimedLock       string             `bson:"timed_lock" json:"timed_lock"`
}

func ValidateCreateRequest(dirtyData *VideoContentCreateRequestIDO) error {
	e := make(map[string]string)

	if dirtyData.Name == "" {
		e["name"] = "missing value"
	}
	if dirtyData.No == 0 {
		e["no"] = "missing value"
	}
	if dirtyData.Description == "" {
		e["description"] = "missing value"
	}
	if dirtyData.AuthorName == "" {
		e["author_name"] = "missing value"
	}
	// if dirtyData.AuthorURL == "" {
	// 	e["author_url"] = "missing value"
	// }
	if dirtyData.Duration == "" {
		e["duration"] = "missing value"
	}
	if dirtyData.Type == 0 {
		e["type"] = "missing value"
	}

	// Video upload
	if dirtyData.VideoType == 0 {
		e["video_type"] = "missing value"
	} else {
		if dirtyData.VideoType == vcon_s.VideoContentVideoTypeS3 && dirtyData.VideoUpload == "" {
			e["video_upload"] = "missing value"
		}
		if dirtyData.VideoType == vcon_s.VideoContentVideoTypeYouTube && dirtyData.VideoURL == "" {
			e["video_url"] = "missing value"
		}
		if dirtyData.VideoType == vcon_s.VideoContentVideoTypeVimeo && dirtyData.VideoURL == "" {
			e["video_url"] = "missing value"
		}
	}

	// Thumbnail upload
	if dirtyData.ThumbnailType == 0 {
		e["thumbnail_type"] = "missing value"
	} else {
		if dirtyData.ThumbnailType == vcon_s.VideoContentThumbnailTypeS3 && dirtyData.ThumbnailUpload == "" {
			e["thumbnail_upload"] = "missing value"
		}
		if dirtyData.ThumbnailType == vcon_s.VideoContentThumbnailTypeExternalURL && dirtyData.ThumbnailURL == "" {
			e["thumbnail_url"] = "missing value"
		}
	}

	if dirtyData.CategoryID.IsZero() {
		e["category_id"] = "missing value"
	}
	if dirtyData.CollectionID.IsZero() {
		e["collection_id"] = "missing value"
	}

	if dirtyData.HasMonetization {
		if dirtyData.OfferID.IsZero() {
			e["offer_id"] = "missing value"
		}
		if dirtyData.HasTimedLock {
			if dirtyData.TimedLock == "" {
				e["timed_lock"] = "missing value"
			}
		}
	}

	if len(e) != 0 {
		return httperror.NewForBadRequest(&e)
	}
	return nil
}

func (c *VideoContentControllerImpl) Create(ctx context.Context, req *VideoContentCreateRequestIDO) (*vcon_s.VideoContent, error) {
	if err := ValidateCreateRequest(req); err != nil {
		c.Logger.Error("validation failure", slog.Any("error", err))
		return nil, err
	}

	// Convert the string to time.Duration
	timedLockDuration, err := time.ParseDuration(req.TimedLock)
	if err != nil && req.TimedLock != "" {
		c.Logger.Error("parse duration err", slog.Any("TimedLock", req.TimedLock))
		return nil, err
	}

	// Extract from our session the following data.
	orgID, _ := ctx.Value(constants.SessionUserOrganizationID).(primitive.ObjectID)
	userID, _ := ctx.Value(constants.SessionUserID).(primitive.ObjectID)
	userName, _ := ctx.Value(constants.SessionUserName).(string)

	e := &vcon_s.VideoContent{}

	// Meta information.
	e.ID = primitive.NewObjectID()
	e.OrganizationID = orgID
	e.Status = vcon_s.VideoContentStatusActive
	e.CreatedByUserName = userName
	e.CreatedByUserID = userID
	e.CreatedAt = time.Now()
	e.ModifiedByUserName = userName
	e.ModifiedByUserID = userID
	e.ModifiedAt = time.Now()

	// Base information.
	e.Name = req.Name
	e.No = req.No
	e.Description = req.Description
	e.AuthorName = req.AuthorName
	e.AuthorURL = req.AuthorURL
	e.Duration = req.Duration
	e.Type = req.Type

	// Fetch offer
	if req.HasMonetization {
		offer, err := c.OfferStorer.GetByID(ctx, req.OfferID)
		if err != nil {
			c.Logger.Error("getting offer error", slog.Any("error", err))
			return nil, err
		}
		if offer == nil {
			c.Logger.Error("offer does not exist", slog.Any("offer_id", req.OfferID))
			return nil, httperror.NewForBadRequestWithSingleField("offer_id", "does not exist")
		}
		e.OfferID = offer.ID
		e.OfferName = offer.Name
		e.OfferMembershipRank = offer.MembershipRank
		e.HasTimedLock = req.HasTimedLock
		e.TimedLock = req.TimedLock
		e.TimedLockDuration = timedLockDuration
	}
	e.HasMonetization = req.HasMonetization

	// Get associate reference. (part 1 of 2)
	cat, err := c.VideoCategoryStorer.GetByID(ctx, req.CategoryID)
	if err != nil {
		c.Logger.Error("video category get id from hex error", slog.Any("error", err))
		return nil, err
	}
	if cat == nil {
		c.Logger.Error("video category does not exist", slog.Any("category_id", req.CategoryID))
		return nil, httperror.NewForBadRequestWithSingleField("category_id", "does not exist")
	}
	e.CategoryID = cat.ID
	e.CategoryName = cat.Name

	// Get associate reference. (part 2 of 2)
	col, err := c.VideoCollectionStorer.GetByID(ctx, req.CollectionID)
	if err != nil {
		c.Logger.Error("video collection get id from hex error", slog.Any("error", err))
		return nil, err
	}
	if col == nil {
		c.Logger.Error("video collection does not exist", slog.Any("collection_id", req.CollectionID))
		return nil, httperror.NewForBadRequestWithSingleField("collection_id", "does not exist")
	}
	e.CollectionID = col.ID
	e.CollectionName = col.Name

	// Get video upload or external video.
	e.VideoType = req.VideoType
	e.VideoURL = req.VideoURL

	//
	// Process video file uploads here...
	//

	if e.VideoType == vcon_s.VideoContentVideoTypeS3 {

		//
		// STEP 1: Fetch the video attachment or error if D.N.E.
		//

		aid, err := primitive.ObjectIDFromHex(req.VideoUpload)
		if err != nil {
			c.Logger.Error("video object id from hex error", slog.Any("error", err))
			return nil, err
		}
		attachment, err := c.AttachmentStorer.GetByID(ctx, aid)
		if err != nil {
			c.Logger.Error("video attachment get by id error", slog.Any("error", err))
			return nil, err
		}
		if attachment == nil {
			c.Logger.Error("video attachment does not exist", slog.Any("aid", aid))
			return nil, httperror.NewForBadRequestWithSingleField("video_upload", "does not exist")
		}

		//
		// STEP 2: Copy attachment into our videocontent.
		//

		// Generate the key of our videocontent video upload.
		videoObjectKey := fmt.Sprintf("org/%v/videocontent/%v/%v", orgID.Hex(), e.ID.Hex(), attachment.Filename)

		// Cut the object from the attachment temporary location into our new videocontent location.
		if cutErr := c.S3.Cut(ctx, attachment.ObjectKey, videoObjectKey); cutErr != nil {
			c.Logger.Error("video s3 cut error", slog.Any("error", err))
			return nil, err
		}

		// Generate a presigned URL for today.
		expiryDur := time.Hour * 12
		videoObjectURL, err := c.S3.GetPresignedURL(ctx, videoObjectKey, expiryDur)

		// Update the videocontent.
		e.VideoObjectKey = videoObjectKey
		e.VideoObjectURL = videoObjectURL
		e.VideoObjectExpiry = time.Now().Add(expiryDur)
		e.VideoAttachmentID = attachment.ID
		e.VideoAttachmentFilename = attachment.Filename

		//
		// STEP 3: Update the attachment.
		//

		attachment.ObjectKey = videoObjectKey
		attachment.ObjectURL = videoObjectURL
		// attachment.OwnershipID = e.ID
		// attachment.OwnershipType = a_d.OwnershipTypeVideoContentVideo
		if updateErr := c.AttachmentStorer.UpdateByID(ctx, attachment); updateErr != nil {
			c.Logger.Error("video attachment update error", slog.Any("updateErr", updateErr))
			return nil, updateErr
		}
	}

	e.ThumbnailType = req.ThumbnailType
	e.ThumbnailURL = req.ThumbnailURL

	//
	// Process thumbnail file uploads here...
	//

	if e.ThumbnailType == vcon_s.VideoContentThumbnailTypeS3 {

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
		if attachment == nil {
			c.Logger.Error("thumbnail attachment does not exist", slog.Any("aid", aid))
			return nil, httperror.NewForBadRequestWithSingleField("thumbnail_upload", "does not exist")
		}

		//
		// STEP 2: Copy attachment into our videocontent.
		//

		// Generate the key of our videocontent video upload.
		thumbnailObjectKey := fmt.Sprintf("org/%v/videocontent/%v/%v", orgID.Hex(), e.ID.Hex(), attachment.Filename)

		// Cut the object from the attachment temporary location into our new videocontent location.
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

		// Update the videocontent.
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
		// attachment.OwnershipID = e.ID
		// attachment.OwnershipType = a_d.OwnershipTypeVideoContentVideo
		if updateErr := c.AttachmentStorer.UpdateByID(ctx, attachment); updateErr != nil {
			c.Logger.Error("video attachment update error", slog.Any("updateErr", updateErr))
			return nil, updateErr
		}
	}

	//
	// Save to our database.
	//

	if err := c.VideoContentStorer.Create(ctx, e); err != nil {
		c.Logger.Error("database create error", slog.Any("error", err))
		return nil, err
	}

	return e, nil
}
