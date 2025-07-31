package impl

import "github.com/fedya-eremin/crypto-api/service/currency"

type Server struct {
	currencyService *currency.CurrencyService
}

func New(c *currency.CurrencyService) *Server {
	return &Server{
		currencyService: c,
	}
}
