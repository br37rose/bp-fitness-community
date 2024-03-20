package google

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	mongo_client "go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/fitness/v1"
	"google.golang.org/api/option"

	c "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/config"
)

type GoogleCloudPlatformAdapter interface {
	Shutdown()
	OAuth2GenerateAuthCodeURL(oauthState string) string
	OAuth2ExchangeCode(code string) (*oauth2.Token, error)
	NewHTTPClientFromToken(token *oauth2.Token, callback func(*oauth2.Token)) (*http.Client, error)
	NewTokenFromExistingToken(token *oauth2.Token) (*oauth2.Token, error)
	NewFitnessStoreFromClient(client *http.Client) (*fitness.Service, error)
	NotAggregatedDatasets(svc *fitness.Service, minTime, maxTime time.Time, dataType string) ([]*fitness.Dataset, error)
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
			fitness.FitnessActivityReadScope,
			fitness.FitnessBloodGlucoseReadScope,
			fitness.FitnessBloodPressureReadScope,
			fitness.FitnessBodyReadScope,
			fitness.FitnessBodyTemperatureReadScope,
			fitness.FitnessHeartRateReadScope,
			fitness.FitnessLocationReadScope,
			fitness.FitnessNutritionReadScope,
			fitness.FitnessOxygenSaturationReadScope,
			fitness.FitnessSleepReadScope,
			// fitness.FitnessReproductiveHealthReadScope, // Not now...
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

func (gcp *gcpAdapter) OAuth2GenerateAuthCodeURL(oauthState string) string {
	// DEVELOPERS NOTE:
	// - `oauth2.AccessTypeOffline` is an important variable to set as it tells google to return a `refresh token` every time, else google will only return the `refresh token` once when our app gets first registered.
	// - If we do not use `oauth2.AccessTypeOffline` then the refresh token is provided to our application ONCE and never again! The only way we can get the refresh token is by the user deleting our app from their google profile and attempting again to register our app with their account.
	// - For more details see on how refresh tokens work see the following link via https://medium.com/starthinker/google-oauth-2-0-access-token-and-refresh-token-explained-cccf2fc0a6d9.
	// - See the following link for an example of refresh tokens working via https://github.com/kbehouse/oauth2/blob/master/google_offline_other_client.go.
	return gcp.GoogleOauthConfig.AuthCodeURL(oauthState, oauth2.AccessTypeOffline)
}

func (gcp *gcpAdapter) OAuth2ExchangeCode(code string) (*oauth2.Token, error) {
	return gcp.GoogleOauthConfig.Exchange(context.Background(), code)
}

// NewHTTPClientFromToken function returns an HTTP client using the provided
// token. The token will auto-refresh as necessary. The underlying
// HTTP transport will be obtained using the provided context.
// The returned client and its Transport should not be modified.
func (gcp *gcpAdapter) NewHTTPClientFromToken(token *oauth2.Token, callback func(*oauth2.Token)) (*http.Client, error) {

	tokenSource := gcp.GoogleOauthConfig.TokenSource(context.Background(), token)
	newToken, err := tokenSource.Token()
	if err != nil {
		return nil, err
	}

	if newToken.AccessToken != token.AccessToken {
		callback(newToken)
	}

	return oauth2.NewClient(context.Background(), tokenSource), nil
}

func (gcp *gcpAdapter) NewTokenFromExistingToken(token *oauth2.Token) (*oauth2.Token, error) {
	tokenSource := gcp.GoogleOauthConfig.TokenSource(context.Background(), token)
	return tokenSource.Token()
}

func (gcp *gcpAdapter) NewFitnessStoreFromClient(client *http.Client) (*fitness.Service, error) {
	ctx := context.Background()
	return fitness.NewService(ctx, option.WithHTTPClient(client))
}

// NotAggregatedDatasets function calls `Google Fit` to lookup data for specific type.
// Special thanks to:
// (1) https://github.com/bronnika/devto-google-fit/blob/a85098882047ff1e647d49905ddedd7e2425e31a/google-api/get.go#L198
// (2) https://dev.to/bronnika/working-with-google-fit-api-using-go-package-fitness-58bn
func (gcp *gcpAdapter) NotAggregatedDatasets(svc *fitness.Service, minTime, maxTime time.Time, dataType string) ([]*fitness.Dataset, error) {
	ds, err := svc.Users.DataSources.List("me").DataTypeName("com.google." + dataType).Do()
	if err != nil {
		// log.Println("Unable to retrieve user's data sources:", err)
		return nil, err
	}
	if len(ds.DataSource) == 0 {
		// log.Println("You have no data sources to explore.")
		return nil, err
	}

	var dataset []*fitness.Dataset

	for _, d := range ds.DataSource {
		setID := fmt.Sprintf("%v-%v", minTime.UnixNano(), maxTime.UnixNano())
		data, err := svc.Users.DataSources.Datasets.Get("me", d.DataStreamId, setID).Do()
		if err != nil {
			// log.Println("Unable to retrieve dataset:", err.Error())
			return nil, err
		}
		dataset = append(dataset, data)
	}

	return dataset, nil

}
