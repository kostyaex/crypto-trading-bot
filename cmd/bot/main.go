package main

import (
	"context"
	"crypto-trading-bot/internal/components"
	"crypto-trading-bot/internal/core/config"
	"crypto-trading-bot/internal/core/logger"
	"crypto-trading-bot/internal/core/repositories"
	"crypto-trading-bot/internal/service/exchange"
	"crypto-trading-bot/internal/service/marketdata"
	"crypto-trading-bot/internal/service/source"
	"crypto-trading-bot/internal/types"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

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

	registry := initRegistry()

	histSourceJson := json.RawMessage(`{
		"symbol": "BTCUSDT",
		"interval": "1s",
		"start_time" : "2025-07-23T00:00:00Z",
		"end_time" : "2025-07-23T01:00:00Z"
	}`)

	comp, err := registry.Build("source", histSourceJson)
	if err != nil {
		fmt.Printf("Ошибка формирования компоненты %s", err)
		return
	}

	_source, err := source.NewHistoricalSource(basicServices.marketDataService, comp)
	if err != nil {
		fmt.Printf("Ошибка создания источника: %s", err)
		return
	}

	_source.Next()

}

func initRegistry() *components.ComponentRegistry {
	reg := components.NewComponentRegistry()

	reg.Register("source", func() types.Component {
		return &components.HistoricalSourceSettings{}
	})

	// Добавляй сюда новые компоненты — система сама их подхватит
	return reg
}

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
