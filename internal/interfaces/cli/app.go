package cli

import (
	"context"
	"crypto-converter/internal/domain/entities"
	"crypto-converter/internal/domain/usecases"
	"fmt"
	"log/slog"
	"os"
	"strings"
)

type CLI struct {
	convertUseCase *usecases.ConvertCurrencyUseCase
	logger         *slog.Logger
}

func NewCLI(convertUseCase *usecases.ConvertCurrencyUseCase, logger *slog.Logger) *CLI {
	return &CLI{
		convertUseCase: convertUseCase,
		logger:         logger,
	}
}

func (a *CLI) Run(args []string) error {
	a.logger.Debug("CLI Run started", "args", args)

	if len(args) != 4 {
		a.logger.Warn("Invalid number of arguments", "args_count", len(args))
		return a.showUsage()
	}

	amountStr := args[1]
	fromCurrency := strings.ToUpper(strings.TrimSpace(args[2]))
	toCurrency := strings.ToUpper(strings.TrimSpace(args[3]))

	a.logger.Info("Starting currency conversion",
		"amount", amountStr,
		"from", fromCurrency,
		"to", toCurrency)

	amount, err := entities.ParseAmount(amountStr)
	if err != nil {
		a.logger.Error("Failed to parse amount", "amount_str", amountStr, "error", err)
		return fmt.Errorf("invalid amount: %w", err)
	}

	ctx := context.Background()
	result, err := a.convertUseCase.Execute(ctx, amount, fromCurrency, toCurrency)
	if err != nil {
		a.logger.Error("Currency conversion failed",
			"amount", amount,
			"from", fromCurrency,
			"to", toCurrency,
			"error", err)
		return fmt.Errorf("conversion failed: %w", err)
	}

	a.logger.Info("Currency conversion successful",
		"from_amount", result.From.Amount,
		"from_currency", result.From.Currency,
		"to_amount", result.To.Amount,
		"to_currency", result.To.Currency,
		"rate", result.Rate)

	a.displayResult(result)
	return nil
}

func (a *CLI) showUsage() error {
	fmt.Fprintf(os.Stderr, "Usage: %s <amount> <from_currency> <to_currency>\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "\nExamples:\n")
	fmt.Fprintf(os.Stderr, "  %s 1000 USD BTC\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "\nNote: This uses CoinMarketCap sandbox API with test data.\n")
	fmt.Fprintf(os.Stderr, "      Crypto symbols are randomly generated for testing.\n")
	return fmt.Errorf("invalid number of arguments")
}

func (a *CLI) displayResult(result *entities.ConversionResult) {
	fmt.Printf("%.8f %s = %.8f %s\n",
		result.From.Amount, result.From.Currency,
		result.To.Amount, result.To.Currency)

	fmt.Printf("Exchange rate: 1 %s = %.8f %s\n",
		result.From.Currency, result.Rate, result.To.Currency)
}
