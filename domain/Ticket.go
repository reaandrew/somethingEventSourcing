package domain

import (
	"errors"

	uuid "github.com/satori/go.uuid"
)

var (
	ErrNoTicketTitle = errors.New("ErrNoTicketTitle")
)

type TicketInfo struct {
	Title   string
	Content string
}

type Ticket struct {
	CommittedEvents []interface{}
	ID              uuid.UUID
	title           string
}

func (ticket *Ticket) handleTicketCreated(event TicketCreated) {
	ticket.ID = event.TicketID
	ticket.title = event.Data.Title
}

func (ticket *Ticket) apply(event interface{}) {
	switch e := event.(type) {
	case TicketCreated:
		ticket.handleTicketCreated(e)
	default:
		panic("Unknown Event")

	}
	ticket.CommittedEvents = append(ticket.CommittedEvents, event)
}

func NewTicket(info TicketInfo) (newTicket *Ticket, err error) {
	newTicket = &Ticket{}
	if info.Title == "" {
		err = ErrNoTicketTitle
	} else {
		newTicket.apply(TicketCreated{
			TicketID: uuid.NewV4(),
			Data:     info,
		})
	}
	return
}

type TicketCreated struct {
	TicketID uuid.UUID
	Data     TicketInfo
}
