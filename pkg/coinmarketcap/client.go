package coinmarketcap

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	DefaultSandboxURL    = "https://sandbox-api.coinmarketcap.com"
	DefaultProductionURL = "https://pro-api.coinmarketcap.com"
	APIVersion           = "v1"
)

type Client struct {
	apiKey     string
	baseURL    string
	httpClient *http.Client
	logger     *slog.Logger
}

type Config struct {
	APIKey  string
	BaseURL string
	Timeout time.Duration
	Logger  *slog.Logger
}

func NewClient(config Config) *Client {
	if config.BaseURL == "" {
		config.BaseURL = DefaultSandboxURL
	}
	if config.Timeout == 0 {
		config.Timeout = 30 * time.Second
	}
	if config.Logger == nil {
		config.Logger = slog.Default()
	}

	return &Client{
		apiKey:  config.APIKey,
		baseURL: config.BaseURL,
		logger:  config.Logger,
		httpClient: &http.Client{
			Timeout: config.Timeout,
		},
	}
}

func NewClientSimple(apiKey, baseURL string, logger *slog.Logger) *Client {
	return NewClient(Config{
		APIKey:  apiKey,
		BaseURL: baseURL,
		Logger:  logger,
	})
}

type QuoteResponse struct {
	Status Status                    `json:"status"`
	Data   map[string]CryptoCurrency `json:"data"`
}

type Status struct {
	Timestamp    time.Time `json:"timestamp"`
	ErrorCode    int       `json:"error_code"`
	ErrorMessage string    `json:"error_message"`
	Elapsed      int       `json:"elapsed"`
	CreditCount  int       `json:"credit_count"`
}

type CryptoCurrency struct {
	ID     int              `json:"id"`
	Name   string           `json:"name"`
	Symbol string           `json:"symbol"`
	Quote  map[string]Quote `json:"quote"`
}

type Quote struct {
	Price            float64   `json:"price"`
	Volume24h        float64   `json:"volume_24h"`
	PercentChange1h  float64   `json:"percent_change_1h"`
	PercentChange24h float64   `json:"percent_change_24h"`
	PercentChange7d  float64   `json:"percent_change_7d"`
	MarketCap        float64   `json:"market_cap"`
	LastUpdated      time.Time `json:"last_updated"`
}

func (c *Client) GetQuotes(ctx context.Context, symbols []string, convert string) (*QuoteResponse, error) {
	endpoint := fmt.Sprintf("%s/%s/cryptocurrency/quotes/latest", c.baseURL, APIVersion)

	params := url.Values{}
	params.Set("symbol", strings.Join(symbols, ","))
	params.Set("convert", convert)

	fullURL := endpoint + "?" + params.Encode()

	c.logger.Debug("Making API request for quotes",
		"endpoint", endpoint,
		"symbols", symbols,
		"convert", convert)

	req, err := http.NewRequestWithContext(ctx, "GET", fullURL, nil)
	if err != nil {
		c.logger.Error("Failed to create HTTP request", "error", err)
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("X-CMC_PRO_API_KEY", c.apiKey)
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		c.logger.Error("HTTP request failed", "error", err, "url", fullURL)
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	c.logger.Debug("Received API response", "status_code", resp.StatusCode)

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		c.logger.Error("API request failed",
			"status_code", resp.StatusCode,
			"response_body", string(body))
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var quoteResponse QuoteResponse
	if err := json.NewDecoder(resp.Body).Decode(&quoteResponse); err != nil {
		c.logger.Error("Failed to decode API response", "error", err)
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if quoteResponse.Status.ErrorCode != 0 {
		c.logger.Error("API returned error",
			"error_code", quoteResponse.Status.ErrorCode,
			"error_message", quoteResponse.Status.ErrorMessage)
		return nil, fmt.Errorf("API error %d: %s", quoteResponse.Status.ErrorCode, quoteResponse.Status.ErrorMessage)
	}

	c.logger.Debug("Successfully retrieved quotes", "data_count", len(quoteResponse.Data), "response", quoteResponse)
	return &quoteResponse, nil
}
