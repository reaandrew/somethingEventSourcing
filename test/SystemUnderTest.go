package test

import (
	"fmt"

	"github.com/icrowley/fake"
	"github.com/reaandrew/eventsourcing-in-go/commands"
	"github.com/reaandrew/eventsourcing-in-go/domain/core"
	"github.com/reaandrew/eventsourcing-in-go/domain/services"
	"github.com/reaandrew/eventsourcing-in-go/infrastructure/inmemory"
	uuid "github.com/satori/go.uuid"
)

type EventRecorder struct {
	Events                []core.DomainEvent
	wrappedEventPublisher core.EventPublisher
}

func (eventRecorder *EventRecorder) Publish(events []core.DomainEvent) (err error) {
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
	var domainEvent = sut.eventRecorder.Events[index]
	value = domainEvent.Data
	return
}

func (sut SystemUnderTest) CreateSampleBoard(name string) (boardID uuid.UUID) {
	boardID = uuid.NewV4()
	var command = commands.CreateBoardCommand{
		BoardID: boardID.String(),
		Name:    name,
		Columns: []string{
			"todo",
			"doing",
			"done",
		},
	}

	var err = sut.CommandExecutor.Execute(command)

	if err != nil {
		panic(err)
	}

	return
}

func (sut SystemUnderTest) CreateSampleTicket(boardID uuid.UUID, column string) (ticketID uuid.UUID) {
	ticketID = uuid.NewV4()
	var command = commands.CreateTicketCommand{
		TicketID: ticketID.String(),
		BoardID:  boardID.String(),
		Column:   column,
		Title:    fake.Title(),
		Content:  fake.Paragraph(),
	}

	var err = sut.CommandExecutor.Execute(command)

	if err != nil {
		panic(err)
	}

	return
}

func NewSystemUnderTest() (sut SystemUnderTest) {
	var eventStore = inmemory.NewInMemoryEventStore()
	var eventPublisher = inmemory.NewInMemoryEventPublisher()
	var eventRecorder = &EventRecorder{
		Events:                []core.DomainEvent{},
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
