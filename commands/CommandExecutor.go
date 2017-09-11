package commands

import (
	"github.com/reaandrew/eventsourcing-in-go/domain/services"
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
	}

	return
}

func NewCommandExecutor(domainRepository services.DomainRepository) (newExecutor CommandExecutor) {
	return CommandExecutor{
		DomainRepository: domainRepository,
	}
}
