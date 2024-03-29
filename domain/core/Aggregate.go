package core

import uuid "github.com/satori/go.uuid"

type Aggregate interface {
	GetUncommittedEvents() []DomainEvent
	GetID() uuid.UUID
	GetVersion() int
	Commit()
}
