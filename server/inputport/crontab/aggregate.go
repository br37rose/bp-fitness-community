package crontab

import (
	"context"
	"log/slog"
)

func (port *crontabInputPort) AggregateThisHour() {
	// port.Logger.Debug("crontab triggered aggregating this hours google fit data")
	if err := port.AggregatePointController.AggregateThisHourForAllActiveGoogleFitApps(context.Background()); err != nil {
		port.Logger.Error("crontab trigger failed aggregating this hours google fit data",
			slog.Any("error", err))
		return
	}
}

func (port *crontabInputPort) AggregateLastHour() {
	// port.Logger.Debug("crontab triggered aggregating yesterdays google fit data")
	if err := port.AggregatePointController.AggregateLastHourForAllActiveGoogleFitApps(context.Background()); err != nil {
		port.Logger.Error("crontab trigger failed aggregating last hours google fit data",
			slog.Any("error", err))
		return
	}
}

func (port *crontabInputPort) AggregateToday() {
	// port.Logger.Debug("crontab triggered aggregating todays google fit data")
	if err := port.AggregatePointController.AggregateTodayForAllActiveGoogleFitApps(context.Background()); err != nil {
		port.Logger.Error("crontab trigger failed aggregating todays google fit data",
			slog.Any("error", err))
		return
	}
}

func (port *crontabInputPort) AggregateYesterday() {
	// port.Logger.Debug("crontab triggered aggregating yesterdays google fit data")
	if err := port.AggregatePointController.AggregateYesterdayForAllActiveGoogleFitApps(context.Background()); err != nil {
		port.Logger.Error("crontab trigger failed aggregating yesterday google fit data",
			slog.Any("error", err))
		return
	}
}

func (port *crontabInputPort) AggregateThisISOWeek() {
	// port.Logger.Debug("crontab triggered aggregating this iso weeks google fit data")
	if err := port.AggregatePointController.AggregateThisISOWeekForAllActiveGoogleFitApps(context.Background()); err != nil {
		port.Logger.Error("crontab trigger failed aggregating this iso weeks google fit data",
			slog.Any("error", err))
		return
	}
}

func (port *crontabInputPort) AggregateLastISOWeek() {
	// port.Logger.Debug("crontab triggered aggregating last iso weeks google fit data")
	if err := port.AggregatePointController.AggregateLastISOWeekForAllActiveGoogleFitApps(context.Background()); err != nil {
		port.Logger.Error("crontab trigger failed aggregating last iso weeks google fit data",
			slog.Any("error", err))
		return
	}
}

func (port *crontabInputPort) AggregateThisMonth() {
	// port.Logger.Debug("crontab triggered aggregating this month google fit data")
	if err := port.AggregatePointController.AggregateThisMonthForAllActiveGoogleFitApps(context.Background()); err != nil {
		port.Logger.Error("crontab trigger failed aggregating this months google fit data",
			slog.Any("error", err))
		return
	}
}

func (port *crontabInputPort) AggregateLastMonth() {
	// port.Logger.Debug("crontab triggered aggregating last month google fit data")
	if err := port.AggregatePointController.AggregateLastMonthForAllActiveGoogleFitApps(context.Background()); err != nil {
		port.Logger.Error("crontab trigger failed aggregating last month google fit data",
			slog.Any("error", err))
		return
	}
}

func (port *crontabInputPort) AggregateThisYear() {
	// port.Logger.Debug("crontab triggered aggregating this year google fit data")
	if err := port.AggregatePointController.AggregateThisYearForAllActiveGoogleFitApps(context.Background()); err != nil {
		port.Logger.Error("crontab trigger failed aggregating this year google fit data",
			slog.Any("error", err))
		return
	}
}

func (port *crontabInputPort) AggregateLastYear() {
	// port.Logger.Debug("crontab triggered aggregating last year google fit data")
	if err := port.AggregatePointController.AggregateLastYearForAllActiveGoogleFitApps(context.Background()); err != nil {
		port.Logger.Error("crontab trigger failed aggregating last year google fit data",
			slog.Any("error", err))
		return
	}
}
