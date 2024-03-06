package datastore

import (
	"context"
	"log"
	"time"

	"log/slog"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	c "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/config"
)

const (
	UserStatusActive                              = 1
	UserStatusArchived                            = 100
	UserRoleRoot                                  = 1
	UserRoleAdmin                                 = 2
	UserRoleTrainer                               = 3
	UserRoleMember                                = 4
	HowDidYouHearAboutUsOther                     = 1
	HowDidYouHearAboutUsDidNotAnswer              = 2
	HowDidYouHearAboutUsFriend                    = 5
	HowDidYouHearAboutUsSocialMedia               = 6
	HowDidYouHearAboutUsBlog                      = 7
	CheckoutSessionPaymentStatusNoPaymentRequired = "no_payment_required"
	CheckoutSessionPaymentStatusPaid              = "paid"
	CheckoutSessionPaymentStatusUnpaid            = "unpaid"
	CheckoutSessionStatusComplete                 = "complete"
	CheckoutSessionStatusExpired                  = "expired"
	CheckoutSessionStatusOpen                     = "open"
	SubscriptionIntervalDay                       = "day"
	SubscriptionIntervalWeek                      = "week"
	SubscriptionIntervalMonth                     = "month"
	SubscriptionIntervalYear                      = "year"
	SubscriptionStatusActive                      = "active"
	SubscriptionStatusAll                         = "all"
	SubscriptionStatusCanceled                    = "canceled"
	SubscriptionStatusIncomplete                  = "incomplete"
	SubscriptionStatusIncompleteExpired           = "incomplete_expired"
	SubscriptionStatusPastDue                     = "past_due"
	SubscriptionStatusTrialing                    = "trialing"
	SubscriptionStatusUnpaid                      = "unpaid"
	UserPrimaryHealthTrackingDeviceTypeNone       = 0
	UserPrimaryHealthTrackingDeviceTypeFitBit     = 1
)

