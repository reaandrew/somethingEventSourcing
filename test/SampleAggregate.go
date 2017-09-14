package test

import (
	"time"

	"github.com/reaandrew/eventsourcing-in-go/domain/core"
	uuid "github.com/satori/go.uuid"
)

type SampleEvent struct {
}

type SampleAggregateCreated struct {
	SampleAggregateID uuid.UUID
}

type SampleAggregate struct {
	CommittedEvents []core.DomainEvent
	ID              uuid.UUID
	version         int
}

func (sample *SampleAggregate) GetCommittedEvents() (events []core.DomainEvent) {
	return sample.CommittedEvents
}

func (sample *SampleAggregate) GetID() (returnID uuid.UUID) {
	return sample.ID
}

func (sample *SampleAggregate) GetVersion() (version int) {
	version = sample.version
	return
}

func (sample *SampleAggregate) Commit() {
	sample.CommittedEvents = []core.DomainEvent{}
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
		Version:   sample.version + len(sample.CommittedEvents) + 1,
		Timestamp: time.Now(),
		Data:      event,
	}
	sample.apply(domainEvent)
	sample.CommittedEvents = append(sample.CommittedEvents, domainEvent)
}

func NewSampleAggregate() (newSampleAggregate *SampleAggregate) {
	newSampleAggregate = &SampleAggregate{}

	newSampleAggregate.publish(SampleAggregateCreated{
		SampleAggregateID: uuid.NewV4(),
	})
	return
}
