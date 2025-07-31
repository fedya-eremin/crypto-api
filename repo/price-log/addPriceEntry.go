package pricelog

import (
	"context"
	"time"

	db_gen "github.com/fedya-eremin/crypto-api/database/gen"
	"github.com/google/uuid"
)

func (r *PriceLogRepo) AddPrice(
	ctx context.Context,
	symbol string,
	price string,
	collectedAt time.Time,
) error {
	err := r.q.AddCurrencyPriceLog(ctx, db_gen.AddCurrencyPriceLogParams{
		Uuid:        uuid.New(),
		PriceUsd:    price,
		Symbol:      symbol,
		CollectedAt: collectedAt,
	})
	return err
}
