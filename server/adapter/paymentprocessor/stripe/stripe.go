package stripe

import (
	"log/slog"

	stripe "github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/checkout/session"
	"github.com/stripe/stripe-go/v72/customer"
	"github.com/stripe/stripe-go/v72/invoice"
	"github.com/stripe/stripe-go/v72/paymentintent"
	"github.com/stripe/stripe-go/v72/price"
	"github.com/stripe/stripe-go/v72/setupintent"
	"github.com/stripe/stripe-go/v72/sub"

	c "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/config"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/provider/uuid"
)

// Special thanks:
// https://github.com/stripe-samples/checkout-single-subscription/blob/main/server/go/server.go
// https://github.com/stripe-samples/checkout-single-subscription/blob/main/client/index.html
// Build a subscriptions integration via https://stripe.com/docs/billing/subscriptions/build-subscriptions?ui=elements
// Prebuilt checkout page via https://stripe.com/docs/checkout/quickstart
// Go & Stripe Subscriptions Quickstart via https://medium.com/@snassr/go-stripe-subscriptions-quickstart-e01db656f2a9
// Webhooks via https://www.youtube.com/watch?v=ZEuYtQ-vnB4

// DEVELOPERS NOTE:
// Run `task stripewebhook` and you can test using the following (https://stripe.com/docs/api/events/types):
// stripe trigger customer.created
// stripe trigger customer.deleted
// stripe trigger customer.subscription.created
// stripe trigger customer.subscription.paused
// stripe trigger customer.subscription.deleted
// stripe trigger customer.subscription.created
// stripe trigger customer.subscription.marked_uncollectible
// stripe trigger customer.subscription.paid
// stripe trigger customer.subscription.payment_failed
// stripe trigger customer.subscription.payment_succeeded
// stripe trigger invoice.paid

// PaymentProcessorProduct Structure represents the product that exists on record
// in the payment processor's database.
type PaymentProcessorProduct struct {
	ProductID string
	PriceID   string
	Price     int64
}

type PaymentProcessor interface {
	GetName() string
	GetProducts() ([]PaymentProcessorProduct, error)
	GetWebhookSecretKey() string
	CreateCustomer(fullName, email, descr, shipName, shipPhone, shipCity, shipCountry, shipLine1, shipLine2, shipPostalCode, shipState, billCity, billCountry, billLine1, billLine2, billPostalCode, billState string) (*string, error)
	UpdateCustomer(customerID, fullName, email, descr, shipName, shipPhone, shipCity, shipCountry, shipLine1, shipLine2, shipPostalCode, shipState, billCity, billCountry, billLine1, billLine2, billPostalCode, billState string) error
	SetupNewCard(customerID string) (secret *string, err error)
	CreateSubscriptionCheckoutSessionURL(domain, successURL, canceledURL, customerID, priceID string, customerHasShippingAddress bool, couponID string) (string, error)
	CreateOneTimeCheckoutSessionURL(domain, successCallbackURL, canceledCallbackURL, customerID, priceID string, customerHasShippingAddress bool, couponID string) (string, error)
	GetCheckoutSession(sessionID string) (*stripe.CheckoutSession, error)
	GetCheckoutSessionLineItems(sessionID string) ([]*stripe.LineItem, error)
	GetCustomer(customerID string) (*stripe.Customer, error)
	GetSubscription(subscriptionID string) (*stripe.Subscription, error)
	CancelSubscription(subscriptionID string) (*stripe.Subscription, error)
	ListInvoicesByCustomerID(customerID string) ([]*stripe.Invoice, error)
	GetPaymentIntent(paymentIntentID string) (*stripe.PaymentIntent, error)
	GetLatestInvoiceByCustomerID(customerID string) (*stripe.Invoice, error)
}

type stripePaymentProcessor struct {
	WebhookSecretKey string
	UUID             uuid.Provider
	Logger           *slog.Logger
}

