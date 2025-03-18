package app

import (
	"crypto-trading-bot/internal/data"
)

// MarketDataLoadedEvent представляет событие загрузки рыночных данных
type MarketDataLoadedEvent struct {
	MarketData []*data.MarketData
}

// Type возвращает тип события
func (e MarketDataLoadedEvent) Type() string {
	return "MarketDataLoaded"
}

// Payload возвращает полезную нагрузку события
func (e MarketDataLoadedEvent) Payload() interface{} {
	return e.MarketData
}
