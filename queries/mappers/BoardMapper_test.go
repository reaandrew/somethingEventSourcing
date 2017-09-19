package mappers_test

import (
	"testing"
	"time"

	"github.com/reaandrew/eventsourcing-in-go/domain/core"
	"github.com/reaandrew/eventsourcing-in-go/domain/models"
	"github.com/reaandrew/eventsourcing-in-go/queries/dtos"
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
	assert.Equal(t, len(board.Columns), len(boardCreatedEvent.Columns))
}
