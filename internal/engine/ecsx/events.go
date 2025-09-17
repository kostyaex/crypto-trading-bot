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
	Entity     string
	Components []ecs.Component // только для Added
}

// EntityLifecycleListener — интерфейс для подписчиков
type EntityLifecycleListener interface {
	OnEntityAdded(e string, components []ecs.Component)
	OnEntityRemoved(e string)
}