// UserLite struct represents the current user but minimal view which is enough
// for the user to utilize the frontend. To get further detailed view then we
// use the `User` struct.
type UserLite struct {
	ID                        primitive.ObjectID `bson:"_id" json:"id"`
	OrganizationID            primitive.ObjectID `bson:"organization_id" json:"organization_id,omitempty"`
	OrganizationName          string             `bson:"organization_name" json:"organization_name"`
	FirstName                 string             `bson:"first_name" json:"first_name"`
	LastName                  string             `bson:"last_name" json:"last_name"`
	Name                      string             `bson:"name" json:"name"`
	LexicalName               string             `bson:"lexical_name" json:"lexical_name"`
	Email                     string             `bson:"email" json:"email"`
	Role                      int8               `bson:"role" json:"role"`
	Phone                     string             `bson:"phone" json:"phone,omitempty"`
	Country                   string             `bson:"country" json:"country,omitempty"`
	Region                    string             `bson:"region" json:"region,omitempty"`
	City                      string             `bson:"city" json:"city,omitempty"`
	PostalCode                string             `bson:"postal_code" json:"postal_code,omitempty"`
	AddressLine1              string             `bson:"address_line_1" json:"address_line_1,omitempty"`
	AddressLine2              string             `bson:"address_line_2" json:"address_line_2,omitempty"`
	HowDidYouHearAboutUs      int8               `bson:"how_did_you_hear_about_us" json:"how_did_you_hear_about_us,omitempty"`
	HowDidYouHearAboutUsOther string             `bson:"how_did_you_hear_about_us_other" json:"how_did_you_hear_about_us_other,omitempty"`
	AgreeTOS                  bool               `bson:"agree_tos" json:"agree_tos,omitempty"`
	AgreePromotionsEmail      bool               `bson:"agree_promotions_email" json:"agree_promotions_email,omitempty"`
	Status                    int8               `bson:"status" json:"status"`

	// Indicates the member is a paying member because they are enrolled in a
	// monthly or annual subscription which is actively occuring. If user
	// cancels their subscription then this value gets set to false.
	IsSubscriber bool `bson:"is_subscriber" json:"is_subscriber"`

	// The name of the payment processor we are using to handle payments with
	// this particular member.
	PaymentProcessorName string `bson:"payment_processor_name" json:"payment_processor_name"`

	// The subscription details from the `Stripe` payment processor, if we
	// enrolled the user through Stripe. All details are pertaining to whatever
	// Stripe requests to support subscription payments. This field is only to
	// be used if the payment processor we are using is "Stripe, Inc.".
	StripeSubscription *StripeSubscription `bson:"stripe_subscription" json:"stripe_subscription"`

	// PrimaryHealthTrackingDeviceType indicates what primary health tracking device the
	// user is using with our system.
	PrimaryHealthTrackingDeviceType int8 `bson:"primary_health_tracking_device_type" json:"primary_health_tracking_device_type"`

	// HeartRateMetricID is the unique identification used to tie the user's
	// heart rate data to.
	PrimaryHealthTrackingDeviceHeartRateMetricID primitive.ObjectID `bson:"primary_health_tracking_device_heart_rate_metric_id" json:"primary_health_tracking_device_heart_rate_metric_id,omitempty"`

	// StepsCountMetricID is the unique identification used to tie the user's
	// steps count data to.
	PrimaryHealthTrackingDeviceStepsCountMetricID primitive.ObjectID `bson:"primary_health_tracking_device_steps_count_metric_id" json:"primary_health_tracking_device_steps_count_metric_id,omitempty"`

	// Tags stores all the user created tags for the user's account.
	Tags []*UserTag `bson:"tags" json:"tags"`

	AvatarObjectExpiry time.Time `bson:"avatar_object_expiry" json:"avatar_object_expiry"`
	AvatarObjectURL    string    `bson:"avatar_object_url" json:"avatar_object_url"`
	AvatarObjectKey    string    `bson:"avatar_object_key" json:"avatar_object_key"`
	AvatarFileType     string    `bson:"avatar_file_type" json:"avatar_file_type"`
	AvatarFileName     string    `bson:"avatar_file_name" json:"avatar_file_name"`
}

