package commands

import (
	"github.com/reaandrew/eventsourcing-in-go/domain"
	"github.com/reaandrew/eventsourcing-in-go/domain/services"
)

type CreateBoardCommandHandler struct {
	DomainRepository services.DomainRepository
}

func (handler CreateBoardCommandHandler) Execute(command CreateBoardCommand) (returnErr error) {
	var board = domain.NewBoard(domain.BoardInfo{
		Name:    command.Name,
		Columns: command.Columns,
	})
	handler.DomainRepository.Save(board)
	return
}
