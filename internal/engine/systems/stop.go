// для остановки выполнения
package systems

import (
	"context"
	"fmt"

	"github.com/andygeiss/ecs"
)

type stopSystem struct {
	cancelContext context.Context
}

func NewStopSystem(cancelContext context.Context) ecs.System {
	return &stopSystem{
		cancelContext: cancelContext,
	}
}

func (s *stopSystem) Process(em ecs.EntityManager) (state int) {

	// Здесь проверяем сигнал на завершение
	select {
	case <-s.cancelContext.Done():
		fmt.Println("\n🛑 (stopSystem) Остановка...")
		return ecs.StateEngineStop
	default:
		return ecs.StateEngineContinue
	}

}

func (s *stopSystem) Setup() {}

func (s *stopSystem) Teardown() {}

// Проверка соответствия интерфейсу
var _ ecs.System = (*stopSystem)(nil)