type User struct {
	Email                     string             `bson:"email" json:"email"`
	FirstName                 string             `bson:"first_name" json:"first_name"`
	LastName                  string             `bson:"last_name" json:"last_name"`
	Status                    int8               `bson:"status" json:"status"`
	ID                        primitive.ObjectID `bson:"_id" json:"id"`
	OrganizationID            primitive.ObjectID `bson:"organization_id" json:"organization_id,omitempty"`
	OrganizationName          string             `bson:"organization_name" json:"organization_name"`
	Name                      string             `bson:"name" json:"name"`
	LexicalName               string             `bson:"lexical_name" json:"lexical_name"`
	PasswordHashAlgorithm     string             `bson:"password_hash_algorithm" json:"password_hash_algorithm,omitempty"`
	PasswordHash              string             `bson:"password_hash" json:"password_hash,omitempty"`
	Role                      int8               `bson:"role" json:"role"`
	WasEmailVerified          bool               `bson:"was_email_verified" json:"was_email_verified"`
	EmailVerificationCode     string             `bson:"email_verification_code,omitempty" json:"email_verification_code,omitempty"`
	EmailVerificationExpiry   time.Time          `bson:"email_verification_expiry,omitempty" json:"email_verification_expiry,omitempty"`
	Phone                     string             `bson:"phone" json:"phone,omitempty"`
	Country                   string             `bson:"country" json:"country,omitempty"`
	Region                    string             `bson:"region" json:"region,omitempty"`
	City                      string             `bson:"city" json:"city,omitempty"`
	PostalCode                string             `bson:"postal_code" json:"postal_code,omitempty"`
	AddressLine1              string             `bson:"address_line_1" json:"address_line_1,omitempty"`
	AddressLine2              string             `bson:"address_line_2" json:"address_line_2,omitempty"`
	HasShippingAddress        bool               `bson:"has_shipping_address" json:"has_shipping_address,omitempty"`
	ShippingName              string             `bson:"shipping_name" json:"shipping_name,omitempty"`
	ShippingPhone             string             `bson:"shipping_phone" json:"shipping_phone,omitempty"`
	ShippingCountry           string             `bson:"shipping_country" json:"shipping_country,omitempty"`
	ShippingRegion            string             `bson:"shipping_region" json:"shipping_region,omitempty"`
	ShippingCity              string             `bson:"shipping_city" json:"shipping_city,omitempty"`
	ShippingPostalCode        string             `bson:"shipping_postal_code" json:"shipping_postal_code,omitempty"`
	ShippingAddressLine1      string             `bson:"shipping_address_line1" json:"shipping_address_line1,omitempty"`
	ShippingAddressLine2      string             `bson:"shipping_address_line2" json:"shipping_address_line2,omitempty"`
	HowDidYouHearAboutUs      int8               `bson:"how_did_you_hear_about_us" json:"how_did_you_hear_about_us,omitempty"`
	HowDidYouHearAboutUsOther string             `bson:"how_did_you_hear_about_us_other" json:"how_did_you_hear_about_us_other,omitempty"`
	AgreeTOS                  bool               `bson:"agree_tos" json:"agree_tos,omitempty"`
	AgreePromotionsEmail      bool               `bson:"agree_promotions_email" json:"agree_promotions_email,omitempty"`
	CreatedAt                 time.Time          `bson:"created_at" json:"created_at,omitempty"`
	CreatedByUserID           primitive.ObjectID `bson:"created_by_user_id" json:"created_by_user_id"`
	CreatedByUserName         string             `bson:"created_by_user_name" json:"created_by_user_name"`
	CreatedFromIPAddress      string             `bson:"created_from_ip_address" json:"created_from_ip_address"`
	ModifiedAt                time.Time          `bson:"modified_at" json:"modified_at,omitempty"`
	ModifiedByUserID          primitive.ObjectID `bson:"modified_by_user_id" json:"modified_by_user_id"`
	ModifiedByUserName        string             `bson:"modified_by_user_name" json:"modified_by_user_name"`
	ModifiedFromIPAddress     string             `bson:"modified_from_ip_address" json:"modified_from_ip_address"`

	// Tags stores all the user created tags for the user's account.
	Tags []*UserTag `bson:"tags" json:"tags"`

	// Comments are all the comments created by admin staff.
	Comments []*UserComment `bson:"comments" json:"comments"`

	Purchases []*UserPurchase `bson:"offer_purchases" json:"offer_purchases"`

	// Indicates the member is a paying member because they are enrolled in a
	// monthly or annual subscription which is actively occuring. If user
	// cancels their subscription then this value gets set to false.
	IsSubscriber bool `bson:"is_subscriber" json:"is_subscriber"`

	// SubscriptionOfferID holds the `ID` value of the subscription `Offer` the
	// user is enrolled in.
	SubscriptionOfferID primitive.ObjectID `bson:"subscription_offer_id" json:"subscription_offer_id,omitempty"`
	// SubscriptionOfferName holds the name of the subscription `Offer` the user
	// us enrolled in.
	SubscriptionOfferName string `bson:"subscription_offer_name" json:"subscription_offer_name"`
	// SubscriptionStatus holds the state of the subscription for the user.
	SubscriptionStatus string `bson:"subscription_status" json:"subscription_status"`
	// SubscriptionStartedAt holds date on when the subscription started.
	SubscriptionStartedAt time.Time `bson:"subscription_started_at" json:"subscription_started_at"`

	// The name of the payment processor we are using to handle payments with
	// this particular member.
	PaymentProcessorName string `bson:"payment_processor_name" json:"payment_processor_name"`

	// The unique identifier used by the payment processor which has a somesort of
	// copy of this member's details saved and we can reference that customer on
	// the payment processor using this `customer_id`.
	PaymentProcessorCustomerID string `bson:"payment_processor_customer_id" json:"payment_processor_customer_id"`

	// The subscription details from the `Stripe` payment processor, if we
	// enrolled the user through Stripe. All details are pertaining to whatever
	// Stripe requests to support subscription payments. This field is only to
	// be used if the payment processor we are using is "Stripe, Inc.".
	StripeSubscription *StripeSubscription `bson:"stripe_subscription" json:"stripe_subscription"`

	// The list of invoices that belong to the user for their continued usage
	// of the subscription plan. This field is only to
	// be used if the payment processor we are using is "Stripe, Inc.".
	StripeInvoices []*StripeInvoice `bson:"stripe_invoices" json:"stripe_invoices"`

	OfferID   primitive.ObjectID `bson:"offer_id" json:"offer_id"`
	OfferName string             `bson:"offer_name" json:"offer_name"`
	// OfferMembershipRank is unique identifier to specify this offer's value in the ranking system, higher is better.
	OfferMembershipRank int `bson:"offer_membership_rank" json:"offer_membership_rank"`

	// FitBitAppID is the unique identifier of the fitbit authorized device in
	// our system. The authorized device information is found in the the
	// `FitBitApp` domain.
	FitBitAppID primitive.ObjectID `bson:"fitbit_app_id" json:"-"`

	// PrimaryHealthTrackingDeviceType indicates what primary health tracking device the
	// user is using with our system.
	PrimaryHealthTrackingDeviceType int8 `bson:"primary_health_tracking_device_type" json:"primary_health_tracking_device_type"`

	// HeartRateMetricID is the unique identification used to tie the user's
	// heart rate data to.
	PrimaryHealthTrackingDeviceHeartRateMetricID primitive.ObjectID `bson:"primary_health_tracking_device_heart_rate_metric_id" json:"primary_health_tracking_device_heart_rate_metric_id,omitempty"`

	// StepsCountMetricID is the unique identification used to tie the user's
	// steps count data to.
	PrimaryHealthTrackingDeviceStepsCountMetricID primitive.ObjectID `bson:"primary_health_tracking_device_steps_count_metric_id" json:"primary_health_tracking_device_steps_count_metric_id,omitempty"`

	AvatarObjectExpiry time.Time `bson:"avatar_object_expiry" json:"avatar_object_expiry"`
	AvatarObjectURL    string    `bson:"avatar_object_url" json:"avatar_object_url"`
	AvatarObjectKey    string    `bson:"avatar_object_key" json:"avatar_object_key"`
	AvatarFileType     string    `bson:"avatar_file_type" json:"avatar_file_type"`
	AvatarFileName     string    `bson:"avatar_file_name" json:"avatar_file_name"`

	// OTPEnabled controls whether we force 2FA or not during login.
	OTPEnabled bool `bson:"otp_enabled" json:"otp_enabled"`

	// OTPVerified indicates user has successfully validated their opt token afer enabling 2FA thus turning it on.
	OTPVerified bool `bson:"otp_verified" json:"otp_verified"`

	// OTPValidated automatically gets set as `false` on successful login and then sets `true` once successfully validated by 2FA.
	OTPValidated bool `bson:"otp_validated" json:"otp_validated"`

	// OTPSecret the unique one-time password secret to be shared between our
	// backend and 2FA authenticator sort of apps that support `TOPT`.
	OTPSecret string `bson:"otp_secret" json:"-"`

	// OTPAuthURL is the URL used to share.
	OTPAuthURL string `bson:"otp_auth_url" json:"-"`
}

