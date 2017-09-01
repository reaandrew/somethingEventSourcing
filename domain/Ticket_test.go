package domain_test

import (
	"testing"

	"github.com/reaandrew/eventsourcing-in-go/domain"
	"github.com/stretchr/testify/assert"
)

func TestCreatingANewTicket(t *testing.T) {
	var ticket = domain.NewTicket()

	assert.Equal(t, len(ticket.CommittedEvents), 1)
	assert.IsType(t, domain.TicketCreated{}, ticket.CommittedEvents[0])
}
