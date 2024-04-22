package scheduler

import (
	"log/slog"
)

func (impl *aggregatePointSchedulerImpl) RunEveryMinuteDeleteAllAnomalousData() error {
	impl.Logger.Debug("scheduled: delete all anomalous data", slog.String("interval", "every minute"))
	err := impl.DistributedScheduler.ScheduleEveryMinuteFunc(func() {
		impl.Logger.Debug("running delete all anomalous data...")
		impl.Logger.Debug("finished running delete all anomalous data")
	})
	if err != nil {
		impl.Logger.Error("error with pinging scheduler", slog.Any("err", err))
	}
	return nil
}

func (impl *aggregatePointSchedulerImpl) deleteAllAnomalousData() error {

	return nil
}
