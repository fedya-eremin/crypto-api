package currency

import (
	"context"
	"time"

	"github.com/hibiken/asynq"
)

type CurrencyStorage interface {
	AddCurrency(ctx context.Context, symbol string, interval int32) error
	RemoveCurrency(ctx context.Context, symbol string) error
	GetPrice(ctx context.Context, symbol string, timestamp int64) (string, time.Time, error)
	GetWatchableCurrencies(ctx context.Context) []Currency
}

type CmcClient interface {
	GetPrice(ctx context.Context, symbol string) (string, error)
	CheckIfExists(ctx context.Context, symbol string) error
}

type Scheduler interface {
	Register(key string, task *asynq.Task, opts ...asynq.Option) (string, error)
	Unregister(key string) error
}

type CurrencyService struct {
	currencyStorage CurrencyStorage
	cmc             CmcClient
	taskScheduler   Scheduler
	taskIds         map[string]string
}

func New(db CurrencyStorage, client CmcClient, taskScheduler Scheduler) *CurrencyService {
	return &CurrencyService{
		currencyStorage: db,
		cmc:             client,
		taskScheduler:   taskScheduler,
		taskIds:         make(map[string]string),
	}
}
