package commands

import (
	"github.com/reaandrew/forora/domain/models"
	"github.com/reaandrew/forora/domain/services"
	uuid "github.com/satori/go.uuid"
)

type AssignTicketCommandHandler struct {
	DomainRepository services.DomainRepository
}

func (handler AssignTicketCommandHandler) Execute(command AssignTicketCommand) (returnErr error) {

	var ticketID, err = uuid.FromString(command.TicketID)
	if err != nil {
		returnErr = models.ErrInvalidTicketID
		return
	}

	var ticket, ticketErr = handler.DomainRepository.GetTicket(ticketID)

	if ticketErr != nil {
		returnErr = ticketErr
		return
	}

	var userID, userErr = uuid.FromString(command.Assignee)
	if userErr != nil {
		returnErr = ErrInvalidAssigneeID
		return
	}

	var ticketAssignmentErr = ticket.AssignTo(userID)

	if ticketAssignmentErr != nil {
		returnErr = ticketAssignmentErr
		return
	}

	handler.DomainRepository.Save(ticket)
	return
}
