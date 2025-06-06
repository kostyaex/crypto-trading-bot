package marketdata

import (
	"crypto-trading-bot/internal/core/config"
	"crypto-trading-bot/internal/core/logger"
	"crypto-trading-bot/internal/core/repositories"
	"crypto-trading-bot/internal/core/utils"
	"crypto-trading-bot/internal/service/exchange"
	"fmt"
	"log"
	"testing"
	"time"
)

type TestSetup struct {
	exchanges         []exchange.Exchange
	exchangeService   exchange.ExchangeService
	marketDataService MarketDataService
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

	marketDataService := NewMarketDataService(cfg, repo, logger, exchanges, exchangeService)

	return &TestSetup{
		exchanges:         exchanges,
		exchangeService:   exchangeService,
		marketDataService: marketDataService,
	}

}

func TestGetIntervalStart(t *testing.T) {
	// Пример использования
	now := time.Now()
	intervals := []string{"1m", "1h", "1d", "1w", "1M"}

	for _, interval := range intervals {
		start, end, nextStart, err := utils.GetIntervalBounds(now, interval)
		if err != nil {
			fmt.Printf("Ошибка для интервала %s: %v\n", interval, err)
			continue
		}
		fmt.Printf("Интервал %s:\n  Начало: %v\n  Конец: %v\n  Следующий: %v\n",
			interval, start.Format(time.RFC3339), end.Format(time.RFC3339), nextStart.Format(time.RFC3339))
	}
}
