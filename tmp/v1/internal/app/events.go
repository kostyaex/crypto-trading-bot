package app

import (
	"crypto-trading-bot/pkg/types"
	"time"
)

// MarketDataLoadedEvent представляет событие загрузки рыночных данных
type MarketDataLoadedEvent struct {
	MarketData []*types.MarketData
}

// Type возвращает тип события
func (e MarketDataLoadedEvent) Type() string {
	return "MarketDataLoaded"
}

// Payload возвращает полезную нагрузку события
func (e MarketDataLoadedEvent) Payload() interface{} {
	return e.MarketData
}

// AnalysisCompletedEvent представляет событие завершения анализа данных
type AnalysisCompletedEvent struct {
	Symbol     string
	Indicators map[string]float64
	Timestamp  time.Time
}

// Type возвращает тип события
func (e AnalysisCompletedEvent) Type() string {
	return "AnalysisCompleted"
}

// Payload возвращает полезную нагрузку события
func (e AnalysisCompletedEvent) Payload() interface{} {
	return e
}
