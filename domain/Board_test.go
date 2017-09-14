package domain_test

import (
	"testing"
	"time"

	"github.com/reaandrew/eventsourcing-in-go/domain"
	"github.com/reaandrew/eventsourcing-in-go/domain/core"
	uuid "github.com/satori/go.uuid"
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
	var board = domain.NewBoard(domain.BoardInfo{
		Columns: columns,
	})

	assert.Equal(t, len(board.CommittedEvents), 1)

	var domainEvent = board.CommittedEvents[0]
	var eventData = domainEvent.Data.(domain.BoardCreated)

	assert.Equal(t, len(eventData.Columns), 3)
	assert.IsType(t, domain.BoardCreated{}, eventData)
	assert.Equal(t, domainEvent.Version, 1)
	assert.NotEqual(t, uuid.Nil, domainEvent.ID)
	assert.False(t, domainEvent.Timestamp.IsZero())
}

func TestAddingATicketToABoard(t *testing.T) {
	var columns = createColumns()
	var board = domain.NewBoard(domain.BoardInfo{
		Columns: columns,
	})

	var ticket, _ = domain.NewTicket(domain.TicketInfo{
		Title: "something",
	})

	var err = board.AddTicket(ticket, "To Do")

	assert.Nil(t, err)

	assert.Equal(t, len(board.CommittedEvents), 2)
	assert.IsType(t, domain.BoardCreated{}, board.CommittedEvents[0].Data)
	assert.IsType(t, domain.TicketAddedToBoard{}, board.CommittedEvents[1].Data)

	var domainEvent = board.CommittedEvents[1]
	var eventData = domainEvent.Data.(domain.TicketAddedToBoard)

	assert.Equal(t, columns[0], eventData.Column.Name)
	assert.NotEmpty(t, uuid.Nil, domainEvent.ID)
	assert.False(t, domainEvent.Timestamp.IsZero())
}

func TestAddingATicketToAColumnWhichDoesNotExistOnABoardReturnsError(t *testing.T) {
	var board = domain.NewBoard(domain.BoardInfo{
		Columns: createColumns(),
	})
	var ticket, _ = domain.NewTicket(domain.TicketInfo{
		Title: "something",
	})
	var err = board.AddTicket(ticket, "Does not exist")
	assert.Equal(t, domain.ErrUnknownColumn, err)
}

func TestLoadingABoardFromEvents(t *testing.T) {
	var events = []core.DomainEvent{
		core.DomainEvent{
			ID:        uuid.NewV4(),
			Timestamp: time.Now(),
			Version:   1,
			Data: domain.BoardCreated{
				BoardID: uuid.NewV4(),
				Columns: []domain.BoardColumn{},
			},
		},
	}

	var board = domain.Board{}
	board.Load(events)

}
