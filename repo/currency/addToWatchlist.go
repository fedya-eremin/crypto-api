package currency

import (
	"context"

	db_gen "github.com/fedya-eremin/crypto-api/database/gen"
	"github.com/google/uuid"
)

func (r *CurrencyRepo) AddCurrency(ctx context.Context, symbol string, interval int32) error {
	err := r.q.AddCurrencyToWatchlist(ctx, db_gen.AddCurrencyToWatchlistParams{
		Uuid:     uuid.New(),
		Symbol:   symbol,
		Interval: interval,
	})
	return err
}
