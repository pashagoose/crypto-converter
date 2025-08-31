package usecases

import (
	"context"
	"crypto-converter/internal/domain/entities"
	"crypto-converter/internal/domain/repositories"
	"fmt"
	"log/slog"
)

type ConvertCurrencyUseCase struct {
	currencyRepo repositories.CurrencyRepository
	logger       *slog.Logger
}

func NewConvertCurrencyUseCase(currencyRepo repositories.CurrencyRepository, logger *slog.Logger) *ConvertCurrencyUseCase {
	return &ConvertCurrencyUseCase{
		currencyRepo: currencyRepo,
		logger:       logger,
	}
}

func (u *ConvertCurrencyUseCase) Execute(ctx context.Context, amount float64, fromCode, toCode string) (*entities.ConversionResult, error) {
	u.logger.Info("Starting currency conversion use case",
		"amount", amount,
		"from", fromCode,
		"to", toCode)

	u.logger.Debug("Creating source money entity")
	fromCurrency, err := entities.NewCurrencyCode(fromCode)
	if err != nil {
		u.logger.Error("Failed to create source currency entity", "error", err)
		return nil, fmt.Errorf("invalid source currency: %w", err)
	}

	fromMoney, err := entities.NewMoney(fromCurrency, amount)
	if err != nil {
		u.logger.Error("Failed to create source money entity", "error", err)
		return nil, fmt.Errorf("invalid source money: %w", err)
	}

	u.logger.Debug("Fetching exchange rate", "from", fromCode, "to", toCode)
	rate, err := u.currencyRepo.GetExchangeRate(ctx, fromCode, toCode)
	if err != nil {
		u.logger.Error("Failed to get exchange rate",
			"from", fromCode,
			"to", toCode,
			"error", err)
		return nil, fmt.Errorf("failed to get exchange rate: %w", err)
	}

	u.logger.Debug("Exchange rate retrieved", "rate", rate)

	convertedAmount := amount * rate
	u.logger.Debug("Calculated converted amount", "converted_amount", convertedAmount)

	toCurrency, err := entities.NewCurrencyCode(toCode)
	if err != nil {
		u.logger.Error("Failed to create target currency entity", "error", err)
		return nil, fmt.Errorf("invalid target currency: %w", err)
	}
	toMoney, err := entities.NewMoney(toCurrency, convertedAmount)
	if err != nil {
		u.logger.Error("Failed to create target money entity", "error", err)
		return nil, fmt.Errorf("invalid target money: %w", err)
	}

	result := entities.NewConversionResult(fromMoney, toMoney, rate)

	u.logger.Info("Currency conversion completed successfully",
		"from_amount", result.From.Amount,
		"from_currency", result.From.Currency,
		"to_amount", result.To.Amount,
		"to_currency", result.To.Currency,
		"exchange_rate", result.Rate)

	return result, nil
}
