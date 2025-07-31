package impl

import (
	"context"
	"errors"

	"github.com/fedya-eremin/crypto-api/api"
	"github.com/fedya-eremin/crypto-api/service/currency"
)

func (s *Server) PostCurrencyPrice(
	ctx context.Context,
	req api.PostCurrencyPriceRequestObject,
) (api.PostCurrencyPriceResponseObject, error) {
	if req.Body.Timestamp <= 0 {
		return api.PostCurrencyPrice400Response{}, nil
	}
	price, collectedAt, err := s.currencyService.GetPriceOffline(
		ctx,
		req.Body.Coin,
		int64(req.Body.Timestamp),
	)
	if err != nil {
		var serviceErr *currency.ServiceError
		if errors.As(err, &serviceErr) {
			switch serviceErr.Code {
			case currency.CodeInvalid:
				return api.PostCurrencyPrice400Response{}, nil
			case currency.CodeNotFound:
				return api.PostCurrencyPrice404Response{}, nil
			default:
				return api.PostCurrencyPrice500Response{}, nil
			}
		}
		return api.PostCurrencyPrice500Response{}, nil
	}
	return api.PostCurrencyPrice200JSONResponse{
		Coin:      req.Body.Coin,
		Price:     price,
		Timestamp: int(collectedAt.Unix()),
	}, nil
}
