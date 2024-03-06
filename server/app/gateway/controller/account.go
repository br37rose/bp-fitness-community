package controller

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	user_s "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/user/datastore"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/config/constants"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/utils/httperror"
)

func (impl *GatewayControllerImpl) Account(ctx context.Context) (*user_s.User, error) {
	// Extract from our session the following data.
	userID := ctx.Value(constants.SessionUserID).(primitive.ObjectID)

	// Lookup the user in our database, else return a `400 Bad Request` error.
	u, err := impl.UserStorer.GetByID(ctx, userID)
	if err != nil {
		impl.Logger.Error("database error", slog.Any("err", err))
		return nil, err
	}
	if u == nil {
		impl.Logger.Warn("user does not exist validation error")
		return nil, httperror.NewForBadRequestWithSingleField("id", "does not exist")
	}

	// Generate most recent avatar URL if it exists and it has expired.
	if u.AvatarObjectKey != "" {
		if time.Now().After(u.AvatarObjectExpiry) {
			// Get preseigned URL.
			avatarObjectExpiry := time.Now().Add(time.Minute * 30)
			avatarFileURL, err := impl.S3.GetPresignedURL(ctx, u.AvatarObjectKey, time.Minute*30)
			if err != nil {
				impl.Logger.Error("s3 failed get presigned url error", slog.Any("error", err))
				return nil, err
			}
			u.AvatarObjectURL = avatarFileURL
			u.AvatarObjectExpiry = avatarObjectExpiry

			// Save to the database the modified associate.
			if err := impl.UserStorer.UpdateByID(ctx, u); err != nil {
				impl.Logger.Error("database update by id error", slog.Any("error", err))
				return nil, err
			}

			// For debugging purposes only.
			impl.Logger.Debug("refreshed avatar object url", slog.String("avatarFileURL", avatarFileURL))
		}
	}

	return u, nil
}

func (impl *GatewayControllerImpl) AccountUpdate(ctx context.Context, nu *user_s.User) error {
	// Extract from our session the following data.
	userID := ctx.Value(constants.SessionUserID).(primitive.ObjectID)
	orgID := ctx.Value(constants.SessionUserOrganizationID).(primitive.ObjectID)

	// Lookup the user in our database, else return a `400 Bad Request` error.
	ou, err := impl.UserStorer.GetByID(ctx, userID)
	if err != nil {
		impl.Logger.Error("database error", slog.Any("err", err))
		return err
	}
	if ou == nil {
		impl.Logger.Warn("user does not exist validation error")
		return httperror.NewForBadRequestWithSingleField("id", "does not exist")
	}

	ou.FirstName = nu.FirstName
	ou.LastName = nu.LastName
	ou.Name = fmt.Sprintf("%s %s", nu.FirstName, nu.LastName)
	ou.LexicalName = fmt.Sprintf("%s, %s", nu.LastName, nu.FirstName)
	ou.Email = nu.Email
	ou.Phone = nu.Phone
	ou.Country = nu.Country
	ou.Region = nu.Region
	ou.City = nu.City
	ou.PostalCode = nu.PostalCode
	ou.AddressLine1 = nu.AddressLine1
	ou.AddressLine2 = nu.AddressLine2
	ou.HowDidYouHearAboutUs = nu.HowDidYouHearAboutUs
	ou.HowDidYouHearAboutUsOther = nu.HowDidYouHearAboutUsOther
	ou.AgreePromotionsEmail = nu.AgreePromotionsEmail

	// Process user tags.
	var modifiedTags []*user_s.UserTag
	for _, tag := range nu.Tags {
		// If no `id` exists then this tag has been recently created so let us
		// finish initializing it by adding our meta information.
		if tag.ID.IsZero() {
			tag.ID = primitive.NewObjectID()
			tag.OrganizationID = orgID
			tag.UserID = userID
			tag.Status = user_s.UserStatusActive
		}
		modifiedTags = append(modifiedTags, tag)
	}
	ou.Tags = modifiedTags

	if err := impl.UserStorer.UpdateByID(ctx, ou); err != nil {
		impl.Logger.Error("user update by id error", slog.Any("error", err))
		return err
	}
	return nil
}

type AccountChangePasswordRequestIDO struct {
	OldPassword      string `json:"old_password"`
	Password         string `json:"password"`
	PasswordRepeated string `json:"password_repeated"`
}

func ValidateAccountChangePassworRequest(dirtyData *AccountChangePasswordRequestIDO) error {
	e := make(map[string]string)

	if dirtyData.OldPassword == "" {
		e["old_password"] = "missing value"
	}
	if dirtyData.Password == "" {
		e["password"] = "missing value"
	}
	if dirtyData.PasswordRepeated == "" {
		e["password_repeated"] = "missing value"
	}
	if dirtyData.PasswordRepeated != dirtyData.Password {
		e["password"] = "does not match"
		e["password_repeated"] = "does not match"
	}

	if len(e) != 0 {
		return httperror.NewForBadRequest(&e)
	}
	return nil
}

func (impl *GatewayControllerImpl) AccountChangePassword(ctx context.Context, req *AccountChangePasswordRequestIDO) error {
	// Extract from our session the following data.
	userID := ctx.Value(constants.SessionUserID).(primitive.ObjectID)

	// Lookup the user in our database, else return a `400 Bad Request` error.
	u, err := impl.UserStorer.GetByID(ctx, userID)
	if err != nil {
		impl.Logger.Error("database error", slog.Any("err", err))
		return err
	}
	if u == nil {
		impl.Logger.Warn("user does not exist validation error")
		return httperror.NewForBadRequestWithSingleField("id", "does not exist")
	}

	if err := ValidateAccountChangePassworRequest(req); err != nil {
		impl.Logger.Warn("user validation failed", slog.Any("err", err))
		return err
	}

	// Verify the inputted password and hashed password match.
	if passwordMatch, _ := impl.Password.ComparePasswordAndHash(req.OldPassword, u.PasswordHash); passwordMatch == false {
		impl.Logger.Warn("password check validation error")
		return httperror.NewForBadRequestWithSingleField("old_password", "old password do not match with record of existing password")
	}

	passwordHash, err := impl.Password.GenerateHashFromPassword(req.Password)
	if err != nil {
		impl.Logger.Error("hashing error", slog.Any("error", err))
		return err
	}
	u.PasswordHash = passwordHash
	u.PasswordHashAlgorithm = impl.Password.AlgorithmName()
	if err := impl.UserStorer.UpdateByID(ctx, u); err != nil {
		impl.Logger.Error("user update by id error", slog.Any("error", err))
		return err
	}
	return nil
}

func (impl *GatewayControllerImpl) AccountListLatestInvoices(ctx context.Context, cursor int64, limit int64) (*user_s.StripeInvoiceListResult, error) {
	// Extract from our session the following data.
	userID := ctx.Value(constants.SessionUserID).(primitive.ObjectID)

	// Lookup the user in our database, else return a `400 Bad Request` error.
	u, err := impl.UserStorer.GetByID(ctx, userID)
	if err != nil {
		impl.Logger.Error("database error", slog.Any("err", err))
		return nil, err
	}
	if u == nil {
		impl.Logger.Warn("user does not exist validation error")
		return nil, httperror.NewForBadRequestWithSingleField("id", "does not exist")
	}

	return impl.UserStorer.ListLatestStripeInvoices(ctx, userID, cursor, limit)
}
