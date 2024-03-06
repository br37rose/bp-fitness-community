package controller

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"log/slog"

	a_d "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/attachment/datastore"
	s_d "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/exercise/datastore"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/config/constants"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/utils/httperror"
)

type ExerciseCreateRequestIDO struct {
	Type int8 `bson:"type" json:"type"`

	VideoType   int8   `bson:"video_type" json:"video_type"`
	VideoUpload string `bson:"video_upload" json:"video_upload"`
	VideoURL    string `bson:"video_url" json:"video_url"`

	ThumbnailType   int8   `bson:"thumbnail_type" json:"thumbnail_type"`
	ThumbnailUpload string `bson:"thumbnail_upload" json:"thumbnail_upload"`
	ThumbnailURL    string `bson:"thumbnail_url" json:"thumbnail_url"`

	Name          string `bson:"name" json:"name"`
	AlternateName string `bson:"alternate_name" json:"alternate_name"`
	Gender        string `bson:"gender" json:"gender"`
	MovementType  int8   `bson:"movement_type" json:"movement_type"`
	Category      int8   `bson:"category" json:"category"`
	Description   string `bson:"description" json:"description"`

	HasMonetization bool               `bson:"has_monetization" json:"has_monetization"`
	OfferID         primitive.ObjectID `bson:"offer_id" json:"offer_id"`
	HasTimedLock    bool               `bson:"has_timed_lock" json:"has_timed_lock"`
	TimedLock       string             `bson:"timed_lock" json:"timed_lock"`
}

func ValidateCreateRequest(dirtyData *ExerciseCreateRequestIDO) error {
	e := make(map[string]string)

	if dirtyData.Type == 0 {
		e["type"] = "missing value"
	}

	// Video upload
	if dirtyData.VideoType == 0 {
		e["video_type"] = "missing value"
	} else {
		if dirtyData.VideoType == s_d.ExerciseVideoTypeS3 && dirtyData.VideoUpload == "" {
			e["video_upload"] = "missing value"
		}
		if dirtyData.VideoType == s_d.ExerciseVideoTypeYouTube && dirtyData.VideoURL == "" {
			e["video_url"] = "missing value"
		}
		if dirtyData.VideoType == s_d.ExerciseVideoTypeVimeo && dirtyData.VideoURL == "" {
			e["video_url"] = "missing value"
		}
	}

	// Thumbnail upload
	if dirtyData.ThumbnailType == 0 {
		e["thumbnail_type"] = "missing value"
	} else {
		if dirtyData.ThumbnailType == s_d.ExerciseThumbnailTypeS3 && dirtyData.ThumbnailUpload == "" {
			e["thumbnail_upload"] = "missing value"
		}
		if dirtyData.ThumbnailType == s_d.ExerciseThumbnailTypeExternalURL && dirtyData.ThumbnailURL == "" {
			e["thumbnail_url"] = "missing value"
		}
	}

	if dirtyData.Name == "" {
		e["name"] = "missing value"
	}
	if dirtyData.AlternateName == "" {
		e["alternate_name"] = "missing value"
	}
	if dirtyData.Gender == "" {
		e["gender"] = "missing value"
	}
	if dirtyData.MovementType == 0 {
		e["movement_type"] = "missing value"
	}
	if dirtyData.Category == 0 {
		e["category"] = "missing value"
	}
	if dirtyData.Description == "" {
		e["description"] = "missing value"
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

func (c *ExerciseControllerImpl) Create(ctx context.Context, req *ExerciseCreateRequestIDO) (*s_d.Exercise, error) {
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

	//
	// Initialize
	//

	e := &s_d.Exercise{}

	e.Status = s_d.ExerciseStatusActive
	e.OrganizationID = orgID
	e.ID = primitive.NewObjectID()
	e.CreatedByUserName = userName
	e.CreatedByUserID = userID
	e.CreatedAt = time.Now()
	e.ModifiedByUserName = userName
	e.ModifiedByUserID = userID
	e.ModifiedAt = time.Now()

	e.Type = req.Type
	e.Name = req.Name
	e.AlternateName = req.AlternateName
	e.Gender = req.Gender
	e.MovementType = req.MovementType
	e.Category = req.Category
	e.Description = req.Description

	e.VideoType = req.VideoType
	e.VideoURL = req.VideoURL

	//
	// Fetch offer
	//

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

	//
	// Process video file uploads here...
	//

	if e.VideoType == s_d.ExerciseVideoTypeS3 {

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
		// STEP 2: Copy attachment into our exercise.
		//

		// Generate the key of our exercise video upload.
		videoObjectKey := fmt.Sprintf("org/%v/exercise/%v/%v", orgID.Hex(), e.ID.Hex(), attachment.Filename)

		// Cut the object from the attachment temporary location into our new exercise location.
		if cutErr := c.S3.Cut(ctx, attachment.ObjectKey, videoObjectKey); cutErr != nil {
			c.Logger.Error("video s3 cut error", slog.Any("error", err))
			return nil, err
		}

		// Generate a presigned URL for today.
		expiryDur := time.Hour * 12
		videoObjectURL, err := c.S3.GetPresignedURL(ctx, videoObjectKey, expiryDur)

		// Update the exercise.
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
		attachment.OwnershipID = e.ID
		attachment.OwnershipType = a_d.OwnershipTypeExerciseVideo
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

	if e.ThumbnailType == s_d.ExerciseThumbnailTypeS3 {

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
		// STEP 2: Copy attachment into our exercise.
		//

		// Generate the key of our exercise video upload.
		thumbnailObjectKey := fmt.Sprintf("org/%v/exercise/%v/%v", orgID.Hex(), e.ID.Hex(), attachment.Filename)

		// Cut the object from the attachment temporary location into our new exercise location.
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

		// Update the exercise.
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
		attachment.OwnershipID = e.ID
		attachment.OwnershipType = a_d.OwnershipTypeExerciseVideo
		if updateErr := c.AttachmentStorer.UpdateByID(ctx, attachment); updateErr != nil {
			c.Logger.Error("video attachment update error", slog.Any("updateErr", updateErr))
			return nil, updateErr
		}
	}

	//
	// Save to our database.
	//

	if err := c.ExerciseStorer.Create(ctx, e); err != nil {
		c.Logger.Error("database create error", slog.Any("error", err))
		return nil, err
	}

	return e, nil
}
