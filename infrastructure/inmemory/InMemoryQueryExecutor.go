package inmemory

import (
	"github.com/reaandrew/eventsourcing-in-go/queries"
	"github.com/reaandrew/eventsourcing-in-go/queries/dtos"
)

type InMemoryQueryExecutor struct {
	Data map[string]interface{}
}

func (executor InMemoryQueryExecutor) Execute(request interface{}) (response interface{}, err error) {
	switch t := request.(type) {
	case queries.GetBoardByIDRequest:
		if _, ok := executor.Data["boards"]; ok {
			var boards = executor.Data["boards"].(map[string]dtos.Board)
			if board, ok := boards[t.BoardID]; ok {
				response = queries.GetBoardByIDResponse{
					Board: board,
				}
			} else {
				err = queries.ErrIDNotFound
			}
		} else {
			err = queries.ErrIDNotFound
		}
	default:
		err = queries.ErrQueryNotSupported
	}

	return
}

func NewInMemoryQueryExecutor(data map[string]interface{}) InMemoryQueryExecutor {
	return InMemoryQueryExecutor{
		Data: data,
	}
}
