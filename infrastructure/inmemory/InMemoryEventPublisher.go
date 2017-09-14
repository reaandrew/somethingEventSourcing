package inmemory

import "github.com/reaandrew/eventsourcing-in-go/domain/core"

type InMemoryEventPublisher struct {
	events []core.DomainEvent
}

func (publisher *InMemoryEventPublisher) Publish(events []core.DomainEvent) (err error) {
	publisher.events = append(publisher.events, events...)
	return
}

func NewInMemoryEventPublisher() (newPublisher *InMemoryEventPublisher) {
	newPublisher = &InMemoryEventPublisher{
		events: []core.DomainEvent{},
	}

	return
}
