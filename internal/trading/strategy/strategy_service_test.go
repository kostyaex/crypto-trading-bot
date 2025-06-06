package strategy

import (
	"crypto-trading-bot/internal/core/config"
	"crypto-trading-bot/internal/core/logger"
	"crypto-trading-bot/internal/core/repositories"
	"crypto-trading-bot/internal/service/exchange"
	"crypto-trading-bot/internal/service/marketdata"
	"fmt"
	"log"
	"testing"
)

type TestSetup struct {
	exchanges         []exchange.Exchange
	exchangeService   exchange.ExchangeService
	marketDataService marketdata.MarketDataService
	strategyService   StrategyService
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

	strategyService := NewStrategyService(repo)

	return &TestSetup{
		exchanges:         exchanges,
		exchangeService:   exchangeService,
		marketDataService: marketDataService,
		strategyService:   strategyService,
	}

}

func Test_strategyService_GetActiveStrategies(t *testing.T) {
	setup := NewTestSetup()

	strategies, err := setup.strategyService.GetActiveStrategies()
	if err != nil {
		t.Errorf("strategyService.GetActiveStrategies() error = %v", err)
		return
	}

	for _, strategy := range strategies {
		settings, err := strategy.Settings()
		if err != nil {
			t.Errorf("strategyService.GetActiveStrategies() error = %v", err)
			return
		}

		fmt.Printf("Strategy: %s %v\n", strategy.Name, settings)
	}
}
