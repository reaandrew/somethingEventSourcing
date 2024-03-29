package commands

import (
	"github.com/reaandrew/forora/domain/services"
)

type CommandExecutor struct {
	DomainRepository services.DomainRepository
}

func (executor CommandExecutor) Execute(command interface{}) (err error) {
	switch c := command.(type) {
	case CreateBoardCommand:
		var handler = CreateBoardCommandHandler{
			DomainRepository: executor.DomainRepository,
		}

		return handler.Execute(c)
	case CreateTicketCommand:
		var handler = CreateTicketCommandHandler{
			DomainRepository: executor.DomainRepository,
		}

		return handler.Execute(c)
	case AssignTicketCommand:
		var handler = AssignTicketCommandHandler{
			DomainRepository: executor.DomainRepository,
		}
		return handler.Execute(c)
	default:
	}
	return
}

func NewCommandExecutor(domainRepository services.DomainRepository) (newExecutor CommandExecutor) {
	return CommandExecutor{
		DomainRepository: domainRepository,
	}
}
