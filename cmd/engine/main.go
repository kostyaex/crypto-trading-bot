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

	em := ecs.NewEntityManager()

	em.Add(ecs.NewEntity("datasource", []ecs.Component{
		components.NewDataSource(components.GenerateTestCandles(10)),
	}))

	em.Add(ecs.NewEntity("trader", []ecs.Component{
		components.NewPosition(time.Time{}),
	}))

	sm := ecs.NewSystemManager()

	fmt.Println("🔄 Запуск...")
	sm.Add(systems.NewStopSystem(ctx))
	sm.Add(systems.NewMovementSystem())

	de := ecsx.NewCustomEngine(em, sm)
	de.Setup()
	defer de.Teardown()

	de.Run()

}
