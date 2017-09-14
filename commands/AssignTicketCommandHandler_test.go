package commands_test

import (
	"fmt"
	"testing"

	"github.com/reaandrew/eventsourcing-in-go/commands"
	"github.com/reaandrew/eventsourcing-in-go/domain/models"
	"github.com/reaandrew/eventsourcing-in-go/test"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

func TestAssignTicketCommandPublishesTicketAssigned(t *testing.T) {
	var sut = test.NewSystemUnderTest()

	var boardID = sut.CreateSampleBoard("fubar")

	var ticketID = sut.CreateSampleTicket(boardID, "todo")

	fmt.Println("TicketID", ticketID)

	sut.ClearRecordedEvents()

	var assignee = uuid.NewV4().String()

	var command = commands.AssignTicketCommand{
		TicketID: ticketID.String(),
		Assignee: assignee,
	}

	var err = sut.CommandExecutor.Execute(command)

	assert.Nil(t, err)
	assert.Equal(t, 1, sut.NumberOfEventsPublished())
	assert.IsType(t, models.TicketAssigned{}, sut.GetEvent(0))
}
