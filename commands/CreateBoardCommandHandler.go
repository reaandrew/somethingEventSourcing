package commands

import (
	"github.com/reaandrew/eventsourcing-in-go/domain/models"
	"github.com/reaandrew/eventsourcing-in-go/domain/services"
)

type CreateBoardCommandHandler struct {
	DomainRepository services.DomainRepository
}

func (handler CreateBoardCommandHandler) Execute(command CreateBoardCommand) (returnErr error) {
	var board = models.NewBoard(models.BoardInfo{
		Name:    command.Name,
		Columns: command.Columns,
	})
	handler.DomainRepository.Save(board)
	return
}
