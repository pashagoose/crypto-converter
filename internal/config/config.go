package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

type Config struct {
	CoinMarketCapAPIKey string
	CoinMarketCapURL    string
	LogLevel            string
}

func Load() (*Config, error) {
	if err := loadEnvFile(); err != nil {
		fmt.Printf("Warning: %v\n", err)
	}

	config := &Config{
		CoinMarketCapAPIKey: os.Getenv("COINMARKETCAP_API_KEY"),
		CoinMarketCapURL:    getEnvWithDefault("COINMARKETCAP_URL", "https://sandbox-api.coinmarketcap.com"),
		LogLevel:            getEnvWithDefault("LOG_LEVEL", "ERROR"),
	}

	if err := config.validate(); err != nil {
		return nil, err
	}

	return config, nil
}

func loadEnvFile() error {
	currentDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}

	for {
		envPath := filepath.Join(currentDir, ".env")
		if _, err := os.Stat(envPath); err == nil {
			return godotenv.Load(envPath)
		}

		parentDir := filepath.Dir(currentDir)
		if parentDir == currentDir {
			break
		}
		currentDir = parentDir
	}

	return fmt.Errorf(".env file not found")
}

func getEnvWithDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func (c *Config) validate() error {
	if c.CoinMarketCapAPIKey == "" {
		return fmt.Errorf("COINMARKETCAP_API_KEY environment variable is required")
	}

	if c.CoinMarketCapURL == "" {
		return fmt.Errorf("COINMARKETCAP_URL cannot be empty")
	}

	return nil
}
