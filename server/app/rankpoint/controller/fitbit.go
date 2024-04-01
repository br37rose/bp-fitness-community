package controller

import (
	"context"
)

func (impl *RankPointControllerImpl) GenerateGlobalRankingForTodayUsingActiveFitBitApps(ctx context.Context) error {
	return nil
	// f := &fba_s.FitBitAppListFilter{
	// 	Cursor:    primitive.NilObjectID,
	// 	PageSize:  1_000_000,
	// 	SortField: "_id",
	// 	SortOrder: -1,
	// 	Status:    fba_s.StatusActive,
	// }
	// fbas, err := impl.FitBitAppStorer.ListByFilter(ctx, f)
	// if err != nil {
	// 	impl.Logger.Error("failed listing by active status",
	// 		slog.Any("date_range", "today"),
	// 		slog.Any("error", err))
	// 	return err
	// }
	// start := timekit.Midnight(time.Now)
	// end := timekit.MidnightTomorrow(time.Now)
	//
	// // --- Heart rate --- //
	//
	// go func(list []*fba_s.FitBitApp, startDT time.Time, endDT time.Time) {
	// 	if err := impl.processGlobalRanksForFitBitApps(context.Background(), list, rp_s.MetricTypeHeartRate, rp_s.FunctionAverage, rp_s.PeriodDay, startDT, endDT); err != nil {
	// 		impl.Logger.Error("failed generating global heart average rate ranking",
	// 			slog.Any("date_range", "today"),
	// 			slog.Any("error", err))
	// 		return
	// 	}
	// }(fbas.Results, start, end)
	//
	// // --- Activity steps --- //
	//
	// go func(list []*fba_s.FitBitApp, startDT time.Time, endDT time.Time) {
	// 	if err := impl.processGlobalRanksForFitBitApps(context.Background(), list, rp_s.MetricTypeActivitySteps, rp_s.FunctionAverage, rp_s.PeriodDay, startDT, endDT); err != nil {
	// 		impl.Logger.Error("failed generating global activity steps average ranking",
	// 			slog.Any("date_range", "today"),
	// 			slog.Any("error", err))
	// 		return
	// 	}
	// }(fbas.Results, start, end)
	//
	// go func(list []*fba_s.FitBitApp, startDT time.Time, endDT time.Time) {
	// 	if err := impl.processGlobalRanksForFitBitApps(context.Background(), list, rp_s.MetricTypeActivitySteps, rp_s.FunctionSum, rp_s.PeriodDay, startDT, endDT); err != nil {
	// 		impl.Logger.Error("failed generating global activity steps sum ranking",
	// 			slog.Any("date_range", "today"),
	// 			slog.Any("error", err))
	// 		return
	// 	}
	// }(fbas.Results, start, end)
	//
	// return nil
}

func (impl *RankPointControllerImpl) GenerateGlobalRankingForThisISOWeekUsingActiveFitBitApps(ctx context.Context) error {
	return nil
	// f := &fba_s.FitBitAppListFilter{
	// 	Cursor:    primitive.NilObjectID,
	// 	PageSize:  1_000_000,
	// 	SortField: "_id",
	// 	SortOrder: -1,
	// 	Status:    fba_s.StatusActive,
	// }
	// fbas, err := impl.FitBitAppStorer.ListByFilter(ctx, f)
	// if err != nil {
	// 	impl.Logger.Error("failed listing by active status",
	// 		slog.Any("date_range", "iso_week"),
	// 		slog.Any("error", err))
	// 	return err
	// }
	// start := timekit.FirstDayOfThisISOWeek(time.Now)
	// end := timekit.FirstDayOfNextISOWeek(time.Now)
	//
	// // --- Heart rate --- //
	//
	// go func(list []*fba_s.FitBitApp, startDT time.Time, endDT time.Time) {
	// 	if err := impl.processGlobalRanksForFitBitApps(context.Background(), list, rp_s.MetricTypeHeartRate, rp_s.FunctionAverage, rp_s.PeriodWeek, startDT, endDT); err != nil {
	// 		impl.Logger.Error("failed generating global heart rate ranking",
	// 			slog.Any("date_range", "iso_week"),
	// 			slog.Any("error", err))
	// 		return
	// 	}
	// }(fbas.Results, start, end)
	//
	// // --- Activity steps --- //
	//
	// go func(list []*fba_s.FitBitApp, startDT time.Time, endDT time.Time) {
	// 	if err := impl.processGlobalRanksForFitBitApps(context.Background(), list, rp_s.MetricTypeActivitySteps, rp_s.FunctionAverage, rp_s.PeriodWeek, startDT, endDT); err != nil {
	// 		impl.Logger.Error("failed generating global activity steps average ranking",
	// 			slog.Any("date_range", "iso_week"),
	// 			slog.Any("error", err))
	// 		return
	// 	}
	// }(fbas.Results, start, end)
	//
	// go func(list []*fba_s.FitBitApp, startDT time.Time, endDT time.Time) {
	// 	if err := impl.processGlobalRanksForFitBitApps(context.Background(), list, rp_s.MetricTypeActivitySteps, rp_s.FunctionSum, rp_s.PeriodWeek, startDT, endDT); err != nil {
	// 		impl.Logger.Error("failed generating global activity steps summation ranking",
	// 			slog.Any("date_range", "iso_week"),
	// 			slog.Any("error", err))
	// 		return
	// 	}
	// }(fbas.Results, start, end)
	//
	// return nil
}

