package test

import (
	"time"

	"github.com/reaandrew/forora/domain/core"
	uuid "github.com/satori/go.uuid"
)

type SampleEvent struct {
}

type SampleAggregateCreated struct {
	SampleAggregateID uuid.UUID
}

type SampleAggregate struct {
	UncommittedEvents []core.DomainEvent
	ID                uuid.UUID
	version           int
}

func (sample *SampleAggregate) GetUncommittedEvents() (events []core.DomainEvent) {
	return sample.UncommittedEvents
}

func (sample *SampleAggregate) GetID() (returnID uuid.UUID) {
	return sample.ID
}

func (sample *SampleAggregate) GetVersion() (version int) {
	version = sample.version
	return
}

func (sample *SampleAggregate) Commit() {
	sample.UncommittedEvents = []core.DomainEvent{}
}

func (sample *SampleAggregate) handleSampleAggregateCreated(event SampleAggregateCreated) {
	sample.ID = event.SampleAggregateID
}

func (sample *SampleAggregate) Load(events []core.DomainEvent) {
	for _, event := range events {
		sample.replay(event)
	}
}

func (sample *SampleAggregate) apply(event core.DomainEvent) {
	switch e := event.Data.(type) {
	case SampleAggregateCreated:
		sample.handleSampleAggregateCreated(e)
	default:

	}
}

func (sample *SampleAggregate) replay(domainEvent core.DomainEvent) {
	sample.apply(domainEvent)
	sample.version = domainEvent.Version
}

func (sample *SampleAggregate) publish(event interface{}) {
	var domainEvent = core.DomainEvent{
		ID:        uuid.NewV4(),
		Version:   sample.version + len(sample.UncommittedEvents) + 1,
		Timestamp: time.Now(),
		Data:      event,
	}
	sample.apply(domainEvent)
	sample.UncommittedEvents = append(sample.UncommittedEvents, domainEvent)
}

func NewSampleAggregate() (newSampleAggregate *SampleAggregate) {
	newSampleAggregate = &SampleAggregate{}

	newSampleAggregate.publish(SampleAggregateCreated{
		SampleAggregateID: uuid.NewV4(),
	})
	return
}
