package controller

import (
	"context"
	"log/slog"

	q_s "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/question/datastore"
)

func (impl *QuestionControllerImpl) ListByFilter(ctx context.Context, f *q_s.QuestionListFilter) (*q_s.QuestionListResult, error) {
	questions, err := impl.QuestionStorer.ListByFilter(ctx, f)
	if err != nil {
		impl.Logger.Error("database list by filter error", slog.Any("error", err))
		return nil, err
	}
	return questions, nil

}
