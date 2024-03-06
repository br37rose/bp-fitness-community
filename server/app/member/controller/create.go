package controller

import (
	"context"
	"fmt"
	"strings"
	"time"

	"log/slog"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	user_s "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/user/datastore"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/config/constants"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/utils/httperror"
)

type MemberCreateRequestIDO struct {
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
}

func (impl *MemberControllerImpl) userFromCreateRequest(requestData *MemberCreateRequestIDO) (*user_s.User, error) {
	passwordHash, err := impl.Password.GenerateHashFromPassword(requestData.Password)
	if err != nil {
		impl.Logger.Error("hashing error", slog.Any("error", err))
		return nil, err
	}

	// Make shipping and billing same if not selected.
	if requestData.HasShippingAddress == false {
		requestData.ShippingCity = requestData.City
		requestData.ShippingCountry = requestData.Country
		requestData.ShippingAddressLine1 = requestData.AddressLine1
		requestData.ShippingAddressLine2 = requestData.AddressLine2
		requestData.ShippingPostalCode = requestData.PostalCode
		requestData.ShippingRegion = requestData.Region
	}

	return &user_s.User{
		OrganizationID:            requestData.OrganizationID,
		FirstName:                 requestData.FirstName,
		LastName:                  requestData.LastName,
		Email:                     requestData.Email,
		PasswordHash:              passwordHash,
		PasswordHashAlgorithm:     impl.Password.AlgorithmName(),
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
		Tags:                      make([]*user_s.UserTag, 0),
		Comments:                  make([]*user_s.UserComment, 0),
		StripeInvoices:            make([]*user_s.StripeInvoice, 0),
	}, nil
}

func (impl *MemberControllerImpl) Create(ctx context.Context, requestData *MemberCreateRequestIDO) (*user_s.User, error) {
	m, err := impl.userFromCreateRequest(requestData)
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
		b, err := impl.OrganizationStorer.GetByID(sessCtx, m.OrganizationID)
		if err != nil {
			impl.Logger.Error("database error", slog.Any("err", err))
			return nil, err
		}
		if b == nil {
			impl.Logger.Warn("organization does not exists validation error")
			return nil, httperror.NewForBadRequestWithSingleField("organization_id", "organization does not exist")
		}

		// Lookup the user in our database, else return a `400 Bad Request` error.
		u, err := impl.UserStorer.GetByEmail(sessCtx, m.Email)
		if err != nil {
			impl.Logger.Error("database error", slog.Any("err", err))
			return nil, err
		}
		if u != nil {
			impl.Logger.Warn("user already exists validation error")
			return nil, httperror.NewForBadRequestWithSingleField("email", "email is not unique")
		}

		// Extract from our session the following data.
		orgID, _ := sessCtx.Value(constants.SessionUserOrganizationID).(primitive.ObjectID)
		orgName, _ := sessCtx.Value(constants.SessionUserOrganizationName).(string)
		userID, _ := sessCtx.Value(constants.SessionUserID).(primitive.ObjectID)
		userName, _ := sessCtx.Value(constants.SessionUserName).(string)
		ipAddress, _ := sessCtx.Value(constants.SessionIPAddress).(string)

		// Add defaults.
		m.OrganizationID = b.ID
		m.OrganizationName = b.Name
		m.Email = strings.ToLower(m.Email)
		m.OrganizationID = orgID
		m.OrganizationName = orgName
		m.ID = primitive.NewObjectID()
		m.CreatedAt = time.Now()
		m.CreatedByUserID = userID
		m.CreatedByUserName = userName
		m.CreatedFromIPAddress = ipAddress
		m.ModifiedAt = time.Now()
		m.ModifiedByUserID = userID
		m.ModifiedByUserName = userName
		m.ModifiedFromIPAddress = ipAddress
		m.Role = user_s.UserRoleMember
		m.Name = fmt.Sprintf("%s %s", m.FirstName, m.LastName)
		m.LexicalName = fmt.Sprintf("%s, %s", m.LastName, m.FirstName)
		m.WasEmailVerified = true

		// Save to our database.
		if err := impl.UserStorer.Create(sessCtx, m); err != nil {
			impl.Logger.Error("database create error", slog.Any("error", err))
			return nil, err
		}

		paymentProcessorCustomerID, err := impl.PaymentProcessor.CreateCustomer(
			fmt.Sprintf("%s %s", m.FirstName, m.LastName),
			m.Email,
			"", // description...
			fmt.Sprintf("%s %s Shipping Address", m.FirstName, m.LastName),
			m.Phone,
			m.ShippingCity, m.ShippingCountry, m.ShippingAddressLine1, m.ShippingAddressLine2, m.ShippingPostalCode, m.ShippingRegion, // Shipping
			m.City, m.Country, m.AddressLine1, m.AddressLine2, m.PostalCode, m.Region, // Billing
		)
		if err != nil {
			impl.Logger.Error("creating customer from payment processor error", slog.Any("error", err))
			return nil, err
		}
		m.PaymentProcessorCustomerID = *paymentProcessorCustomerID
		if err := impl.UserStorer.UpdateByID(sessCtx, m); err != nil {
			impl.Logger.Error("database update error", slog.Any("error", err))
			return nil, err
		}

		return m, nil
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
