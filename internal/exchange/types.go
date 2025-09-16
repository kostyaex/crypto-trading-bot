package exchange

import "time"

type Candle struct {
	Symbol    string
	Timestamp time.Time
	Open      float64
	High      float64
	Low       float64
	Close     float64
	Volume    float64
}

type Order struct {
	ID     string
	Symbol string
	Side   string // "buy", "sell"
	Type   string // "limit", "market"
	Price  float64
	Amount float64
	Status string // "open", "filled", "canceled"
}

type Position struct {
	Symbol        string
	Side          string // "long", "short"
	Size          float64
	Entry         float64
	UnrealizedPnL float64
}

type Balance struct {
	Asset  string
	Free   float64
	Locked float64
}

type CommandID string // Уникальный ID команды для получения результата позже
