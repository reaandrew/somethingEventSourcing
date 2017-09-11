package domain

import (
	"errors"
	"time"

	"github.com/satori/go.uuid"
)

var (
	ErrUnknownColumn = errors.New("ErrUnknownColumn")
	ErrBoardNotExist = errors.New("ErrBoardNoExist")
)

type BoardInfo struct {
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
	CommittedEvents []interface{}
	ID              uuid.UUID
	columns         []BoardColumn
	version         int
}

func (board *Board) AddTicket(ticket *Ticket, columnName string) (err error) {
	var column, findErr = board.findColumn(columnName)
	if findErr != nil {
		err = findErr
	}
	var event = TicketAddedToBoard{
		Version:   board.version + 1,
		TicketID:  ticket.ID,
		Column:    column,
		EventID:   uuid.NewV4(),
		Timestamp: time.Now(),
	}
	board.apply(event)
	return
}

func (board *Board) GetCommittedEvents() (events []interface{}) {
	return board.CommittedEvents
}

func (board *Board) GetID() (returnID uuid.UUID) {
	return board.ID
}

func (board *Board) GetVersion() (version int) {
	version = board.version
	return
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

func (board *Board) Load(events []interface{}) {
	for _, event := range events {
		board.applyEvent(event)
	}
}

func (board *Board) applyEvent(event interface{}) {
	switch e := event.(type) {
	case BoardCreated:
		board.handleBoardCreated(e)
		board.version = e.Version
	case TicketAddedToBoard:
		board.handleTicketAddedToBoard(e)
		board.version = e.Version
	default:
		panic("Unknown Event")
	}
}

func (board *Board) apply(event interface{}) {
	board.applyEvent(event)
	board.CommittedEvents = append(board.CommittedEvents, event)
}

func NewBoard(info BoardInfo) (newBoard *Board) {
	newBoard = &Board{}
	var boardColumns = []BoardColumn{}
	for _, column := range info.Columns {
		boardColumns = append(boardColumns, NewBoardColumn(column))
	}
	newBoard.apply(BoardCreated{
		Version:   1,
		BoardID:   uuid.NewV4(),
		Columns:   boardColumns,
		EventID:   uuid.NewV4(),
		Timestamp: time.Now(),
	})
	return
}

type TicketAddedToBoard struct {
	EventID   uuid.UUID
	Timestamp time.Time
	Version   int
	TicketID  uuid.UUID
	Column    BoardColumn
}

type BoardCreated struct {
	EventID   uuid.UUID
	Timestamp time.Time
	Version   int
	BoardID   uuid.UUID
	Columns   []BoardColumn
}
