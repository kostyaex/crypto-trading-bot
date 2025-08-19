package main

import (
	"crypto-trading-bot/internal/core/config"
	"crypto-trading-bot/internal/core/logger"
	"crypto-trading-bot/internal/core/repositories"
	"crypto-trading-bot/internal/models"
	"crypto-trading-bot/internal/service/exchange"
	"crypto-trading-bot/internal/service/marketdata"
	"crypto-trading-bot/internal/service/marketdata/sources"
	"crypto-trading-bot/internal/service/series"
	"crypto-trading-bot/internal/trading/dispatcher"
	"crypto-trading-bot/internal/trading/trader"
	"fmt"
	"log"
	"time"
)

type basicServices struct {
	conf              *config.Config
	logger            *logger.Logger
	marketDataService marketdata.MarketDataService
}

func NewBasicServices() basicServices {

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

	return basicServices{
		conf:              cfg,
		logger:            logger,
		marketDataService: marketDataService,
	}
}

func main() {

	basicServices := NewBasicServices()

	basicServices.logger.Debugf("Запуск бектеста...")

	// ---------------------------------------------------------------------

	// Пример конфигурации построения сериий
	config := map[string]interface{}{
		"type":         "simple",
		"value_factor": 1.0,
		"time_factor":  0.001,
	}

	// ---------------------------------------------------------------------
	strategySettings := models.StrategySettings{
		Symbol:              "BTCUSDT",
		Interval:            "1s",
		Cluster:             models.ClusterSettings{NumClusters: 5, Block: 300, Interval: "5m"},
		SeriesBuilderConfig: config,
	}
	strategy, err := models.NewStrategy("test-strategy", "", strategySettings)

	if err != nil {
		basicServices.logger.Errorf("Ошибка создания новой стратегии %v", err)
		return
	}

	// ---------------------------------------------------------------------

	disp := dispatcher.NewDispatcher(
		&dispatcher.VolumeTrendRule{MinVolumeChangePercent: 10},
	)
	disp.Register(dispatcher.SignalBuy, &dispatcher.LoggerHandler{})
	disp.Register(dispatcher.SignalSell, &dispatcher.LoggerHandler{})

	// ---------------------------------------------------------------------

	// Подготовка тестовых данных
	startTime, _ := time.Parse(time.DateTime, "2025-07-23 00:00:00")
	stopTime, _ := time.Parse(time.DateTime, "2025-07-23 23:59:59") //
	// Получить данные за период
	marketData, err := basicServices.marketDataService.GetMarketDataPeriod(strategySettings.Symbol, strategySettings.Interval, startTime, stopTime)
	// marketData, err := s.GetMarketData("BTCUSDT", 1000)
	if err != nil {
		basicServices.logger.Errorf("Ошибка получения торговых данных: %s\n", err)
		return
	}

	basicServices.logger.Debugf("Получено данных: %d\n", len(marketData))

	//testData := sources.GenerateTestMarketData(10)

	source := sources.NewMockMarketDataSource(marketData)

	// ---------------------------------------------------------------------

	seriesBuilder, err := series.NewSeriesBuilder(strategy.SeriesBuilderConfig)
	if err != nil {
		basicServices.logger.Errorf("Ошибка формирования построителя серий: %s\n", err)
		return
	}

	// ---------------------------------------------------------------------

	pipeline := trader.Pipeline{
		Conf:          basicServices.conf,
		Logger:        basicServices.logger,
		Mode:          "backtest",
		DataSource:    source,
		Strategy:      *strategy,
		SeriesBuilder: seriesBuilder,
		Dispatcher:    disp,
	}

	err = pipeline.Run()
	if err != nil {
		fmt.Printf("Ошибка при запуске пайплайна: %v\n", err)
	}
}
