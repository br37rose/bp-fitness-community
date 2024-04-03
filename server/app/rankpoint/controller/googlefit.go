package controller

import (
	"context"
	"log/slog"
	"time"

	"github.com/bartmika/timekit"
	"go.mongodb.org/mongo-driver/bson/primitive"

	gfa_ds "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/googlefitapp/datastore"
	rp_s "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/rankpoint/datastore"
)

func (impl *RankPointControllerImpl) GenerateGlobalRankingForTodayUsingActiveGoogleFitApps(ctx context.Context) error {
	f := &gfa_ds.GoogleFitAppListFilter{
		Cursor:    primitive.NilObjectID,
		PageSize:  1_000_000,
		SortField: "_id",
		SortOrder: -1,
		Status:    gfa_ds.StatusActive,
	}
	gfas, err := impl.GoogleFitAppStorer.ListByFilter(ctx, f)
	if err != nil {
		impl.Logger.Error("failed listing by active status",
			slog.Any("date_range", "today"),
			slog.Any("error", err))
		return err
	}
	start := timekit.Midnight(time.Now)
	end := timekit.MidnightTomorrow(time.Now)

	// --- Heart rate --- //

	go func(list []*gfa_ds.GoogleFitApp, startDT time.Time, endDT time.Time) {
		if err := impl.processGlobalRanksForGoogleFitApps(context.Background(), list, rp_s.MetricTypeHeartRate, rp_s.FunctionAverage, rp_s.PeriodDay, startDT, endDT); err != nil {
			impl.Logger.Error("failed generating global heart average rate ranking",
				slog.Any("date_range", "today"),
				slog.Any("error", err))
			return
		}
	}(gfas.Results, start, end)

	// TODO: Add more metric types below...

	// // --- Activity steps --- //
	//
	// go func(list []*gfa_ds.GoogleFitApp, startDT time.Time, endDT time.Time) {
	// 	if err := impl.processGlobalRanksForGoogleFitApps(context.Background(), list, rp_s.MetricTypeActivitySteps, rp_s.FunctionAverage, rp_s.PeriodDay, startDT, endDT); err != nil {
	// 		impl.Logger.Error("failed generating global activity steps average ranking",
	// 			slog.Any("date_range", "today"),
	// 			slog.Any("error", err))
	// 		return
	// 	}
	// }(gfas.Results, start, end)
	//
	// go func(list []*gfa_ds.GoogleFitApp, startDT time.Time, endDT time.Time) {
	// 	if err := impl.processGlobalRanksForGoogleFitApps(context.Background(), list, rp_s.MetricTypeActivitySteps, rp_s.FunctionSum, rp_s.PeriodDay, startDT, endDT); err != nil {
	// 		impl.Logger.Error("failed generating global activity steps sum ranking",
	// 			slog.Any("date_range", "today"),
	// 			slog.Any("error", err))
	// 		return
	// 	}
	// }(gfas.Results, start, end)
	//
	return nil
}

func (impl *RankPointControllerImpl) GenerateGlobalRankingForThisISOWeekUsingActiveGoogleFitApps(ctx context.Context) error {
	f := &gfa_ds.GoogleFitAppListFilter{
		Cursor:    primitive.NilObjectID,
		PageSize:  1_000_000,
		SortField: "_id",
		SortOrder: -1,
		Status:    gfa_ds.StatusActive,
	}
	gfas, err := impl.GoogleFitAppStorer.ListByFilter(ctx, f)
	if err != nil {
		impl.Logger.Error("failed listing by active status",
			slog.Any("date_range", "iso_week"),
			slog.Any("error", err))
		return err
	}
	start := timekit.FirstDayOfThisISOWeek(time.Now)
	end := timekit.FirstDayOfNextISOWeek(time.Now)

	// --- Heart rate --- //

	go func(list []*gfa_ds.GoogleFitApp, startDT time.Time, endDT time.Time) {
		if err := impl.processGlobalRanksForGoogleFitApps(context.Background(), list, rp_s.MetricTypeHeartRate, rp_s.FunctionAverage, rp_s.PeriodWeek, startDT, endDT); err != nil {
			impl.Logger.Error("failed generating global heart rate ranking",
				slog.Any("date_range", "iso_week"),
				slog.Any("error", err))
			return
		}
	}(gfas.Results, start, end)

	// // --- Activity steps --- //
	//
	// go func(list []*gfa_ds.GoogleFitApp, startDT time.Time, endDT time.Time) {
	// 	if err := impl.processGlobalRanksForGoogleFitApps(context.Background(), list, rp_s.MetricTypeActivitySteps, rp_s.FunctionAverage, rp_s.PeriodWeek, startDT, endDT); err != nil {
	// 		impl.Logger.Error("failed generating global activity steps average ranking",
	// 			slog.Any("date_range", "iso_week"),
	// 			slog.Any("error", err))
	// 		return
	// 	}
	// }(gfas.Results, start, end)
	//
	// go func(list []*gfa_ds.GoogleFitApp, startDT time.Time, endDT time.Time) {
	// 	if err := impl.processGlobalRanksForGoogleFitApps(context.Background(), list, rp_s.MetricTypeActivitySteps, rp_s.FunctionSum, rp_s.PeriodWeek, startDT, endDT); err != nil {
	// 		impl.Logger.Error("failed generating global activity steps summation ranking",
	// 			slog.Any("date_range", "iso_week"),
	// 			slog.Any("error", err))
	// 		return
	// 	}
	// }(gfas.Results, start, end)

	return nil
}

