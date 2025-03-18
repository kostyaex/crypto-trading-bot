package app

// Publisher представляет издателя событий
type Publisher interface {
	Subscribe(subscriber Subscriber)
	Unsubscribe(subscriber Subscriber)
	Publish(event Event)
}

// Subscriber представляет подписчика событий
type Subscriber interface {
	Handle(event Event)
}

// Event представляет событие, которое может быть опубликовано
type Event interface {
	Type() string
	Payload() interface{}
}
