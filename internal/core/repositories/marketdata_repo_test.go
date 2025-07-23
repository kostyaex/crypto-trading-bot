package repositories

import (
	"crypto-trading-bot/internal/core/config"
	"crypto-trading-bot/internal/core/logger"
	"fmt"
	"log"
	"testing"
)

type TestSetup struct {
	cfg  *config.Config
	db   *DB
	repo *Repository
}

func NewTestSetup() *TestSetup {
	cfg := config.LoadConfig()

	db, err := NewDB(cfg)

	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	logger := logger.NewLogger(cfg.Logging.Level)

	repo := NewRepository(db, logger)

	return &TestSetup{
		cfg:  cfg,
		db:   db,
		repo: repo,
	}

}

func Test_marketDataRepository_GetMarketDataStatusList(t *testing.T) {
	testSetup := NewTestSetup()
	list, err := testSetup.repo.MarketData.GetMarketDataStatusList()
	if err != nil {
		t.Errorf("GetMarketDataStatusList() error = %v", err)
		return
	}

	for _, item := range list {
		fmt.Printf("ActualTime %s - location: %s\n", item.ActualTime, item.ActualTime.Location().String())
	}
}
