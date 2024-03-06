package crontab

import (
	"context"
	"log/slog"
)

func (port *crontabInputPort) pullFitBitAppRawData() {
	port.Logger.Debug("crontab trigger pulling latest data from fitbit web-service")
	if err := port.FitBitAppController.PullAllActiveDevices(context.Background()); err != nil {
		port.Logger.Error("crontab trigger failed pulling latest data from fitbit web-service",
			slog.Any("error", err))
		return
	}
}

func (port *crontabInputPort) processAllQueuedData() {
	port.Logger.Debug("crontab trigger processing queued fitbit data")
	if err := port.FitBitAppController.ProcessAllQueuedData(context.Background()); err != nil {
		port.Logger.Error("crontab trigger failed processing queued fitbit data",
			slog.Any("error", err))
		return
	}
}

func (port *crontabInputPort) processAllActiveSimulators() {
	port.Logger.Debug("crontab trigger processing active fitbit simulators")
	if err := port.FitBitAppController.ProcessAllActiveSimulators(context.Background()); err != nil {
		port.Logger.Error("crontab trigger failed processing active fitbit simulators",
			slog.Any("error", err))
		return
	}
}