func NewPaymentProcessor(cfg *c.Conf, logger *slog.Logger, uuidp uuid.Provider) PaymentProcessor {
	// Defensive code: Make sure we have access to the file before proceeding any further with the code.
	logger.Debug("payment processor initializing...")

	// Initialize our secret key for the stripe payment processor which is required.
	stripe.Key = cfg.PaymentProcessor.SecretKey

	logger.Debug("payment processor was initialized with stripe.")

	return &stripePaymentProcessor{
		UUID:             uuidp,
		Logger:           logger,
		WebhookSecretKey: cfg.PaymentProcessor.WebhookSecretKey,
	}
}

// Return the name of the payment processor of this adapter.
func (pm *stripePaymentProcessor) GetName() string {
	return "Stripe, Inc."
}

func (pm *stripePaymentProcessor) GetWebhookSecretKey() string {
	return pm.WebhookSecretKey
}

// GetProducts Function will pull the latest products on record in the stripe
// for our account and return the product details.
func (pm *stripePaymentProcessor) GetProducts() ([]PaymentProcessorProduct, error) {
	// Special thanks:
	// https://medium.com/@snassr/go-stripe-subscriptions-quickstart-e01db656f2a9

	products := make([]PaymentProcessorProduct, 0)
	priceParams := &stripe.PriceListParams{}
	priceIterator := price.List(priceParams)
	for priceIterator.Next() {
		products = append(products, PaymentProcessorProduct{
			ProductID: priceIterator.Price().Product.ID,
			PriceID:   priceIterator.Price().ID,
			Price:     priceIterator.Price().UnitAmount,
		})
	}
	return products, nil
}

// func (pm *stripePaymentProcessor) CreateInvoice() error {
//  https://stripe.com/docs/receipts#:~:text=To%20generate%20invoices%2C%20first%2C%20in,invoice%20template%20in%20the%20Dashboard.
//  https://stripe.com/docs/invoicing/integration
//  https://stripe.com/docs/api/invoices
//  https://stripe.com/docs/invoicing/overview?dashboard-or-api=api
//
// 	params := &stripe.InvoiceParams{
// 		Customer:         stripe.String("{{CUSTOMER_ID}}"),
// 		CollectionMethod: stripe.String(string(stripe.InvoiceCollectionMethodSendInvoice)),
// 		DaysUntilDue:     stripe.Int64(30),
// 	}
// 	result, err := invoice.New(params)
//
// 	//
//
// 	params := &stripe.InvoiceItemParams{
// 		Customer: stripe.String("{{CUSTOMER_ID}}"),
// 		Price:    stripe.String("{{PRICE_ID}}"),
// 		Invoice:  stripe.String("{{INVOICE_ID}}"),
// 	}
// 	result, err := invoiceitem.New(params)
//
// 	//
//
// 	params := &stripe.InvoiceFinalizeInvoiceParams{}
// 	result, err := invoice.FinalizeInvoice("{{INVOICE_ID}}", params)
//
// 	return nil
//
// }

