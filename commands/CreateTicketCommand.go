package commands

type CreateTicketCommand struct {
	BoardID  string
	Title    string
	Content  string
	Assignee string
}
