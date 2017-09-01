package domain

import (
	"errors"

	uuid "github.com/satori/go.uuid"
)

var (
	ErrNoTicketTitle = errors.New("ErrNoTicketTitle")
)

type TicketInfo struct {
	Title    string
	Content  string
	Assignee uuid.UUID
}

type Ticket struct {
	CommittedEvents []interface{}
	ID              uuid.UUID
	title           string
	content         string
	assignee        uuid.UUID
}

func (ticket *Ticket) handleTicketCreated(event TicketCreated) {
	ticket.ID = event.TicketID
	ticket.title = event.Data.Title
	ticket.content = event.Data.Content
}

func (ticket *Ticket) handleTicketAssigned(event TicketAssigned) {
	ticket.assignee = event.Assignee
}

func (ticket *Ticket) AssignTo(userID uuid.UUID) {
	ticket.apply(TicketAssigned{
		TicketID: ticket.ID,
		Assignee: userID,
	})
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

		if info.Assignee != uuid.Nil {
			newTicket.AssignTo(info.Assignee)
		}
	}
	return
}

func (ticket *Ticket) apply(event interface{}) {
	switch e := event.(type) {
	case TicketAssigned:
		ticket.handleTicketAssigned(e)
	case TicketCreated:
		ticket.handleTicketCreated(e)
	default:
		panic("Unknown Event")

	}
	ticket.CommittedEvents = append(ticket.CommittedEvents, event)
}

type TicketCreated struct {
	TicketID uuid.UUID
	Data     TicketInfo
}

type TicketAssigned struct {
	TicketID uuid.UUID
	Assignee uuid.UUID
}
