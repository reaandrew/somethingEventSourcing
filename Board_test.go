package main_test

import (
	"testing"

	es "github.com/reaandrew/eventsourcing-in-go"
	"github.com/stretchr/testify/assert"
)

func createColumns() (columns []string) {
	columns = []string{
		"To Do",
		"In Progress",
		"Completed",
	}
	return
}

func TestCreatingABoard(t *testing.T) {
	var columns = createColumns()
	var board = es.NewBoard(columns)

	assert.Equal(t, len(board.CommittedEvents), 1)
	assert.Equal(t, len(board.CommittedEvents[0].(es.BoardCreated).Columns), 3)
	assert.IsType(t, es.BoardCreated{}, board.CommittedEvents[0])
}

func TestAddingATicketToABoard(t *testing.T) {
	var columns = createColumns()
	var board = es.NewBoard(columns)

	var err = board.AddTicket(es.NewTicket(), "To Do")

	assert.Nil(t, err)

	assert.Equal(t, len(board.CommittedEvents), 2)
	assert.IsType(t, es.BoardCreated{}, board.CommittedEvents[0])
	assert.IsType(t, es.TicketAddedToBoard{}, board.CommittedEvents[1])

	var event = board.CommittedEvents[1].(es.TicketAddedToBoard)
	assert.Equal(t, columns[0], event.Column.Name)
}

func TestAddingATicketToAColumnWhichDoesNotExistOnABoardReturnsError(t *testing.T) {
	var board = es.NewBoard(createColumns())
	var err = board.AddTicket(es.NewTicket(), "Does not exist")
	assert.Equal(t, es.ErrUnknownColumn, err)
}
