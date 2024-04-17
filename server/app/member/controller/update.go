package controller

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	user_s "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/user/datastore"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/config/constants"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/utils/httperror"
)

type MemberUpdateRequestIDO struct {
	ID                        primitive.ObjectID `json:"id"`
	OrganizationID            primitive.ObjectID `json:"organization_id"`
	FirstName                 string             `json:"first_name"`
	LastName                  string             `json:"last_name"`
	Email                     string             `json:"email"`
	Password                  string             `json:"password"`
	PasswordRepeated          string             `json:"password_repeated"`
	Phone                     string             `json:"phone,omitempty"`
	Country                   string             `json:"country,omitempty"`
	Region                    string             `json:"region,omitempty"`
	City                      string             `json:"city,omitempty"`
	PostalCode                string             `json:"postal_code,omitempty"`
	AddressLine1              string             `json:"address_line_1,omitempty"`
	AddressLine2              string             `json:"address_line_2,omitempty"`
	HowDidYouHearAboutUs      int8               `json:"how_did_you_hear_about_us,omitempty"`
	HowDidYouHearAboutUsOther string             `json:"how_did_you_hear_about_us_other,omitempty"`
	AgreeTOS                  bool               `json:"agree_tos,omitempty"`
	AgreePromotionsEmail      bool               `json:"agree_promotions_email,omitempty"`
	Status                    int8               `json:"status,omitempty"`
	HasShippingAddress        bool               `bson:"has_shipping_address" json:"has_shipping_address,omitempty"`
	ShippingName              string             `bson:"shipping_name" json:"shipping_name,omitempty"`
	ShippingPhone             string             `bson:"shipping_phone" json:"shipping_phone,omitempty"`
	ShippingCountry           string             `bson:"shipping_country" json:"shipping_country,omitempty"`
	ShippingRegion            string             `bson:"shipping_region" json:"shipping_region,omitempty"`
	ShippingCity              string             `bson:"shipping_city" json:"shipping_city,omitempty"`
	ShippingPostalCode        string             `bson:"shipping_postal_code" json:"shipping_postal_code,omitempty"`
	ShippingAddressLine1      string             `bson:"shipping_address_line1" json:"shipping_address_line1,omitempty"`
	ShippingAddressLine2      string             `bson:"shipping_address_line2" json:"shipping_address_line2,omitempty"`
	Tags                      []*user_s.UserTag  `bson:"tags" json:"tags"`
	OnboardingAnswers         []*user_s.Answer   `bson:"onboarding_answers" json:"onboarding_answers,omitempty"`
	OnboardingCompleted       bool               `bson:"onboarding_completed" json:"onboarding_completed"`
}

func (impl *MemberControllerImpl) userFromUpdateRequest(requestData *MemberUpdateRequestIDO) (*user_s.User, error) {
	// Make shipping and billing same if not selected.
	if requestData.HasShippingAddress == false {
		requestData.ShippingCity = requestData.City
		requestData.ShippingCountry = requestData.Country
		requestData.ShippingAddressLine1 = requestData.AddressLine1
		requestData.ShippingAddressLine2 = requestData.AddressLine2
		requestData.ShippingPostalCode = requestData.PostalCode
		requestData.ShippingRegion = requestData.Region
	}

	u := &user_s.User{
		ID:                        requestData.ID,
		OrganizationID:            requestData.OrganizationID,
		FirstName:                 requestData.FirstName,
		LastName:                  requestData.LastName,
		Email:                     requestData.Email,
		Phone:                     requestData.Phone,
		Country:                   requestData.Country,
		Region:                    requestData.Region,
		City:                      requestData.City,
		PostalCode:                requestData.PostalCode,
		AddressLine1:              requestData.AddressLine1,
		AddressLine2:              requestData.AddressLine2,
		HowDidYouHearAboutUs:      requestData.HowDidYouHearAboutUs,
		HowDidYouHearAboutUsOther: requestData.HowDidYouHearAboutUsOther,
		AgreeTOS:                  requestData.AgreeTOS,
		AgreePromotionsEmail:      requestData.AgreePromotionsEmail,
		Status:                    requestData.Status,
		HasShippingAddress:        requestData.HasShippingAddress,
		ShippingName:              requestData.ShippingName,
		ShippingPhone:             requestData.ShippingPhone,
		ShippingCountry:           requestData.ShippingCountry,
		ShippingRegion:            requestData.ShippingRegion,
		ShippingCity:              requestData.ShippingCity,
		ShippingPostalCode:        requestData.ShippingPostalCode,
		ShippingAddressLine1:      requestData.ShippingAddressLine1,
		ShippingAddressLine2:      requestData.ShippingAddressLine2,
		Tags:                      requestData.Tags,
		OnboardingAnswers:         requestData.OnboardingAnswers,
		OnboardingCompleted:       requestData.OnboardingCompleted,
	}

	// Set new password if user added new password.
	if requestData.Password != "" {
		passwordHash, err := impl.Password.GenerateHashFromPassword(requestData.Password)
		if err != nil {
			impl.Logger.Error("hashing error", slog.Any("error", err))
			return nil, err
		}
		u.PasswordHash = passwordHash
		u.PasswordHashAlgorithm = impl.Password.AlgorithmName()
	}

	return u, nil
}

