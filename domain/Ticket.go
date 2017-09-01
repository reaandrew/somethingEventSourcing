package domain

import uuid "github.com/satori/go.uuid"

type Ticket struct {
	CommittedEvents []interface{}
	ID              uuid.UUID
	title           string
}

func (ticket *Ticket) handleTicketCreated(event TicketCreated) {
	ticket.ID = event.TicketID
	ticket.title = event.Title
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

func NewTicket(title string) (newTicket *Ticket) {
	newTicket = &Ticket{}
	newTicket.apply(TicketCreated{
		TicketID: uuid.NewV4(),
		Title:    title,
	})
	return
}

type TicketCreated struct {
	TicketID uuid.UUID
	Title    string
}
