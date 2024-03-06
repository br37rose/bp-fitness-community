package controller

import (
	"context"
	"fmt"
	"time"

	a_d "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/attachment/datastore"
	domain "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/exercise/datastore"
	s_d "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/exercise/datastore"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/config/constants"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/utils/httperror"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log/slog"
)

type ExerciseUpdateRequestIDO struct {
	ID primitive.ObjectID `bson:"id" json:"id"`

	Type int8 `bson:"type" json:"type"`

	VideoType   int8   `bson:"video_type" json:"video_type"`
	VideoUpload string `bson:"video_upload" json:"video_upload"`
	VideoURL    string `bson:"video_url" json:"video_url"`

	ThumbnailType   int8   `bson:"thumbnail_type" json:"thumbnail_type"`
	ThumbnailUpload string `bson:"thumbnail_upload" json:"thumbnail_upload"`
	ThumbnailURL    string `bson:"thumbnail_url" json:"thumbnail_url"`

	Name            string             `bson:"name" json:"name"`
	AlternateName   string             `bson:"alternate_name" json:"alternate_name"`
	Gender          string             `bson:"gender" json:"gender"`
	MovementType    int8               `bson:"movement_type" json:"movement_type"`
	Category        int8               `bson:"category" json:"category"`
	Description     string             `bson:"description" json:"description"`
	Status          int8               `bson:"status" json:"status"`
	HasMonetization bool               `bson:"has_monetization" json:"has_monetization"`
	OfferID         primitive.ObjectID `bson:"offer_id" json:"offer_id"`
	HasTimedLock    bool               `bson:"has_timed_lock" json:"has_timed_lock"`
	TimedLock       string             `bson:"timed_lock" json:"timed_lock"`
}

func ValidateUpdateRequest(dirtyData *ExerciseUpdateRequestIDO) error {
	e := make(map[string]string)

	if dirtyData.ID.IsZero() {
		e["id"] = "missing value"
	}

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
	if dirtyData.Status == 0 {
		e["status"] = "missing value"
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

func (c *ExerciseControllerImpl) UpdateByID(ctx context.Context, req *ExerciseUpdateRequestIDO) (*domain.Exercise, error) {
	if err := ValidateUpdateRequest(req); err != nil {
		c.Logger.Error("validation failure", slog.Any("error", err))
		return nil, err
	}

	// Convert the string to time.Duration
	timedLockDuration, err := time.ParseDuration(req.TimedLock)
	if err != nil && req.TimedLock != "" {
		c.Logger.Error("parse duration err", slog.Any("TimedLock", req.TimedLock))
		return nil, err
	}

	// Lock this database record when we are updating it so in case we don't cause any data inconsistency when system refreshes the presigned URL in the detail and list functions.
	c.Kmutex.Lockf("exercise-%v", req.ID.Hex()) // Step 1
	defer func() {
		c.Kmutex.Unlockf("exercise-%v", req.ID.Hex()) // Step 2
	}()

	// Extract from our session the following data.
	orgID, _ := ctx.Value(constants.SessionUserOrganizationID).(primitive.ObjectID)
	userID, _ := ctx.Value(constants.SessionUserID).(primitive.ObjectID)
	userName, _ := ctx.Value(constants.SessionUserName).(string)

	// Fetch the original exercise.
	os, err := c.ExerciseStorer.GetByID(ctx, req.ID)
	if err != nil {
		c.Logger.Error("database get by id error", slog.Any("error", err))
		return nil, err
	}
	if os == nil {
		return nil, httperror.NewForBadRequestWithSingleField("id", "exercise does not exist")
	}

	// Defensive code: Do not edit system.
	if os.Type == domain.ExerciseTypeSystem {
		c.Logger.Error("system exercise edit error")
		return nil, httperror.NewForBadRequestWithSingleField("message", "the `system` exercises are read-only and thus access is denied")
	}

	// Modify our original exercise.
	os.ModifiedByUserName = userName
	os.ModifiedByUserID = userID
	os.ModifiedAt = time.Now()
	os.OrganizationID = orgID
	os.Status = req.Status

	os.Type = req.Type
	os.Name = req.Name
	os.AlternateName = req.AlternateName
	os.Gender = req.Gender
	os.MovementType = req.MovementType
	os.Category = req.Category
	os.Description = req.Description

	os.VideoType = req.VideoType
	os.VideoURL = req.VideoURL

	os.HasMonetization = req.HasMonetization
	os.OfferID = req.OfferID
	os.HasTimedLock = req.HasTimedLock
	os.TimedLock = req.TimedLock
	os.TimedLockDuration = timedLockDuration

	//
	// Process video file uploads here...
	//

	if req.VideoType == s_d.ExerciseVideoTypeS3 {

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
		videoObjectKey := fmt.Sprintf("org/%v/exercise/%v/%v", orgID.Hex(), req.ID.Hex(), attachment.Filename)

		// Cut the object from the attachment temporary location into our new exercise location.
		if cutErr := c.S3.Cut(ctx, attachment.ObjectKey, videoObjectKey); cutErr != nil {
			c.Logger.Error("video s3 cut error", slog.Any("error", err))
			return nil, err
		}

		// Generate a presigned URL for today.
		expiryDur := time.Hour * 12
		videoObjectURL, err := c.S3.GetPresignedURL(ctx, videoObjectKey, expiryDur)

		// Update the exercise.
		os.VideoObjectKey = videoObjectKey
		os.VideoObjectURL = videoObjectURL
		os.VideoObjectExpiry = time.Now().Add(expiryDur)
		os.VideoAttachmentID = attachment.ID
		os.VideoAttachmentFilename = attachment.Filename

		//
		// STEP 3: Update the attachment.
		//

		attachment.ObjectKey = videoObjectKey
		attachment.ObjectURL = videoObjectURL
		attachment.OwnershipID = os.ID
		attachment.OwnershipType = a_d.OwnershipTypeExerciseVideo
		if updateErr := c.AttachmentStorer.UpdateByID(ctx, attachment); updateErr != nil {
			c.Logger.Error("video attachment update error", slog.Any("updateErr", updateErr))
			return nil, updateErr
		}
	}

	os.ThumbnailType = req.ThumbnailType
	os.ThumbnailURL = req.ThumbnailURL

	//
	// Process thumbnail file uploads here...
	//

	if os.ThumbnailType == s_d.ExerciseThumbnailTypeS3 {

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
		thumbnailObjectKey := fmt.Sprintf("org/%v/exercise/%v/%v", orgID.Hex(), os.ID.Hex(), attachment.Filename)

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
		attachment.OwnershipID = os.ID
		attachment.OwnershipType = a_d.OwnershipTypeExerciseVideo
		if updateErr := c.AttachmentStorer.UpdateByID(ctx, attachment); updateErr != nil {
			c.Logger.Error("video attachment update error", slog.Any("updateErr", updateErr))
			return nil, updateErr
		}
	}

	//
	// Save to the database the modified exercise.
	//

	if err := c.ExerciseStorer.UpdateByID(ctx, os); err != nil {
		c.Logger.Error("database update by id error", slog.Any("error", err))
		return nil, err
	}

	return os, nil
}
