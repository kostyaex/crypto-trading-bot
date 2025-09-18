package ecsx

import "github.com/andygeiss/ecs"

// EntityEventType — тип события
type EntityEventType int

const (
	EntityAdded EntityEventType = iota
	EntityRemoved
)

// EntityEvent — структура события
type EntityEvent struct {
	Type       EntityEventType
	Entity     *ecs.Entity
	Components []ecs.Component // только для Added
}

// EntityLifecycleListener — интерфейс для подписчиков
type EntityLifecycleListener interface {
	OnEntityAdded(Entity *ecs.Entity, components []ecs.Component)
	OnEntityRemoved(Entity *ecs.Entity)
}
