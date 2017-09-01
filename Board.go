package main

import (
	"errors"

	"github.com/satori/go.uuid"
)

var (
	ErrUnknownColumn = errors.New("Unknown Column")
)

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
	CommittedEvents []interface{}
	ID              uuid.UUID
	Columns         []BoardColumn
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
	board.apply(event)
	return
}

func (board *Board) findColumn(columnName string) (matchingBoard BoardColumn, err error) {
	for _, column := range board.Columns {
		if columnName == column.Name {
			return column, nil
		}
	}
	err = ErrUnknownColumn
	return
}

func (board *Board) handleBoardCreated(event BoardCreated) {
	board.ID = event.BoardID
	board.Columns = event.Columns
}

func (board *Board) handleTicketAddedToBoard(event TicketAddedToBoard) {
	board.Columns[0].Tickets = append(board.Columns[0].Tickets, event.TicketID)
}

func (board *Board) apply(event interface{}) {
	switch e := event.(type) {
	case BoardCreated:
		board.handleBoardCreated(e)
	case TicketAddedToBoard:
		board.handleTicketAddedToBoard(e)
	default:
		panic("Unknown Event")

	}

	board.CommittedEvents = append(board.CommittedEvents, event)
}

func NewBoard(columns []string) (newBoard *Board) {
	newBoard = &Board{}
	var boardColumns = []BoardColumn{}
	for _, column := range columns {
		boardColumns = append(boardColumns, NewBoardColumn(column))
	}
	newBoard.apply(BoardCreated{
		BoardID: uuid.NewV4(),
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
	Columns []BoardColumn
}
