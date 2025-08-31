package main

import (
	"crypto-converter/internal/config"
	"crypto-converter/internal/domain/usecases"
	"crypto-converter/internal/infrastructure/repositories"
	"crypto-converter/internal/interfaces/cli"
	"crypto-converter/pkg/logger"
	"fmt"
	"os"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load configuration: %v\n", err)
		os.Exit(1)
	}

	appLogger := logger.New(cfg.LogLevel)

	appLogger.Info("Starting crypto-converter application")

	currencyRepo := repositories.NewCoinMarketCapRepository(cfg.CoinMarketCapAPIKey, cfg.CoinMarketCapURL, appLogger)
	convertUseCase := usecases.NewConvertCurrencyUseCase(currencyRepo, appLogger)
	app := cli.NewCLI(convertUseCase, appLogger)

	appLogger.Info("Dependencies initialized, starting CLI")

	if err := app.Run(os.Args); err != nil {
		appLogger.Error("Application failed", "error", err)
		os.Exit(1)
	}

	appLogger.Info("Application completed successfully")
}
