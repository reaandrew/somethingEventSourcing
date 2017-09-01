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
	assert.IsType(t, domain.TicketCreated{}, ticket.CommittedEvents[0])

	var event = ticket.CommittedEvents[0].(domain.TicketCreated)
	assert.Equal(t, event.Data.Title, expectedTitle)
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

	var event = ticket.CommittedEvents[0].(domain.TicketCreated)
	assert.Equal(t, event.Data.Content, expectedContent)
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
	assert.IsType(t, domain.TicketAssigned{}, ticket.CommittedEvents[1])
	var event = ticket.CommittedEvents[1].(domain.TicketAssigned)
	assert.Equal(t, event.Assignee, expectedAssignee)
}
