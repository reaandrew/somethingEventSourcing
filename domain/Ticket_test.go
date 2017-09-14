package domain_test

import (
	"testing"

	"github.com/reaandrew/eventsourcing-in-go/domain"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

func TestCreatingANewTicket(t *testing.T) {
	var expectedTitle = "something"
	var ticket, _ = domain.NewTicket(domain.TicketInfo{
		Title: expectedTitle,
	})

	assert.Equal(t, len(ticket.CommittedEvents), 1)

	var domainEvent = ticket.CommittedEvents[0]
	var eventData = domainEvent.Data.(domain.TicketCreated)

	assert.IsType(t, domain.TicketCreated{}, eventData)
	assert.Equal(t, eventData.Info.Title, expectedTitle)
	assert.Equal(t, domainEvent.Version, 1)
	assert.NotEqual(t, uuid.Nil, domainEvent.ID)
	assert.False(t, domainEvent.Timestamp.IsZero())
}

func TestCreatingANewTicketWithoutATitleReturnsAnError(t *testing.T) {
	var _, err = domain.NewTicket(domain.TicketInfo{})

	assert.Equal(t, domain.ErrNoTicketTitle, err)
}

func TestCreatingANewTicketWithBody(t *testing.T) {
	var expectedContent = "stuff"

	var ticket, _ = domain.NewTicket(domain.TicketInfo{
		Title:   "Something",
		Content: expectedContent,
	})

	var domainEvent = ticket.CommittedEvents[0]
	var eventData = domainEvent.Data.(domain.TicketCreated)

	assert.Equal(t, eventData.Info.Content, expectedContent)
}

func TestCreatigANewTicketWithAssignee(t *testing.T) {
	var expectedAssignee = uuid.NewV4()
	var ticketInfo = domain.TicketInfo{
		Title:    "Something",
		Content:  "stuff",
		Assignee: expectedAssignee,
	}
	var ticket, _ = domain.NewTicket(ticketInfo)

	assert.Equal(t, len(ticket.CommittedEvents), 2)

	var domainEvent = ticket.CommittedEvents[1]
	var eventData = domainEvent.Data.(domain.TicketAssigned)

	assert.IsType(t, domain.TicketAssigned{}, eventData)
	assert.Equal(t, eventData.Assignee, expectedAssignee)
	assert.False(t, domainEvent.Timestamp.IsZero())
	assert.Equal(t, domainEvent.Version, 2)
	assert.NotEmpty(t, uuid.Nil, domainEvent.ID)
}

func TestAssigningAnEmptyUserIDReturnsAnError(t *testing.T) {
	var ticketInfo = domain.TicketInfo{
		Title:   "Something",
		Content: "stuff",
	}
	var ticket, _ = domain.NewTicket(ticketInfo)
	var err = ticket.AssignTo(uuid.UUID{})

	assert.Equal(t, domain.ErrCannotAssignToEmptyUserID, err)
}
