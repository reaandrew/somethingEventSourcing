package core

type EventPublisher interface {
	Publish(events []DomainEvent) error
}
