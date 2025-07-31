package cmc

import "encoding/json"

type Quote struct {
	Price json.Number `json:"price"`
}

type CryptoData struct {
	Symbol string           `json:"symbol"`
	Quote  map[string]Quote `json:"quote"`
}

type Response struct {
	Data map[string]CryptoData `json:"data"`
}
