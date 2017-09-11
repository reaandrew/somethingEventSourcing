package commands_test

import (
	"testing"

	"github.com/reaandrew/eventsourcing-in-go/commands"
	"github.com/reaandrew/eventsourcing-in-go/domain"
	"github.com/reaandrew/eventsourcing-in-go/domain/services"
	"github.com/reaandrew/eventsourcing-in-go/infrastructure/inmemory"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

func TestCreateTicketCommandPublishesTicketCreated(t *testing.T) {
	var eventStore = inmemory.NewInMemoryEventStore()
	var eventPublisher = inmemory.NewInMemoryEventPublisher()
	var domainRepository = services.NewDomainRepository(eventStore, eventPublisher)
	var commandExecutor = commands.NewCommandExecutor(domainRepository)

	var command = commands.CreateBoardCommand{
		Name: "some board",
		Columns: []string{
			"todo",
			"doing",
			"done",
		},
	}

	var err = commandExecutor.Execute(command)
	assert.Nil(t, err)

	var boardCreatedEvent domain.BoardCreated
	boardCreatedEvent = eventPublisher.GetEvent(0).(domain.BoardCreated)

	var createTicketCommand = commands.CreateTicketCommand{
		BoardID: boardCreatedEvent.BoardID.String(),
		Title:   "some ticket",
	}

	var createErr = commandExecutor.Execute(createTicketCommand)
	assert.Nil(t, createErr)

}

func TestCreateTicketCommandReturnsErrorWhenBoardDoesNotExist(t *testing.T) {
	var eventStore = inmemory.NewInMemoryEventStore()
	var eventPublisher = inmemory.NewInMemoryEventPublisher()
	var domainRepository = services.NewDomainRepository(eventStore, eventPublisher)
	var commandExecutor = commands.NewCommandExecutor(domainRepository)
	var createTicketCommand = commands.CreateTicketCommand{
		BoardID: uuid.NewV4().String(),
	}

	var createErr = commandExecutor.Execute(createTicketCommand)
	assert.Equal(t, domain.ErrBoardNotExist, createErr)
}

func TestCreateTicketCommandReturnsErrorWhenAssigneeNotUUID(t *testing.T) {

}

func TestCreateTicketCommandReturnsErrorWhenTitleIsEmpty(t *testing.T) {

}
