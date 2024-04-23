package config

import (
	"log"
	"os"
	"strconv"
	"strings"
)

type Conf struct {
	AppServer           serverConf
	DB                  dbConfig
	Redis               redisConfig
	Cache               cacheConfig
	AWS                 awsConfig
	PDFBuilder          pdfBuilderConfig
	Emailer             mailgunConfig
	PaymentProcessor    paymentProcessorConfig
	AI                  openAIConfig
	GoogleCloudPlatform googleCloudPlatformAppConfig
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

type redisConfig struct {
	IsClusterConfiguration bool     `env:"BP8_BACKEND_REDIS_IS_CLUSTER_CONFIGURATION,required"`
	Addresses              []string `env:"BP8_BACKEND_REDIS_ADDRESSES,required"`
	Username               string   `env:"BP8_BACKEND_REDIS_USERNAME"`
	Password               string   `env:"BP8_BACKEND_REDIS_PASSWORD"`
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
	APIKey                 string `env:"BP8_BACKEND_OPENAI_API_KEY,required"`
	OrganizationKey        string `env:"BP8_BACKEND_OPENAI_ORGANIZATION_KEY,required"`
	FitnessPlanAssistantID string `env:"BP8_BACKEND_OPENAI_API_FITNESS_PLAN_ASSISTANT_ID,required"`
}

type googleCloudPlatformAppConfig struct {
	ClientID                 string `env:"BP8_BACKEND_GOOGLE_CLOUD_PLATFORM_CLIENT_ID,required"`
	ClientSecret             string `env:"BP8_BACKEND_GOOGLE_CLOUD_PLATFORM_CLIENT_SECRET,required"`
	AuthorizationRedirectURI string `env:"BP8_BACKEND_GOOGLE_CLOUD_PLATFORM_AUTHORIZATION_REDIRECT_URI,required"`
	SuccessRedirectURI       string `env:"BP8_BACKEND_GOOGLE_CLOUD_PLATFORM_SUCCESS_REDIRECT_URL,required"`
	// ProjectID                string `env:"BP8_BACKEND_GOOGLE_CLOUD_PLATFORM_PROJECT_ID,required"`
	// AuthURI                  string `env:"BP8_BACKEND_GOOGLE_CLOUD_PLATFORM_AUTH_URI,required"`
	// TokenURI                 string `env:"BP8_BACKEND_GOOGLE_CLOUD_PLATFORM_TOKEN_URI,required"`
	// AuthProviderX509CertURL  string `env:"BP8_BACKEND_GOOGLE_CLOUD_PLATFORM_AUTH_PROVIDER_X509_CERT_URL,required"`
}

func New() *Conf {
	var c Conf

	// Server configuration
	c.AppServer.Port = getEnvString("BP8_BACKEND_PORT", true)
	c.AppServer.IP = getEnvString("BP8_BACKEND_IP", false)
	c.AppServer.HMACSecret = []byte(getEnvString("BP8_BACKEND_HMAC_SECRET", true))
	c.AppServer.HasDebugging = getEnvBool("BP8_BACKEND_HAS_DEBUGGING", false, true)
	c.AppServer.IsDeveloperMode = getEnvBool("BP8_BACKEND_IS_DEVELOPER_MODE", false, true)
	c.AppServer.InitialRootAdminID = getEnvString("BP8_BACKEND_INITIAL_ROOT_ADMIN_ID", true)
	c.AppServer.InitialRootAdminEmail = getEnvString("BP8_BACKEND_INITIAL_ROOT_ADMIN_EMAIL", true)
	c.AppServer.InitialRootAdminPassword = getEnvString("BP8_BACKEND_INITIAL_ROOT_ADMIN_PASSWORD", true)
	c.AppServer.InitialOrgID = getEnvString("BP8_BACKEND_INITIAL_ORG_ID", true)
	c.AppServer.InitialOrgName = getEnvString("BP8_BACKEND_INITIAL_ORG_NAME", true)
	c.AppServer.InitialOrgBranchID = getEnvString("BP8_BACKEND_INITIAL_ORG_BRANCH_ID", true)
	c.AppServer.InitialOrgAdminEmail = getEnvString("BP8_BACKEND_INITIAL_ORG_ADMIN_EMAIL", true)
	c.AppServer.InitialOrgAdminPassword = getEnvString("BP8_BACKEND_INITIAL_ORG_ADMIN_PASSWORD", true)
	c.AppServer.InitialOrgTrainerEmail = getEnvString("BP8_BACKEND_INITIAL_ORG_TRAINER_EMAIL", true)
	c.AppServer.InitialOrgTrainerPassword = getEnvString("BP8_BACKEND_INITIAL_ORG_TRAINER_PASSWORD", true)
	c.AppServer.InitialOrgMemberEmail = getEnvString("BP8_BACKEND_INITIAL_ORG_MEMBER_EMAIL", true)
	c.AppServer.InitialOrgMemberPassword = getEnvString("BP8_BACKEND_INITIAL_ORG_MEMBER_PASSWORD", true)
	c.AppServer.APIDomainName = getEnvString("BP8_BACKEND_API_DOMAIN_NAME", true)
	c.AppServer.AppDomainName = getEnvString("BP8_BACKEND_APP_DOMAIN_NAME", true)
	c.AppServer.Enable2FAOnRegistration = getEnvBool("BP8_BACKEND_APP_ENABLE_2FA_ON_REGISTRATION", false, false)

	// Database configuration
	c.DB.URI = getEnvString("BP8_BACKEND_DB_URI", true)
	c.DB.Name = getEnvString("BP8_BACKEND_DB_NAME", true)

	// Redis configuration
	c.Redis.IsClusterConfiguration = getEnvBool("BP8_BACKEND_REDIS_IS_CLUSTER_CONFIGURATION", true, false)
	c.Redis.Addresses = getEnvStrings("BP8_BACKEND_REDIS_ADDRESSES", true)
	c.Redis.Username = getEnvString("BP8_BACKEND_REDIS_USERNAME", false)
	c.Redis.Password = getEnvString("BP8_BACKEND_REDIS_PASSWORD", false)

	// Cache configuration
	c.Cache.URI = getEnvString("BP8_BACKEND_CACHE_URI", true)

	// AWS configuration
	c.AWS.AccessKey = getEnvString("BP8_BACKEND_AWS_ACCESS_KEY", true)
	c.AWS.SecretKey = getEnvString("BP8_BACKEND_AWS_SECRET_KEY", true)
	c.AWS.Endpoint = getEnvString("BP8_BACKEND_AWS_ENDPOINT", true)
	c.AWS.Region = getEnvString("BP8_BACKEND_AWS_REGION", true)
	c.AWS.BucketName = getEnvString("BP8_BACKEND_AWS_BUCKET_NAME", true)

	// PDF builder configuration
	c.PDFBuilder.CBFFTemplatePath = getEnvString("BP8_BACKEND_PDF_BUILDER_CBFF_TEMPLATE_FILE_PATH", true)
	c.PDFBuilder.DataDirectoryPath = getEnvString("BP8_BACKEND_PDF_BUILDER_DATA_DIRECTORY_PATH", true)

	// Mailgun configuration
	c.Emailer.APIKey = getEnvString("BP8_BACKEND_MAILGUN_API_KEY", true)
	c.Emailer.Domain = getEnvString("BP8_BACKEND_MAILGUN_DOMAIN", true)
	c.Emailer.APIBase = getEnvString("BP8_BACKEND_MAILGUN_API_BASE", true)
	c.Emailer.SenderEmail = getEnvString("BP8_BACKEND_MAILGUN_SENDER_EMAIL", true)
	c.Emailer.MaintenanceEmail = getEnvString("BP8_BACKEND_MAILGUN_MAINTENANCE_EMAIL", true)

	// Payment processor configuration
	c.PaymentProcessor.SecretKey = getEnvString("BP8_BACKEND_PAYMENT_PROCESSOR_SECRET_KEY", true)
	c.PaymentProcessor.PublicKey = getEnvString("BP8_BACKEND_PAYMENT_PROCESSOR_PUBLIC_KEY", true)
	c.PaymentProcessor.WebhookSecretKey = getEnvString("BP8_BACKEND_PAYMENT_PROCESSOR_WEBHOOK_SECRET_KEY", true)

	// OpenAI configuration
	c.AI.APIKey = getEnvString("BP8_BACKEND_OPENAI_API_KEY", true)
	c.AI.OrganizationKey = getEnvString("BP8_BACKEND_OPENAI_ORGANIZATION_KEY", true)
	c.AI.FitnessPlanAssistantID = getEnvString("BP8_BACKEND_OPENAI_API_FITNESS_PLAN_ASSISTANT_ID", true)

	// Google Cloud Platform application configuration
	c.GoogleCloudPlatform.ClientID = getEnvString("BP8_BACKEND_GOOGLE_CLOUD_PLATFORM_CLIENT_ID", true)
	c.GoogleCloudPlatform.ClientSecret = getEnvString("BP8_BACKEND_GOOGLE_CLOUD_PLATFORM_CLIENT_SECRET", true)
	c.GoogleCloudPlatform.AuthorizationRedirectURI = getEnvString("BP8_BACKEND_GOOGLE_CLOUD_PLATFORM_AUTHORIZATION_REDIRECT_URI", true)
	c.GoogleCloudPlatform.SuccessRedirectURI = getEnvString("BP8_BACKEND_GOOGLE_CLOUD_PLATFORM_SUCCESS_REDIRECT_URL", true)

	return &c
}

func getEnvString(key string, required bool) string {
	value := os.Getenv(key)
	if required && value == "" {
		log.Fatalf("Environment variable not found: %s", key)
	}
	return value
}

func getEnvStrings(key string, required bool) []string {
	value := os.Getenv(key)
	if required && value == "" {
		log.Fatalf("Environment variable not found: %s", key)
	}

	// Split the comma-separated values into a slice
	values := strings.Split(value, ",")

	// Remove any leading or trailing whitespaces from each value
	for i := range values {
		values[i] = strings.TrimSpace(values[i])
	}

	return values
}

func getEnvBool(key string, required bool, defaultValue bool) bool {
	valueStr := getEnvString(key, required)
	if valueStr == "" {
		return defaultValue
	}
	value, err := strconv.ParseBool(valueStr)
	if err != nil {
		log.Fatalf("Invalid boolean value for environment variable %s", key)
	}
	return value
}

func getEnvInt64(key string, required bool) int64 {
	valueStr := os.Getenv(key)
	if required && valueStr == "" {
		log.Fatalf("Environment variable not found: %s", key)
	}

	// Convert string value to int64
	valueInt64, err := strconv.ParseInt(valueStr, 10, 64)
	if err != nil {
		log.Fatalf("Error converting %s to int64: %v", valueStr, err)
	}

	return valueInt64
}

func getEnvInt64s(key string, required bool) []int64 {
	value := os.Getenv(key)
	if required && value == "" {
		log.Fatalf("Environment variable not found: %s", key)
	}

	// Split the comma-separated values into a slice
	values := strings.Split(value, ",")

	// Convert string values to int64
	var intValues []int64
	for _, v := range values {
		trimmedValue := strings.TrimSpace(v)
		intValue, err := strconv.ParseInt(trimmedValue, 10, 64)
		if err != nil {
			log.Fatalf("Error converting %s to int64: %v", trimmedValue, err)
		}
		intValues = append(intValues, intValue)
	}

	return intValues
}
