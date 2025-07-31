package currency

import (
	"context"
	"time"

	db_gen "github.com/fedya-eremin/crypto-api/database/gen"
)

func (r *CurrencyRepo) GetPrice(
	ctx context.Context,
	symbol string,
	timestamp int64,
) (string, time.Time, error) {
	price, err := r.q.GetNearestPrice(ctx, db_gen.GetNearestPriceParams{
		Symbol:    symbol,
		Timestamp: timestamp,
	})
	if err != nil {
		return "", time.Now(), err
	}
	return price.PriceUsd, price.CollectedAt, nil
}
