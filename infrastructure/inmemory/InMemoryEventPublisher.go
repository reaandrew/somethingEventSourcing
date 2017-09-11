package inmemory

type InMemoryEventPublisher struct {
	events []interface{}
}

func (publisher *InMemoryEventPublisher) Publish(events []interface{}) (err error) {
	publisher.events = append(publisher.events, events...)
	return
}

func NewInMemoryEventPublisher() (newPublisher *InMemoryEventPublisher) {
	newPublisher = &InMemoryEventPublisher{
		events: []interface{}{},
	}

	return
}
