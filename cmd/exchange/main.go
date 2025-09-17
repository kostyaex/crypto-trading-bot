package main

import (
	"context"
	"crypto-trading-bot/internal/exchange"
	"crypto-trading-bot/internal/exchange/exchanges/binance"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// Создаем экземпляр Binance
	ex := binance.NewBinanceExchange()

	// Подписываемся на свечи в реальном времени через WebSocket
	log.Println("Подписываемся на свечи BTCUSDT 1m...")
	ex.SubscribeCandles("BTCUSDT", "1m", func(candle exchange.Candle) {
		log.Printf("📈 [WS] Свеча: %s | %.2f - %.2f (объём: %.2f)",
			candle.Timestamp.Format("15:04:05"),
			candle.Open, candle.Close, candle.Volume)
	})

	// Асинхронно запрашиваем исторические свечи
	log.Println("Запрашиваем последние 5 свечей...")
	cmdID := ex.FetchCandlesAsync("BTCUSDT", "1m", 5)

	// Проверяем результат каждую секунду
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				if result, ok := ex.GetResult(cmdID); ok {
					if err, isError := result.(error); isError {
						log.Printf("❌ Ошибка при получении свечей: %v", err)
					} else if candles, isCandles := result.([]exchange.Candle); isCandles {
						log.Printf("✅ Получено %d свечей:", len(candles))
						for i, c := range candles {
							log.Printf("  %d) %s: O=%.2f H=%.2f L=%.2f C=%.2f V=%.2f",
								i+1, c.Timestamp.Format("15:04:05"), c.Open, c.High, c.Low, c.Close, c.Volume)
						}
						// После получения — отменяем тикер
						cancel()
					}
				}
			}
		}
	}()

	// Продолжаем работать, не блокируя основной поток
	log.Println("Бот запущен. Нажмите Ctrl+C для остановки...")

	// Ожидаем SIGINT (Ctrl+C)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	log.Println("Остановка бота...")
}