// CreateCustomer Function registers our member with the payment processor's
// customer database so we can use for billing purposes. Function will return
// the `customer_id` value that the payment processor generated in their database
// for the saved customer record.
func (pm *stripePaymentProcessor) CreateCustomer(fullName, email, descr, shipName, shipPhone, shipCity, shipCountry, shipLine1, shipLine2, shipPostalCode, shipState, billCity, billCountry, billLine1, billLine2, billPostalCode, billState string) (*string, error) {
	pm.Logger.Debug("Beginning to create customer account on stripe...")
	pm.Logger.Debug("First lookup by email to see if previously created")

	// The following code will lookup the customer's email and if it already
	// exists then return the customer ID here after we update the customer.
	// Special thanks to: https://stripe.com/docs/invoicing/customer
	params1 := &stripe.CustomerListParams{Email: stripe.String(email)}
	result := customer.List(params1)
	for result.Next() {
		c := result.Customer()
		if c.Email == email {
			pm.Logger.Debug("Found customer, proceeding to update...", slog.String("customerID", c.ID))

			// Special thanks:
			// https://medium.com/@snassr/go-stripe-subscriptions-quickstart-e01db656f2a9
			// https://stripe.com/docs/billing/subscriptions/build-subscriptions?ui=elements

			params := &stripe.CustomerParams{
				Name:        &fullName,
				Email:       &email,
				Description: &descr,
				// Shipping: &stripe.CustomerShippingDetailsParams{
				// 	Address: &stripe.AddressParams{
				// 		City:       stripe.String(shipCity),
				// 		Country:    stripe.String(shipCountry),
				// 		Line1:      stripe.String(shipLine1),
				// 		Line2:      stripe.String(shipLine2),
				// 		PostalCode: stripe.String(shipPostalCode),
				// 		State:      stripe.String(shipState),
				// 	},
				// 	Name:  &shipName,
				// 	Phone: &shipPhone,
				// },
				Address: &stripe.AddressParams{
					City:       stripe.String(billCity),
					Country:    stripe.String(billCountry),
					Line1:      stripe.String(billLine1),
					Line2:      stripe.String(billLine2),
					PostalCode: stripe.String(billPostalCode),
					State:      stripe.String(billState),
				},
			}

			// Only attach shipping address if the user has one.
			if shipCountry != "" && shipState != "" && shipCity != "" && shipLine1 != "" {
				pm.Logger.Debug("attaching shipping address to new customer on stripe")
				params.Shipping = &stripe.CustomerShippingDetailsParams{
					Address: &stripe.AddressParams{
						City:       stripe.String(shipCity),
						Country:    stripe.String(shipCountry),
						Line1:      stripe.String(shipLine1),
						Line2:      stripe.String(shipLine2),
						PostalCode: stripe.String(shipPostalCode),
						State:      stripe.String(shipState),
					},
					Name:  &shipName,
					Phone: &shipPhone,
				}
			} else {
				params.Shipping = nil
			}

			_, err := customer.Update(c.ID, params)
			if err != nil {
				return stripe.String(""), err
			}

			pm.Logger.Debug("Updated customer on stripe", slog.String("customerID", c.ID))

			// Return our ID
			return &c.ID, nil
		}
	}

	// Special thanks:
	// https://medium.com/@snassr/go-stripe-subscriptions-quickstart-e01db656f2a9
	// https://stripe.com/docs/billing/subscriptions/build-subscriptions?ui=elements

	pm.Logger.Debug("creating new customer on stripe")

	params := &stripe.CustomerParams{
		Name:        &fullName,
		Email:       &email,
		Description: &descr,
		// Shipping: &stripe.CustomerShippingDetailsParams{
		// 	Address: &stripe.AddressParams{
		// 		City:       stripe.String(shipCity),
		// 		Country:    stripe.String(shipCountry),
		// 		Line1:      stripe.String(shipLine1),
		// 		Line2:      stripe.String(shipLine2),
		// 		PostalCode: stripe.String(shipPostalCode),
		// 		State:      stripe.String(shipState),
		// 	},
		// 	Name:  &shipName,
		// 	Phone: &shipPhone,
		// },
		Address: &stripe.AddressParams{
			City:       stripe.String(billCity),
			Country:    stripe.String(billCountry),
			Line1:      stripe.String(billLine1),
			Line2:      stripe.String(billLine2),
			PostalCode: stripe.String(billPostalCode),
			State:      stripe.String(billState),
		},
	}

	// Only attach shipping address if the user has one.
	if shipCountry != "" && shipState != "" && shipCity != "" && shipLine1 != "" {
		pm.Logger.Debug("attaching shipping address to new customer on stripe")
		params.Shipping = &stripe.CustomerShippingDetailsParams{
			Address: &stripe.AddressParams{
				City:       stripe.String(shipCity),
				Country:    stripe.String(shipCountry),
				Line1:      stripe.String(shipLine1),
				Line2:      stripe.String(shipLine2),
				PostalCode: stripe.String(shipPostalCode),
				State:      stripe.String(shipState),
			},
			Name:  &shipName,
			Phone: &shipPhone,
		}
	} else {
		params.Shipping = nil
	}

	c, err := customer.New(params)
	if err != nil {
		return nil, err
	}

	pm.Logger.Debug("created new customer on stripe", slog.String("customerID", c.ID))
	return &c.ID, nil
}

