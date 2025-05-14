package services

import (
	"crypto-trading-bot/internal/models"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_GroupingAndBroadcasting(t *testing.T) {
	// Создаем тестовые данные
	strategiesSettings := []models.StrategySettings{
		{Symbol: "BTCUSDT", Interval: "1m", Waves: models.StrategyWavesSettings{BlockSize: 5, Overlap: 4}},
		{Symbol: "BTCUSDT", Interval: "1m", Waves: models.StrategyWavesSettings{BlockSize: 10, Overlap: 2}},
		{Symbol: "ETHUSDT", Interval: "5m", Waves: models.StrategyWavesSettings{BlockSize: 8, Overlap: 3}},
		{Symbol: "ETHUSDT", Interval: "5m", Waves: models.StrategyWavesSettings{BlockSize: 6, Overlap: 1}},
		{Symbol: "BTCUSDT", Interval: "5m", Waves: models.StrategyWavesSettings{BlockSize: 7, Overlap: 3}},
	}

	strategies := make([]models.Strategy, 0)
	for n, settings := range strategiesSettings {
		strategy, err := models.NewStrategy("strat"+strconv.Itoa(n), "", settings)
		if err != nil {
			t.Errorf("Ошибка создания новой стратегии %v", err)
		}
		strategies = append(strategies, *strategy)
	}

	grouped := groupStrategiesBySymbolInterval(strategies)

	expectedGroups := map[string]int{
		"BTCUSDT|1m": 2,
		"ETHUSDT|5m": 2,
		"BTCUSDT|5m": 1,
	}

	assert.Equal(t, len(expectedGroups), len(grouped))

	for key, expectedCount := range expectedGroups {
		t.Run("Group_"+key, func(t *testing.T) {
			group, ok := grouped[key]
			assert.True(t, ok)
			assert.Equal(t, expectedCount, len(group))
		})
	}

	// Проверяем работу мультикастера
	testData := []*models.MarketData{
		{Timestamp: time.Now(), OpenPrice: 30000, Volume: 100},
		{Timestamp: time.Now().Add(time.Minute), OpenPrice: 30100, Volume: 150},
	}

	source := NewHistoricalSource(testData)
	broadcaster := NewBroadcaster(source.GetMarketDataCh())
	broadcaster.Start()

	sub1 := broadcaster.Subscribe()
	sub2 := broadcaster.Subscribe()

	count1 := 0
	for range sub1 {
		count1++
	}

	count2 := 0
	for range sub2 {
		count2++
	}

	broadcaster.Wait()

	assert.Equal(t, len(testData), count1)
	assert.Equal(t, len(testData), count2)

}
