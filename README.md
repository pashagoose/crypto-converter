# Crypto Converter

A simple command-line utility for currency conversion using the CoinMarketCap API.

## Quick Start

### 1. Install Dependencies

```bash
go mod download
```

### 2. API Key Setup

Copy the example configuration file:
```bash
cp env.example .env
```

Edit the `.env` file and add your CoinMarketCap API key:
```bash
# For testing (sandbox)
COINMARKETCAP_API_KEY=b54bcf4d-1bca-4e8e-9a24-22ff2c3d462c
COINMARKETCAP_URL=https://sandbox-api.coinmarketcap.com

# Log level (DEBUG, INFO, WARN, ERROR)
LOG_LEVEL=INFO
```

> **Note**: The provided API key is for sandbox testing. For production, get your own key at [coinmarketcap.com/api](https://coinmarketcap.com/api/).

### 3. Build Application

```bash
# Using Taskfile (recommended)
task build

# Or directly with Go
go build -o bin/crypto-converter cmd/main.go
```

### 4. Usage

```bash
# Basic syntax
./bin/crypto-converter <amount> <from_currency> <to_currency>

# Usage examples
./bin/crypto-converter 1000 USD BTC    # $1000 to Bitcoin
./bin/crypto-converter 0.5 BTC USD     # 0.5 Bitcoin to USD
./bin/crypto-converter 100 USD ETH     # $100 to Ethereum
```

## Output Examples

```bash
$ ./bin/crypto-converter 1000 USD BTC
1000.00000000 USD = 0.02285156 BTC
Exchange rate: 1 USD = 0.00002285 BTC
```

```bash
$ ./bin/crypto-converter 0.1 BTC USD
0.10000000 BTC = 4375.23000000 USD
Exchange rate: 1 BTC = 43752.30000000 USD
```

## Supported Operations

- ✅ **Fiat → Crypto**: USD → BTC, EUR → ETH, etc.
- ✅ **Crypto → Fiat**: BTC → USD, ETH → EUR, etc.  
- ✅ **Crypto → Crypto**: BTC → ETH, ETH → BTC, etc.
- ❌ **Fiat → Fiat**: USD → EUR (not supported)
