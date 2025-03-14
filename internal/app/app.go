package app

import (
	"context"
	"crypto-trading-bot/internal/data"
	"crypto-trading-bot/internal/exchange"
	"crypto-trading-bot/internal/utils"
	"crypto-trading-bot/internal/web"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	_ "github.com/lib/pq"
)

type App struct {
	cfg       *Config
	db        *data.DB
	repo      *data.PostgresRepository
	exchanges []exchange.Exchange
	//trader    *trading.Trader
	webServer *web.Server
	scheduler *Scheduler
	logger    *utils.Logger
}

func NewApp() *App {
	cfg := LoadConfig()

	db, err := data.NewDB(cfg.PostgresDSN())

	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	logger := utils.NewLogger(cfg.Logging.Level)

	exchanges := []exchange.Exchange{
		exchange.NewBinance(cfg.Binance.APIKey, cfg.Binance.APISecret, logger),
		//exchange.NewHuobi(cfg.Huobi.APIKey, cfg.Huobi.APISecret, logger),
	}

	repo := data.NewPostgresRepository(db, logger)
	//trader := trading.NewTrader(repo, exchanges, logger)
	webServer := web.NewServer(strconv.Itoa(cfg.Web.Port), repo)
	scheduler := NewScheduler(repo, exchanges, logger)

	return &App{
		cfg:       cfg,
		db:        db,
		repo:      repo,
		exchanges: exchanges,
		//trader:    trader,
		webServer: webServer,
		scheduler: scheduler,
		logger:    logger,
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

	// Добавление задачи для загрузки данных с бирж каждые 5 минут
	task := NewDataFetchingTask(a.repo, a.exchanges, a.logger)
	_, err := a.scheduler.AddJob("@every 1m", task)
	if err != nil {
		return err
	}

	if err := a.webServer.Start(ctx); err != nil {
		return err
	}

	return nil
}
