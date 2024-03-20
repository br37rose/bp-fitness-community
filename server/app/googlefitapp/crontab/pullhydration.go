package crontab

import (
	"context"
	"fmt"
	"log/slog"
	"strconv"
	"time"

	"google.golang.org/api/fitness/v1"

	gfa_ds "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/googlefitapp/datastore"
)

// HydrationStruct defines the hydration data type provided by `Google Fit`. Special thanks to https://github.com/bronnika/devto-google-fit/blob/main/models/models.go#L52C1-L56C2
type HydrationStruct struct {
	Amount    int       `json:"amount"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
}

// parseHydration function converts the `Google Fit` hydration data into usable format for our app. Special thanks to: https://github.com/bronnika/devto-google-fit/blob/main/google-api/parse.go#L103
func parseHydration(datasets []*fitness.Dataset) []HydrationStruct {
	var data []HydrationStruct

	for _, ds := range datasets {
		var value float64
		for _, p := range ds.Point {
			for _, v := range p.Value {
				valueString := fmt.Sprintf("%.3f", v.FpVal)
				value, _ = strconv.ParseFloat(valueString, 64)
			}
			var row HydrationStruct
			row.StartTime = NanosToTime(p.StartTimeNanos)
			row.EndTime = NanosToTime(p.EndTimeNanos)
			// liters to milliliters
			row.Amount = int(value * 1000)
			data = append(data, row)
		}
	}
	return data
}

func (impl *googleFitAppCrontaberImpl) pullHydrationDataFromGoogleWithGfaAndFitnessStore(ctx context.Context, gfa *gfa_ds.GoogleFitApp, svc *fitness.Service) error {
	impl.Logger.Debug("pulling hydration data",
		slog.String("gfa_id", gfa.ID.Hex()))

	////
	//// Get `Google Fit` data
	////

	maxTime := time.Now()
	minTime := gfa.LastFetchedAt
	dataType := "hydration" // Note: This is a `Google Fit` specific constant.
	dataset, err := impl.GCP.NotAggregatedDatasets(svc, minTime, maxTime, dataType)
	if err != nil {
		impl.Logger.Error("failed listing hydration",
			slog.Any("error", err))
		return err
	}

	impl.Logger.Debug("pulled hydration data",
		slog.String("gfa_id", gfa.ID.Hex()))

	if len(dataset) == 0 {
		return nil
	}

	////
	//// Convert from `Google Fit` format into our apps format.
	////

	hydration := parseHydration(dataset)
	impl.Logger.Debug("parsed hydration data",
		slog.Any("data", hydration))

	////
	//// Save into our database.
	////

	return nil
}
