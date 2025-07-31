package currency

import db_gen "github.com/fedya-eremin/crypto-api/database/gen"

type CurrencyRepo struct {
	q *db_gen.Queries
}

func New(q *db_gen.Queries) *CurrencyRepo {
	return &CurrencyRepo{
		q: q,
	}
}
