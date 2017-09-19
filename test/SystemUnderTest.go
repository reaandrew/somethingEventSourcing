package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
	"github.com/icrowley/fake"
	"github.com/reaandrew/eventsourcing-in-go/api/http/rest"
	"github.com/reaandrew/eventsourcing-in-go/commands"
	"github.com/reaandrew/eventsourcing-in-go/domain/core"
	"github.com/reaandrew/eventsourcing-in-go/domain/services"
	"github.com/reaandrew/eventsourcing-in-go/infrastructure/inmemory"
	"github.com/reaandrew/eventsourcing-in-go/queries"
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

func (eventRecorder *EventRecorder) Clear() {
	eventRecorder.Events = []core.DomainEvent{}
}

type SystemUnderTest struct {
	EventStore       core.EventStore
	EventPublisher   core.EventPublisher
	DomainRepository services.DomainRepository
	CommandExecutor  commands.CommandExecutor
	QueryExecutor    queries.QueryExecutor
	eventRecorder    *EventRecorder
}

func (sut SystemUnderTest) NumberOfEventsPublished() (value int) {
	value = len(sut.eventRecorder.Events)
	return
}

func (sut SystemUnderTest) GetEvent(index int) (value interface{}) {
	fmt.Println("Events", sut.eventRecorder.Events)
	if len(sut.eventRecorder.Events) == 0 {
		return nil
	}
	var domainEvent = sut.eventRecorder.Events[index]
	value = domainEvent.Data
	return
}

func (sut SystemUnderTest) CreateSampleBoard(name string) (boardID uuid.UUID) {
	defer sut.DomainRepository.Commit()
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
	defer sut.DomainRepository.Commit()
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

func (sut SystemUnderTest) ClearRecordedEvents() {
	sut.eventRecorder.Clear()
}

func (sut SystemUnderTest) Post(obj interface{}, url string) *httptest.ResponseRecorder {
	var jsonBytes, _ = json.Marshal(obj)
	var reader = bytes.NewReader(jsonBytes)
	var request, _ = http.NewRequest("POST", url, reader)
	var resp = httptest.NewRecorder()
	var router *gin.Engine
	router = rest.SetupRouter(sut.CommandExecutor, sut.QueryExecutor)

	router.ServeHTTP(resp, request)
	return resp
}

func NewSystemUnderTest() (sut SystemUnderTest) {
	var queryStore = map[string]interface{}{}
	var queryExecutor = inmemory.NewInMemoryQueryExecutor(queryStore)

	var eventStore = inmemory.NewInMemoryEventStore()
	var eventPublisher = inmemory.NewInMemoryEventPublisher(queryStore)
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
		QueryExecutor:    queryExecutor,
		eventRecorder:    eventRecorder,
	}
	return
}