func (impl *RankPointControllerImpl) GenerateGlobalRankingForThisMonthUsingActiveGoogleFitApps(ctx context.Context) error {
	f := &gfa_ds.GoogleFitAppListFilter{
		Cursor:    primitive.NilObjectID,
		PageSize:  1_000_000,
		SortField: "_id",
		SortOrder: -1,
		Status:    gfa_ds.StatusActive,
	}
	gfas, err := impl.GoogleFitAppStorer.ListByFilter(ctx, f)
	if err != nil {
		impl.Logger.Error("failed listing by active status",
			slog.Any("date_range", "month"),
			slog.Any("error", err))
		return err
	}
	start := timekit.FirstDayOfThisMonth(time.Now)
	end := timekit.FirstDayOfNextMonth(time.Now)

	// --- Heart rate --- //

	go func(list []*gfa_ds.GoogleFitApp, startDT time.Time, endDT time.Time) {
		if err := impl.processGlobalRanksForGoogleFitApps(context.Background(), list, rp_s.MetricTypeHeartRate, rp_s.FunctionAverage, rp_s.PeriodMonth, startDT, endDT); err != nil {
			impl.Logger.Error("failed generating global heart rate ranking",
				slog.Any("date_range", "month"),
				slog.Any("error", err))
			return
		}
	}(gfas.Results, start, end)

	// // --- Activity steps --- //
	//
	// go func(list []*gfa_ds.GoogleFitApp, startDT time.Time, endDT time.Time) {
	// 	if err := impl.processGlobalRanksForGoogleFitApps(context.Background(), list, rp_s.MetricTypeActivitySteps, rp_s.FunctionAverage, rp_s.PeriodMonth, startDT, endDT); err != nil {
	// 		impl.Logger.Error("failed generating global activity steps average ranking",
	// 			slog.Any("date_range", "month"),
	// 			slog.Any("error", err))
	// 		return
	// 	}
	// }(gfas.Results, start, end)
	//
	// go func(list []*gfa_ds.GoogleFitApp, startDT time.Time, endDT time.Time) {
	// 	if err := impl.processGlobalRanksForGoogleFitApps(context.Background(), list, rp_s.MetricTypeActivitySteps, rp_s.FunctionSum, rp_s.PeriodMonth, startDT, endDT); err != nil {
	// 		impl.Logger.Error("failed generating global activity steps summation ranking",
	// 			slog.Any("date_range", "month"),
	// 			slog.Any("error", err))
	// 		return
	// 	}
	// }(gfas.Results, start, end)

	return nil
}

func (impl *RankPointControllerImpl) GenerateGlobalRankingForThisYearUsingActiveGoogleFitApps(ctx context.Context) error {
	f := &gfa_ds.GoogleFitAppListFilter{
		Cursor:    primitive.NilObjectID,
		PageSize:  1_000_000,
		SortField: "_id",
		SortOrder: -1,
		Status:    gfa_ds.StatusActive,
	}
	gfas, err := impl.GoogleFitAppStorer.ListByFilter(ctx, f)
	if err != nil {
		impl.Logger.Error("failed listing by active status",
			slog.Any("date_range", "year"),
			slog.Any("error", err))
		return err
	}
	start := timekit.FirstDayOfThisYear(time.Now)
	end := timekit.FirstDayOfNextYear(time.Now)

	// --- Heart rate --- //

	go func(list []*gfa_ds.GoogleFitApp, startDT time.Time, endDT time.Time) {
		if err := impl.processGlobalRanksForGoogleFitApps(context.Background(), list, rp_s.MetricTypeHeartRate, rp_s.FunctionAverage, rp_s.PeriodYear, startDT, endDT); err != nil {
			impl.Logger.Error("failed generating global heart rate ranking",
				slog.Any("date_range", "year"),
				slog.Any("error", err))
			return
		}
	}(gfas.Results, start, end)
	//
	// // --- Activity steps --- //
	//
	// go func(list []*gfa_ds.GoogleFitApp, startDT time.Time, endDT time.Time) {
	// 	if err := impl.processGlobalRanksForGoogleFitApps(context.Background(), list, rp_s.MetricTypeActivitySteps, rp_s.FunctionAverage, rp_s.PeriodYear, startDT, endDT); err != nil {
	// 		impl.Logger.Error("failed generating global activity steps average ranking",
	// 			slog.Any("date_range", "year"),
	// 			slog.Any("error", err))
	// 		return
	// 	}
	// }(gfas.Results, start, end)
	//
	// go func(list []*gfa_ds.GoogleFitApp, startDT time.Time, endDT time.Time) {
	// 	if err := impl.processGlobalRanksForGoogleFitApps(context.Background(), list, rp_s.MetricTypeActivitySteps, rp_s.FunctionSum, rp_s.PeriodYear, startDT, endDT); err != nil {
	// 		impl.Logger.Error("failed generating global activity steps summation ranking",
	// 			slog.Any("date_range", "year"),
	// 			slog.Any("error", err))
	// 		return
	// 	}
	// }(gfas.Results, start, end)

	return nil
}
