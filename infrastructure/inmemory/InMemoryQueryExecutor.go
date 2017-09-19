package inmemory

import (
	"github.com/reaandrew/forora/queries"
	"github.com/reaandrew/forora/queries/dtos"
)

const (
	BoardsKey = "boards"
)

type InMemoryQueryExecutor struct {
	Data map[string]interface{}
}

func (executor InMemoryQueryExecutor) Execute(request interface{}) (response interface{}, err error) {
	switch t := request.(type) {
	case queries.GetBoardByIDRequest:
		var boards = executor.Data[BoardsKey].(map[string]dtos.Board)
		if board, ok := boards[t.BoardID]; ok {
			response = queries.GetBoardByIDResponse{
				Board: board,
			}
		} else {
			err = queries.ErrIDNotFound
		}
	case queries.GetAllBoardsRequest:
		var boards = executor.Data[BoardsKey].(map[string]dtos.Board)
		var boardArray = []dtos.Board{}
		for _, board := range boards {
			boardArray = append(boardArray, board)
		}
		response = queries.GetAllBoardsResponse{
			Boards: boardArray,
		}
	default:
		err = queries.ErrQueryNotSupported
	}

	return
}

func (executor InMemoryQueryExecutor) setup() (next InMemoryQueryExecutor) {
	next = executor
	next.Data[BoardsKey] = map[string]dtos.Board{}
	return
}

func NewInMemoryQueryExecutor(data map[string]interface{}) InMemoryQueryExecutor {
	return InMemoryQueryExecutor{
		Data: data,
	}.setup()
}