// Address describes common properties for an Address hash.
type Address struct {
	City       string `json:"city"`
	Country    string `json:"country"`
	Line1      string `json:"line1"`
	Line2      string `json:"line2"`
	PostalCode string `json:"postal_code"`
	State      string `json:"state"`
}

type UserTag struct {
	ID             primitive.ObjectID `bson:"_id" json:"id"`
	OrganizationID primitive.ObjectID `bson:"organization_id" json:"organization_id,omitempty"`
	UserID         primitive.ObjectID `bson:"user_id" json:"user_id"`
	Text           string             `bson:"text" json:"text"`
	Description    string             `bson:"description" json:"description"`
	Status         int8               `bson:"status" json:"status"`
}

// Mailing and shipping address for the customer. Appears on invoices emailed to this customer.
type CustomerShippingDetails struct {
	Address Address `json:"address"`
	// Recipient name.
	Name string `json:"name"`
	// Recipient phone (including extension).
	Phone string `json:"phone"`
}

type UserComment struct {
	ID               primitive.ObjectID `bson:"_id" json:"id"`
	OrganizationID   primitive.ObjectID `bson:"organization_id" json:"organization_id"`
	CreatedAt        time.Time          `bson:"created_at,omitempty" json:"created_at,omitempty"`
	CreatedByUserID  primitive.ObjectID `bson:"created_by_user_id" json:"created_by_user_id"`
	CreatedByName    string             `bson:"created_by_name" json:"created_by_name"`
	ModifiedAt       time.Time          `bson:"modified_at,omitempty" json:"modified_at,omitempty"`
	ModifiedByUserID primitive.ObjectID `bson:"modified_by_user_id" json:"modified_by_user_id"`
	ModifiedByName   string             `bson:"modified_by_name" json:"modified_by_name"`
	Content          string             `bson:"content" json:"content"`
}

