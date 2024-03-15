package config

import (
	"log"

	"github.com/joeshaw/envdecode"
)

type Conf struct {
	AppServer        serverConf
	DB               dbConfig
	Cache            cacheConfig
	AWS              awsConfig
	PDFBuilder       pdfBuilderConfig
	Emailer          mailgunConfig
	PaymentProcessor paymentProcessorConfig
	AI               openAIConfig
	FitBitApp        fitbitAppConfig // Deprecated
	// OURARingApp  // Deprecated
	GoogleCloudPlatform googleCloudPlatformAppConfig // Deprecated
}

type serverConf struct {
	Port                      string `env:"BP8_BACKEND_PORT,required"`
	IP                        string `env:"BP8_BACKEND_IP"`
	HMACSecret                []byte `env:"BP8_BACKEND_HMAC_SECRET,required"`
	HasDebugging              bool   `env:"BP8_BACKEND_HAS_DEBUGGING,default=true"`
	IsDeveloperMode           bool   `env:"BP8_BACKEND_IS_DEVELOPER_MODE,default=true"`
	InitialRootAdminID        string `env:"BP8_BACKEND_INITIAL_ROOT_ADMIN_ID,required"`
	InitialRootAdminEmail     string `env:"BP8_BACKEND_INITIAL_ROOT_ADMIN_EMAIL,required"`
	InitialRootAdminPassword  string `env:"BP8_BACKEND_INITIAL_ROOT_ADMIN_PASSWORD,required"`
	InitialOrgID              string `env:"BP8_BACKEND_INITIAL_ORG_ID,required"`
	InitialOrgName            string `env:"BP8_BACKEND_INITIAL_ORG_NAME,required"`
	InitialOrgBranchID        string `env:"BP8_BACKEND_INITIAL_ORG_BRANCH_ID,required"`
	InitialOrgAdminEmail      string `env:"BP8_BACKEND_INITIAL_ORG_ADMIN_EMAIL,required"`
	InitialOrgAdminPassword   string `env:"BP8_BACKEND_INITIAL_ORG_ADMIN_PASSWORD,required"`
	InitialOrgTrainerEmail    string `env:"BP8_BACKEND_INITIAL_ORG_TRAINER_EMAIL,required"`
	InitialOrgTrainerPassword string `env:"BP8_BACKEND_INITIAL_ORG_TRAINER_PASSWORD,required"`
	InitialOrgMemberEmail     string `env:"BP8_BACKEND_INITIAL_ORG_MEMBER_EMAIL,required"`
	InitialOrgMemberPassword  string `env:"BP8_BACKEND_INITIAL_ORG_MEMBER_PASSWORD,required"`
	APIDomainName             string `env:"BP8_BACKEND_API_DOMAIN_NAME,required"`
	AppDomainName             string `env:"BP8_BACKEND_APP_DOMAIN_NAME,required"`
	Enable2FAOnRegistration   bool   `env:"BP8_BACKEND_APP_ENABLE_2FA_ON_REGISTRATION"`
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
	SecretKey        string `env:"BP8_BACKEND_PAYMENT_PROCESSOR_SECRET_KEY,required"`
	PublicKey        string `env:"BP8_BACKEND_PAYMENT_PROCESSOR_PUBLIC_KEY,required"`
	WebhookSecretKey string `env:"BP8_BACKEND_PAYMENT_PROCESSOR_WEBHOOK_SECRET_KEY,required"`
}

type openAIConfig struct {
	APIKey          string `env:"BP8_BACKEND_OPENAI_API_KEY,required"`
	OrganizationKey string `env:"BP8_BACKEND_OPENAI_ORGANIZATION_KEY,required"`
}

type fitbitAppConfig struct { // DEPRECATED
	ClientID                       string `env:"BP8_BACKEND_FITBIT_APP_CLIENT_ID,required"`
	ClientSecret                   string `env:"BP8_BACKEND_FITBIT_APP_CLIENT_SECRET,required"`
	SubscriberVerificationCode     string `env:"BP8_BACKEND_FITBIT_APP_SUBSCRIBER_VERIFICATION_CODE,required"`
	RegistrationSuccessRedirectURL string `env:"BP8_BACKEND_FITBIT_APP_REGISTRATION_SUCCESS_REDIRECT_URL,required"`
}

type ouraRingAppConfig struct { // DEPRECATED
	ClientID     string `env:"BP8_BACKEND_OURA_RING_APP_CLIENT_ID,required"`
	ClientSecret string `env:"BP8_BACKEND_OURA_RING_APP_CLIENT_SECRET,required"`
}

type googleCloudPlatformAppConfig struct {
	ClientID                 string `env:"BP8_BACKEND_GOOGLE_CLOUD_PLATFORM_CLIENT_ID,required"`
	ProjectID                string `env:"BP8_BACKEND_GOOGLE_CLOUD_PLATFORM_PROJECT_ID,required"`
	AuthURI                  string `env:"BP8_BACKEND_GOOGLE_CLOUD_PLATFORM_AUTH_URI,required"`
	TokenURI                 string `env:"BP8_BACKEND_GOOGLE_CLOUD_PLATFORM_TOKEN_URI,required"`
	AuthProviderX509CertURL  string `env:"BP8_BACKEND_GOOGLE_CLOUD_PLATFORM_AUTH_PROVIDER_X509_CERT_URL,required"`
	ClientSecret             string `env:"BP8_BACKEND_GOOGLE_CLOUD_PLATFORM_CLIENT_SECRET,required"`
	AuthorizationRedirectURI string `env:"BP8_BACKEND_GOOGLE_CLOUD_PLATFORM_AUTHORIZATION_REDIRECT_URI,required"`
	SuccessRedirectURI       string `env:"BP8_BACKEND_GOOGLE_CLOUD_PLATFORM_SUCCESS_REDIRECT_URL,required"`
}

func New() *Conf {
	var c Conf
	if err := envdecode.StrictDecode(&c); err != nil {
		log.Fatalf("Failed to decode: %s", err)
	}
	return &c
}
