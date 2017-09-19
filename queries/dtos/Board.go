package dtos

import "time"

type Board struct {
	ID      string
	Name    string
	Columns []BoardColumn
	Version int
	Updated time.Time
	Created time.Time
}
