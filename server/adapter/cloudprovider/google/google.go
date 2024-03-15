package google

import (
	"context"
	"log/slog"

	mongo_client "go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	c "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/config"
)

type GoogleCloudPlatformAdapter interface {
	Shutdown()
	GenerateAuthCodeURL(oauthState string) string
	ExchangeCode(code string) (*oauth2.Token, error)
}

type gcpAdapter struct {
	Logger            *slog.Logger
	GoogleOauthConfig oauth2.Config
}

func NewAdapter(cfg *c.Conf, logger *slog.Logger, dbClient *mongo_client.Client) GoogleCloudPlatformAdapter {
	logger.Debug("google cloud platform connecting...")

	googleLoginConfig := oauth2.Config{
		RedirectURL:  cfg.GoogleCloudPlatform.AuthorizationRedirectURI,
		ClientID:     cfg.GoogleCloudPlatform.ClientID,
		ClientSecret: cfg.GoogleCloudPlatform.ClientSecret,
		Scopes: []string{
			"https://www.googleapis.com/auth/fitness.activity.read",
			"https://www.googleapis.com/auth/fitness.blood_glucose.read",
			"https://www.googleapis.com/auth/fitness.blood_pressure.read",
			"https://www.googleapis.com/auth/fitness.body.read",
			"https://www.googleapis.com/auth/fitness.heart_rate.read",
			"https://www.googleapis.com/auth/fitness.body_temperature.read",
			"https://www.googleapis.com/auth/fitness.location.read",
			"https://www.googleapis.com/auth/fitness.nutrition.read",
			"https://www.googleapis.com/auth/fitness.oxygen_saturation.read",
			"https://www.googleapis.com/auth/fitness.sleep.read",
		},
		Endpoint: google.Endpoint,
	}

	logger.Debug("connected with google cloud platform")
	return &gcpAdapter{
		Logger:            logger,
		GoogleOauthConfig: googleLoginConfig,
	}
}

func (gcp *gcpAdapter) Shutdown() {
	// Do nothing...
}

func (gcp *gcpAdapter) GenerateAuthCodeURL(oauthState string) string {
	return gcp.GoogleOauthConfig.AuthCodeURL(oauthState)
}

func (gcp *gcpAdapter) ExchangeCode(code string) (*oauth2.Token, error) {
	return gcp.GoogleOauthConfig.Exchange(context.Background(), code)
}
