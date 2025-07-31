package currency

import (
	"context"

	currency_svc "github.com/fedya-eremin/crypto-api/service/currency"
)

func (r *CurrencyRepo) GetWatchableCurrencies(ctx context.Context) []currency_svc.Currency {
	curs, err := r.q.BootstrapWatchingEntries(ctx)
	if err != nil {
		return make([]currency_svc.Currency, 0)
	}
	res := make([]currency_svc.Currency, 0)

	for _, cur := range curs {
		res = append(res, currency_svc.Currency{
			Symbol:   cur.Symbol,
			Interval: int(cur.Interval),
			Watching: true,
		})
	}
	return res
}
