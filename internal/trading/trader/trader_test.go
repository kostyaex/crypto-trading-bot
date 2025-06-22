package trader

import (
	"crypto-trading-bot/internal/models"
	"crypto-trading-bot/internal/service/marketdata"
	"crypto-trading-bot/internal/service/marketdata/sources"
	"crypto-trading-bot/internal/trading/dispatcher"
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

	source := sources.NewHistoricalSource(testData)
	broadcaster := marketdata.NewBroadcaster(source.GetMarketDataCh())
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

func Test_runStrategyForSource(t *testing.T) {

	//setup := NewTestSetup()

	// Подготовка тестовых данных
	now, _ := time.Parse(time.RFC3339, "2025-01-05T00:00:00Z") //time.Now()
	testData := []*models.MarketData{
		{Timestamp: now, OpenPrice: 100, Volume: 10},
		{Timestamp: now.Add(time.Minute), OpenPrice: 105, Volume: 20},
		{Timestamp: now.Add(2 * time.Minute), OpenPrice: 103, Volume: 15},
		{Timestamp: now.Add(3 * time.Minute), OpenPrice: 107, Volume: 25},
		{Timestamp: now.Add(4 * time.Minute), OpenPrice: 108, Volume: 30},
	}

	source := NewMockMarketDataSource(testData)

	strategy, err := models.NewStrategy("test-strategy", "",
		models.StrategySettings{
			Symbol:   "BTCUSDT",
			Interval: "1m",
			Waves:    models.StrategyWavesSettings{BlockSize: 3, Overlap: 2, NumClusters: 2}})
	if err != nil {
		t.Errorf("Ошибка создания новой стратегии %v", err)
		return
	}

	disp := dispatcher.NewSignalDispatcher(
		&dispatcher.VolumeTrendRule{MinVolumeChangePercent: 10},
	)

	// Вызов тестируемой функции
	err = runStrategyForSource(*strategy, source, disp)
	if err != nil {
		t.Errorf("Ошибка выполнения RunStrategyForSource %v", err)
		return
	}

	// // Проверяем, что результаты содержат ожидаемое количество волн
	// assert.NotEmpty(t, results)
	// assert.Contains(t, results[0].Log, "Waves:")
}
