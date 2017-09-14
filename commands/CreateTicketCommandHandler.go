package commands

import (
	"fmt"

	"github.com/reaandrew/eventsourcing-in-go/domain/models"
	"github.com/reaandrew/eventsourcing-in-go/domain/services"
	uuid "github.com/satori/go.uuid"
)

type CreateTicketCommandHandler struct {
	DomainRepository services.DomainRepository
}

func (handler CreateTicketCommandHandler) Execute(command CreateTicketCommand) (returnErr error) {
	var boardID, idErr = uuid.FromString(command.BoardID)
	if idErr != nil {
		returnErr = ErrInvalidBoardID
		return
	}

	var board, err = handler.DomainRepository.GetBoard(boardID)
	if err != nil {
		returnErr = err
		return
	}

	var ticketInfo = models.TicketInfo{}
	var assigneeID, assigneeErr = uuid.FromString(command.Assignee)
	if command.Assignee != "" && assigneeErr != nil {
		returnErr = ErrInvalidAssigneeID
		return
	}
	if command.Assignee != "" {
		ticketInfo.Assignee = assigneeID
	}

	var ticketID, ticketIDErr = uuid.FromString(command.TicketID)
	if ticketIDErr != nil {
		returnErr = models.ErrInvalidTicketID
		return
	}

	ticketInfo.TicketID = ticketID
	ticketInfo.Title = command.Title
	ticketInfo.Content = command.Content

	var ticket, ticketErr = models.NewTicket(ticketInfo)
	if ticketErr != nil {
		returnErr = ticketErr
		return
	}

	board.AddTicket(ticket, command.Column)

	fmt.Println("Saving Ticket", ticket.GetID())
	handler.DomainRepository.Save(ticket)
	handler.DomainRepository.Save(board)

	return
}
