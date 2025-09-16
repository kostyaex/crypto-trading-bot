package exchange

import (
	"context"
	"crypto-trading-bot/internal/engine"
	"fmt"
	"testing"
	"time"
)

func TestFetchDataExecutor(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	executor := NewFetchDataExecutor(ctx)

	commands := []FetchDataCommand{
		{Symbol: "BTCUSDT", Interval: "1m"},
		{Symbol: "ETHUSDT", Interval: "1m"},
	}

	for _, cmd := range commands {
		executor.SubmitCommand(cmd)
		fmt.Printf("📤 Submitted: %v\n", cmd)
	}

	// Основной цикл стратегии — НЕ БЛОКИРУЕТСЯ!
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			fmt.Println("\n--- Main loop tick ---")

			// Проверяем результаты — НЕ БЛОКИРУЕМСЯ!
			results := executor.GetAllResults()
			for cmd, res := range results {
				if res.Error != nil {
					fmt.Printf("❌ Error: %s %s | %v\n", cmd.Symbol, cmd.Interval, res.Error)
				} else {
					fmt.Printf("✅ %s: %s %s\n",
						res.Command.Symbol, res.Command.Interval, res.Timestamp.Format(engine.TimeFormat))
				}
			}

			// Можно здесь принять решение: например, если ордер заполнился — запустить следующий
			// или сбросить позицию, если что-то не так

		case <-time.After(8 * time.Second):
			fmt.Println("\n🛑 Stopping after 8 seconds...")
			executor.Close()
			return
		}
	}
}
