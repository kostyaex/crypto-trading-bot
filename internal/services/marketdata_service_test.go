package services

import (
	"crypto-trading-bot/internal/models"
	"crypto-trading-bot/internal/utils"
	"fmt"
	"testing"
	"time"
)

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

func Test_marketDataService_RunBacktesting(t *testing.T) {

	setup := NewTestSetup()

	startTime := time.Now().Add(-time.Hour)
	// пока не важно какое время
	err := setup.marketDataService.RunBacktesting(startTime, startTime)
	if err != nil {
		t.Errorf("Ошибка: %s", err)
	}
}

func Test_marketDataService_GetIntervals(t *testing.T) {

	setup := NewTestSetup()
	marketDataCh := make(chan *models.MarketData)
	//defer close(marketDataCh)

	now := time.Now().Truncate(5 * time.Minute)

	marketData := []*models.MarketData{
		&models.MarketData{
			Timestamp: now.Add(-10 * time.Minute),
			Exchange:  "binance",
			Symbol:    "BTCUSDT",
		},
		&models.MarketData{
			Timestamp: now.Add(-7 * time.Minute),
			Exchange:  "binance",
			Symbol:    "BTCUSDT",
		},
		&models.MarketData{
			Timestamp: now.Add(-3 * time.Minute),
			Exchange:  "binance",
			Symbol:    "BTCUSDT",
		},
	}

	// Передаем тестовые данные в канал биржевых данных
	go func() {
		for _, marketDataItem := range marketData {
			marketDataCh <- marketDataItem
		}
		close(marketDataCh)
	}()

	intervalsCh := setup.marketDataService.GetIntervals(marketDataCh)

	// получаем интервалы из канала
	for interval := range intervalsCh {
		fmt.Printf("Interval: %v - %d\n", interval.Start, len(interval.Records))
	}
}