func (impl *RankPointControllerImpl) GenerateGlobalRankingForThisMonthUsingActiveFitBitApps(ctx context.Context) error {
	return nil
	// f := &fba_s.FitBitAppListFilter{
	// 	Cursor:    primitive.NilObjectID,
	// 	PageSize:  1_000_000,
	// 	SortField: "_id",
	// 	SortOrder: -1,
	// 	Status:    fba_s.StatusActive,
	// }
	// fbas, err := impl.FitBitAppStorer.ListByFilter(ctx, f)
	// if err != nil {
	// 	impl.Logger.Error("failed listing by active status",
	// 		slog.Any("date_range", "month"),
	// 		slog.Any("error", err))
	// 	return err
	// }
	// start := timekit.FirstDayOfThisMonth(time.Now)
	// end := timekit.FirstDayOfNextMonth(time.Now)
	//
	// // --- Heart rate --- //
	//
	// go func(list []*fba_s.FitBitApp, startDT time.Time, endDT time.Time) {
	// 	if err := impl.processGlobalRanksForFitBitApps(context.Background(), list, rp_s.MetricTypeHeartRate, rp_s.FunctionAverage, rp_s.PeriodMonth, startDT, endDT); err != nil {
	// 		impl.Logger.Error("failed generating global heart rate ranking",
	// 			slog.Any("date_range", "month"),
	// 			slog.Any("error", err))
	// 		return
	// 	}
	// }(fbas.Results, start, end)
	//
	// // --- Activity steps --- //
	//
	// go func(list []*fba_s.FitBitApp, startDT time.Time, endDT time.Time) {
	// 	if err := impl.processGlobalRanksForFitBitApps(context.Background(), list, rp_s.MetricTypeActivitySteps, rp_s.FunctionAverage, rp_s.PeriodMonth, startDT, endDT); err != nil {
	// 		impl.Logger.Error("failed generating global activity steps average ranking",
	// 			slog.Any("date_range", "month"),
	// 			slog.Any("error", err))
	// 		return
	// 	}
	// }(fbas.Results, start, end)
	//
	// go func(list []*fba_s.FitBitApp, startDT time.Time, endDT time.Time) {
	// 	if err := impl.processGlobalRanksForFitBitApps(context.Background(), list, rp_s.MetricTypeActivitySteps, rp_s.FunctionSum, rp_s.PeriodMonth, startDT, endDT); err != nil {
	// 		impl.Logger.Error("failed generating global activity steps summation ranking",
	// 			slog.Any("date_range", "month"),
	// 			slog.Any("error", err))
	// 		return
	// 	}
	// }(fbas.Results, start, end)
	//
	// return nil
}

func (impl *RankPointControllerImpl) GenerateGlobalRankingForThisYearUsingActiveFitBitApps(ctx context.Context) error {
	return nil
	// f := &fba_s.FitBitAppListFilter{
	// 	Cursor:    primitive.NilObjectID,
	// 	PageSize:  1_000_000,
	// 	SortField: "_id",
	// 	SortOrder: -1,
	// 	Status:    fba_s.StatusActive,
	// }
	// fbas, err := impl.FitBitAppStorer.ListByFilter(ctx, f)
	// if err != nil {
	// 	impl.Logger.Error("failed listing by active status",
	// 		slog.Any("date_range", "year"),
	// 		slog.Any("error", err))
	// 	return err
	// }
	// start := timekit.FirstDayOfThisYear(time.Now)
	// end := timekit.FirstDayOfNextYear(time.Now)
	//
	// // --- Heart rate --- //
	//
	// go func(list []*fba_s.FitBitApp, startDT time.Time, endDT time.Time) {
	// 	if err := impl.processGlobalRanksForFitBitApps(context.Background(), list, rp_s.MetricTypeHeartRate, rp_s.FunctionAverage, rp_s.PeriodYear, startDT, endDT); err != nil {
	// 		impl.Logger.Error("failed generating global heart rate ranking",
	// 			slog.Any("date_range", "year"),
	// 			slog.Any("error", err))
	// 		return
	// 	}
	// }(fbas.Results, start, end)
	//
	// // --- Activity steps --- //
	//
	// go func(list []*fba_s.FitBitApp, startDT time.Time, endDT time.Time) {
	// 	if err := impl.processGlobalRanksForFitBitApps(context.Background(), list, rp_s.MetricTypeActivitySteps, rp_s.FunctionAverage, rp_s.PeriodYear, startDT, endDT); err != nil {
	// 		impl.Logger.Error("failed generating global activity steps average ranking",
	// 			slog.Any("date_range", "year"),
	// 			slog.Any("error", err))
	// 		return
	// 	}
	// }(fbas.Results, start, end)
	//
	// go func(list []*fba_s.FitBitApp, startDT time.Time, endDT time.Time) {
	// 	if err := impl.processGlobalRanksForFitBitApps(context.Background(), list, rp_s.MetricTypeActivitySteps, rp_s.FunctionSum, rp_s.PeriodYear, startDT, endDT); err != nil {
	// 		impl.Logger.Error("failed generating global activity steps summation ranking",
	// 			slog.Any("date_range", "year"),
	// 			slog.Any("error", err))
	// 		return
	// 	}
	// }(fbas.Results, start, end)
	//
	// return nil
}
