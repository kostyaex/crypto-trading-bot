package types

import (
	"time"
)

// ЗДЕСЬ НЕ ДОЛЖНО БЫТЬ ДИРЕКТИВ db!!

// MarketDataStatus представляет модель MarketDataStatus.
type MarketDataStatus struct {
	ID         int       `db:"id"`
	Exchange   string    `db:"exchange"`
	Symbol     string    `db:"symbol"`
	TimeFrame  string    `db:"time_frame"`
	Active     bool      `db:"active"`
	ActualTime time.Time `db:"actual_time"`
	Status     string    `db:"status"`
}
