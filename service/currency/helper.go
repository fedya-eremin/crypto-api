package currency

func (s *CurrencyService) generateTaskId(symbol string) string {
	return updateCurrencyTaskPrefix + symbol
}
