package commands_test

import (
	"testing"

	"github.com/reaandrew/eventsourcing-in-go/commands"
	"github.com/reaandrew/eventsourcing-in-go/domain/models"
	"github.com/reaandrew/eventsourcing-in-go/test"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

func TestCreateTicketCommandPublishesEvents(t *testing.T) {
	var sut = test.NewSystemUnderTest()

	var command = commands.CreateBoardCommand{
		Name: "some board",
		Columns: []string{
			"todo",
			"doing",
			"done",
		},
	}

	var err = sut.CommandExecutor.Execute(command)
	assert.Nil(t, err)

	var boardCreatedEvent models.BoardCreated
	boardCreatedEvent = sut.GetEvent(0).(models.BoardCreated)

	var createTicketCommand = commands.CreateTicketCommand{
		BoardID: boardCreatedEvent.BoardID.String(),
		Column:  "todo",
		Title:   "some ticket",
	}

	var createErr = sut.CommandExecutor.Execute(createTicketCommand)
	assert.Nil(t, createErr)

	assert.IsType(t, models.TicketCreated{}, sut.GetEvent(1))
	assert.IsType(t, models.TicketAddedToBoard{}, sut.GetEvent(2))
}

func TestCreateTicketCommandReturnsErrorWhenBoardDoesNotExist(t *testing.T) {
	var sut = test.NewSystemUnderTest()
	var createTicketCommand = commands.CreateTicketCommand{
		BoardID: uuid.NewV4().String(),
	}

	var createErr = sut.CommandExecutor.Execute(createTicketCommand)
	assert.Equal(t, models.ErrBoardNotExist, createErr)
}

func TestCreateTicketCommandReturnsErrorWhenBoardIDNotUUID(t *testing.T) {
	var sut = test.NewSystemUnderTest()
	var createTicketCommand = commands.CreateTicketCommand{
		BoardID: "something",
	}

	var createErr = sut.CommandExecutor.Execute(createTicketCommand)
	assert.Equal(t, commands.ErrInvalidBoardID, createErr)
}

func TestCreateTicketCommandReturnsErrorWhenAssigneeNotUUID(t *testing.T) {
	var sut = test.NewSystemUnderTest()
	var command = commands.CreateBoardCommand{
		Name: "some board",
		Columns: []string{
			"todo",
			"doing",
			"done",
		},
	}

	var err = sut.CommandExecutor.Execute(command)
	assert.Nil(t, err)

	var boardCreatedEvent models.BoardCreated
	boardCreatedEvent = sut.GetEvent(0).(models.BoardCreated)

	var createTicketCommand = commands.CreateTicketCommand{
		BoardID:  boardCreatedEvent.BoardID.String(),
		Assignee: "something",
	}
	var createErr = sut.CommandExecutor.Execute(createTicketCommand)
	assert.Equal(t, commands.ErrInvalidAssigneeID, createErr)
}

func TestCreateTicketCommandReturnsErrorWhenTitleIsEmpty(t *testing.T) {
	var sut = test.NewSystemUnderTest()
	var command = commands.CreateBoardCommand{
		Name: "some board",
		Columns: []string{
			"todo",
			"doing",
			"done",
		},
	}

	var err = sut.CommandExecutor.Execute(command)
	assert.Nil(t, err)

	var boardCreatedEvent models.BoardCreated
	boardCreatedEvent = sut.GetEvent(0).(models.BoardCreated)

	var createTicketCommand = commands.CreateTicketCommand{
		BoardID:  boardCreatedEvent.BoardID.String(),
		Assignee: uuid.NewV4().String(),
	}
	var createErr = sut.CommandExecutor.Execute(createTicketCommand)
	assert.Equal(t, models.ErrNoTicketTitle, createErr)
}