func (pm *stripePaymentProcessor) UpdateCustomer(
	customerID,
	fullName,
	email,
	descr,
	shipName,
	shipPhone,
	shipCity,
	shipCountry,
	shipLine1,
	shipLine2,
	shipPostalCode,
	shipState,
	billCity,
	billCountry,
	billLine1,
	billLine2,
	billPostalCode,
	billState string,
) error {
	// Special thanks:
	// https://medium.com/@snassr/go-stripe-subscriptions-quickstart-e01db656f2a9
	// https://stripe.com/docs/billing/subscriptions/build-subscriptions?ui=elements

	params := &stripe.CustomerParams{
		Name:        &fullName,
		Email:       &email,
		Description: &descr,
		Address: &stripe.AddressParams{
			City:       stripe.String(billCity),
			Country:    stripe.String(billCountry),
			Line1:      stripe.String(billLine1),
			Line2:      stripe.String(billLine2),
			PostalCode: stripe.String(billPostalCode),
			State:      stripe.String(billState),
		},
	}

	if shipCountry != "" && shipState != "" && shipCity != "" && shipLine1 != "" {
		pm.Logger.Debug("attaching shipping address to customer update on stripe")
		params.Shipping = &stripe.CustomerShippingDetailsParams{
			Address: &stripe.AddressParams{
				City:       stripe.String(shipCity),
				Country:    stripe.String(shipCountry),
				Line1:      stripe.String(shipLine1),
				Line2:      stripe.String(shipLine2),
				PostalCode: stripe.String(shipPostalCode),
				State:      stripe.String(shipState),
			},
			Name:  &shipName,
			Phone: &shipPhone,
		}
	} else {
		params.Shipping = nil
	}

	_, err := customer.Update(customerID, params)
	if err != nil {
		return err
	}
	pm.Logger.Debug("updated customer on stripe", slog.String("customerID", customerID))
	return nil
}

func (pm *stripePaymentProcessor) SetupNewCard(customerID string) (secret *string, err error) {
	// Special thanks:
	// https://medium.com/@snassr/go-stripe-subscriptions-quickstart-e01db656f2a9

	params := &stripe.SetupIntentParams{
		PaymentMethodTypes: []*string{
			stripe.String("card"),
		},
		Customer: &customerID,
	}
	si, err := setupintent.New(params)
	if err != nil {
		return nil, err
	}
	return &si.ClientSecret, nil
}

func (pm *stripePaymentProcessor) createCheckoutSessionURL(
	domain string,
	successCallbackURL string,
	canceledCallbackURL string,
	mode stripe.CheckoutSessionMode,
	customerID string,
	priceID string,
	quantity int64,
	customerHasShippingAddress bool,
	couponID string,
) (string, error) {
	params := &stripe.CheckoutSessionParams{
		SuccessURL: stripe.String("https://" + domain + successCallbackURL + "?session_id={CHECKOUT_SESSION_ID}"),
		CancelURL:  stripe.String("https://" + domain + canceledCallbackURL),
		Mode:       stripe.String(string(mode)),
		// AllowPromotionCodes: stripe.Bool(true),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				Price:    stripe.String(priceID),
				Quantity: stripe.Int64(quantity),
			},
		},
		// AutomaticTax: &stripe.CheckoutSessionAutomaticTaxParams{Enabled: stripe.Bool(true)},
		CustomerUpdate: &stripe.CheckoutSessionCustomerUpdateParams{
			Address: stripe.String("auto"),
			Name:    stripe.String("auto"),
			// Shipping: stripe.String("never"),
		},
		// ShippingAddressCollection: &stripe.CheckoutSessionShippingAddressCollectionParams{
		// 	AllowedCountries: []*string{stripe.String("US"), stripe.String("CA")},
		// },
		BillingAddressCollection: stripe.String("auto"),
	}

	// If the `customer id` was inputted as a parameter then include in the session.
	if customerID != "" {
		params.Customer = &customerID
	}

	// If customer used shipping address then our checkout will require the
	// following changes.
	if customerHasShippingAddress {
		params.CustomerUpdate.Shipping = stripe.String("auto")
		params.ShippingAddressCollection = &stripe.CheckoutSessionShippingAddressCollectionParams{
			AllowedCountries: []*string{stripe.String("US"), stripe.String("CA")},
		}
	}

	// If the `coupon_id` was inputted as a parameter then include in the session.
	// If no coupon coude was inputted then grant access to users to add
	// promotional codes to a checkout session.
	if couponID != "" {
		params.Discounts = []*stripe.CheckoutSessionDiscountParams{
			&stripe.CheckoutSessionDiscountParams{Coupon: stripe.String(couponID)},
		}
	} else {
		params.AllowPromotionCodes = stripe.Bool(true)
	}

	s, err := session.New(params)
	if err != nil {
		return "", err
	}
	return s.URL, nil
}

