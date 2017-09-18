package commands

type CreateBoardCommand struct {
	BoardID string   `json:"board_id"`
	Name    string   `json:"name"`
	Columns []string `json:"columns"`
}
