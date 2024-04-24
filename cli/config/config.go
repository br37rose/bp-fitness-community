package config

import (
	"log"

	"github.com/joeshaw/envdecode"
)

type Conf struct {
	AppServer serverConf
	DB        dbConfig
	// Cache            cacheConfig
	AWS awsConfig
	// PDFBuilder       pdfBuilderConfig
	// Emailer          mailgunConfig
	// PaymentProcessor paymentProcessorConfig
	Ngrok ngrokConfig
}

type serverConf struct {
	// Port                     string `env:"BP8_BACKEND_PORT,required"`
	// IP                       string `env:"BP8_BACKEND_IP"`
	// HMACSecret               []byte `env:"BP8_BACKEND_HMAC_SECRET,required"`
	// HasDebugging             bool   `env:"BP8_BACKEND_HAS_DEBUGGING,default=true"`
	InitialRootAdminID       string `env:"BP8_BACKEND_INITIAL_ROOT_ADMIN_ID,required"`
	InitialRootAdminEmail    string `env:"BP8_BACKEND_INITIAL_ROOT_ADMIN_EMAIL,required"`
	InitialRootAdminPassword string `env:"BP8_BACKEND_INITIAL_ROOT_ADMIN_PASSWORD,required"`
	InitialOrgID             string `env:"BP8_BACKEND_INITIAL_ORG_ID,required"`
	InitialOrgName           string `env:"BP8_BACKEND_INITIAL_ORG_NAME,required"`
	InitialOrgBranchID       string `env:"BP8_BACKEND_INITIAL_ORG_BRANCH_ID,required"`
	InitialOrgAdminEmail     string `env:"BP8_BACKEND_INITIAL_ORG_ADMIN_EMAIL,required"`
	InitialOrgAdminPassword  string `env:"BP8_BACKEND_INITIAL_ORG_ADMIN_PASSWORD,required"`
	// InitialOrgTrainerEmail    string `env:"BP8_BACKEND_INITIAL_ORG_TRAINER_EMAIL,required"`
	// InitialOrgTrainerPassword string `env:"BP8_BACKEND_INITIAL_ORG_TRAINER_PASSWORD,required"`
	// InitialOrgMemberEmail     string `env:"BP8_BACKEND_INITIAL_ORG_MEMBER_EMAIL,required"`
	// InitialOrgMemberPassword  string `env:"BP8_BACKEND_INITIAL_ORG_MEMBER_PASSWORD,required"`
	// APIDomainName             string `env:"BP8_BACKEND_API_DOMAIN_NAME,required"`
	// AppDomainName             string `env:"BP8_BACKEND_APP_DOMAIN_NAME,required"`
}

type dbConfig struct {
	URI  string `env:"BP8_BACKEND_DB_URI,required"`
	Name string `env:"BP8_BACKEND_DB_NAME,required"`
}

type cacheConfig struct {
	URI string `env:"BP8_BACKEND_CACHE_URI,required"`
}

type awsConfig struct {
	AccessKey  string `env:"BP8_BACKEND_AWS_ACCESS_KEY,required"`
	SecretKey  string `env:"BP8_BACKEND_AWS_SECRET_KEY,required"`
	Endpoint   string `env:"BP8_BACKEND_AWS_ENDPOINT,required"`
	Region     string `env:"BP8_BACKEND_AWS_REGION,required"`
	BucketName string `env:"BP8_BACKEND_AWS_BUCKET_NAME,required"`
}

type pdfBuilderConfig struct {
	CBFFTemplatePath  string `env:"BP8_BACKEND_PDF_BUILDER_CBFF_TEMPLATE_FILE_PATH,required"`
	DataDirectoryPath string `env:"BP8_BACKEND_PDF_BUILDER_DATA_DIRECTORY_PATH,required"`
}

type mailgunConfig struct {
	APIKey           string `env:"BP8_BACKEND_MAILGUN_API_KEY,required"`
	Domain           string `env:"BP8_BACKEND_MAILGUN_DOMAIN,required"`
	APIBase          string `env:"BP8_BACKEND_MAILGUN_API_BASE,required"`
	SenderEmail      string `env:"BP8_BACKEND_MAILGUN_SENDER_EMAIL,required"`
	MaintenanceEmail string `env:"BP8_BACKEND_MAILGUN_MAINTENANCE_EMAIL,required"`
}

type paymentProcessorConfig struct {
	MonthlyMembershipProductID string `env:"BP8_BACKEND_PAYMENT_PROCESSOR_MONTHLY_MEMBERSHIP_PRODUCT_ID,required"`
	MonthlyMembershipPriceID   string `env:"BP8_BACKEND_PAYMENT_PROCESSOR_MONTHLY_MEMBERSHIP_PRICE_ID,required"`
	AnnualMembershipProductID  string `env:"BP8_BACKEND_PAYMENT_PROCESSOR_ANNUAL_MEMBERSHIP_PRODUCT_ID,required"`
	AnnualyMembershipPriceID   string `env:"BP8_BACKEND_PAYMENT_PROCESSOR_ANNUAL_MEMBERSHIP_PRICE_ID,required"`
	SecretKey                  string `env:"BP8_BACKEND_PAYMENT_PROCESSOR_SECRET_KEY,required"`
	PublicKey                  string `env:"BP8_BACKEND_PAYMENT_PROCESSOR_PUBLIC_KEY,required"`
	WebhookSecretKey           string `env:"BP8_BACKEND_PAYMENT_PROCESSOR_WEBHOOK_SECRET_KEY,required"`
}

type ngrokConfig struct {
	AuthToken string `env:"BP8_BACKEND_NGROK_AUTH_TOKEN,required"`
}

func New() *Conf {
	var c Conf
	if err := envdecode.StrictDecode(&c); err != nil {
		log.Fatalf("Failed to decode: %s", err)
	}
	return &c
}
