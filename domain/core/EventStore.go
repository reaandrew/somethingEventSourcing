package core

import uuid "github.com/satori/go.uuid"

type EventStore interface {
	Save(aggregate Aggregate) error
	GetEvents(id uuid.UUID) ([]interface{}, error)
}
