package domain_test

import (
	"testing"

	"github.com/reaandrew/eventsourcing-in-go/domain"
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
	var board = domain.NewBoard(columns)

	assert.Equal(t, len(board.CommittedEvents), 1)
	assert.Equal(t, len(board.CommittedEvents[0].(domain.BoardCreated).Columns), 3)
	assert.IsType(t, domain.BoardCreated{}, board.CommittedEvents[0])
}

func TestAddingATicketToABoard(t *testing.T) {
	var columns = createColumns()
	var board = domain.NewBoard(columns)

	var ticket, _ = domain.NewTicket(domain.TicketInfo{
		Title: "something",
	})

	var err = board.AddTicket(ticket, "To Do")

	assert.Nil(t, err)

	assert.Equal(t, len(board.CommittedEvents), 2)
	assert.IsType(t, domain.BoardCreated{}, board.CommittedEvents[0])
	assert.IsType(t, domain.TicketAddedToBoard{}, board.CommittedEvents[1])

	var event = board.CommittedEvents[1].(domain.TicketAddedToBoard)
	assert.Equal(t, columns[0], event.Column.Name)
}

func TestAddingATicketToAColumnWhichDoesNotExistOnABoardReturnsError(t *testing.T) {
	var board = domain.NewBoard(createColumns())
	var ticket, _ = domain.NewTicket(domain.TicketInfo{
		Title: "something",
	})
	var err = board.AddTicket(ticket, "Does not exist")
	assert.Equal(t, domain.ErrUnknownColumn, err)
}