type UserPurchase struct {
	ID                 primitive.ObjectID `bson:"_id" json:"id"`
	OrganizationID     primitive.ObjectID `bson:"organization_id" json:"organization_id"`
	CreatedAt          time.Time          `bson:"created_at,omitempty" json:"created_at,omitempty"`
	ModifiedAt         time.Time          `bson:"modified_at,omitempty" json:"modified_at,omitempty"`
	OfferID            primitive.ObjectID `bson:"offer_id" json:"offer_id"`                         // Copied from `Offer`.
	OfferName          string             `bson:"offer_name" json:"offer_name"`                     // Copied from `Offer`.
	OfferDescription   string             `bson:"offer_description" json:"offer_description"`       // Copied from `Offer`.
	OfferPrice         float64            `bson:"offer_price" json:"offer_price"`                   // Copied from `Offer`.
	OfferPriceCurrency string             `bson:"offer_price_currency" json:"offer_price_currency"` // Copied from `Offer`.
	OfferPayFrequency  int8               `bson:"offer_pay_frequency" json:"offer_pay_frequency"`   // Copied from `Offer`.ç
	// Controls how the user is able to book in our system. Special thanks to http://www.heppnetz.de/ontologies/goodrelations/v1#BusinessFunction.
	OfferBusinessFunction int8 `bson:"offer_business_function" json:"offer_business_function"` // Copied from `Offer`.
	OfferType             int8 `bson:"offer_type" json:"offer_type"`
	// IncludesOfferIDs controls what benefits this offer will include with other offers. For example, if this offer is "Gold" subscription, then inside this variable would include the IDs of the "Silver" and "Bronze" subscriptions.
	IncludesOfferIDs []primitive.ObjectID `bson:"includes_offer_ids" json:"includes_offer_ids"`
}

type StripeSubscription struct {
	// The unique identification of the specific subscription product we want to assign.
	PriceID string `bson:"price_id" json:"price_id"`

	// The unique identification created by Stripe after the intent for subscription has been submitted.
	SubscriptionID string `bson:"subscription_id"`

	// The frequency at which a subscription is billed. One of `day`, `week`, `month` or `year`.
	Interval string `bson:"interval" json:"interval"`

	// The current state of the subscription.
	Status string `bson:"status" json:"status"`

	// OfferPurchase indicates the purchased offer we have for this subscription.
	OfferPurchase *UserPurchase `bson:"offer_purchase" json:"offer_purchase"`
}

