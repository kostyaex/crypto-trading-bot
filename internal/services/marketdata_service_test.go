package services

import (
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
	setup.marketDataService.RunBacktesting(startTime, startTime)

}
