package sources

import (
	"context"
	"crypto-trading-bot/internal/engine"
	"crypto-trading-bot/pkg/types"
	"fmt"
	"log"
	"math"
	"math/rand"
	"time"
)

type MockMarketDataSourceConfig struct {
	// Interval string `json:"interval"`
	// Source   string `json:"source"`
}

func (c MockMarketDataSourceConfig) GetComponentType() string { return "mock_marketdata_source" }

type MockMarketDataSourceComponent struct {
	config MockMarketDataSourceConfig
	logger *log.Logger
	data   []*types.MarketData
}

func NewMockMarketDataSourceComponent(
	config types.ComponentConfig,
	logger *log.Logger,
	data []*types.MarketData,
) (*MockMarketDataSourceComponent, error) {

	typedConfig, ok := config.(MockMarketDataSourceConfig)
	if !ok {
		return nil, fmt.Errorf("invalid config type for sampling")
	}

	return &MockMarketDataSourceComponent{
		config: typedConfig,
		logger: logger,
		data:   data,
	}, nil
}

// Фабрика — возвращаем функцию, а не сразу создаем
func NewMockMarketDataSourceFactory(data []*types.MarketData) engine.ComponentFactory {
	return func(config types.ComponentConfig, logger *log.Logger) (*engine.Component, error) {
		comp, err := NewMockMarketDataSourceComponent(config, logger, data)
		if err != nil {
			return nil, err
		}
		return &engine.Component{MarketDataSource: comp}, nil
	}
}

func (c *MockMarketDataSourceComponent) Run(ctx context.Context, output chan<- *types.MarketData) error {

	for _, item := range c.data {
		output <- item
	}
	close(output)

	return nil

	// ticker := time.NewTicker(parseInterval(c.config.Interval))
	// defer ticker.Stop()

	// for {
	// 	select {
	// 	case <-ctx.Done():
	// 		return ctx.Err()
	// 	case <-ticker.C:
	// 		candle := c.fetchCandle() // условная функция
	// 		if c.output != nil {
	// 			select {
	// 			case c.output <- candle:
	// 			case <-ctx.Done():
	// 				return ctx.Err()
	// 			}
	// 		}
	// 	}
	// }
}

// func NewMockMarketDataSource(data []*types.MarketData) *MockMarketDataSource {
// 	return &MockMarketDataSource{data: data}
// }

// func (m *MockMarketDataSource) GetMarketDataCh() <-chan *types.MarketData {
// 	ch := make(chan *types.MarketData)
// 	go func() {
// 		for _, item := range m.data {
// 			ch <- item
// 		}
// 		close(ch)
// 	}()
// 	return ch
// }

// func (m *MockMarketDataSource) Close() {}

// Генерация тестовых данных: тренд + шум
func GenerateTestMarketData(n int) []*types.MarketData {
	candles := make([]*types.MarketData, n)
	base := 40000.0
	trend := 10.0
	noise := 2000.0

	rand.New(rand.NewSource(time.Now().UnixNano()))

	now := time.Now()

	for i := 0; i < n; i++ {
		priceBase := base + trend*float64(i) + noise*math.Sin(float64(i)/20)
		price1 := priceBase + (rand.Float64()-0.5)*1000 // шум
		//price2 := priceBase + (rand.Float64()-0.5)*1000 // шум
		candles[i] = &types.MarketData{
			Timestamp:  now,
			ClosePrice: price1,
		}

		now = now.Add(time.Second)
	}
	return candles
}
