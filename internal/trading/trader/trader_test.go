package trader

import (
	"crypto-trading-bot/internal/core/config"
	"crypto-trading-bot/internal/core/logger"
	"crypto-trading-bot/internal/core/repositories"
	"crypto-trading-bot/internal/models"
	"crypto-trading-bot/internal/service/exchange"
	"crypto-trading-bot/internal/service/marketdata"
	"crypto-trading-bot/internal/service/marketdata/sources"
	"crypto-trading-bot/internal/trading/dispatcher"
	"log"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type TestSetup struct {
	//exchanges         []exchange.Exchange
	//exchangeService   exchange.ExchangeService
	//marketDataService marketdata.MarketDataService
	//strategyService   StrategyService
	traderService TraderService
}

// NewTestSetup инициализирует все необходимые зависимости для тестов
func NewTestSetup() *TestSetup {
	cfg := config.LoadConfig()

	db, err := repositories.NewDB(cfg)

	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	logger := logger.NewLogger(cfg.Logging.Level)

	exchanges := []exchange.Exchange{
		exchange.NewBinance(cfg.Binance.APIKey, cfg.Binance.APISecret, logger),
		//exchange.NewHuobi(cfg.Huobi.APIKey, cfg.Huobi.APISecret, logger),
	}

	repo := repositories.NewRepository(db, logger)

	exchangeService := exchange.NewEchangeService(repo, logger, exchanges)

	marketDataService := marketdata.NewMarketDataService(cfg, repo, logger, exchanges, exchangeService)

	//strategyService := NewStrategyService(repo)

	traderService := NewTraderService(logger, marketDataService)

	return &TestSetup{
		//exchanges:         exchanges,
		//exchangeService:   exchangeService,
		//marketDataService: marketDataService,
		//strategyService:   strategyService,
		traderService: traderService,
	}

}

func Test_GroupingAndBroadcasting(t *testing.T) {
	// Создаем тестовые данные
	strategiesSettings := []models.StrategySettings{
		{Symbol: "BTCUSDT", Interval: "1m"},
		{Symbol: "BTCUSDT", Interval: "1m"},
		{Symbol: "ETHUSDT", Interval: "5m"},
		{Symbol: "ETHUSDT", Interval: "5m"},
		{Symbol: "BTCUSDT", Interval: "5m"},
	}

	// Пример конфигурации
	config := map[string]interface{}{
		"type":         "simple",
		"value_factor": 1.0,
		"time_factor":  100.0,
	}

	strategies := make([]models.Strategy, 0)
	for n, settings := range strategiesSettings {
		strategy, err := models.NewStrategy("strat"+strconv.Itoa(n), "", settings, config)
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

func Test_runStrategyForSource(t *testing.T) {

	//setup := NewTestSetup()

	// Подготовка тестовых данных
	now, _ := time.Parse(time.RFC3339, "2025-01-05T00:00:00Z") //time.Now()
	testData := []*models.MarketData{
		{Timestamp: now, ClosePrice: 100, Volume: 10, BuyVolume: 5, SellVolume: 5, Symbol: "TEST"},
		{Timestamp: now, ClosePrice: 100, Volume: 10, BuyVolume: 5, SellVolume: 5, Symbol: "TEST"},
		{Timestamp: now, ClosePrice: 100, Volume: 10, BuyVolume: 5, SellVolume: 5, Symbol: "TEST"},
		{Timestamp: now, ClosePrice: 100, Volume: 10, BuyVolume: 5, SellVolume: 5, Symbol: "TEST"},
		{Timestamp: now, ClosePrice: 100, Volume: 10, BuyVolume: 5, SellVolume: 5, Symbol: "TEST"},
		{Timestamp: now.Add(time.Minute), ClosePrice: 100, Volume: 10, BuyVolume: 5, SellVolume: 5, Symbol: "TEST"},
		{Timestamp: now.Add(time.Minute), ClosePrice: 100, Volume: 10, BuyVolume: 5, SellVolume: 5, Symbol: "TEST"},
		{Timestamp: now.Add(time.Minute), ClosePrice: 100, Volume: 10, BuyVolume: 5, SellVolume: 5, Symbol: "TEST"},
		{Timestamp: now.Add(time.Minute), ClosePrice: 100, Volume: 10, BuyVolume: 5, SellVolume: 5, Symbol: "TEST"},
		{Timestamp: now.Add(time.Minute), ClosePrice: 100, Volume: 10, BuyVolume: 5, SellVolume: 5, Symbol: "TEST"},
		{Timestamp: now.Add(2 * time.Minute), ClosePrice: 100, Volume: 20, BuyVolume: 5, SellVolume: 15, Symbol: "TEST"},
		{Timestamp: now.Add(2 * time.Minute), ClosePrice: 100, Volume: 20, BuyVolume: 5, SellVolume: 15, Symbol: "TEST"},
		{Timestamp: now.Add(2 * time.Minute), ClosePrice: 100, Volume: 20, BuyVolume: 5, SellVolume: 15, Symbol: "TEST"},
		{Timestamp: now.Add(2 * time.Minute), ClosePrice: 100, Volume: 20, BuyVolume: 5, SellVolume: 15, Symbol: "TEST"},
		{Timestamp: now.Add(2 * time.Minute), ClosePrice: 100, Volume: 20, BuyVolume: 5, SellVolume: 15, Symbol: "TEST"},
	}

	source := NewMockMarketDataSource(testData)

	// Пример конфигурации
	config := map[string]interface{}{
		"type":         "simple",
		"value_factor": 1.0,
		"time_factor":  100.0,
	}

	strategy, err := models.NewStrategy("test-strategy", "",
		models.StrategySettings{
			Symbol:   "BTCUSDT",
			Interval: "1m",
			Cluster:  models.ClusterSettings{NumClusters: 1, Block: 5, Interval: "5m"}},
		config)

	if err != nil {
		t.Errorf("Ошибка создания новой стратегии %v", err)
		return
	}

	disp := dispatcher.NewDispatcher(
		&dispatcher.VolumeTrendRule{MinVolumeChangePercent: 10},
	)

	disp.Register(dispatcher.SignalBuy, &dispatcher.LoggerHandler{})
	disp.Register(dispatcher.SignalSell, &dispatcher.LoggerHandler{})
	//disp.Register(dispatcher.SignalHold, &dispatcher.LoggerHandler{})

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

func TestTraderService_RunBacktesting(t *testing.T) {
	setup := NewTestSetup()

	// Пример конфигурации
	config := map[string]interface{}{
		"type":         "simple",
		"value_factor": 200.0,
		"time_factor":  10.0,
	}

	strategy, err := models.NewStrategy("test-strategy", "",
		models.StrategySettings{
			Symbol:   "BTCUSDT",
			Interval: "1s",
			Cluster:  models.ClusterSettings{NumClusters: 10, Block: 300, Interval: "5m"}},
		config)

	if err != nil {
		t.Errorf("Ошибка создания новой стратегии %v", err)
		return
	}

	startTime, _ := time.Parse(time.DateTime, "2025-07-01 00:00:00")
	stopTime, _ := time.Parse(time.DateTime, "2025-07-01 23:59:59") //

	setup.traderService.RunBacktesting(strategy, startTime, stopTime)
}
