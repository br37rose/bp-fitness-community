package controller

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/bartmika/timekit"

	ap_s "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/aggregatepoint/datastore"
	fba_s "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/fitbitapp/datastore"
)

func (impl *AggregatePointControllerImpl) AggregateThisYearForAllActiveFitBitApps(ctx context.Context) error {
	res, err := impl.FitBitAppStorer.ListIDsByStatus(ctx, fba_s.StatusActive)
	if err != nil {
		impl.Logger.Error("failed listing by active status",
			slog.Any("error", err))
		return err
	}
	for _, fbaID := range res {
		// Lock this fitbit device for modification.
		impl.Kmutex.Lockf("fitbitapp_%v", fbaID.Hex())
		defer impl.Kmutex.Unlockf("fitbitapp_%v", fbaID.Hex())

		fba, err := impl.FitBitAppStorer.GetByID(ctx, fbaID)
		if err != nil {
			impl.Logger.Error("failed getting fitbit by id",
				slog.Any("fitbit_app_id", fbaID),
				slog.Any("error", err))
			return err
		}
		if fba == nil {
			err := fmt.Errorf("fitbit does not exist for id %v", fbaID.Hex())
			impl.Logger.Error("fitbit does not exist",
				slog.Any("fitbit_app_id", fbaID),
				slog.Any("error", err))
			return err
		}

		start := timekit.FirstDayOfThisYear(time.Now)
		end := timekit.FirstDayOfNextYear(time.Now)

		if err := impl.aggregateForMetric(ctx, fba.HeartRateMetricID, ap_s.PeriodYear, start, end); err != nil {
			impl.Logger.Error("failed aggregating",
				slog.Any("fitbit_app_id", fbaID),
				slog.Any("error", err))
		}
		if err := impl.aggregateForMetric(ctx, fba.StepsCountMetricID, ap_s.PeriodYear, start, end); err != nil {
			impl.Logger.Error("failed aggregating",
				slog.Any("fitbit_app_id", fbaID),
				slog.Any("error", err))
		}
	}
	return nil
}
