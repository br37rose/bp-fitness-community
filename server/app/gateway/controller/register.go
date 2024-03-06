package controller

import (
	"context"
	"fmt"
	"strings"
	"time"

	"log/slog"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	gateway_s "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/gateway/datastore"
	user_s "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/user/datastore"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/config/constants"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/utils/httperror"
)

type MemberRegisterRequestIDO struct {
	OrganizationID            primitive.ObjectID `json:"organization_id"`
	FirstName                 string             `json:"first_name"`
	LastName                  string             `json:"last_name"`
	Email                     string             `json:"email"`
	EmailRepeated             string             `json:"email_repeated"`
	Password                  string             `json:"password"`
	PasswordRepeated          string             `json:"password_repeated"`
	Phone                     string             `json:"phone,omitempty"`
	Country                   string             `json:"country,omitempty"`
	Region                    string             `json:"region,omitempty"`
	City                      string             `json:"city,omitempty"`
	PostalCode                string             `json:"postal_code,omitempty"`
	AddressLine1              string             `json:"address_line_1,omitempty"`
	AddressLine2              string             `json:"address_line_2,omitempty"`
	StoreLogo                 string             `json:"store_logo,omitempty"`
	HowDidYouHearAboutUs      int8               `json:"how_did_you_hear_about_us,omitempty"`
	HowDidYouHearAboutUsOther string             `json:"how_did_you_hear_about_us_other,omitempty"`
	AgreeTOS                  bool               `json:"agree_tos,omitempty"`
	AgreePromotionsEmail      bool               `json:"agree_promotions_email,omitempty"`
	HasShippingAddress        bool               `bson:"has_shipping_address" json:"has_shipping_address,omitempty"`
	ShippingName              string             `bson:"shipping_name" json:"shipping_name,omitempty"`
	ShippingPhone             string             `bson:"shipping_phone" json:"shipping_phone,omitempty"`
	ShippingCountry           string             `bson:"shipping_country" json:"shipping_country,omitempty"`
	ShippingRegion            string             `bson:"shipping_region" json:"shipping_region,omitempty"`
	ShippingCity              string             `bson:"shipping_city" json:"shipping_city,omitempty"`
	ShippingPostalCode        string             `bson:"shipping_postal_code" json:"shipping_postal_code,omitempty"`
	ShippingAddressLine1      string             `bson:"shipping_address_line1" json:"shipping_address_line1,omitempty"`
	ShippingAddressLine2      string             `bson:"shipping_address_line2" json:"shipping_address_line2,omitempty"`
}

type MemberRegisterResponseIDO struct {
	User                   *user_s.User `json:"user"`
	AccessToken            string       `json:"access_token"`
	AccessTokenExpiryTime  time.Time    `json:"access_token_expiry_time"`
	RefreshToken           string       `json:"refresh_token"`
	RefreshTokenExpiryTime time.Time    `json:"refresh_token_expiry_time"`
}

func (impl *GatewayControllerImpl) MemberRegister(ctx context.Context, req *MemberRegisterRequestIDO) (*gateway_s.LoginResponseIDO, error) {
	// Defensive Code: For security purposes we need to remove all whitespaces from the email and lower the characters.
	req.Email = strings.ToLower(req.Email)
	req.Password = strings.ReplaceAll(req.Password, " ", "")

	impl.Kmutex.Lockf("REGISTRATION-EMAIL-%v", req.Email)
	defer impl.Kmutex.Unlockf("REGISTRATION-EMAIL-%v", req.Email)

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
		u, err := impl.UserStorer.GetByEmail(sessCtx, req.Email)
		if err != nil {
			impl.Logger.Error("database error", slog.Any("err", err))
			return nil, err
		}
		if u != nil {
			impl.Logger.Warn("user already exists validation error",
				slog.String("Email", req.Email))
			return nil, httperror.NewForBadRequestWithSingleField("email", "email is not unique")
		}

		// Create our user.
		u, err = impl.createUserForRequest(sessCtx, req)
		if err != nil {
			return nil, err
		}

		// Send our verification email.
		if err := impl.TemplatedEmailer.SendMemberVerificationEmail(u.Email, u.EmailVerificationCode, u.FirstName, u.OrganizationName); err != nil {
			impl.Logger.Error("failed sending verification email with error from registration",
				slog.Any("err", err),
				slog.String("Email", u.Email),
				slog.Any("UserID", u.ID))
			// Do not send error message to user nor abort the registration process.
			// Just simply log an error message and continue.
		}

		return u, nil
	}

	// Start a transaction
	if _, err := session.WithTransaction(ctx, transactionFunc); err != nil {
		impl.Logger.Error("session failed error",
			slog.Any("error", err))
		return nil, err
	}

	return impl.Login(ctx, req.Email, req.Password)
}

