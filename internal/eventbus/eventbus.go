package eventbus

import "sync"

type EventType string

const (
	EventUserCreated EventType = "user.created"
	ProductCreated   EventType = "product.created"
)

type Event struct {
	Type EventType
	Data interface{}
}

type EventBus struct {
	subscribers []chan Event
	mu          sync.Mutex
}

func NewEventBus() *EventBus {
	return &EventBus{
		subscribers: make([]chan Event, 0),
	}
}

func (eb *EventBus) Subscribe() <-chan Event {
	ch := make(chan Event, 10)
	eb.mu.Lock()
	defer eb.mu.Unlock()
	eb.subscribers = append(eb.subscribers, ch)
	return ch
}

func (eb *EventBus) Publish(e Event) {
	eb.mu.Lock()
	defer eb.mu.Unlock()
	for _, ch := range eb.subscribers {
		ch <- e
	}
}
