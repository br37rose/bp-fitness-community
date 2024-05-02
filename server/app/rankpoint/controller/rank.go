package controller

import (
	"context"
	"log/slog"
	"sync"

	gfa_ds "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/googlefitapp/datastore"
	rp_s "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/rankpoint/datastore"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (impl *RankPointControllerImpl) GenerateGlobalRankingForActiveGoogleFitApps(ctx context.Context) error {
	f := &gfa_ds.GoogleFitAppListFilter{
		Cursor:    primitive.NilObjectID,
		PageSize:  1_000_000,
		SortField: "created_at",
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

	// MetricDataTypeNames := []int8{
	// 	gcp_a.DataTypeKeyCaloriesBurned,
	// 	gcp_a.DataTypeKeyStepCountDelta,
	// 	gcp_a.DataTypeKeyHeartRateBPM,
	// 	gcp_a.DataTypeKeyDistanceDelta,
	// 	//TODO: Add more health sensors here...
	// }
	//
	// funcTypes := []int8{
	// 	rp_s.FunctionAverage,
	// 	rp_s.FunctionSum,
	// 	rp_s.FunctionCount,
	// 	rp_s.FunctionMin,
	// 	rp_s.FunctionMax,
	// }

	// Variable stores the number of goroutines we expect to wait for. We
	// set value of `1` because we have the following functions we want to
	// process in the background as goroutines:
	// - Today
	// - ISO Week
	// - Month
	// - Year
	numWorkers := 4

	// Create a channel to collect errors from goroutines.
	errCh := make(chan error, numWorkers)

	// Variable used to synchronize all the go routines running in
	// background outside of this function.
	var wg sync.WaitGroup

	// // Variable used to lock / unlock access when the goroutines want to
	// // perform writes to our output response.
	// var mu sync.Mutex

	// Load up the number of workers our waitgroup will need to handle.
	wg.Add(numWorkers)

	// ------ TODAY ------ //
	go func() {
		if err := impl.processGlobalRanksForGoogleFitAppsV2(context.Background(), gfas.Results, rp_s.PeriodDay); err != nil {
			impl.Logger.Error("failed generating global rate ranking for today",
				slog.Any("period", "day"),
				slog.Any("error", err))
			return
		}
		wg.Done() // We are done this background task.
	}()

	// ------ ISO Week ------ //
	go func() {
		if err := impl.processGlobalRanksForGoogleFitAppsV2(context.Background(), gfas.Results, rp_s.PeriodWeek); err != nil {
			impl.Logger.Error("failed generating global rate ranking for iso week",
				slog.Any("period", "iso week"),
				slog.Any("error", err))
			return
		}
		wg.Done() // We are done this background task.
	}()

	// ------ Month ------ //
	go func() {
		if err := impl.processGlobalRanksForGoogleFitAppsV2(context.Background(), gfas.Results, rp_s.PeriodMonth); err != nil {
			impl.Logger.Error("failed generating global rate ranking for month",
				slog.Any("period", "month"),
				slog.Any("error", err))
			return
		}
		wg.Done() // We are done this background task.
	}()

	// ------ Year ------ //
	go func() {
		go func() {
			if err := impl.processGlobalRanksForGoogleFitAppsV2(context.Background(), gfas.Results, rp_s.PeriodYear); err != nil {
				impl.Logger.Error("failed generating global rate ranking for year",
					slog.Any("period", "year"),
					slog.Any("error", err))
				return
			}
			wg.Done() // We are done this background task.
		}()
	}()

	// Create a goroutine to close the error channel when all workers are done
	go func() {
		wg.Wait()
		close(errCh)
	}()

	// Iterate over the error channel to collect any errors from workers
	for err := range errCh {
		impl.Logger.Error("failed executing in goroutine",
			slog.Any("error", err))
		return err
	}

	impl.Logger.Debug("ranking completed")
	return nil
}
