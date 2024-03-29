package models

import (
	"errors"
	"time"

	"github.com/reaandrew/forora/domain/core"
	uuid "github.com/satori/go.uuid"
)

var (
	ErrNoTicketTitle             = errors.New("ErrNoTicketTitle")
	ErrCannotAssignToEmptyUserID = errors.New("ErrCannotAssignToEmptyUserID")
	ErrInvalidAssigneeID         = errors.New("ErrInvalidAssigneeID")
	ErrInvalidTicketID           = errors.New("ErrInvalidTicketID")
	ErrTicketNotExist            = errors.New("ErrTicketNotExist")
)

type TicketInfo struct {
	TicketID uuid.UUID
	Title    string
	Content  string
	Assignee uuid.UUID
}

type Ticket struct {
	UncommittedEvents []core.DomainEvent
	ID                uuid.UUID
	version           int
	title             string
	content           string
	assignee          uuid.UUID
}

func (ticket *Ticket) handleTicketCreated(event TicketCreated) {
	ticket.ID = event.TicketID
	ticket.title = event.Info.Title
	ticket.content = event.Info.Content
}

func (ticket *Ticket) handleTicketAssigned(event TicketAssigned) {
	ticket.assignee = event.Assignee
}

func (ticket *Ticket) AssignTo(userID uuid.UUID) (err error) {
	if userID == uuid.Nil {
		err = ErrCannotAssignToEmptyUserID
	} else {
		ticket.publish(TicketAssigned{
			TicketID: ticket.ID,
			Assignee: userID,
		})
	}
	return
}

func (ticket *Ticket) GetUncommittedEvents() (events []core.DomainEvent) {
	return ticket.UncommittedEvents
}

func (ticket *Ticket) GetID() (returnID uuid.UUID) {
	return ticket.ID
}

func (ticket *Ticket) GetVersion() (version int) {
	version = ticket.version
	return
}

func (ticket *Ticket) Commit() {
	ticket.UncommittedEvents = []core.DomainEvent{}
}

func NewTicket(info TicketInfo) (newTicket *Ticket, err error) {
	newTicket = &Ticket{}
	if info.Title == "" {
		err = ErrNoTicketTitle
	} else {
		newTicket.publish(TicketCreated{
			TicketID: info.TicketID,
			Info:     info,
		})

		if info.Assignee != uuid.Nil {
			newTicket.AssignTo(info.Assignee)
		}
	}
	return
}

func (ticket *Ticket) Load(events []core.DomainEvent) {
	for _, event := range events {
		ticket.replay(event)
	}
}

func (ticket *Ticket) apply(event core.DomainEvent) {
	switch e := event.Data.(type) {
	case TicketAssigned:
		ticket.handleTicketAssigned(e)
	case TicketCreated:
		ticket.handleTicketCreated(e)
	default:
		panic("Unknown Event")

	}
}

func (ticket *Ticket) replay(domainEvent core.DomainEvent) {
	ticket.apply(domainEvent)
	ticket.version = domainEvent.Version
}

func (ticket *Ticket) publish(event interface{}) {
	var domainEvent = core.DomainEvent{
		ID:        uuid.NewV4(),
		Version:   ticket.version + len(ticket.UncommittedEvents) + 1,
		Timestamp: time.Now(),
		Data:      event,
	}
	ticket.apply(domainEvent)
	ticket.UncommittedEvents = append(ticket.UncommittedEvents, domainEvent)
}

type TicketCreated struct {
	TicketID uuid.UUID
	Info     TicketInfo
}

type TicketAssigned struct {
	TicketID uuid.UUID
	Assignee uuid.UUID
}
