package controller

import (
	"context"

	"log/slog"

	c_s "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/aggregatepoint/datastore"
)

func (c *AggregatePointControllerImpl) ListByFilter(ctx context.Context, f *c_s.AggregatePointPaginationListFilter) (*c_s.AggregatePointPaginationListResult, error) {
	c.Logger.Debug("listing using filter options:",
		slog.Any("Cursor", f.Cursor),
		slog.Int64("PageSize", f.PageSize),
		slog.String("SortField", f.SortField),
		slog.Int("SortOrder", int(f.SortOrder)),
		slog.Any("MetricIDs", f.MetricIDs),
		slog.Time("CreatedAtGTE", f.CreatedAtGTE))

	m, err := c.AggregatePointStorer.ListByFilter(ctx, f)
	if err != nil {
		c.Logger.Error("database list by filter error", slog.Any("error", err))
		return nil, err
	}
	return m, err
}
