package ecsx

import (
	"fmt"

	"github.com/andygeiss/ecs"
)

type CustomEntitityManager struct {
	ecs.EntityManager
	listeners  []EntityLifecycleListener
	eventQueue []EntityEvent // <-- очередь событий
}

func NewEntityManager() ecs.EntityManager {
	return &CustomEntitityManager{
		EntityManager: ecs.NewEntityManager(),
		listeners:     make([]EntityLifecycleListener, 0),
		eventQueue:    make([]EntityEvent, 0),
	}
}

func (e *CustomEntitityManager) Add(entities ...*ecs.Entity) {

	// Накапливаем событие
	for _, entity := range entities {
		fmt.Printf("Add Entity %v\n", entity.Id)
		e.eventQueue = append(e.eventQueue, EntityEvent{
			Type:       EntityAdded,
			Entity:     entity.Id,
			Components: entity.Components,
		})
	}

	e.EntityManager.Add(entities...)
}
