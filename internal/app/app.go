package app

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/lib/pq"
)

type App struct {
	cfg *Config
	// db         *sqlx.DB
	// exchanges  []exchange.Exchange
	// trader     *trading.Trader
	// webServer  *web.Server
}

func NewApp() *App {
	cfg := LoadConfig()

	fmt.Printf("cfg: %s\n", cfg.Postgres.Host)

	// db, err := sqlx.Connect("postgres", cfg.PostgresDSN)
	// if err != nil {
	//     log.Fatalf("Failed to connect to database: %v", err)
	// }

	// exchanges := []exchange.Exchange{
	//     exchange.NewBinance(cfg.BinanceAPIKey, cfg.BinanceAPISecret),
	//     exchange.NewHuobi(cfg.HuobiAPIKey, cfg.HuobiAPISecret),
	// }

	// repo := data.NewPostgresRepository(db)
	// trader := trading.NewTrader(repo, exchanges)
	// webServer := web.NewServer(cfg.WebPort, repo, trader)

	return &App{
		cfg: cfg,
		// db:         db,
		// exchanges:  exchanges,
		// trader:     trader,
		// webServer:  webServer,
	}

}

func (a *App) Run() error {
	//ctx, cancel := context.WithCancel(context.Background())
	_, cancel := context.WithCancel(context.Background())
	defer cancel()

	fmt.Println("App is running")

	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan
		cancel()
	}()

	// if err := a.webServer.Start(ctx); err != nil {
	//     return err
	// }

	return nil
}
