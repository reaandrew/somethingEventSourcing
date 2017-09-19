package dtos

import (
	"time"

	"github.com/reaandrew/eventsourcing-in-go/domain/core"
	"github.com/reaandrew/eventsourcing-in-go/domain/models"
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
	mappedBoard.ID = boardCreatedEvent.BoardID.String()
	mappedBoard.Columns = []BoardColumn{}
	for _, col := range boardCreatedEvent.Columns {
		mappedBoard.Columns = append(mappedBoard.Columns, BoardColumn{
			ID:   col.ID.String(),
			Name: col.Name,
		})
	}
	return
}
