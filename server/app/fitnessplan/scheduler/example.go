package scheduler

import (
	"log/slog"
)

// THIS IS AN EXAMPLE - PLEASE WRITE YOUR OWN FUNCTIONS AS YOU NEED USING THIS
// EXAMPLE TO HELP GUIDE YOU IN CODE IMPLEMENTATION.
func (impl *fitnessPlanSchedulerImpl) RunEveryMinuteDeleteAllAnomalousData() error {
	impl.Logger.Debug("scheduled: delete all anomalous data", slog.String("interval", "every minute"))
	err := impl.DistributedScheduler.ScheduleEveryMinuteFunc(func() {
		impl.Logger.Debug("running delete all anomalous data...")

		// Add your code here...

		impl.Logger.Debug("finished running delete all anomalous data")
	})
	if err != nil {
		impl.Logger.Error("error with pinging scheduler", slog.Any("err", err))
	}
	return nil
}

func (impl *fitnessPlanSchedulerImpl) deleteAllAnomalousData() error {

	return nil
}
