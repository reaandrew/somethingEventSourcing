package services

import (
	"github.com/reaandrew/eventsourcing-in-go/domain"
	"github.com/reaandrew/eventsourcing-in-go/domain/core"
	uuid "github.com/satori/go.uuid"
)

type DomainRepository struct {
	eventStore     core.EventStore
	eventPublisher core.EventPublisher
}

func (repository DomainRepository) Save(aggregate core.Aggregate) {
	repository.eventStore.Save(aggregate)
	repository.eventPublisher.Publish(aggregate.GetCommittedEvents())
}

func (repository DomainRepository) GetBoard(id uuid.UUID) (newBoard *domain.Board, returnErr error) {
	newBoard = &domain.Board{}
	var events, err = repository.eventStore.GetEvents(id)
	if err != nil {
		returnErr = err
		return
	}
	if len(events) == 0 {
		returnErr = domain.ErrBoardNotExist
		return
	}
	newBoard.Load(events)
	return
}

func NewDomainRepository(eventStore core.EventStore, eventPublisher core.EventPublisher) (domainRepository DomainRepository) {
	return DomainRepository{
		eventStore:     eventStore,
		eventPublisher: eventPublisher,
	}
}
