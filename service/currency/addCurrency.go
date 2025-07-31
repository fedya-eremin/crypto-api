package currency

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hibiken/asynq"
)

func (s *CurrencyService) AddCurrency(ctx context.Context, symbol string, interval int32) error {
	if err := s.cmc.CheckIfExists(ctx, symbol); err != nil {
		return NewServiceError(err, cmcOrigin, "symbol not found", CodeNotFound)
	}
	if err := s.currencyStorage.AddCurrency(ctx, symbol, interval); err != nil {
		return NewServiceError(err, currencyRepoOrigin, "error in repo", CodeUnknown)
	}
	payload := TaskPayload{Symbol: symbol}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return NewServiceError(err, serviceOrigin, "error marshaling payload", CodeUnknown)
	}
	taskId := s.generateTaskId(symbol)
	task := asynq.NewTask(
		TypeUpdateCurrencyTask,
		payloadBytes,
		asynq.TaskID(taskId),
	)
	entryId, err := s.taskScheduler.Register(fmt.Sprintf("@every %ds", interval), task)
	s.taskIds[taskId] = entryId
	if err != nil {
		return NewServiceError(err, serviceOrigin, "error adding task", CodeUnknown)
	}
	return nil
}
