package currency

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/hibiken/asynq"
)

func (s *CurrencyService) BootstrapTasks(ctx context.Context) error {
	curs := s.currencyStorage.GetWatchableCurrencies(ctx)
	for _, cur := range curs {
		slog.Info("Bootstraping task for symbol", "symbol", cur.Symbol)
		payload := TaskPayload{Symbol: cur.Symbol}
		payloadBytes, err := json.Marshal(payload)
		if err != nil {
			return NewServiceError(err, serviceOrigin, "error marshaling payload", CodeUnknown)
		}
		taskId := s.generateTaskId(cur.Symbol)
		task := asynq.NewTask(
			TypeUpdateCurrencyTask,
			payloadBytes,
			asynq.TaskID(taskId),
		)
		entryId, err := s.taskScheduler.Register(fmt.Sprintf("@every %ds", cur.Interval), task)
		s.taskIds[taskId] = entryId
	}
	return nil
}
