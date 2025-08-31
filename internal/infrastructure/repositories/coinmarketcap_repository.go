package repositories

import (
	"context"
	"crypto-converter/pkg/coinmarketcap"
	"errors"
	"fmt"
	"log/slog"
	"strings"
)

type CoinMarketCapRepository struct {
	client *coinmarketcap.Client
	logger *slog.Logger
}

func NewCoinMarketCapRepository(apiKey, baseURL string, logger *slog.Logger) *CoinMarketCapRepository {
	return &CoinMarketCapRepository{
		client: coinmarketcap.NewClientSimple(apiKey, baseURL, logger),
		logger: logger,
	}
}

func (r *CoinMarketCapRepository) GetExchangeRate(ctx context.Context, from, to string) (float64, error) {
	from = strings.ToUpper(from)
	to = strings.ToUpper(to)

	r.logger.Debug("Getting exchange rate", "from", from, "to", to)
	if from == to {
		r.logger.Debug("Same currency conversion, returning rate 1.0")
		return 1.0, nil
	}

	if r.isFiatCurrency(from) && r.isFiatCurrency(to) {
		return 0, errors.New("fiat to fiat conversion is not supported by coinmarket cap")
	}

	var symbols []string
	var convert string

	if r.isFiatCurrency(to) {
		symbols = []string{from}
		convert = to
	} else if r.isFiatCurrency(from) {
		symbols = []string{to}
		convert = from
	} else {
		symbols = []string{from, to}
		convert = "USD"
	}

	quotes, err := r.client.GetQuotes(ctx, symbols, convert)
	if err != nil {
		return 0, fmt.Errorf("failed to get quotes: %w", err)
	}

	return r.calculateRate(quotes, from, to, convert)
}

func (r *CoinMarketCapRepository) calculateRate(quotes *coinmarketcap.QuoteResponse, from, to, convert string) (float64, error) {
	getPrice := func(currency string) (float64, error) {
		cryptoData, exists := quotes.Data[currency]
		if !exists {
			return 0, fmt.Errorf("no data found for currency %s", currency)
		}

		quoteData, exists := cryptoData.Quote[convert]
		if !exists {
			return 0, fmt.Errorf("no quote found for %s in %s", currency, convert)
		}

		return quoteData.Price, nil
	}

	fromIsFiat := r.isFiatCurrency(from)
	toIsFiat := r.isFiatCurrency(to)

	switch {
	case fromIsFiat && toIsFiat:
		return 0, errors.New("fiat-to-fiat conversion not supported")

	case fromIsFiat && !toIsFiat:
		cryptoPrice, err := getPrice(to)
		if err != nil {
			return 0, err
		}
		return 1.0 / cryptoPrice, nil

	case !fromIsFiat && toIsFiat:
		return getPrice(from)

	default:
		fromPrice, err := getPrice(from)
		if err != nil {
			return 0, err
		}

		toPrice, err := getPrice(to)
		if err != nil {
			return 0, err
		}

		return fromPrice / toPrice, nil
	}
}

func (r *CoinMarketCapRepository) isFiatCurrency(code string) bool {
	fiatCurrencies := map[string]bool{
		"USD": true, "EUR": true, "GBP": true, "JPY": true,
		"AUD": true, "CAD": true, "CHF": true, "CNY": true,
		"KRW": true, "RUB": true, "BRL": true, "INR": true,
	}

	return fiatCurrencies[strings.ToUpper(code)]
}
