package impl

import (
	"context"
	"errors"
	"fmt"

	"github.com/fedya-eremin/crypto-api/api"
	"github.com/fedya-eremin/crypto-api/service/currency"
)

func (s *Server) PostCurrencyAdd(
	ctx context.Context,
	req api.PostCurrencyAddRequestObject,
) (api.PostCurrencyAddResponseObject, error) {
	if req.Body.Interval <= 0 {
		return api.PostCurrencyAdd400Response{}, nil
	}
	err := s.currencyService.AddCurrency(ctx, req.Body.Coin, int32(req.Body.Interval))
	if err != nil {
		var serviceErr *currency.ServiceError
		if errors.As(err, &serviceErr) {
			switch serviceErr.Code {
			case currency.CodeInvalid:
				return api.PostCurrencyAdd400Response{}, nil
			case currency.CodeNotFound:
				return api.PostCurrencyAdd404Response{}, nil
			default:
				return api.PostCurrencyAdd500Response{}, nil
			}
		}
		return api.PostCurrencyAdd500Response{}, nil
	}
	return api.PostCurrencyAdd200JSONResponse{
		Message: fmt.Sprintf(
			"currency %v added with interval %v",
			req.Body.Coin,
			req.Body.Interval,
		),
		Success: true,
	}, nil
}
