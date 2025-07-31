package pricelog

import db_gen "github.com/fedya-eremin/crypto-api/database/gen"

type PriceLogRepo struct {
	q *db_gen.Queries
}

func New(q *db_gen.Queries) *PriceLogRepo {
	return &PriceLogRepo{
		q: q,
	}
}
