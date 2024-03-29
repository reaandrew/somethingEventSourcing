package dtos

import (
	"time"

	"github.com/reaandrew/forora/domain/core"
	"github.com/reaandrew/forora/domain/models"
)

type Board struct {
	ID      string
	Name    string
	Columns []BoardColumn
	Version int
	Updated time.Time
	Created time.Time
}

func (board Board) MapDomainCreated(event core.DomainEvent) (mappedBoard Board) {
	var boardCreatedEvent = event.Data.(models.BoardCreated)
	mappedBoard.Name = boardCreatedEvent.Name
	mappedBoard.ID = boardCreatedEvent.BoardID.String()
	mappedBoard.Columns = []BoardColumn{}
	for _, col := range boardCreatedEvent.Columns {
		mappedBoard.Columns = append(mappedBoard.Columns, BoardColumn{
			ID:   col.ID.String(),
			Name: col.Name,
		})
	}
	mappedBoard.Created = event.Timestamp
	mappedBoard.Updated = event.Timestamp
	mappedBoard.Version = event.Version
	return
}
