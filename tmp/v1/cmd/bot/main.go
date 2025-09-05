package main

import (
	"context"
	"crypto-trading-bot/internal/components/sources"
	"crypto-trading-bot/internal/engine"
)

func main() {

	data := sources.GenerateTestMarketData(10)

	engine.RegisterComponent("source", sources.NewMockMarketDataSourceFactory(data))

	s := &sources.TickSource{}
	s.Run(context.Background())
}
