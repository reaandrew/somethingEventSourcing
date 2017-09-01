package domain_test

import (
	"testing"

	"github.com/reaandrew/eventsourcing-in-go/domain"
	"github.com/stretchr/testify/assert"
)

func TestCreatingANewTicket(t *testing.T) {
	var expectedTitle = "something"
	var ticket = domain.NewTicket(expectedTitle)

	assert.Equal(t, len(ticket.CommittedEvents), 1)
	assert.IsType(t, domain.TicketCreated{}, ticket.CommittedEvents[0])

	var event = ticket.CommittedEvents[0].(domain.TicketCreated)
	assert.Equal(t, event.Title, expectedTitle)
}
