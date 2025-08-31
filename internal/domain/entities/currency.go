package entities

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type CurrencyCode string

func NewCurrencyCode(code string) (CurrencyCode, error) {
	if code == "" {
		return "", errors.New("currency code cannot be empty")
	}

	normalized := strings.ToUpper(strings.TrimSpace(code))
	return CurrencyCode(normalized), nil
}

func (c CurrencyCode) String() string {
	return string(c)
}

type Money struct {
	Currency CurrencyCode
	Amount   float64
}

func NewMoney(currency CurrencyCode, amount float64) (*Money, error) {
	if currency == "" {
		return nil, errors.New("currency code cannot be empty")
	}

	if amount < 0 {
		return nil, errors.New("amount cannot be negative")
	}

	return &Money{
		Currency: currency,
		Amount:   amount,
	}, nil
}

func (m *Money) String() string {
	return fmt.Sprintf("%.8f %s", m.Amount, m.Currency)
}

func ParseAmount(amountStr string) (float64, error) {
	amount, err := strconv.ParseFloat(amountStr, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid amount format: %s", amountStr)
	}
	return amount, nil
}
