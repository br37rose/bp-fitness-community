package scheduler

import (
	"context"
	"log/slog"
)

func (impl *googleFitAppSchedulerImpl) RunEveryMinuteDeleteAllAnomalousData() error {
	impl.Logger.Debug("scheduled every minute: delete all anomalous data")
	err := impl.DistributedScheduler.ScheduleEveryMinuteFunc(func() {
		impl.Logger.Debug("running delete all anomalous data...")
		impl.Logger.Debug("finished running delete all anomalous data")
	})
	if err != nil {
		impl.Logger.Error("error with pinging scheduler", slog.Any("err", err))
	}
	return nil
}

func (impl *googleFitAppSchedulerImpl) deleteAllAnomalousData() error {
	ctx := context.Background()
	res, err := impl.GoogleFitDataPointStorer.ListByAnomalousDetection(ctx)
	if err != nil {
		impl.Logger.Error("database error", slog.Any("error", err))
		return err
	}
	for _, dp := range res.Results {
		impl.Logger.Debug("deleted anomalous datapoint", slog.String("id", dp.ID.Hex()))
		if err := impl.GoogleFitDataPointStorer.DeleteByID(ctx, dp.ID); err != nil {
			impl.Logger.Error("database error", slog.Any("error", err))
			return err
		}
	}
	return nil
}
