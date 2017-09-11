package commands

import (
	"github.com/reaandrew/eventsourcing-in-go/domain"
	"github.com/reaandrew/eventsourcing-in-go/domain/services"
)

type CommandExecutor struct {
	DomainRepository services.DomainRepository
}

func (executor CommandExecutor) executeCreateBoardCommand(command CreateBoardCommand) {
	var board = domain.NewBoard(domain.BoardInfo{
		Name:    command.Name,
		Columns: command.Columns,
	})
	executor.DomainRepository.Save(board)
}

func (executor CommandExecutor) Execute(command interface{}) (err error) {
	switch c := command.(type) {
	case CreateBoardCommand:
		executor.executeCreateBoardCommand(c)
	case CreateTicketCommand:
		var handler = CreateTicketCommandHandler{
			DomainRepository: executor.DomainRepository,
		}

		return handler.Execute(c)
	}

	return
}

func NewCommandExecutor(domainRepository services.DomainRepository) (newExecutor CommandExecutor) {
	return CommandExecutor{
		DomainRepository: domainRepository,
	}
}