type StripeInvoice struct {
	// The unique identification created by Stripe to present this particular invoice.
	InvoiceID string `bson:"invoice_id" json:"invoice_id"`
	// Time at which the object was created. Measured in seconds since the Unix epoch.
	Created int64 `bson:"created" json:"created"`
	// Whether payment was successfully collected for this invoice. An invoice can be paid (most commonly) with a charge or with credit from the customer's account balance.
	Paid bool `bson:"paid" json:"paid"`
	// The URL for the hosted invoice page, which allows customers to view and pay an invoice. If the invoice has not been finalized yet, this will be null.
	HostedInvoiceURL string `bson:"hosted_invoice_url" json:"hosted_invoice_url"`
	// The link to download the PDF for the invoice. If the invoice has not been finalized yet, this will be null.
	InvoicePDF string `bson:"invoice_pdf" json:"invoice_pdf"`
	// The integer amount in %s representing the subtotal of the invoice before any invoice level discount or tax is applied. Item discounts are already incorporated
	SubtotalExcludingTax float64 `bson:"subtotal_excluding_tax" json:"subtotal_excluding_tax"`
	// The amount of tax on this invoice. This is the sum of all the tax amounts on this invoice.
	Tax float64 `bson:"tax" json:"tax"`
	// Total after discounts and taxes.
	Total float64 `bson:"total" json:"total"`
	// A unique, identifying string that appears on emails sent to the customer for this invoice. This starts with the customer's unique invoice_prefix if it is specified.
	Number string `bson:"number" json:"number"`
	// Three-letter [ISO currency code](https://www.iso.org/iso-4217-currency-codes.html), in lowercase. Must be a [supported currency](https://stripe.com/docs/currencies).
	Currency string `bson:"currency" json:"currency"`
}

type UserListFilter struct {
	// Pagination related.
	Cursor    primitive.ObjectID
	PageSize  int64
	SortField string
	SortOrder int8 // 1=ascending | -1=descending

	// Filter related.
	OrganizationID  primitive.ObjectID `bson:"organization_id" json:"organization_id,omitempty"`
	Role            int8               `bson:"role" json:"role"`
	Status          int8               `json:"status"`
	ExcludeArchived bool               `json:"exclude_archived"`
	SearchText      string             `json:"search_text"`
	FirstName       string             `json:"first_name"`
	LastName        string             `json:"last_name"`
	Email           string             `json:"email"`
	Phone           string             `json:"phone"`
}

type UserListResult struct {
	Results     []*User            `json:"results"`
	NextCursor  primitive.ObjectID `json:"next_cursor"`
	HasNextPage bool               `json:"has_next_page"`
}

type UserAsSelectOption struct {
	Value primitive.ObjectID `bson:"_id" json:"value"` // Extract from the database `_id` field and output through API as `value`.
	Label string             `bson:"name" json:"label"`
}

type StripeInvoiceListResult struct {
	Results     []*StripeInvoice `json:"results"`
	NextCursor  int64            `json:"next_cursor"`
	HasNextPage bool             `json:"has_next_page"`
}

