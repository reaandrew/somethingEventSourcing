package dtos_test

import (
	"testing"
	"time"

	"github.com/reaandrew/forora/domain/core"
	"github.com/reaandrew/forora/domain/models"
	"github.com/reaandrew/forora/queries/dtos"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

func TestMappingBoardCreatedToBoard(t *testing.T) {
	var event = core.DomainEvent{
		Data: models.BoardCreated{
			BoardID: uuid.NewV4(),
			Name:    "some name",
			Columns: []models.BoardColumnInfo{
				models.BoardColumnInfo{
					ID:   uuid.NewV4(),
					Name: "todo",
				},
			},
		},
		ID:        uuid.NewV4(),
		Timestamp: time.Now(),
		Version:   1,
	}

	var board = dtos.Board{}.MapDomainCreated(event)

	var boardCreatedEvent = event.Data.(models.BoardCreated)

	assert.Equal(t, board.ID, boardCreatedEvent.BoardID.String())
	assert.Equal(t, board.Name, boardCreatedEvent.Name)
	assert.Equal(t, len(board.Columns), len(boardCreatedEvent.Columns))
	assert.Equal(t, board.Columns[0].ID, boardCreatedEvent.Columns[0].ID.String())
	assert.Equal(t, board.Created, event.Timestamp)
	assert.Equal(t, board.Updated, event.Timestamp)
	assert.Equal(t, board.Version, event.Version)
}
