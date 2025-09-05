package main

import (
	"context"
	"crypto-trading-bot/internal/components/sources"
	"crypto-trading-bot/internal/core/config"
	"crypto-trading-bot/internal/core/logger"
	"crypto-trading-bot/internal/core/repositories"
	"crypto-trading-bot/internal/engine"
	"crypto-trading-bot/internal/models"
	"crypto-trading-bot/internal/service/exchange"
	"crypto-trading-bot/internal/service/marketdata"
	"crypto-trading-bot/internal/service/series"
	"crypto-trading-bot/internal/trading/dispatcher"
	"crypto-trading-bot/internal/trading/trader"
	"crypto-trading-bot/pkg/types"
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"syscall"
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

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// // === Запускаем менеджер ===
	// if err := manager.LoadAndStartAll(); err != nil {
	//     log.Printf("Warning: failed to start some strategies: %v", err)
	// }

	// === Перехватываем сигналы ===
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		sig := <-c
		log.Printf("Received signal: %s. Shutting down gracefully...", sig)
		cancel() // ← отменяем контекст → остановка всех Runner'ов

		// Дополнительная задержка на завершение (опционально)
		time.Sleep(1 * time.Second)
		os.Exit(0)
	}()

	// ===  ===

	run(ctx, cancel)

	// // === Запускаем веб-сервер ===
	// r := gin.Default()
	// web.SetupRoutes(r, manager)

	// go func() {
	//     if err := r.Run(":8080"); err != nil {
	//         log.Printf("Server error: %v", err)
	//     }
	// }()

	// === Ждём отмены контекста ===
	<-ctx.Done()

	log.Println("Graceful shutdown...")

	// // Можно добавить дополнительную очистку
	// manager.StopAll() // явная остановка всех стратегий (если нужно)

	log.Println("Bye!")
}

func run(ctx context.Context, cancel context.CancelFunc) {
	basicServices := NewBasicServices()

	basicServices.logger.Debugf("Запуск бектеста...")

	// ---------------------------------------------------------------------

	strategyConf := `{
		"symbol": "BTCUSDT",
		"interval": "1s",
		"cluster": {
			"num_clusters": 5,
			"block": 300,
			"interval": "5m"
		},
		"series_builder": {
			"type": "simple",
			"value_factor": 1.0,
			"time_factor": 0.001
		}
	}`

	strategy := models.Strategy{
		ID:     1,
		Name:   "test-strategy",
		Config: json.RawMessage(strategyConf),
	}

	err := strategy.UpdateSettingsFromConf()

	if err != nil {
		basicServices.logger.Errorf("Ошибка создания новой стратегии %v", err)
		return
	}

	// ---------------------------------------------------------------------

	// disp := dispatcher.NewDispatcher(
	// 	&dispatcher.VolumeTrendRule{MinVolumeChangePercent: 10},
	// )
	// disp.Register(dispatcher.SignalBuy, &dispatcher.LoggerHandler{})
	// disp.Register(dispatcher.SignalSell, &dispatcher.LoggerHandler{})

	// Построение компоненты из конфига JSON:
	dispatcherConfigJSON := `{
  "rules": [
    {
      "type": "volume_trand",
      "config": {
        "min_volume_change_percent": 10.0
      }
    }
  ],
  "handlers": {
    "buy": [
      {
        "type": "logger",
        "config": {}
      }
    ],
    "sell": [
      {
        "type": "logger",
        "config": {}
      }
    ],
    "hold": []
  }
}`

	componentType := "dispatcher"
	engine.RegisterComponent(componentType, dispatcher.NewDispatcherFactory())

	//disp, err := dispatcher.NewDispatcherFromJSON([]byte(dispatcherCondig))
	dispatcherConfig, err := types.DeserializeConfig(componentType, []byte(dispatcherConfigJSON))
	if err != nil {
		basicServices.logger.Errorf("Ошибка разбора конфигурации диспетчера: %v", err)
		return
	}

	dispatcherComponent, err := engine.New(componentType, dispatcherConfig, log.Default())
	if err != nil {
		basicServices.logger.Errorf("Ошибка создания компоненты диспетчера: %v", err)
		return
	}

	// ---------------------------------------------------------------------

	strategySettings := strategy.Settings

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

	//source := sources.NewMockMarketDataSource(marketData)

	source := sources.NewHistoricalSource(marketData, ctx)

	// ---------------------------------------------------------------------

	seriesBuilder, err := series.NewSeriesBuilder(strategySettings.SeriesBuilderConfig)
	if err != nil {
		basicServices.logger.Errorf("Ошибка формирования построителя серий: %s\n", err)
		return
	}

	// ---------------------------------------------------------------------

	pipeline := &trader.Pipeline{
		Conf:          basicServices.conf,
		Logger:        basicServices.logger,
		Mode:          "backtest",
		DataSource:    source,
		Strategy:      strategy,
		SeriesBuilder: seriesBuilder,
		Dispatcher:    dispatcherComponent.Dispatcher,
	}

	// ---------------------------------------------------------------------

	// err = pipeline.Run(ctx)
	// if err != nil {
	// 	fmt.Printf("Ошибка при запуске пайплайна: %v\n", err)
	// }

	// Теперь запускаем пайплайн через Runner
	runner := trader.NewRunner(&strategy, pipeline)

	runner.Start(ctx)
	log.Printf("Runner запущен.")

	// Ожидаем завершения runner'а
	<-runner.Done()
	log.Printf("Runner выполнился. Завершаем работу.")
	cancel()
}
