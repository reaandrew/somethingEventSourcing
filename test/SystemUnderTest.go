package test

import (
	"fmt"

	"github.com/reaandrew/eventsourcing-in-go/commands"
	"github.com/reaandrew/eventsourcing-in-go/domain/core"
	"github.com/reaandrew/eventsourcing-in-go/domain/services"
	"github.com/reaandrew/eventsourcing-in-go/infrastructure/inmemory"
)

type EventRecorder struct {
	Events                []interface{}
	wrappedEventPublisher core.EventPublisher
}

func (eventRecorder *EventRecorder) Publish(events []interface{}) (err error) {
	eventRecorder.Events = append(eventRecorder.Events, events...)
	eventRecorder.wrappedEventPublisher.Publish(events)
	return
}

type SystemUnderTest struct {
	EventStore       core.EventStore
	EventPublisher   core.EventPublisher
	DomainRepository services.DomainRepository
	CommandExecutor  commands.CommandExecutor
	eventRecorder    *EventRecorder
}

func (sut SystemUnderTest) NumberOfEventsPublished() (value int) {
	fmt.Println(sut.eventRecorder.Events)
	value = len(sut.eventRecorder.Events)
	return
}

func (sut SystemUnderTest) GetEvent(index int) (value interface{}) {
	value = sut.eventRecorder.Events[index]
	return
}

func NewSystemUnderTest() (sut SystemUnderTest) {
	var eventStore = inmemory.NewInMemoryEventStore()
	var eventPublisher = inmemory.NewInMemoryEventPublisher()
	var eventRecorder = &EventRecorder{
		Events:                []interface{}{},
		wrappedEventPublisher: eventPublisher,
	}
	var domainRepository = services.NewDomainRepository(eventStore, eventRecorder)
	var commandExecutor = commands.NewCommandExecutor(domainRepository)

	sut = SystemUnderTest{
		EventStore:       eventStore,
		EventPublisher:   eventRecorder,
		DomainRepository: domainRepository,
		CommandExecutor:  commandExecutor,
		eventRecorder:    eventRecorder,
	}
	return
}