func ToUserLite(u *User) *UserLite {
	return &UserLite{
		ID:                              u.ID,
		OrganizationID:                  u.OrganizationID,
		OrganizationName:                u.OrganizationName,
		FirstName:                       u.FirstName,
		LastName:                        u.LastName,
		Name:                            u.Name,
		LexicalName:                     u.LexicalName,
		Email:                           u.Email,
		Role:                            u.Role,
		Phone:                           u.Phone,
		Country:                         u.Country,
		Region:                          u.Region,
		City:                            u.City,
		PostalCode:                      u.PostalCode,
		AddressLine1:                    u.AddressLine1,
		AddressLine2:                    u.AddressLine2,
		HowDidYouHearAboutUs:            u.HowDidYouHearAboutUs,
		HowDidYouHearAboutUsOther:       u.HowDidYouHearAboutUsOther,
		AgreeTOS:                        u.AgreeTOS,
		AgreePromotionsEmail:            u.AgreePromotionsEmail,
		Status:                          u.Status,
		IsSubscriber:                    u.IsSubscriber,
		PaymentProcessorName:            u.PaymentProcessorName,
		StripeSubscription:              u.StripeSubscription,
		PrimaryHealthTrackingDeviceType: u.PrimaryHealthTrackingDeviceType,
		PrimaryHealthTrackingDeviceHeartRateMetricID:  u.PrimaryHealthTrackingDeviceHeartRateMetricID,
		PrimaryHealthTrackingDeviceStepsCountMetricID: u.PrimaryHealthTrackingDeviceStepsCountMetricID,
		Tags:               u.Tags,
		AvatarObjectExpiry: u.AvatarObjectExpiry,
		AvatarObjectURL:    u.AvatarObjectURL,
		AvatarObjectKey:    u.AvatarObjectKey,
		AvatarFileType:     u.AvatarFileType,
		AvatarFileName:     u.AvatarFileName,
	}
}

// UserStorer Interface for user.
type UserStorer interface {
	Create(ctx context.Context, m *User) error
	GetByID(ctx context.Context, id primitive.ObjectID) (*User, error)
	GetLiteByID(ctx context.Context, id primitive.ObjectID) (*UserLite, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	GetLiteByEmail(ctx context.Context, email string) (*UserLite, error)
	GetByVerificationCode(ctx context.Context, verificationCode string) (*User, error)
	GetByPaymentProcessorCustomerID(ctx context.Context, paymentProcessorCustomerID string) (*User, error)
	GetStripeInvoiceByPaymentProcessorInvoiceID(ctx context.Context, paymentProcessorInvoiceID string) (*StripeInvoice, error)
	CheckIfExistsByEmail(ctx context.Context, email string) (bool, error)
	UpdateByID(ctx context.Context, m *User) error
	UpdateLiteByID(ctx context.Context, m *UserLite) error
	UpdateStripeInvoiceByPaymentProcessorInvoiceID(ctx context.Context, newInvoice *StripeInvoice) error
	ListByFilter(ctx context.Context, f *UserListFilter) (*UserListResult, error)
	ListAsSelectOptionByFilter(ctx context.Context, f *UserListFilter) ([]*UserAsSelectOption, error)
	ListLatestStripeInvoices(ctx context.Context, userID primitive.ObjectID, cursor int64, limit int64) (*StripeInvoiceListResult, error)
	DeleteByID(ctx context.Context, id primitive.ObjectID) error
	// //TODO: Add more...
}

type UserStorerImpl struct {
	Logger     *slog.Logger
	DbClient   *mongo.Client
	Collection *mongo.Collection
}

func NewDatastore(appCfg *c.Conf, loggerp *slog.Logger, client *mongo.Client) UserStorer {
	// ctx := context.Background()
	uc := client.Database(appCfg.DB.Name).Collection("users")

	// The following few lines of code will create the index for our app for this
	// colleciton.
	indexModel := mongo.IndexModel{
		Keys: bson.D{
			{"organization_name", "text"},
			{"branch_name", "text"},
			{"name", "text"},
			{"lexical_name", "text"},
			{"email", "text"},
			{"phone", "text"},
			{"country", "text"},
			{"region", "text"},
			{"city", "text"},
			{"postal_code", "text"},
			{"address_line_1", "text"},
		},
	}
	_, err := uc.Indexes().CreateOne(context.TODO(), indexModel)
	if err != nil {
		// It is important that we crash the app on startup to meet the
		// requirements of `google/wire` framework.
		log.Fatal(err)
	}

	s := &UserStorerImpl{
		Logger:     loggerp,
		DbClient:   client,
		Collection: uc,
	}

	return s
}
