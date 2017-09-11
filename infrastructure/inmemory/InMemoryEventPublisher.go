package inmemory

type InMemoryEventPublisher struct {
	events []interface{}
}

func (publisher *InMemoryEventPublisher) Publish(events []interface{}) (err error) {
	publisher.events = append(publisher.events, events...)
	return
}

func (publisher *InMemoryEventPublisher) NumberOfEventsPublished() (value int) {
	value = len(publisher.events)
	return
}

func (publisher *InMemoryEventPublisher) GetEvent(index int) (value interface{}) {
	value = publisher.events[index]
	return
}

func NewInMemoryEventPublisher() (newPublisher *InMemoryEventPublisher) {
	newPublisher = &InMemoryEventPublisher{
		events: []interface{}{},
	}

	return
}
