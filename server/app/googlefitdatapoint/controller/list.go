package controller

import (
	"context"

	"log/slog"

	c_s "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/googlefitdatapoint/datastore"
)

func (c *GoogleFitDataPointControllerImpl) ListByFilter(ctx context.Context, f *c_s.GoogleFitDataPointPaginationListFilter) (*c_s.GoogleFitDataPointPaginationListResult, error) {
	c.Logger.Debug("listing using filter options:",
		slog.Any("Cursor", f.Cursor),
		slog.Int64("PageSize", f.PageSize),
		slog.String("SortField", f.SortField),
		slog.Int("SortOrder", int(f.SortOrder)),
		slog.Any("MetricIDs", f.MetricIDs),
		slog.Time("GTE", f.GTE),
		slog.Time("GT", f.GT),
		slog.Time("LTE", f.LTE),
		slog.Time("LT", f.LT))

	m, err := c.GoogleFitDataPointStorer.ListByFilter(ctx, f)
	if err != nil {
		c.Logger.Error("database list by filter error", slog.Any("error", err))
		return nil, err
	}
	return m, err
}