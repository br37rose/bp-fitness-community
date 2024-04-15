package crontab

import (
	"context"
	"log/slog"
)

func (impl *googleFitAppCrontaberImpl) DeleteAllAnomalousData() error {
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
