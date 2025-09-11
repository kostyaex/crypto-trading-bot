// Пример использования ECS
package main

import (
	"context"
	"crypto-trading-bot/internal/engine/components"
	"crypto-trading-bot/internal/engine/ecsx"
	"crypto-trading-bot/internal/engine/systems"
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/andygeiss/ecs"
)

const ModeLive = "live"
const ModeBacktest = "backtest"

func main() {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// // === Запускаем менеджер ===
	// if err := manager.LoadAndStartAll(); err != nil {
	//     log.Printf("Warning: failed to start some strategies: %v", err)
	// }

	// === Перехватываем сигналы ===
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		sig := <-c
		log.Printf("Получен сигнал: %s. Отправляем сигнал на завершение приложения...", sig)
		cancel() // ← отправка сигнала Done

		// Дополнительная задержка на завершение (опционально)
		time.Sleep(1 * time.Second)
		os.Exit(0)
	}()

	// =====================================================================

	em := ecs.NewEntityManager()

	// // компонент с текущим временем
	// em.Add(ecs.NewEntity("currenttime", []ecs.Component{
	// 	components.NewCurrentTime(),
	// }))

	sm := ecs.NewSystemManager()

	//currentTime := resources.NewCurrentTime()

	fmt.Println("🔄 Запуск...")
	sm.Add(systems.NewStopSystem(ctx))
	//candles := generateTestCandles(10)
	//sm.Add(systems.NewHistoricalTimeUpdateSystem(currentTime, candles))
	//sm.Add(systems.NewTimeUpdateSystem(currentTime))
	//sm.Add(systems.NewMonitoringSystem(currentTime))

	de := ecsx.NewCustomEngine(em, sm)
	de.Setup()
	defer de.Teardown()

	de.Run()

}

// Генерация тестовых данных: тренд + шум
func generateTestCandles(n int) []components.Candle {
	candles := make([]components.Candle, n)
	base := 40000.0
	trend := 10.0
	noise := 2000.0

	rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 0; i < n; i++ {
		price := base + trend*float64(i) + noise*math.Sin(float64(i)/20)
		price += (rand.Float64() - 0.5) * 1000 // шум
		candles[i] = components.Candle{
			Timestamp: time.Now().Unix() + int64(i)*60,
			Close:     price,
		}
	}
	return candles
}
