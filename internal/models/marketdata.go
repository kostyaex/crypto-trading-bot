package models

import (
	"time"
)

type MarketData struct {
	Symbol    string    `db:"symbol"`
	Price     float64   `db:"price"`
	Timestamp time.Time `db:"timestamp"`
}
