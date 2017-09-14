package models_test

import (
	"testing"

	"github.com/reaandrew/eventsourcing-in-go/domain/models"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

func TestCreatingANewTicket(t *testing.T) {
	var expectedTitle = "something"
	var ticket, _ = models.NewTicket(models.TicketInfo{
		Title: expectedTitle,
	})

	assert.Equal(t, len(ticket.UncommittedEvents), 1)

	var domainEvent = ticket.UncommittedEvents[0]
	var eventData = domainEvent.Data.(models.TicketCreated)

	assert.IsType(t, models.TicketCreated{}, eventData)
	assert.Equal(t, eventData.Info.Title, expectedTitle)
	assert.Equal(t, domainEvent.Version, 1)
	assert.NotEqual(t, uuid.Nil, domainEvent.ID)
	assert.False(t, domainEvent.Timestamp.IsZero())
}

func TestCreatingANewTicketWithoutATitleReturnsAnError(t *testing.T) {
	var _, err = models.NewTicket(models.TicketInfo{})

	assert.Equal(t, models.ErrNoTicketTitle, err)
}

func TestCreatingANewTicketWithBody(t *testing.T) {
	var expectedContent = "stuff"

	var ticket, _ = models.NewTicket(models.TicketInfo{
		Title:   "Something",
		Content: expectedContent,
	})

	var domainEvent = ticket.UncommittedEvents[0]
	var eventData = domainEvent.Data.(models.TicketCreated)

	assert.Equal(t, eventData.Info.Content, expectedContent)
}

func TestCreatigANewTicketWithAssignee(t *testing.T) {
	var expectedAssignee = uuid.NewV4()
	var ticketInfo = models.TicketInfo{
		Title:    "Something",
		Content:  "stuff",
		Assignee: expectedAssignee,
	}
	var ticket, _ = models.NewTicket(ticketInfo)

	assert.Equal(t, len(ticket.UncommittedEvents), 2)

	var domainEvent = ticket.UncommittedEvents[1]
	var eventData = domainEvent.Data.(models.TicketAssigned)

	assert.IsType(t, models.TicketAssigned{}, eventData)
	assert.Equal(t, eventData.Assignee, expectedAssignee)
	assert.False(t, domainEvent.Timestamp.IsZero())
	assert.Equal(t, domainEvent.Version, 2)
	assert.NotEmpty(t, uuid.Nil, domainEvent.ID)
}

func TestAssigningAnEmptyUserIDReturnsAnError(t *testing.T) {
	var ticketInfo = models.TicketInfo{
		Title:   "Something",
		Content: "stuff",
	}
	var ticket, _ = models.NewTicket(ticketInfo)
	var err = ticket.AssignTo(uuid.UUID{})

	assert.Equal(t, models.ErrCannotAssignToEmptyUserID, err)
}
