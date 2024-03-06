package controller

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log/slog"
	"net/http"
	"strings"
	"time"

	fb_api "github.com/Thomas2500/go-fitbit"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	fitbitapp_s "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/fitbitapp/datastore"
	fitbitdatum_s "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/fitbitdatum/datastore"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/config/constants"
)

func (c *FitBitAppControllerImpl) fetchHeartRateDataForFitBit(
	sessCtx mongo.SessionContext,
	fba *fitbitapp_s.FitBitApp,
	startAt time.Time,
	resource string,
	typeID int64,
) error {
	c.Logger.Debug(fmt.Sprintf("pulling fitbit device ID %v for resouce %v", fba.ID.Hex(), resource))
	defer c.Logger.Debug(fmt.Sprintf("pulled fitbit device ID %v for resouce %v", fba.ID.Hex(), resource))
	// https://cloud.fitbit.com/v2/docs#tag/Activity-Rate

	// startDateStr := timekit.ToISO8601String(startDate)
	// endDateStr := timekit.ToISO8601String(endDate)
	//
	url := fmt.Sprintf(constants.FitBitGetHeartRateIntradayByDateURL,
		fba.FitBitUserID, // [user-id]
		"today",
		"1min",
	)

	var bearer = "Bearer " + fba.AccessToken

	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)

	// add authorization header to the req
	req.Header.Add("Authorization", bearer)

	// Set data format.
	req.Header.Set("Content-Type", "application/json")

	// Send req using http Client
	resp, err := client.Do(req)
	if err != nil {
		c.Logger.Error("failed fetching from fitbit web-service",
			slog.Any("device_id", fba.ID),
			slog.Any("resource", resource),
			slog.Any("error", err))
		return err
	}

	defer resp.Body.Close()

	// Read the response body
	responseBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		c.Logger.Error("failed reading bytes",
			slog.Any("device_id", fba.ID),
			slog.Any("resource", resource),
			slog.Any("error", err))
		return err
	}

	responseStr := string(responseBytes)

	// If the token has expired then we will attempt to fetch the token again from
	// the FitBit API and then attempt to fetch the data again.
	if strings.Contains(responseStr, "Access token expired") {
		c.Logger.Warn("access token expired",
			slog.Any("device_id", fba.ID),
			slog.Any("resource", resource))

		if err := c.refreshAccessTokenForFitBit(sessCtx, fba); err != nil {
			return err
		}
		return c.fetchHeartRateDataForFitBit(sessCtx, fba, startAt, resource, typeID)
	}

	// Convert into the data-structure set by FitBit inc.
	content := &fb_api.HeartIntraday{}

	// Read the JSON string and convert it into our golang stuct else we get an error.
	if err := json.Unmarshal([]byte(responseStr), &content); err != nil {
		c.Logger.Error("failed unmarshalling steps fitbit data",
			slog.Any("fitbit_app_id", fba.ID),
			slog.Any("error", err))
	}

	formatted := fitbitdatum_s.FitBitFormattedDatum{
		HeartIntraday: content,
	}

	c.Logger.Debug("unmarshalled heart  datum",
		slog.Any("content", content))

	datum := &fitbitdatum_s.FitBitDatum{
		ID:             primitive.NewObjectID(),
		UserID:         fba.UserID,
		FitBitAppID:    fba.ID,
		Path:           url,
		StartAt:        fba.LastFetchedAt,
		EndAt:          time.Now(),
		Raw:            responseStr,
		Formatted:      formatted,
		Errors:         "",
		Type:           typeID,
		Status:         fitbitdatum_s.StatusQueued,
		CreatedAt:      time.Now(),
		ModifiedAt:     time.Now(),
		OrganizationID: fba.OrganizationID,
		MetricID:       fba.HeartRateMetricID,
	}

	if strings.Contains(datum.Raw, "errors") {
		datum.Raw = ""
		datum.Errors = responseStr
		datum.Status = fitbitdatum_s.StatusError
	}

	err = c.FitBitDatumStorer.Create(sessCtx, datum)
	if err != nil {
		// bc.Logger.Err(err).Caller().Str("taskID", taskID).
		// 	Str("taskName", "FitBitFetchData").
		// 	Msg("database error")
		return err
	}

	// // bc.Logger.Info().Str("taskID", taskID).
	// // 	Str("taskName", "FitBitFetchData").
	// // 	Str("func", "fetchActivityDataForFitBit").
	// // 	Msg("finished")
	return nil
}
