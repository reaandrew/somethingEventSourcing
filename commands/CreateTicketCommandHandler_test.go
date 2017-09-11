package commands_test

import (
	"testing"

	"github.com/reaandrew/eventsourcing-in-go/commands"
	"github.com/reaandrew/eventsourcing-in-go/domain"
	"github.com/reaandrew/eventsourcing-in-go/test"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

func TestCreateTicketCommandPublishesTicketCreated(t *testing.T) {
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

	var boardCreatedEvent domain.BoardCreated
	boardCreatedEvent = sut.GetEvent(0).(domain.BoardCreated)

	var createTicketCommand = commands.CreateTicketCommand{
		BoardID: boardCreatedEvent.BoardID.String(),
		Title:   "some ticket",
	}

	var createErr = sut.CommandExecutor.Execute(createTicketCommand)
	assert.Nil(t, createErr)

}

func TestCreateTicketCommandReturnsErrorWhenBoardDoesNotExist(t *testing.T) {
	var sut = test.NewSystemUnderTest()
	var createTicketCommand = commands.CreateTicketCommand{
		BoardID: uuid.NewV4().String(),
	}

	var createErr = sut.CommandExecutor.Execute(createTicketCommand)
	assert.Equal(t, domain.ErrBoardNotExist, createErr)
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

	var boardCreatedEvent domain.BoardCreated
	boardCreatedEvent = sut.GetEvent(0).(domain.BoardCreated)

	var createTicketCommand = commands.CreateTicketCommand{
		BoardID:  boardCreatedEvent.BoardID.String(),
		Assignee: "something",
	}
	var createErr = sut.CommandExecutor.Execute(createTicketCommand)
	assert.Equal(t, domain.ErrInvalidAssigneeID, createErr)
}

func TestCreateTicketCommandReturnsErrorWhenTitleIsEmpty(t *testing.T) {

}
