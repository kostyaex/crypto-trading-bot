package engine

import (
	"crypto-trading-bot/pkg/types"
	"log"
)

// Component — универсальная обёртка
type Component struct {
	Runnable         types.Runnable         `json:"-"`
	Processor        types.Processor        `json:"-"`
	MarketDataSource types.MarketDataSource `json:"-"`
	Dispatcher       types.Dispatcher       `json:"-"`
}

// ComponentFactory — функция, создающая компонент
type ComponentFactory func(config types.ComponentConfig, logger *log.Logger) (*Component, error)
