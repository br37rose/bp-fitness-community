package crontab

import (
	"log/slog"
	"time"

	gcp_a "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/adapter/cloudprovider/google"
	googlefitapp_c "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/googlefitapp/controller"
	gfa_ds "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/googlefitapp/datastore"
	user_s "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/user/datastore"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/provider/kmutex"
)

type GoogleFitAppCrontaber interface {
	RefreshTokensFromGoogleJob() error
	PullDataFromGoogleJob() error
}

// Handler Creates http request handler
type googleFitAppCrontaberImpl struct {
	Logger                 *slog.Logger
	Kmutex                 kmutex.Provider
	GCP                    gcp_a.GoogleCloudPlatformAdapter
	GoogleFitAppStorer     gfa_ds.GoogleFitAppStorer
	GoogleFitAppController googlefitapp_c.GoogleFitAppController
	UserStorer             user_s.UserStorer
}

// NewHandler Constructor
func NewCrontab(
	loggerp *slog.Logger,
	kmutexp kmutex.Provider,
	gcpa gcp_a.GoogleCloudPlatformAdapter,
	gfa_storer gfa_ds.GoogleFitAppStorer,
	c googlefitapp_c.GoogleFitAppController,
	usr_storer user_s.UserStorer,
) GoogleFitAppCrontaber {
	return &googleFitAppCrontaberImpl{
		Logger:                 loggerp,
		Kmutex:                 kmutexp,
		GCP:                    gcpa,
		GoogleFitAppStorer:     gfa_storer,
		GoogleFitAppController: c,
		UserStorer:             usr_storer,
	}
}

// NanosToTime converts Unix nanos to time.Time. Special thanks to: https://github.com/bronnika/devto-google-fit/blob/main/google-api/init.go#L49
func NanosToTime(t int64) time.Time {
	return time.Unix(0, t)
}

// TimeToNanos coverts time.Time to Unix nanos. Special thanks to: https://github.com/bronnika/devto-google-fit/blob/main/google-api/init.go#L54
func TimeToNanos(time2 time.Time) int64 {
	return time2.UnixNano()
}
