// источник торговых данных
package components

import (
	"crypto-trading-bot/internal/engine"
	"math"
	"math/rand"
	"time"
)

type DataSource struct {
	data           []engine.MarketData
	indexTimestamp map[time.Time]int // индекс записи по времени
}

func NewDataSource(data []engine.MarketData) *DataSource {
	s := &DataSource{
		data: data,
	}

	s.indexTimestamp = make(map[time.Time]int, len(data))

	for n, v := range data {
		s.indexTimestamp[v.Timestamp] = n
	}

	return s
}

func (c *DataSource) Mask() uint64 {
	return MaskDatasource
}

// Получить следующую отметку времени, если она есть
func (c *DataSource) NextPosition(timestamp time.Time) (time.Time, bool) {

	if timestamp.IsZero() {
		if len(c.data) > 0 {
			return c.data[0].Timestamp, true
		} else {
			return time.Time{}, false
		}
	}

	ind := c.indexTimestamp[timestamp]
	if ind < len(c.data)-1 {
		return c.data[ind+1].Timestamp, true
	} else {
		return time.Time{}, false
	}

}

// Генерация тестовых данных: тренд + шум
func GenerateTestCandles(n int) []engine.MarketData {
	candles := make([]engine.MarketData, n)
	base := 40000.0
	trend := 10.0
	noise := 2000.0

	rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 0; i < n; i++ {
		price := base + trend*float64(i) + noise*math.Sin(float64(i)/20)
		price += (rand.Float64() - 0.5) * 1000 // шум
		candles[i] = engine.MarketData{
			Timestamp:  time.Now().Add(time.Second * time.Duration(i)),
			ClosePrice: price,
		}
	}
	return candles
}
