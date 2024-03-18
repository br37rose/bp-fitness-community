package google

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"strconv"
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

	NewHTTPClientFromToken(token *oauth2.Token) *http.Client
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

func (gcp *gcpAdapter) OAuth2GenerateAuthCodeURL(oauthState string) string {
	return gcp.GoogleOauthConfig.AuthCodeURL(oauthState)
}

func (gcp *gcpAdapter) OAuth2ExchangeCode(code string) (*oauth2.Token, error) {
	return gcp.GoogleOauthConfig.Exchange(context.Background(), code)
}

// NewHTTPClientFromToken function returns an HTTP client using the provided
// token. The token will auto-refresh as necessary. The underlying
// HTTP transport will be obtained using the provided context.
// The returned client and its Transport should not be modified.
func (gcp *gcpAdapter) NewHTTPClientFromToken(token *oauth2.Token) *http.Client {
	return gcp.GoogleOauthConfig.Client(context.Background(), token)
}

func (gcp *gcpAdapter) NewFitnessStoreFromClient(client *http.Client) (*fitness.Service, error) {
	ctx := context.Background()
	return fitness.NewService(ctx, option.WithHTTPClient(client))
}

func (gcp *gcpAdapter) NotAggregatedDatasets(svc *fitness.Service, minTime, maxTime time.Time, dataType string) ([]*fitness.Dataset, error) {
	ds, err := svc.Users.DataSources.List("me").DataTypeName("com.google." + dataType).Do()
	if err != nil {
		log.Println("Unable to retrieve user's data sources:", err)
		return nil, err
	}
	if len(ds.DataSource) == 0 {
		log.Println("You have no data sources to explore.")
		return nil, err
	}

	var dataset []*fitness.Dataset

	for _, d := range ds.DataSource {
		setID := fmt.Sprintf("%v-%v", minTime.UnixNano(), maxTime.UnixNano())
		data, err := svc.Users.DataSources.Datasets.Get("me", d.DataStreamId, setID).Do()
		if err != nil {
			log.Println("Unable to retrieve dataset:", err.Error())
			return nil, err
		}
		dataset = append(dataset, data)
	}

	return dataset, nil

}

const (
	layout        = "Jan 2, 2006 at 3:04pm" // for time.Format
	nanosPerMilli = 1e6
) // Special thanks: https://github.com/googleapis/google-api-go-client/blob/main/examples/fitness.go#L18C1-L21C2

// millisToTime converts Unix millis to time.Time.
func millisToTime(t int64) time.Time {
	// Special thanks: https://github.com/googleapis/google-api-go-client/blob/main/examples/fitness.go#L36
	return time.Unix(0, t*nanosPerMilli)
}

type HydrationStruct struct {
	Amount    int       `json:"amount"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
}

func ParseHydration(datasets []*fitness.Dataset) []HydrationStruct {
	var data []HydrationStruct

	for _, ds := range datasets {
		var value float64
		for _, p := range ds.Point {
			for _, v := range p.Value {
				valueString := fmt.Sprintf("%.3f", v.FpVal)
				value, _ = strconv.ParseFloat(valueString, 64)
			}
			var row HydrationStruct
			row.StartTime = millisToTime(p.StartTimeNanos)
			row.EndTime = millisToTime(p.EndTimeNanos)
			// liters to milliliters
			row.Amount = int(value * 1000)
			data = append(data, row)
		}
	}
	return data
}
