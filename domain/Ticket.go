package domain

import uuid "github.com/satori/go.uuid"

type Ticket struct {
	CommittedEvents []interface{}
	ID              uuid.UUID
}

func (ticket *Ticket) apply(event interface{}) {
	switch e := event.(type) {
	case TicketCreated:
		ticket.ID = e.TicketID
	default:
		panic("Unknown Event")

	}
	ticket.CommittedEvents = append(ticket.CommittedEvents, event)
}

func NewTicket() (newTicket *Ticket) {
	newTicket = &Ticket{}
	newTicket.apply(TicketCreated{
		TicketID: uuid.NewV4(),
	})
	return
}

type TicketCreated struct {
	TicketID uuid.UUID
}
