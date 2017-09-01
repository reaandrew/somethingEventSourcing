package domain

import (
	"errors"
	"time"

	uuid "github.com/satori/go.uuid"
)

var (
	ErrNoTicketTitle             = errors.New("ErrNoTicketTitle")
	ErrCannotAssignToEmptyUserID = errors.New("ErrCannotAssignToEmptyUserID")
)

type TicketInfo struct {
	Title    string
	Content  string
	Assignee uuid.UUID
}

type Ticket struct {
	CommittedEvents []interface{}
	ID              uuid.UUID
	version         int
	title           string
	content         string
	assignee        uuid.UUID
}

func (ticket *Ticket) handleTicketCreated(event TicketCreated) {
	ticket.ID = event.TicketID
	ticket.title = event.Data.Title
	ticket.content = event.Data.Content
	ticket.version = event.Version
}

func (ticket *Ticket) handleTicketAssigned(event TicketAssigned) {
	ticket.assignee = event.Assignee
	ticket.version = event.Version
}

func (ticket *Ticket) AssignTo(userID uuid.UUID) (err error) {
	if userID == uuid.Nil {
		err = ErrCannotAssignToEmptyUserID
	} else {
		ticket.apply(TicketAssigned{
			TicketID:  ticket.ID,
			Assignee:  userID,
			Version:   ticket.version + 1,
			Timestamp: time.Now(),
			EventID:   uuid.NewV4(),
		})
	}
	return
}

func NewTicket(info TicketInfo) (newTicket *Ticket, err error) {
	newTicket = &Ticket{}
	if info.Title == "" {
		err = ErrNoTicketTitle
	} else {
		newTicket.apply(TicketCreated{
			TicketID:  uuid.NewV4(),
			Data:      info,
			Version:   1,
			EventID:   uuid.NewV4(),
			Timestamp: time.Now(),
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
	EventID   uuid.UUID
	Timestamp time.Time
	Version   int
	TicketID  uuid.UUID
	Data      TicketInfo
}

type TicketAssigned struct {
	EventID   uuid.UUID
	Timestamp time.Time
	Version   int
	TicketID  uuid.UUID
	Assignee  uuid.UUID
}
