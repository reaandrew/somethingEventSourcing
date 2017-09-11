package core

import uuid "github.com/satori/go.uuid"

type Aggregate interface {
	GetCommittedEvents() []interface{}
	GetID() uuid.UUID
	GetVersion() int
}
