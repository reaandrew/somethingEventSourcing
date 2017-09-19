package models

import (
	"errors"
	"time"

	"github.com/reaandrew/eventsourcing-in-go/domain/core"
	"github.com/satori/go.uuid"
)

var (
	ErrUnknownColumn  = errors.New("ErrUnknownColumn")
	ErrBoardNotExist  = errors.New("ErrBoardNotExist")
	ErrInvalidBoardID = errors.New("ErrInvalidBoardID")
)

type BoardInfo struct {
	BoardID uuid.UUID
	Name    string
	Columns []string
}

type BoardColumn struct {
	ID      uuid.UUID
	Name    string
	Tickets []uuid.UUID
}

func NewBoardColumn(name string) (newBoardColumn BoardColumn) {
	newBoardColumn = BoardColumn{
		ID:      uuid.NewV4(),
		Name:    name,
		Tickets: []uuid.UUID{},
	}
	return
}

type Board struct {
	UncommittedEvents []core.DomainEvent
	ID                uuid.UUID
	columns           []BoardColumn
	version           int
}

func (board *Board) AddTicket(ticket *Ticket, columnName string) (err error) {
	var column, findErr = board.findColumn(columnName)
	if findErr != nil {
		err = findErr
	}
	var event = TicketAddedToBoard{
		TicketID: ticket.ID,
		Column:   column,
	}
	board.publish(event)
	return
}

func (board *Board) GetUncommittedEvents() (events []core.DomainEvent) {
	return board.UncommittedEvents
}

func (board *Board) GetID() (returnID uuid.UUID) {
	return board.ID
}

func (board *Board) GetVersion() (version int) {
	version = board.version
	return
}

func (board *Board) Commit() {
	board.UncommittedEvents = []core.DomainEvent{}
}

func (board *Board) findColumn(columnName string) (matchingBoard BoardColumn, err error) {
	for _, column := range board.columns {
		if columnName == column.Name {
			return column, nil
		}
	}
	err = ErrUnknownColumn
	return
}

func (board *Board) handleBoardCreated(event BoardCreated) {
	board.ID = event.BoardID
	board.columns = event.Columns
}

func (board *Board) handleTicketAddedToBoard(event TicketAddedToBoard) {
	board.columns[0].Tickets = append(board.columns[0].Tickets, event.TicketID)
}

func (board *Board) Load(events []core.DomainEvent) {
	for _, event := range events {
		board.replay(event)
	}
}

func (board *Board) apply(event core.DomainEvent) {
	switch e := event.Data.(type) {
	case BoardCreated:
		board.handleBoardCreated(e)
	case TicketAddedToBoard:
		board.handleTicketAddedToBoard(e)
	default:
		panic("Unknown Event")
	}
}

func (board *Board) replay(domainEvent core.DomainEvent) {
	board.apply(domainEvent)
	board.version = domainEvent.Version
}

func (board *Board) publish(event interface{}) {
	var domainEvent = core.DomainEvent{
		ID:        uuid.NewV4(),
		Version:   board.version + len(board.UncommittedEvents) + 1,
		Timestamp: time.Now(),
		Data:      event,
	}
	board.apply(domainEvent)
	board.UncommittedEvents = append(board.UncommittedEvents, domainEvent)
}

func NewBoard(info BoardInfo) (newBoard *Board) {
	newBoard = &Board{}
	var boardColumns = []BoardColumn{}
	for _, column := range info.Columns {
		boardColumns = append(boardColumns, NewBoardColumn(column))
	}
	newBoard.publish(BoardCreated{
		BoardID: info.BoardID,
		Name:    info.Name,
		Columns: boardColumns,
	})
	return
}

type TicketAddedToBoard struct {
	TicketID uuid.UUID
	Column   BoardColumn
}

type BoardCreated struct {
	BoardID uuid.UUID
	Name    string
	Columns []BoardColumn
}
