package commands

import (
	"github.com/reaandrew/eventsourcing-in-go/domain"
	"github.com/reaandrew/eventsourcing-in-go/domain/services"
	uuid "github.com/satori/go.uuid"
)

type CreateTicketCommandHandler struct {
	DomainRepository services.DomainRepository
}

func (handler CreateTicketCommandHandler) Execute(command CreateTicketCommand) (returnErr error) {
	var boardID, idErr = uuid.FromString(command.BoardID)
	if idErr != nil {
		returnErr = idErr
		return
	}

	var _, err = handler.DomainRepository.GetBoard(boardID)
	if err != nil {
		returnErr = err
		return
	}

	var ticketInfo = domain.TicketInfo{}
	var assigneeID, assigneeErr = uuid.FromString(command.Assignee)
	if command.Assignee != "" && assigneeErr != nil {
		returnErr = domain.ErrInvalidAssigneeID
		return
	}
	if command.Assignee != "" {
		ticketInfo.Assignee = assigneeID
	}

	ticketInfo.Title = command.Title
	ticketInfo.Content = command.Content

	var ticket, ticketErr = domain.NewTicket(ticketInfo)
	if ticketErr != nil {
		returnErr = ticketErr
		return
	}

	handler.DomainRepository.Save(ticket)
	return
}