func (impl *GatewayControllerImpl) createUserForRequest(sessCtx mongo.SessionContext, req *MemberRegisterRequestIDO) (*user_s.User, error) {
	// Lookup the user in our database, else return a `400 Bad Request` error.
	org, err := impl.OrganizationStorer.GetByID(sessCtx, req.OrganizationID)
	if err != nil {
		impl.Logger.Error("database error", slog.Any("err", err))
		return nil, err
	}
	if org == nil {
		impl.Logger.Warn("organization does not exists validation error")
		return nil, httperror.NewForBadRequestWithSingleField("organization_id", "organization does not exist")
	}

	// Hash the password for security purposes.
	passwordHash, err := impl.Password.GenerateHashFromPassword(req.Password)
	if err != nil {
		impl.Logger.Error("hashing error", slog.Any("error", err))
		return nil, err
	}

	// //TODO: UNCOMMENT THIS CODE WHEN PROGRAMMER IS READY TO IMPLEMENT PAYMENT STUFF.
	// paymentProcessorCustomerID, err := impl.PaymentProcessor.CreateCustomer(
	// 	fmt.Sprintf("%s %s", req.FirstName, req.LastName),
	// 	req.Email,
	// 	"", // description...
	// 	fmt.Sprintf("%s %s Shipping Address", req.FirstName, req.LastName),
	// 	req.Phone,
	// 	req.ShippingCity, req.ShippingCountry, req.ShippingAddressLine1, req.ShippingAddressLine2, req.ShippingPostalCode, req.ShippingRegion, // Shipping
	// 	req.City, req.Country, req.AddressLine1, req.AddressLine2, req.PostalCode, req.Region, // Billing
	// )
	// if err != nil {
	// 	impl.Logger.Error("creating customer from payment processor error", slog.Any("error", err))
	// 	return nil, err
	// }

	// Generate the unique identifier used for MongoDB.
	userID := primitive.NewObjectID()

	ipAddress, _ := sessCtx.Value(constants.SessionIPAddress).(string)

	// Create our user record in MongoDB.
	u := &user_s.User{
		OrganizationID:            org.ID,
		OrganizationName:          org.Name,
		ID:                        userID,
		FirstName:                 req.FirstName,
		LastName:                  req.LastName,
		Name:                      fmt.Sprintf("%s %s", req.FirstName, req.LastName),
		LexicalName:               fmt.Sprintf("%s, %s", req.LastName, req.FirstName),
		Email:                     req.Email,
		PasswordHash:              passwordHash,
		PasswordHashAlgorithm:     impl.Password.AlgorithmName(),
		Role:                      user_s.UserRoleMember,
		Phone:                     req.Phone,
		Country:                   req.Country,
		Region:                    req.Region,
		City:                      req.City,
		PostalCode:                req.PostalCode,
		AddressLine1:              req.AddressLine1,
		HowDidYouHearAboutUs:      req.HowDidYouHearAboutUs,
		HowDidYouHearAboutUsOther: req.HowDidYouHearAboutUsOther,
		AgreeTOS:                  req.AgreeTOS,
		AgreePromotionsEmail:      req.AgreePromotionsEmail,
		CreatedByUserID:           userID,
		CreatedAt:                 time.Now(),
		CreatedByUserName:         fmt.Sprintf("%s %s", req.FirstName, req.LastName),
		CreatedFromIPAddress:      ipAddress,
		ModifiedByUserID:          userID,
		ModifiedAt:                time.Now(),
		ModifiedByUserName:        fmt.Sprintf("%s %s", req.FirstName, req.LastName),
		ModifiedFromIPAddress:     ipAddress,
		WasEmailVerified:          false,
		EmailVerificationCode:     impl.UUID.NewUUID(),
		EmailVerificationExpiry:   time.Now().Add(72 * time.Hour),
		Status:                    user_s.UserStatusActive,
		HasShippingAddress:        req.HasShippingAddress,
		ShippingName:              req.ShippingName,
		ShippingPhone:             req.ShippingPhone,
		ShippingCountry:           req.ShippingCountry,
		ShippingRegion:            req.ShippingRegion,
		ShippingCity:              req.ShippingCity,
		ShippingPostalCode:        req.ShippingPostalCode,
		ShippingAddressLine1:      req.ShippingAddressLine1,
		ShippingAddressLine2:      req.ShippingAddressLine2,
		Tags:                      make([]*user_s.UserTag, 0),
		Comments:                  make([]*user_s.UserComment, 0),
		StripeInvoices:            make([]*user_s.StripeInvoice, 0),
		PaymentProcessorName:      impl.PaymentProcessor.GetName(),
		// PaymentProcessorCustomerID: *paymentProcessorCustomerID,  //TODO: UNCOMMENT THIS CODE WHEN PROGRAMMER IS READY TO IMPLEMENT PAYMENT STUFF.
		OTPEnabled:   impl.Config.AppServer.Enable2FAOnRegistration,
		OTPVerified:  false,
		OTPValidated: false,
		OTPSecret:    "",
		OTPAuthURL:   "",
	}
	err = impl.UserStorer.Create(sessCtx, u)
	if err != nil {
		impl.Logger.Error("database create error", slog.Any("error", err))
		return nil, err
	}
	impl.Logger.Info("Member created.",
		slog.Any("organization_id", u.OrganizationID),
		slog.Any("id", u.ID),
		slog.String("full_name", u.Name),
		slog.String("email", u.Email),
		slog.String("password_hash_algorithm", u.PasswordHashAlgorithm),
		slog.String("password_hash", u.PasswordHash))

	return u, nil
}
