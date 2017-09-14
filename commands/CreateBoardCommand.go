package commands

type CreateBoardCommand struct {
	BoardID string
	Name    string
	Columns []string
}
