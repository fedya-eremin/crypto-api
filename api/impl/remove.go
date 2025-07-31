package impl

import (
	"context"
	"errors"
	"fmt"

	"github.com/fedya-eremin/crypto-api/api"
	"github.com/fedya-eremin/crypto-api/service/currency"
)

func (s *Server) PostCurrencyRemove(
	ctx context.Context,
	req api.PostCurrencyRemoveRequestObject,
) (api.PostCurrencyRemoveResponseObject, error) {
	err := s.currencyService.RemoveCurrency(ctx, req.Body.Coin)
	if err != nil {
		var serviceErr *currency.ServiceError
		if errors.As(err, &serviceErr) {
			switch serviceErr.Code {
			case currency.CodeInvalid:
				return api.PostCurrencyRemove400Response{}, nil
			case currency.CodeNotFound:
				return api.PostCurrencyRemove404Response{}, nil
			default:
				return api.PostCurrencyRemove500Response{}, nil
			}
		}
		return api.PostCurrencyRemove500Response{}, nil
	}
	return api.PostCurrencyRemove200JSONResponse{
		Message: fmt.Sprintf(
			"currency %v removed from watchlist",
			req.Body.Coin,
		),
		Success: true,
	}, nil
}
