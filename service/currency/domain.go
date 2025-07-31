package currency

const (
	TypeUpdateCurrencyTask   string = "currency:update"
	updateCurrencyTaskPrefix string = "update_currency_"
)

type TaskPayload struct {
	Symbol string `json:"symbol"`
}

type Currency struct {
	Symbol   string
	Watching bool
	Interval int
}
