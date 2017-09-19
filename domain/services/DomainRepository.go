package services

import (
	"github.com/reaandrew/eventsourcing-in-go/domain/core"
	"github.com/reaandrew/eventsourcing-in-go/domain/models"
	uuid "github.com/satori/go.uuid"
)

type DomainRepository struct {
	eventStore        core.EventStore
	eventPublisher    core.EventPublisher
	trackedAggregates map[uuid.UUID]core.Aggregate
}

func (repository DomainRepository) Save(aggregate core.Aggregate) {
	if _, ok := repository.trackedAggregates[aggregate.GetID()]; !ok {
		repository.trackedAggregates[aggregate.GetID()] = aggregate
	}

	repository.eventStore.Save(aggregate)
	repository.eventPublisher.Publish(aggregate.GetUncommittedEvents())
}

func (repository DomainRepository) Commit() {
	for _, aggregate := range repository.trackedAggregates {
		aggregate.Commit()
	}
}

func (repository DomainRepository) GetBoard(id uuid.UUID) (newBoard *models.Board, returnErr error) {
	newBoard = &models.Board{}
	var events, err = repository.eventStore.GetEvents(id)
	if err != nil {
		returnErr = err
		return
	}
	if len(events) == 0 {
		returnErr = models.ErrBoardNotExist
		return
	}
	newBoard.Load(events)
	return
}

func (repository DomainRepository) GetTicket(id uuid.UUID) (newTicket *models.Ticket, returnErr error) {
	newTicket = &models.Ticket{}
	var events, err = repository.eventStore.GetEvents(id)

	if err != nil {
		returnErr = err
		return
	}
	if len(events) == 0 {
		returnErr = models.ErrTicketNotExist
		return
	}
	newTicket.Load(events)
	return
}

func NewDomainRepository(eventStore core.EventStore, eventPublisher core.EventPublisher) (domainRepository DomainRepository) {
	return DomainRepository{
		eventStore:        eventStore,
		eventPublisher:    eventPublisher,
		trackedAggregates: map[uuid.UUID]core.Aggregate{},
	}
}
