package currency

import (
	"context"
)

func (r *CurrencyRepo) RemoveCurrency(ctx context.Context, symbol string) error {
	err := r.q.UnwatchCurrency(ctx, symbol)
	return err
}
