package commands

import "errors"

var (
	ErrInvalidAssigneeID = errors.New("ErrInvalidAssigneeID")
	ErrInvalidBoardID    = errors.New("ErrInvalidBoardID")
)

type CreateTicketCommand struct {
	TicketID string
	BoardID  string
	Column   string
	Title    string
	Content  string
	Assignee string
}
