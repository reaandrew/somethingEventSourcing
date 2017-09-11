package inmemory

import (
	"github.com/reaandrew/eventsourcing-in-go/domain/core"
	uuid "github.com/satori/go.uuid"
)

type InMemoryEventStore struct {
	events map[uuid.UUID][]interface{}
}

func (store *InMemoryEventStore) Save(aggregate core.Aggregate) (err error) {
	var id = aggregate.GetID()
	var events = aggregate.GetCommittedEvents()
	if _, ok := store.events[id]; !ok {
		store.events[id] = []interface{}{}
	}

	store.events[id] = append(store.events[id], events...)
	return
}

func (store *InMemoryEventStore) GetEvents(id uuid.UUID) (events []interface{}, err error) {
	events = store.events[id]
	return
}

func NewInMemoryEventStore() (newEventStore *InMemoryEventStore) {
	return &InMemoryEventStore{
		events: map[uuid.UUID][]interface{}{},
	}
}