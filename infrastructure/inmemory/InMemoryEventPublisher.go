package inmemory

import (
	"github.com/reaandrew/eventsourcing-in-go/domain/core"
	"github.com/reaandrew/eventsourcing-in-go/domain/models"
	"github.com/reaandrew/eventsourcing-in-go/queries/dtos"
)

type InMemoryEventPublisher struct {
	events     []core.DomainEvent
	queryStore map[string]interface{}
}

func (publisher *InMemoryEventPublisher) Publish(events []core.DomainEvent) (err error) {

	for _, event := range events {
		publisher.updateQueryStore(event)
	}

	publisher.events = append(publisher.events, events...)
	return
}

func NewInMemoryEventPublisher(queryStore map[string]interface{}) (newPublisher *InMemoryEventPublisher) {
	newPublisher = &InMemoryEventPublisher{
		events:     []core.DomainEvent{},
		queryStore: queryStore,
	}

	return
}

func (publisher *InMemoryEventPublisher) updateQueryStore(event core.DomainEvent) {
	switch event.Data.(type) {
	case models.BoardCreated:
		if _, ok := publisher.queryStore["boards"]; !ok {
			publisher.queryStore["boards"] = map[string]dtos.Board{}
		}

		var data = publisher.queryStore["boards"].(map[string]dtos.Board)
		var board = dtos.Board{}.MapDomainCreated(event)
		data[board.ID] = board
	}
}
