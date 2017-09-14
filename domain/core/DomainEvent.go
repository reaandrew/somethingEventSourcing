package core

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

type DomainEvent struct {
	ID        uuid.UUID
	Timestamp time.Time
	Version   int
	Data      interface{}
}