func (impl *MemberControllerImpl) UpdateByID(ctx context.Context, requestData *MemberUpdateRequestIDO) (*user_s.User, error) {
	nu, err := impl.userFromUpdateRequest(requestData)
	if err != nil {
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

		// Lookup the user in our database, else return a `400 Bad Request` error.
		b, err := impl.OrganizationStorer.GetByID(sessCtx, nu.OrganizationID)
		if err != nil {
			impl.Logger.Error("database error", slog.Any("err", err))
			return nil, err
		}
		if b == nil {
			impl.Logger.Warn("organization does not exists validation error")
			return nil, httperror.NewForBadRequestWithSingleField("organization_id", "organization does not exist")
		}

		// Extract from our session the following data.
		orgID, _ := sessCtx.Value(constants.SessionUserOrganizationID).(primitive.ObjectID)
		orgName, _ := sessCtx.Value(constants.SessionUserOrganizationName).(string)
		userID, _ := sessCtx.Value(constants.SessionUserID).(primitive.ObjectID)
		userName, _ := sessCtx.Value(constants.SessionUserName).(string)
		ipAddress, _ := sessCtx.Value(constants.SessionIPAddress).(string)

		// Lookup the user in our database, else return a `400 Bad Request` error.
		ou, err := impl.UserStorer.GetByID(sessCtx, nu.ID)
		if err != nil {
			impl.Logger.Error("database error", slog.Any("err", err))
			return nil, err
		}
		if ou == nil {
			impl.Logger.Warn("user does not exist validation error")
			return nil, httperror.NewForBadRequestWithSingleField("id", "does not exist")
		}

		if ou.Email != nu.Email { // Defensive code: Only pick unique emails.
			// Lookup the user in our database, else return a `400 Bad Request` error.
			u, err := impl.UserStorer.GetByEmail(sessCtx, nu.Email)
			if err != nil {
				impl.Logger.Error("database error", slog.Any("err", err))
				return nil, err
			}
			if u != nil {
				impl.Logger.Warn("user already exists validation error")
				return nil, httperror.NewForBadRequestWithSingleField("email", "email is not unique")
			}
		}

		ou.OrganizationID = b.ID
		ou.OrganizationName = b.Name
		ou.FirstName = nu.FirstName
		ou.LastName = nu.LastName
		ou.Name = fmt.Sprintf("%s %s", nu.FirstName, nu.LastName)
		ou.LexicalName = fmt.Sprintf("%s, %s", nu.LastName, nu.FirstName)
		ou.Email = nu.Email
		ou.OrganizationID = orgID
		ou.OrganizationName = orgName
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
		ou.ModifiedAt = time.Now()
		ou.ModifiedByUserID = userID
		ou.ModifiedByUserName = userName
		ou.ModifiedFromIPAddress = ipAddress
		ou.HasShippingAddress = nu.HasShippingAddress
		ou.ShippingName = nu.ShippingName
		ou.ShippingPhone = nu.ShippingPhone
		ou.ShippingCountry = nu.ShippingCountry
		ou.ShippingRegion = nu.ShippingRegion
		ou.ShippingCity = nu.ShippingCity
		ou.ShippingPostalCode = nu.ShippingPostalCode
		ou.ShippingAddressLine1 = nu.ShippingAddressLine1
		ou.ShippingAddressLine2 = nu.ShippingAddressLine2
		ou.OnboardingAnswers = nu.OnboardingAnswers
		ou.OnboardingCompleted = nu.OnboardingCompleted

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

		if err := impl.UserStorer.UpdateByID(sessCtx, ou); err != nil {
			impl.Logger.Error("user update by id error", slog.Any("error", err))
			return nil, err
		}

		// Defensive Code: In case user does not have an account with the payment processor.
		if ou.PaymentProcessorCustomerID != "" {
			err = impl.PaymentProcessor.UpdateCustomer(
				ou.PaymentProcessorCustomerID,
				fmt.Sprintf("%s %s", ou.FirstName, ou.LastName),
				ou.Email,
				"", // description...
				fmt.Sprintf("%s %s Shipping Address", ou.FirstName, ou.LastName),
				ou.Phone,
				ou.ShippingCity, ou.ShippingCountry, ou.ShippingAddressLine1, ou.ShippingAddressLine2, ou.ShippingPostalCode, ou.ShippingRegion, // Shipping
				ou.City, ou.Country, ou.AddressLine1, ou.AddressLine2, ou.PostalCode, ou.Region, // Billing
			)
			if err != nil {
				impl.Logger.Error("update customer in payment processor error", slog.Any("error", err))
				return nil, err
			}
		}

		return ou, nil
	}

	// Start a transaction
	result, err := session.WithTransaction(ctx, transactionFunc)
	if err != nil {
		impl.Logger.Error("session failed error",
			slog.Any("error", err))
		return nil, err
	}

	return result.(*user_s.User), nil
}
