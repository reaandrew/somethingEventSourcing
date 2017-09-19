package commands

import (
	"github.com/reaandrew/forora/domain/models"
	"github.com/reaandrew/forora/domain/services"
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
	var boardInfo = models.BoardInfo{
		BoardID: boardID,
		Name:    command.Name,
		Columns: []models.BoardColumnInfo{},
	}
	for _, colName := range command.Columns {
		boardInfo.Columns = append(boardInfo.Columns, models.BoardColumnInfo{
			ID:   uuid.NewV4(),
			Name: colName,
		})
	}
	var board = models.NewBoard(boardInfo)
	handler.DomainRepository.Save(board)
	return
}
