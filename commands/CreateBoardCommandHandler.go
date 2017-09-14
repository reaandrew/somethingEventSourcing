package commands

import (
	"github.com/reaandrew/eventsourcing-in-go/domain/models"
	"github.com/reaandrew/eventsourcing-in-go/domain/services"
	uuid "github.com/satori/go.uuid"
)

type CreateBoardCommandHandler struct {
	DomainRepository services.DomainRepository
}

func (handler CreateBoardCommandHandler) Execute(command CreateBoardCommand) (returnErr error) {

	var boardID, err = uuid.FromString(command.BoardID)
	if err != nil {
		returnErr = models.ErrInvalidBoardID
		return
	}

	var board = models.NewBoard(models.BoardInfo{
		BoardID: boardID,
		Name:    command.Name,
		Columns: command.Columns,
	})
	handler.DomainRepository.Save(board)
	return
}
