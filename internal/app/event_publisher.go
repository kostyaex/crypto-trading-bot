package app

import (
	"sync"
)

// EventPublisher представляет конкретную реализацию издателя
type EventPublisher struct {
	subscribers map[Subscriber]struct{}
	mu          sync.Mutex
}

// NewEventPublisher создает новый издатель событий
func NewEventPublisher() *EventPublisher {
	return &EventPublisher{
		subscribers: make(map[Subscriber]struct{}),
	}
}

// Subscribe добавляет нового подписчика
func (ep *EventPublisher) Subscribe(subscriber Subscriber) {
	ep.mu.Lock()
	defer ep.mu.Unlock()
	ep.subscribers[subscriber] = struct{}{}
}

// Unsubscribe удаляет подписчика
func (ep *EventPublisher) Unsubscribe(subscriber Subscriber) {
	ep.mu.Lock()
	defer ep.mu.Unlock()
	delete(ep.subscribers, subscriber)
}

// Publish отправляет событие всем подписчикам
func (ep *EventPublisher) Publish(event Event) {
	ep.mu.Lock()
	defer ep.mu.Unlock()
	for subscriber := range ep.subscribers {
		go subscriber.Handle(event)
	}
}
