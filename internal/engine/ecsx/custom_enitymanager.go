package ecsx

import (
	"github.com/andygeiss/ecs"
)

type EntityManager interface {
	ecs.EntityManager
	AddListener(listener EntityLifecycleListener)
	ProcessEvents()
}

type CustomEntitityManager struct {
	ecs.EntityManager
	listeners  []EntityLifecycleListener
	eventQueue []EntityEvent // <-- очередь событий
}

func NewEntityManager() EntityManager {
	return &CustomEntitityManager{
		EntityManager: ecs.NewEntityManager(),
		listeners:     make([]EntityLifecycleListener, 0),
		eventQueue:    make([]EntityEvent, 0),
	}
}

func (e *CustomEntitityManager) Add(entities ...*ecs.Entity) {

	// Накапливаем событие
	for _, entity := range entities {
		// fmt.Printf("Add Entity %v\n", entity.Id)
		e.eventQueue = append(e.eventQueue, EntityEvent{
			Type:       EntityAdded,
			Entity:     entity,
			Components: entity.Components,
		})
	}

	e.EntityManager.Add(entities...)
}

func (e *CustomEntitityManager) Remove(entity *ecs.Entity) {

	// Накапливаем событие
	e.eventQueue = append(e.eventQueue, EntityEvent{
		Type:   EntityRemoved,
		Entity: entity,
	})

	//e.Remove(entity)
	//здесь не удаляем. Перенес ниже, сначала должно обработаться событие удаления, потом уже удаляем
}

// ProcessEvents — обрабатывает все накопленные события пачкой
func (e *CustomEntitityManager) ProcessEvents() {
	for _, event := range e.eventQueue {
		switch event.Type {
		case EntityAdded:
			//	fmt.Printf("Add Entity %v\n", event.Entity.Id)
			for _, listener := range e.listeners {
				listener.OnEntityAdded(event.Entity, event.Components)
			}
		case EntityRemoved:
			// Сначала уведомляем системы
			for _, listener := range e.listeners {
				listener.OnEntityRemoved(event.Entity)

				e.Remove(event.Entity)
			}
		}
	}

	// Очищаем очередь
	e.eventQueue = e.eventQueue[:0] // reset slice, keep capacity
}

// AddListener подписывает систему на события
func (e *CustomEntitityManager) AddListener(listener EntityLifecycleListener) {
	e.listeners = append(e.listeners, listener)
}
