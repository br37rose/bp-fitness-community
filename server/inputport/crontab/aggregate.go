package crontab

import (
	"context"
	"log/slog"
)

func (port *crontabInputPort) AggregateThisHour() {
	// port.Logger.Debug("crontab triggered aggregating this hours fitbit data")
	if err := port.AggregatePointController.AggregateThisHourForAllActiveFitBitApps(context.Background()); err != nil {
		port.Logger.Error("crontab trigger failed aggregating this hours fitbit data",
			slog.Any("error", err))
		return
	}
}

func (port *crontabInputPort) AggregateLastHour() {
	// port.Logger.Debug("crontab triggered aggregating yesterdays fitbit data")
	if err := port.AggregatePointController.AggregateLastHourForAllActiveFitBitApps(context.Background()); err != nil {
		port.Logger.Error("crontab trigger failed aggregating last hours fitbit data",
			slog.Any("error", err))
		return
	}
}

func (port *crontabInputPort) AggregateToday() {
	// port.Logger.Debug("crontab triggered aggregating todays fitbit data")
	if err := port.AggregatePointController.AggregateTodayForAllActiveFitBitApps(context.Background()); err != nil {
		port.Logger.Error("crontab trigger failed aggregating todays fitbit data",
			slog.Any("error", err))
		return
	}
}

func (port *crontabInputPort) AggregateYesterday() {
	// port.Logger.Debug("crontab triggered aggregating yesterdays fitbit data")
	if err := port.AggregatePointController.AggregateYesterdayForAllActiveFitBitApps(context.Background()); err != nil {
		port.Logger.Error("crontab trigger failed aggregating yesterday fitbit data",
			slog.Any("error", err))
		return
	}
}

func (port *crontabInputPort) AggregateThisISOWeek() {
	// port.Logger.Debug("crontab triggered aggregating this iso weeks fitbit data")
	if err := port.AggregatePointController.AggregateThisISOWeekForAllActiveFitBitApps(context.Background()); err != nil {
		port.Logger.Error("crontab trigger failed aggregating this iso weeks fitbit data",
			slog.Any("error", err))
		return
	}
}

func (port *crontabInputPort) AggregateLastISOWeek() {
	// port.Logger.Debug("crontab triggered aggregating last iso weeks fitbit data")
	if err := port.AggregatePointController.AggregateLastISOWeekForAllActiveFitBitApps(context.Background()); err != nil {
		port.Logger.Error("crontab trigger failed aggregating last iso weeks fitbit data",
			slog.Any("error", err))
		return
	}
}

func (port *crontabInputPort) AggregateThisMonth() {
	// port.Logger.Debug("crontab triggered aggregating this month fitbit data")
	if err := port.AggregatePointController.AggregateThisMonthForAllActiveFitBitApps(context.Background()); err != nil {
		port.Logger.Error("crontab trigger failed aggregating this months fitbit data",
			slog.Any("error", err))
		return
	}
}

func (port *crontabInputPort) AggregateLastMonth() {
	// port.Logger.Debug("crontab triggered aggregating last month fitbit data")
	if err := port.AggregatePointController.AggregateLastMonthForAllActiveFitBitApps(context.Background()); err != nil {
		port.Logger.Error("crontab trigger failed aggregating last month fitbit data",
			slog.Any("error", err))
		return
	}
}

func (port *crontabInputPort) AggregateThisYear() {
	// port.Logger.Debug("crontab triggered aggregating this year fitbit data")
	if err := port.AggregatePointController.AggregateThisYearForAllActiveFitBitApps(context.Background()); err != nil {
		port.Logger.Error("crontab trigger failed aggregating this year fitbit data",
			slog.Any("error", err))
		return
	}
}

func (port *crontabInputPort) AggregateLastYear() {
	// port.Logger.Debug("crontab triggered aggregating last year fitbit data")
	if err := port.AggregatePointController.AggregateLastYearForAllActiveFitBitApps(context.Background()); err != nil {
		port.Logger.Error("crontab trigger failed aggregating last year fitbit data",
			slog.Any("error", err))
		return
	}
}
