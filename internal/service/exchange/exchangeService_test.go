package exchange

import (
	"crypto-trading-bot/internal/core/config"
	"crypto-trading-bot/internal/core/logger"
	"crypto-trading-bot/internal/core/repositories"
	"log"
	"testing"
	"time"
)

type TestSetup struct {
	exchanges       []Exchange
	exchangeService ExchangeService
}

// NewTestSetup инициализирует все необходимые зависимости для тестов
func NewTestSetup() *TestSetup {
	cfg := config.LoadConfig()

	db, err := repositories.NewDB(cfg)

	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	logger := logger.NewLogger(cfg.Logging.Level)

	exchanges := []Exchange{
		NewBinance(cfg.Binance.APIKey, cfg.Binance.APISecret, logger),
		//exchange.NewHuobi(cfg.Huobi.APIKey, cfg.Huobi.APISecret, logger),
	}

	repo := repositories.NewRepository(db, logger)

	exchangeService := NewEchangeService(repo, logger, exchanges)

	return &TestSetup{
		exchanges:       exchanges,
		exchangeService: exchangeService,
	}

}

func Test_exchangeService_LoadData(t *testing.T) {

	setup := NewTestSetup()

	// Выбираем данные час назад - они должны быть

	exchange := setup.exchanges[0]
	symbol := "BTCUSDT"
	timeFrame := "5m"
	startTime := time.Now().Add(-time.Hour)

	marketData, lastTime, err := setup.exchangeService.LoadData(exchange, symbol, timeFrame, startTime)
	if err != nil {
		t.Errorf("exchangeService.LoadData() error = %v", err)
		return
	}

	if len(marketData) == 0 {
		t.Error("exchangeService.LoadData() не получены данные")
		return
	}

	t.Logf("Loaded data %v to %v", len(marketData), lastTime)

}

func Test_exchangeService_LoadData_Empty(t *testing.T) {
	// Здесь случай когда данных нет за указанное время. Время указываем на час вперед.

	setup := NewTestSetup()

	exchange := setup.exchanges[0]
	symbol := "BTCUSDT"
	timeFrame := "5m"
	startTime := time.Now().Add(time.Hour)

	marketData, lastTime, err := setup.exchangeService.LoadData(exchange, symbol, timeFrame, startTime)
	if err != nil {
		t.Errorf("exchangeService.LoadData() error = %v", err)
		return
	}

	if len(marketData) != 0 {
		t.Error("exchangeService.LoadData() не должно возвращать данные")
		return
	}

	t.Logf("Loaded data %v to %v", len(marketData), lastTime)

}