func (pm *stripePaymentProcessor) CreateOneTimeCheckoutSessionURL(domain, successCallbackURL, canceledCallbackURL, customerID, priceID string, customerHasShippingAddress bool, couponID string) (string, error) {
	return pm.createCheckoutSessionURL(domain, successCallbackURL, canceledCallbackURL, stripe.CheckoutSessionModePayment, customerID, priceID, 1, customerHasShippingAddress, couponID)
}

func (pm *stripePaymentProcessor) CreateSubscriptionCheckoutSessionURL(domain, successCallbackURL, canceledCallbackURL, customerID, priceID string, customerHasShippingAddress bool, couponID string) (string, error) {
	return pm.createCheckoutSessionURL(domain, successCallbackURL, canceledCallbackURL, stripe.CheckoutSessionModeSubscription, customerID, priceID, 1, customerHasShippingAddress, couponID)
}

func (pm *stripePaymentProcessor) GetCheckoutSession(sessionID string) (*stripe.CheckoutSession, error) {
	s, err := session.Get(sessionID, nil)
	if err != nil {
		return nil, err
	}
	return s, nil
}

func (pm *stripePaymentProcessor) GetCheckoutSessionLineItems(sessionID string) ([]*stripe.LineItem, error) {
	params := &stripe.CheckoutSessionParams{}
	params.AddExpand("line_items")

	sessionWithLineItems, err := session.Get(sessionID, params)
	if err != nil {
		return nil, err
	}
	lineItems := sessionWithLineItems.LineItems
	return lineItems.Data, nil
}

func (pm *stripePaymentProcessor) GetPaymentIntent(paymentIntentID string) (*stripe.PaymentIntent, error) {
	s, err := paymentintent.Get(paymentIntentID, nil)
	if err != nil {
		return nil, err
	}
	return s, nil
}

func (pm *stripePaymentProcessor) GetCustomer(customerID string) (*stripe.Customer, error) {
	s, err := customer.Get(customerID, nil)
	if err != nil {
		return nil, err
	}
	return s, nil
}

func (pm *stripePaymentProcessor) GetSubscription(subscriptionID string) (*stripe.Subscription, error) {
	s, err := sub.Get(subscriptionID, nil)
	if err != nil {
		return nil, err
	}
	return s, nil
}
func (pm *stripePaymentProcessor) CancelSubscription(subscriptionID string) (*stripe.Subscription, error) {
	s, err := sub.Cancel(subscriptionID, nil)
	if err != nil {
		return nil, err
	}
	return s, nil
}

func (pm *stripePaymentProcessor) ListInvoicesByCustomerID(customerID string) ([]*stripe.Invoice, error) {
	params := &stripe.InvoiceListParams{}
	params.Filters.AddFilter("customer", "", customerID)
	params.Filters.AddFilter("limit", "", "3")
	i := invoice.List(params)

	ii := []*stripe.Invoice{}
	for i.Next() {
		in := i.Invoice()
		ii = append(ii, in)
	}
	return ii, nil
}

func (pm *stripePaymentProcessor) GetLatestInvoiceByCustomerID(customerID string) (*stripe.Invoice, error) {
	params := &stripe.InvoiceListParams{}
	params.Filters.AddFilter("customer", "", customerID)
	params.Filters.AddFilter("limit", "", "1")
	params.Filters.AddFilter("status", "", "paid")
	i := invoice.List(params)

	for i.Next() {
		in := i.Invoice()
		if in != nil {
			return in, nil
		}
	}
	return nil, nil
}
