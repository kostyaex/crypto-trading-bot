package sources

import (
	"crypto-trading-bot/internal/models"
	"math"
	"math/rand"
	"time"
)

type MockMarketDataSource struct {
	data []*models.MarketData
}

func NewMockMarketDataSource(data []*models.MarketData) *MockMarketDataSource {
	return &MockMarketDataSource{data: data}
}

func (m *MockMarketDataSource) GetMarketDataCh() <-chan *models.MarketData {
	ch := make(chan *models.MarketData)
	go func() {
		for _, item := range m.data {
			ch <- item
		}
		close(ch)
	}()
	return ch
}

func (m *MockMarketDataSource) Close() {}

// Генерация тестовых данных: тренд + шум
func GenerateTestMarketData(n int) []*models.MarketData {
	candles := make([]*models.MarketData, n)
	base := 40000.0
	trend := 10.0
	noise := 2000.0

	rand.New(rand.NewSource(time.Now().UnixNano()))

	now := time.Now()

	for i := 0; i < n; i++ {
		priceBase := base + trend*float64(i) + noise*math.Sin(float64(i)/20)
		price1 := priceBase + (rand.Float64()-0.5)*1000 // шум
		//price2 := priceBase + (rand.Float64()-0.5)*1000 // шум
		candles[i] = &models.MarketData{
			Timestamp:  now,
			ClosePrice: price1,
		}

		now = now.Add(time.Second)
	}
	return candles
}
