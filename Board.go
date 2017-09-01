package main

import (
	"github.com/satori/go.uuid"
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

func (board *Board) AddTicket(ticket *Ticket) (err error) {
	board.apply(TicketAddedToBoard{
		TicketID: ticket.ID,
	})
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
}

type BoardCreated struct {
	BoardID uuid.UUID
	Columns []BoardColumn
}
