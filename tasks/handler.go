package tasks

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/fedya-eremin/crypto-api/clients/cmc"
	pricelog "github.com/fedya-eremin/crypto-api/repo/price-log"
	currency_svc "github.com/fedya-eremin/crypto-api/service/currency"
	"github.com/hibiken/asynq"
)

type Handler struct {
	logRepo   *pricelog.PriceLogRepo
	cmcClient *cmc.Client
}

func NewHandler(repo *pricelog.PriceLogRepo, client *cmc.Client) *Handler {
	return &Handler{
		logRepo:   repo,
		cmcClient: client,
	}
}

func (h *Handler) HandleCurrencyUpdateTask(ctx context.Context, task *asynq.Task) error {
	var payload currency_svc.TaskPayload
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", err)
	}
	slog.Info("starting task for symbol", "symbol", payload.Symbol)
	price, err := h.cmcClient.GetPrice(ctx, payload.Symbol)
	if err != nil {
		return err
	}
	return h.logRepo.AddPrice(ctx, payload.Symbol, price, time.Now())
}
