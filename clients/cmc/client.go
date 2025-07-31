package cmc

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"
)

type Client struct {
	apiKey          string
	baseURL         string
	httpClient      *http.Client
	convertCurrency string
}

func New(apiKey string, defaultCurrency string) *Client {
	return &Client{
		apiKey:  apiKey,
		baseURL: "https://pro-api.coinmarketcap.com/v1/cryptocurrency",
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		convertCurrency: defaultCurrency,
	}
}

func (c *Client) GetPrice(ctx context.Context, symbol string) (string, error) {
	data, err := c.getCryptoData(ctx, symbol, c.convertCurrency)
	if err != nil {
		return "0", err
	}

	if quote, exists := data.Quote[c.convertCurrency]; exists {
		return quote.Price.String(), nil
	}

	return "0", fmt.Errorf("default currency %s not found for symbol %s", c.convertCurrency, symbol)
}

func (c *Client) CheckIfExists(ctx context.Context, symbol string) error {
	_, err := c.getCryptoData(ctx, symbol, c.convertCurrency)
	return err
}

func (c *Client) getCryptoData(
	ctx context.Context,
	symbol string,
	currency string,
) (*CryptoData, error) {
	url := fmt.Sprintf("%s/quotes/latest", c.baseURL)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	q := req.URL.Query()
	q.Add("symbol", symbol)
	q.Add("convert", currency)
	req.URL.RawQuery = q.Encode()

	req.Header.Add("X-CMC_PRO_API_KEY", c.apiKey)
	req.Header.Add("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		slog.Error("network error when fetching api", "error", err)
		return nil, fmt.Errorf("API request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		slog.Error("api wrong status", "status", resp.StatusCode, "body", body)
		return nil, NewCmcClientError(fmt.Sprintf("body: %v", body), resp.StatusCode)
	}

	var apiResponse Response
	if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	if cryptoData, exists := apiResponse.Data[symbol]; exists {
		return &cryptoData, nil
	}

	return nil, fmt.Errorf("symbol %s not found in response", symbol)
}
