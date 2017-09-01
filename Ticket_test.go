package main_test

import (
	"testing"

	es "github.com/reaandrew/eventsourcing"
	"github.com/stretchr/testify/assert"
)

func TestCreatingANewTicket(t *testing.T) {
	var ticket = es.NewTicket()

	assert.Equal(t, len(ticket.CommittedEvents), 1)
	assert.IsType(t, es.TicketCreated{}, ticket.CommittedEvents[0])
}
