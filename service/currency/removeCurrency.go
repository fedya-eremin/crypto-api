package currency

import (
	"context"
	"errors"
	"log/slog"

	"github.com/jackc/pgx/v5"
)

func (s *CurrencyService) RemoveCurrency(ctx context.Context, symbol string) error {
	err := s.currencyStorage.RemoveCurrency(ctx, symbol)
	if errors.Is(err, pgx.ErrNoRows) {
		return NewServiceError(err, currencyRepoOrigin, "cannot remove currency", CodeNotFound)
	}
	if err != nil {
		return NewServiceError(err, currencyRepoOrigin, "cannot remove currency", CodeUnknown)
	}
	taskId := s.generateTaskId(symbol)
	err = s.taskScheduler.Unregister(s.taskIds[taskId])
	delete(s.taskIds, taskId)
	if err != nil {
		slog.Error("cannot stop task", "error", err)
	}
	return nil
}
