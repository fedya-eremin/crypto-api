package currency

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
)

func (s *CurrencyService) GetCurrencyPriceOnline(
	ctx context.Context,
	symbol string,
) (string, error) {
	price, err := s.cmc.GetPrice(ctx, symbol)
	if err != nil {
		return "", NewServiceError(err, cmcOrigin, "cannot get price", CodeUnknown)
	}
	return price, nil
}

func (s *CurrencyService) GetPriceOffline(
	ctx context.Context,
	symbol string,
	timestamp int64,
) (string, time.Time, error) {
	price, collectedAt, err := s.currencyStorage.GetPrice(ctx, symbol, timestamp)
	if errors.Is(err, pgx.ErrNoRows) {
		return "", time.Now(), NewServiceError(err, cmcOrigin, "cannot get price", CodeNotFound)
	}
	if err != nil {
		return "", time.Now(), NewServiceError(
			err,
			currencyRepoOrigin,
			"error in repo",
			CodeUnknown,
		)
	}
	return price, collectedAt, nil
}
