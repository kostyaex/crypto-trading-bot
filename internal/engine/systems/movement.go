// Обновление текущего времени
package systems

import (
	"crypto-trading-bot/internal/engine/components"
	"fmt"

	"github.com/andygeiss/ecs"
)

type movementSystem struct {
}

func NewMovementSystem() ecs.System {
	return &movementSystem{}
}

func (s *movementSystem) Process(em ecs.EntityManager) (state int) {

	// сдвигаем позиции в источнике данных

	for _, dataComp := range em.FilterByMask(components.MaskDatasource) {

		datasource := dataComp.Get(components.MaskDatasource).(*components.DataSource)

		for _, posComp := range em.FilterByMask(components.MaskPosition) {
			position := posComp.Get(components.MaskPosition).(*components.Position)

			if next, ok := datasource.NextPosition(position.Timestamp); ok {
				position.SetPosition(next)
				fmt.Printf("Позиция перемещена: %s\n", next.Format("2006-01-02 15:04:05"))
			}

		}
	}

	return ecs.StateEngineContinue
}

func (s *movementSystem) Setup() {}

func (s *movementSystem) Teardown() {}

// Проверка соответствия интерфейсу
var _ ecs.System = (*movementSystem)(nil)
