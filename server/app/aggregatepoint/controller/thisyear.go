package controller

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/bartmika/timekit"
	ap_s "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/aggregatepoint/datastore"
	gfa_ds "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/googlefitapp/datastore"
)

func (impl *AggregatePointControllerImpl) AggregateThisYearForAllActiveGoogleFitApps(ctx context.Context) error {
	res, err := impl.GoogleFitAppStorer.ListIDsByStatus(ctx, gfa_ds.StatusActive)
	if err != nil {
		impl.Logger.Error("failed listing by active status",
			slog.Any("error", err))
		return err
	}
	for _, gfaID := range res {
		// Lock this Google Fit App for modification.
		impl.Kmutex.Lockf("gfa_%v", gfaID.Hex())
		defer impl.Kmutex.Unlockf("gfa_%v", gfaID.Hex())

		gfa, err := impl.GoogleFitAppStorer.GetByID(ctx, gfaID)
		if err != nil {
			impl.Logger.Error("failed getting google fit app by id",
				slog.Any("google_fit_app_id", gfaID),
				slog.Any("error", err))
			return err
		}
		if gfa == nil {
			err := fmt.Errorf("google fit app does not exist for id %v", gfaID.Hex())
			impl.Logger.Error("google fit app does not exist",
				slog.Any("google_fit_app_id", gfaID),
				slog.Any("error", err))
			return err
		}

		start := timekit.FirstDayOfThisYear(time.Now)
		end := timekit.FirstDayOfNextYear(time.Now)

		// impl.Logger.Debug("aggregate this year",
		// 	slog.Any("start", start),
		// 	slog.Any("end", end))

		if err := impl.aggregateForMetric(ctx, gfa.HeartRateBPMMetricID, ap_s.PeriodYear, start, end); err != nil {
			impl.Logger.Error("failed aggregating",
				slog.Any("google_fit_app_id", gfaID),
				slog.Any("error", err))
		}
	}
	return nil
}
