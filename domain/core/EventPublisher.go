package core

type EventPublisher interface {
	Publish(events []interface{}) error
}
