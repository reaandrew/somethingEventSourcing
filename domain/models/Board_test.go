package models_test

import (
	"testing"
	"time"

	"github.com/reaandrew/eventsourcing-in-go/domain/core"
	"github.com/reaandrew/eventsourcing-in-go/domain/models"
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
	var board = models.NewBoard(models.BoardInfo{
		Columns: columns,
	})

	assert.Equal(t, len(board.UncommittedEvents), 1)

	var domainEvent = board.UncommittedEvents[0]
	var eventData = domainEvent.Data.(models.BoardCreated)

	assert.Equal(t, len(eventData.Columns), 3)
	assert.IsType(t, models.BoardCreated{}, eventData)
	assert.Equal(t, domainEvent.Version, 1)
	assert.NotEqual(t, uuid.Nil, domainEvent.ID)
	assert.False(t, domainEvent.Timestamp.IsZero())
}

func TestAddingATicketToABoard(t *testing.T) {
	var columns = createColumns()
	var board = models.NewBoard(models.BoardInfo{
		Columns: columns,
	})

	var ticket, _ = models.NewTicket(models.TicketInfo{
		Title: "something",
	})

	var err = board.AddTicket(ticket, "To Do")

	assert.Nil(t, err)

	assert.Equal(t, len(board.UncommittedEvents), 2)
	assert.IsType(t, models.BoardCreated{}, board.UncommittedEvents[0].Data)
	assert.IsType(t, models.TicketAddedToBoard{}, board.UncommittedEvents[1].Data)

	var domainEvent = board.UncommittedEvents[1]
	var eventData = domainEvent.Data.(models.TicketAddedToBoard)

	assert.Equal(t, columns[0], eventData.Column.Name)
	assert.NotEmpty(t, uuid.Nil, domainEvent.ID)
	assert.False(t, domainEvent.Timestamp.IsZero())
}

func TestAddingATicketToAColumnWhichDoesNotExistOnABoardReturnsError(t *testing.T) {
	var board = models.NewBoard(models.BoardInfo{
		Columns: createColumns(),
	})
	var ticket, _ = models.NewTicket(models.TicketInfo{
		Title: "something",
	})
	var err = board.AddTicket(ticket, "Does not exist")
	assert.Equal(t, models.ErrUnknownColumn, err)
}

func TestLoadingABoardFromEvents(t *testing.T) {
	var events = []core.DomainEvent{
		core.DomainEvent{
			ID:        uuid.NewV4(),
			Timestamp: time.Now(),
			Version:   1,
			Data: models.BoardCreated{
				BoardID: uuid.NewV4(),
				Columns: []models.BoardColumn{},
			},
		},
	}

	var board = models.Board{}
	board.Load(events)

}
