package notifications

import (
	"fmt"
	"log"
	"modular-task/internal/eventbus"
	"modular-task/internal/products"
	"modular-task/internal/users"
	"sync"
)

type NotificationsService struct {
	name           string
	eventBus       *eventbus.EventBus
	stopChan       chan struct{}
	wg             sync.WaitGroup
	receivedEvents []eventbus.Event
	mu             sync.Mutex
}

func NewNotificationsService(name string, eventBus *eventbus.EventBus) *NotificationsService {
	return &NotificationsService{
		name:           name,
		eventBus:       eventBus,
		stopChan:       make(chan struct{}),
		receivedEvents: make([]eventbus.Event, 0),
	}
}

func (ns *NotificationsService) Start() {
	ch := ns.eventBus.Subscribe()

	ns.wg.Add(1)
	go func() {
		defer ns.wg.Done()
		for {
			select {
			case e := <-ch:
				ns.handleEvent(e)
			case <-ns.stopChan:
				fmt.Println("Notification service stopped")
				return
			}
		}
	}()
}

func (ns *NotificationsService) Stop() {
	close(ns.stopChan)
	ns.wg.Wait()
}
func (ns *NotificationsService) handleEvent(e eventbus.Event) {
	ns.mu.Lock()
	defer ns.mu.Unlock()
	ns.receivedEvents = append(ns.receivedEvents, e)

	switch e.Type {
	case eventbus.EventUserCreated:
		user, ok := e.Data.(users.User)
		if !ok {
			log.Printf("[%s] Received user.created event with invalid data", ns.name)
			return
		}
		log.Printf("[%s] Received user.created event: %+v", ns.name, user)
	case eventbus.ProductCreated:
		product, ok := e.Data.(products.Product)
		if !ok {
			log.Printf("[%s] Received product.created event with invalid data", ns.name)
			return
		}
		log.Printf("[%s] Received product.created event: %+v", ns.name, product)
	default:
		log.Printf("[%s] Received unknown event type: %s", ns.name, e.Type)
	}
}
