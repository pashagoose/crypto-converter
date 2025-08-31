package repositories

import (
	"context"
)

type CurrencyRepository interface {
	GetExchangeRate(ctx context.Context, from, to string) (float64, error)
}
