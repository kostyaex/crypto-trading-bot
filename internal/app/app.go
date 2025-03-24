package app

import (
	"context"
	"crypto-trading-bot/internal/config"
	"crypto-trading-bot/internal/exchange"
	"crypto-trading-bot/internal/repositories"
	"crypto-trading-bot/internal/services"
	"crypto-trading-bot/internal/utils"
	"crypto-trading-bot/internal/web"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

type App struct {
	cfg       *config.Config
	db        *repositories.DB
	repo      *repositories.Repository
	exchanges []exchange.Exchange
	//trader    *trading.Trader
	webServer      *web.Server
	scheduler      *Scheduler
	logger         *utils.Logger
	eventPublisher *EventPublisher
}

func NewApp() *App {
	cfg := config.LoadConfig()

	db, err := repositories.NewDB(cfg)

	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// if err := data.applyMigrations(cfg.PostgresDSN()); err != nil {
	//     log.Fatalf("Failed to apply migrations: %v", err)
	// }

	logger := utils.NewLogger(cfg.Logging.Level)

	exchanges := []exchange.Exchange{
		exchange.NewBinance(cfg.Binance.APIKey, cfg.Binance.APISecret, logger),
		//exchange.NewHuobi(cfg.Huobi.APIKey, cfg.Huobi.APISecret, logger),
	}

	repo := repositories.NewRepository(db, logger)

	strategyService := services.NewStrategyService(repo)
	behaviorTreeService := services.NewBehaviorTree(repo)
	indicatorService := services.NewIndicatorService(repo)
	marketDataService := services.NewMarketDataService(repo, logger)
	marketDataStatusService := services.NewMarketDataStatusService(repo, logger)
	excahngeService := services.NewEchangeService(repo, logger, exchanges, marketDataStatusService, marketDataService)

	//trader := trading.NewTrader(repo, exchanges, logger)
	webServer := web.NewServer(strconv.Itoa(cfg.Web.Port), repo, logger, excahngeService, strategyService, marketDataService)
	scheduler := NewScheduler(exchanges, logger)
	eventPublisher := NewEventPublisher()

	// Создание подписчика для анализа данных
	dataAnalysisSubscriber := NewDataAnalysisSubscriber(logger, eventPublisher, marketDataService, indicatorService)
	eventPublisher.Subscribe(dataAnalysisSubscriber)

	// Создание подписчика для обновления состояния стратегий
	strategyUpdateSubscriber := NewStrategyUpdateSubscriber(logger, strategyService, behaviorTreeService)
	eventPublisher.Subscribe(strategyUpdateSubscriber)

	return &App{
		cfg:       cfg,
		db:        db,
		repo:      repo,
		exchanges: exchanges,
		//trader:    trader,
		webServer:      webServer,
		scheduler:      scheduler,
		logger:         logger,
		eventPublisher: eventPublisher,
	}
}

func (a *App) Run() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan
		cancel()
	}()

	// Запуск планировщика
	a.scheduler.Start()
	defer a.scheduler.Stop()

	// // Добавление задачи для загрузки данных с бирж каждые 5 минут
	// task := NewDataFetchingTask(a.repo, a.exchanges, a.logger, a.eventPublisher)
	// _, err := a.scheduler.AddJob("@every 5m", task)
	// if err != nil {
	// 	return err
	// }

	if err := a.webServer.Start(ctx); err != nil {
		return err
	}

	return nil
}
