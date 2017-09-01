package main_test

import (
	"testing"

	es "github.com/reaandrew/eventsourcing"
	"github.com/stretchr/testify/assert"
)

func TestAddingATicketToABoard(t *testing.T) {
	var board = es.NewBoard()

	board.AddTicket(es.NewTicket())

	assert.Equal(t, len(board.CommittedEvents), 2)
	assert.IsType(t, es.BoardCreated{}, board.CommittedEvents[0])
	assert.IsType(t, es.TicketAddedToBoard{}, board.CommittedEvents[1])
}
