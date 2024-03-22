package crontab

import (
	"context"
	"log/slog"
)

func (port *crontabInputPort) RankToday() {
	// port.Logger.Debug("crontab triggered ranking today")
	if err := port.RankPointController.GenerateGlobalRankingForTodayUsingActiveFitBitApps(context.Background()); err != nil {
		port.Logger.Error("crontab trigger failed ranking today",
			slog.Any("error", err))
		return
	}
}

func (port *crontabInputPort) RankThisISOWeek() {
	// port.Logger.Debug("crontab triggered ranking for this iso weeks")
	if err := port.RankPointController.GenerateGlobalRankingForThisISOWeekUsingActiveFitBitApps(context.Background()); err != nil {
		port.Logger.Error("crontab trigger failed ranking this iso weeks",
			slog.Any("error", err))
		return
	}
}

func (port *crontabInputPort) RankThisMonth() {
	// port.Logger.Debug("crontab triggered ranking for this month")
	if err := port.RankPointController.GenerateGlobalRankingForThisMonthUsingActiveFitBitApps(context.Background()); err != nil {
		port.Logger.Error("crontab trigger failed ranking this month",
			slog.Any("error", err))
		return
	}
}

func (port *crontabInputPort) RankThisYear() {
	// port.Logger.Debug("crontab triggered ranking for this year")
	if err := port.RankPointController.GenerateGlobalRankingForThisYearUsingActiveFitBitApps(context.Background()); err != nil {
		port.Logger.Error("crontab trigger failed ranking this year",
			slog.Any("error", err))
		return
	}
}