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

	var boardID = sut.CreateSampleBoard("something")

	var createTicketCommand = commands.CreateTicketCommand{
		TicketID: uuid.NewV4().String(),
		BoardID:  boardID.String(),
		Column:   "todo",
		Title:    "some ticket",
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
	var boardID = sut.CreateSampleBoard("something")
	var createTicketCommand = commands.CreateTicketCommand{
		BoardID:  boardID.String(),
		Assignee: "something",
	}
	var createErr = sut.CommandExecutor.Execute(createTicketCommand)
	assert.Equal(t, commands.ErrInvalidAssigneeID, createErr)
}

func TestCreateTicketCommandReturnsErrorWhenTitleIsEmpty(t *testing.T) {
	var sut = test.NewSystemUnderTest()
	var boardID = sut.CreateSampleBoard("something")

	var createTicketCommand = commands.CreateTicketCommand{
		TicketID: uuid.NewV4().String(),
		BoardID:  boardID.String(),
		Assignee: uuid.NewV4().String(),
	}
	var createErr = sut.CommandExecutor.Execute(createTicketCommand)
	assert.Equal(t, models.ErrNoTicketTitle, createErr)
}
