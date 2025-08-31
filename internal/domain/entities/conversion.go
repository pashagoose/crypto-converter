package entities

import "time"

type ConversionResult struct {
	From      *Money
	To        *Money
	Rate      float64
	Timestamp time.Time
}

func NewConversionResult(from, to *Money, rate float64) *ConversionResult {
	return &ConversionResult{
		From:      from,
		To:        to,
		Rate:      rate,
		Timestamp: time.Now(),
	}
}
