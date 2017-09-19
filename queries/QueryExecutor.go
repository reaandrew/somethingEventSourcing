package queries

import "errors"

var (
	ErrQueryNotSupported = errors.New("QueryNotSupported")
	ErrIDNotFound        = errors.New("IDNotFound")
)

type QueryExecutor interface {
	Execute(request interface{}) (response interface{}, err error)
}
