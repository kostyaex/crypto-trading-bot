// Пример использования ECS
package main

import (
	"context"
	"crypto-trading-bot/internal/engine/components"
	"crypto-trading-bot/internal/engine/ecsx"
	"crypto-trading-bot/internal/engine/systems"
	"fmt"
	"log"
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

	//em := ecs.NewEntityManager()
	em := ecsx.NewEntityManager()

	em.Add(ecs.NewEntity("datasource1", []ecs.Component{
		components.NewDataSource("BTCUSDT", "1m", components.GenerateTestCandles(10)),
	}))

	em.Add(ecs.NewEntity("datasource2", []ecs.Component{
		components.NewDataSource("ETHUSDT", "1m", components.GenerateTestCandles(10)),
	}))

	// em.Add(ecs.NewEntity("trader", []ecs.Component{
	// 	components.NewPosition(time.Time{}),
	// }))

	sm := ecs.NewSystemManager()

	fmt.Println("🔄 Запуск...")
	sm.Add(systems.NewStopSystem(ctx))
	sm.Add(systems.NewMovementSystem())

	fetchdataSystem := systems.NewFetchDataSystem(ctx, em)
	sm.Add(fetchdataSystem)

	em.AddListener(fetchdataSystem)

	_engine := ecsx.NewCustomEngine(em, sm)
	_engine.Setup()

	go func() {
		time.Sleep(time.Second)
		_entity := ecs.NewEntity("datasource3", []ecs.Component{
			components.NewDataSource("ETHUSDT", "1m", components.GenerateTestCandles(10)),
		})
		em.Add(_entity)
		time.Sleep(time.Second)
		em.Remove(_entity)
	}()

	defer _engine.Teardown()
	_engine.Run()

}
