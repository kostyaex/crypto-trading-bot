// Пример использования ECS
package main

import (
	"crypto-trading-bot/internal/engine/components"
	"crypto-trading-bot/internal/engine/ecsx"
	"crypto-trading-bot/internal/engine/resources"
	"crypto-trading-bot/internal/engine/systems"
	"fmt"
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

	mode := os.Getenv("MODE")
	if mode == "" {
		mode = ModeBacktest // по умолчанию
	}

	em := ecs.NewEntityManager()

	// // компонент с текущим временем
	// em.Add(ecs.NewEntity("currenttime", []ecs.Component{
	// 	components.NewCurrentTime(),
	// }))

	sm := ecs.NewSystemManager()

	currentTime := resources.NewCurrentTime()

	// Режим: бэктест или live
	if mode == ModeBacktest {
		fmt.Println("🔄 Запуск в режиме бэктеста...")
		candles := generateTestCandles(10)
		sm.Add(systems.NewHistoricalTimeUpdateSystem(currentTime, candles))
	} else {
		sm.Add(systems.NewTimeUpdateSystem(currentTime))
	}
	sm.Add(systems.NewMonitoringSystem(currentTime))

	de := ecsx.NewCustomEngine(em, sm)
	de.Setup()
	defer de.Teardown()

	//de.Run()
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	if mode == ModeBacktest {

		for {
			select {
			case <-sig:
				fmt.Println("\n🛑 Остановка...")
				return
			default:
				de.Tick()
				if de.IsDone() {
					fmt.Println("\n✅Выполнено")
					return
				}
			}
		}

	} else {
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				de.Tick()
			case <-sig:
				fmt.Println("\n🛑 Остановка...")
				//liveSys.Shutdown()
				//goto exit
				return
			}
		}
	}

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
