package controller

import (
	"context"
	"fmt"
	"log/slog"
	"mime/multipart"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	rp_s "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/rankpoint/datastore"
	u_s "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/user/datastore"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/config/constants"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/utils/httperror"
)

type AccountOperationAvatarRequest struct {
	AccountID primitive.ObjectID `bson:"account_id" json:"account_id"`
	FileName  string
	FileType  string
	File      multipart.File
}

func (impl *GatewayControllerImpl) validateOperationAvatarRequest(ctx context.Context, dirtyData *AccountOperationAvatarRequest) error {
	e := make(map[string]string)

	if dirtyData.AccountID.IsZero() {
		e["account_id"] = "missing value"
	}
	if dirtyData.File == nil {
		e["file"] = "missing value"
	}

	if len(e) != 0 {
		return httperror.NewForBadRequest(&e)
	}
	return nil
}

func (impl *GatewayControllerImpl) Avatar(ctx context.Context, req *AccountOperationAvatarRequest) (*u_s.User, error) {

	// Get variables from our user authenticated session.
	//

	orgID := ctx.Value(constants.SessionUserOrganizationID).(primitive.ObjectID)
	// role, _ := ctx.Value(constants.SessionUserRole).(int8)
	userID, _ := ctx.Value(constants.SessionUserID).(primitive.ObjectID)
	userName, _ := ctx.Value(constants.SessionUserName).(string)
	ipAddress, _ := ctx.Value(constants.SessionIPAddress).(string)

	//
	// Perform our validation and return validation error on any issues detected.
	//

	if err := impl.validateOperationAvatarRequest(ctx, req); err != nil {
		impl.Logger.Error("validation error", slog.Any("error", err))
		return nil, err
	}

	////
	//// Start the transaction.
	////

	session, err := impl.DbClient.StartSession()
	if err != nil {
		impl.Logger.Error("start session error",
			slog.Any("error", err))
		return nil, err
	}
	defer session.EndSession(ctx)

	// Define a transaction function with a series of operations
	transactionFunc := func(sessCtx mongo.SessionContext) (interface{}, error) {

		//
		// Fetch the original account.
		//

		s, err := impl.UserStorer.GetByID(sessCtx, req.AccountID)
		if err != nil {
			impl.Logger.Error("database get by id error", slog.Any("error", err))
			return nil, err
		}
		if s == nil {
			return nil, nil
		}
		if s.OrganizationID != orgID {
			return nil, httperror.NewForForbiddenWithSingleField("security", "you do not belong to this organization")
		}

		// Update the file if the user uploaded a new file.
		if req.File != nil {
			// Proceed to delete the physical files from AWS s3.
			if err := impl.S3.DeleteByKeys(sessCtx, []string{s.AvatarObjectKey}); err != nil {
				impl.Logger.Warn("s3 delete by keys error", slog.Any("error", err))
				// Do not return an error, simply continue this function as there might
				// be a case were the file was removed on the s3 bucket by ourselves
				// or some other reason.
			}

			// Generate the key of our upload.
			objectKey := fmt.Sprintf("org/%v/uploads/user/%v/avatar/%v", orgID.Hex(), req.AccountID.Hex(), req.FileName)

			// For debugging purposes only.
			impl.Logger.Debug("pre-upload meta",
				slog.String("AvatarFileName", req.FileName),
				slog.String("AvatarFileType", req.FileType),
				slog.String("AvatarObjectKey", objectKey),
			)

			go func(file multipart.File, objkey string) {
				impl.Logger.Debug("beginning private s3 image upload...")
				if err := impl.S3.UploadContentFromMulipart(context.Background(), objkey, file); err != nil {
					impl.Logger.Error("private s3 upload error", slog.Any("error", err))
					// Do not return an error, simply continue this function as there might
					// be a case were the file was removed on the s3 bucket by ourselves
					// or some other reason.
				}
				impl.Logger.Debug("Finished private s3 image upload")
			}(req.File, objectKey)

			// Update file.
			s.AvatarObjectKey = objectKey
			s.AvatarFileName = req.FileName
			s.AvatarFileType = req.FileType

			// Get preseigned URL.
			avatarObjectExpiry := time.Now().Add(time.Minute * 30)
			avatarFileURL, err := impl.S3.GetPresignedURL(sessCtx, objectKey, time.Minute*30)
			if err != nil {
				impl.Logger.Error("s3 failed get presigned url error", slog.Any("error", err))
				return nil, err
			}
			s.AvatarObjectURL = avatarFileURL
			s.AvatarObjectExpiry = avatarObjectExpiry

			// Update metainformation.
			s.ModifiedByUserID = userID
			s.ModifiedAt = time.Now()
			s.ModifiedByUserName = userName
			s.ModifiedFromIPAddress = ipAddress
		}

		// Save to the database the modified account.
		if err := impl.UserStorer.UpdateByID(sessCtx, s); err != nil {
			impl.Logger.Error("database update by id error", slog.Any("error", err))
			return nil, err
		}

		// UPDATE RELATED
		rpf := &rp_s.RankPointPaginationListFilter{
			Cursor:    "",
			PageSize:  1_000_000_000,
			SortField: "", // No need for sorting here b/c of repeated keys
			SortOrder: 1,
			UserID:    userID,
		}
		rplist, err := impl.RankPointStorer.ListByFilter(sessCtx, rpf)
		if err != nil {
			impl.Logger.Error("database get by id error", slog.Any("error", err))
			return nil, err
		}
		for _, rp := range rplist.Results {
			rp.UserAvatarObjectExpiry = s.AvatarObjectExpiry
			rp.UserAvatarObjectURL = s.AvatarObjectURL
			rp.UserAvatarObjectKey = s.AvatarObjectKey
			rp.UserAvatarFileType = s.AvatarFileType
			rp.UserAvatarFileName = s.AvatarFileName
			if err := impl.RankPointStorer.UpdateByID(sessCtx, rp); err != nil {
				impl.Logger.Error("database update by id error", slog.Any("error", err))
				return nil, err
			}
		}

		return s, nil
	}

	// Start a transaction
	result, err := session.WithTransaction(ctx, transactionFunc)
	if err != nil {
		impl.Logger.Error("session failed error",
			slog.Any("error", err))
		return nil, err
	}

	return result.(*u_s.User), nil
}
