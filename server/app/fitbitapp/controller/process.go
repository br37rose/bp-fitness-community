package controller

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/bartmika/timekit"
	"go.mongodb.org/mongo-driver/bson/primitive"

	dp_s "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/datapoint/datastore"
	fitbitdatum_s "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/fitbitdatum/datastore"
)

func (c *FitBitAppControllerImpl) ProcessAllQueuedData(ctx context.Context) error {
	c.Logger.Debug("processing all queued fitbit data")
	res, err := c.FitBitDatumStorer.ListByQueuedStatus(ctx)
	if err != nil {
		c.Logger.Error("failed listing by queued status",
			slog.Any("error", err))
		return err
	}
	for _, datum := range res.Results {
		if err := c.processQueuedDatum(ctx, datum); err != nil {
			c.Logger.Error("failed processing queued data for fitbit device",
				slog.Any("fitbit_app_id", datum.FitBitAppID),
				slog.Any("fitbit_datum_id", datum.ID),
				slog.Any("error", err))
			return err
		}
	}
	return nil
}

func (c *FitBitAppControllerImpl) processQueuedDatum(ctx context.Context, datum *fitbitdatum_s.FitBitDatum) error {
	switch datum.Type {
	case fitbitdatum_s.TypeActivitySteps:
		return c.processQueuedStepsDatum(ctx, datum)
	case fitbitdatum_s.TypeHeartRate:
		return c.processQueuedHeartRateDatum(ctx, datum)
	default:
		err := fmt.Errorf("queued datum is not supported for type %v", datum.Type)
		c.Logger.Error("failed processing queued data for fitbit device",
			slog.Any("fitbit_app_id", datum.FitBitAppID),
			slog.Any("fitbit_datum_id", datum.ID),
			slog.Any("error", err))
		return err
	}
}

func (c *FitBitAppControllerImpl) processQueuedStepsDatum(ctx context.Context, datum *fitbitdatum_s.FitBitDatum) error {
	c.Logger.Debug(fmt.Sprintf("processing queued fitbit steps datum #%s", datum.ID.Hex()),
		slog.Any("metric_id", datum.MetricID))

	for _, raw := range datum.Formatted.ActivitiesInterdayLog.ActivitiesStepsIntraday.Dataset {
		dur, err := timekit.ParseHourMinuteSecondDurationString(raw.Time)
		if err != nil {
			return err
		}

		// DEVELOPERS NOTE: The reason why we have to calculate the date-time
		// is because that's how the API is structured.
		midnight := timekit.Midnight(time.Now)
		dt := midnight.Add(dur)

		// Check to see if this was already submitted and if not then
		// proceed to submit.
		if exists, err := c.DataPointStorer.CheckIfExistsByCompositeKey(ctx, datum.MetricID, dt); err != nil || exists == true {
			if err != nil {
				return err
			}

			c.Logger.Warn(fmt.Sprintf("skipped processing datum #%s", datum.ID.Hex()),
				slog.Any("metric_id", datum.MetricID))
			return nil
		}

		// Generate our time-series datum.
		dp := &dp_s.DataPoint{
			ID:        primitive.NewObjectID(),
			MetricID:  datum.MetricID,
			Timestamp: dt,
			Value:     raw.Value,
		}
		if err := c.DataPointStorer.Create(ctx, dp); err != nil {
			return err
		}
		c.Logger.Warn(fmt.Sprintf("created data point #%s", dp.ID.Hex()),
			slog.Any("metric_id", datum.MetricID))

		// Update our datum that we datum has been processed.
		datum.Status = fitbitdatum_s.StatusActive
		datum.Errors = ""
		datum.ModifiedAt = time.Now()
		if err := c.FitBitDatumStorer.UpdateByID(ctx, datum); err != nil {
			return err
		}

		c.Logger.Warn(fmt.Sprintf("updated datum #%s", datum.ID.Hex()),
			slog.Any("metric_id", datum.MetricID))
	}

	fba, err := c.FitBitAppStorer.GetByID(ctx, datum.FitBitAppID)
	if err != nil {
		c.Logger.Error("failed getting fitbit app",
			slog.Any("fitbit_app_id", datum.FitBitAppID),
			slog.Any("fitbit_datum_id", datum.ID),
			slog.Any("error", err))
		return err
	}
	fba.LastFetchedAt = time.Now()
	if err := c.FitBitAppStorer.UpdateByID(ctx, fba); err != nil {
		return err
	}

	return nil
}

func (c *FitBitAppControllerImpl) processQueuedHeartRateDatum(ctx context.Context, datum *fitbitdatum_s.FitBitDatum) error {
	c.Logger.Debug(fmt.Sprintf("processing queued fitbit heart rate datum #%s", datum.ID.Hex()),
		slog.Any("metric_id", datum.MetricID))

	for _, raw := range datum.Formatted.HeartIntraday.ActivitiesHeartIntraday.Dataset {
		dur, err := timekit.ParseHourMinuteSecondDurationString(raw.Time)
		if err != nil {
			return err
		}

		// DEVELOPERS NOTE: The reason why we have to calculate the date-time
		// is because that's how the API is structured.
		midnight := timekit.Midnight(time.Now)
		dt := midnight.Add(dur)

		// Check to see if this was already submitted and if not then
		// proceed to submit.
		if exists, err := c.DataPointStorer.CheckIfExistsByCompositeKey(ctx, datum.MetricID, dt); err != nil || exists == true {
			if err != nil {
				return err
			}

			c.Logger.Warn(fmt.Sprintf("skipped processing datum #%s", datum.ID.Hex()),
				slog.Any("metric_id", datum.MetricID))
			return nil
		}

		// Generate our time-series datum.
		dp := &dp_s.DataPoint{
			ID:        primitive.NewObjectID(),
			MetricID:  datum.MetricID,
			Timestamp: dt,
			Value:     float64(raw.Value),
		}
		if err := c.DataPointStorer.Create(ctx, dp); err != nil {
			return err
		}
		c.Logger.Warn(fmt.Sprintf("created data point #%s", dp.ID.Hex()),
			slog.Any("metric_id", datum.MetricID))

		// Update our datum that we datum has been processed.
		datum.Status = fitbitdatum_s.StatusActive
		datum.Errors = ""
		datum.ModifiedAt = time.Now()
		if err := c.FitBitDatumStorer.UpdateByID(ctx, datum); err != nil {
			return err
		}

		c.Logger.Warn(fmt.Sprintf("updated datum #%s", datum.ID.Hex()),
			slog.Any("metric_id", datum.MetricID))
	}

	fba, err := c.FitBitAppStorer.GetByID(ctx, datum.FitBitAppID)
	if err != nil {
		c.Logger.Error("failed getting fitbit app",
			slog.Any("fitbit_app_id", datum.FitBitAppID),
			slog.Any("fitbit_datum_id", datum.ID),
			slog.Any("error", err))
		return err
	}
	fba.LastFetchedAt = time.Now()
	if err := c.FitBitAppStorer.UpdateByID(ctx, fba); err != nil {
		return err
	}

	return nil
}
