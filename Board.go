package main

import (
	"github.com/satori/go.uuid"
)

type Board struct {
	CommittedEvents []interface{}
	ID              uuid.UUID
	Tickets         []uuid.UUID
}

func (board *Board) AddTicket(ticket *Ticket) (err error) {
	board.apply(TicketAddedToBoard{
		TicketID: ticket.ID,
	})
	return
}

func (board *Board) handleBoardCreated(event BoardCreated) {
	board.ID = event.BoardID
}

func (board *Board) handleTicketAddedToBoard(event TicketAddedToBoard) {
	board.Tickets = append(board.Tickets, event.TicketID)
}

func (board *Board) apply(event interface{}) {
	switch e := event.(type) {
	case BoardCreated:
	case TicketAddedToBoard:
		board.handleTicketAddedToBoard(e)
	default:
		panic("Unknown Event")

	}

	board.CommittedEvents = append(board.CommittedEvents, event)
}

func NewBoard() (newBoard *Board) {
	newBoard = &Board{}
	newBoard.apply(BoardCreated{
		BoardID: uuid.NewV4(),
	})
	return
}

type TicketAddedToBoard struct {
	TicketID uuid.UUID
}

type BoardCreated struct {
	BoardID uuid.UUID
}
