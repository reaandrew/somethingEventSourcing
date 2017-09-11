package commands_test

import (
	"testing"

	"github.com/reaandrew/eventsourcing-in-go/commands"
	"github.com/reaandrew/eventsourcing-in-go/domain"
	"github.com/reaandrew/eventsourcing-in-go/domain/services"
	"github.com/reaandrew/eventsourcing-in-go/infrastructure/inmemory"
	"github.com/stretchr/testify/assert"
)

func TestCreateBoardCommandPublishesBoardCreated(t *testing.T) {
	var eventStore = inmemory.NewInMemoryEventStore()
	var eventPublisher = inmemory.NewInMemoryEventPublisher()
	var domainRepository = services.NewDomainRepository(eventStore, eventPublisher)
	var commandExecutor = commands.NewCommandExecutor(domainRepository)

	var command = commands.CreateBoardCommand{
		Name: "some board",
		Columns: []string{
			"todo",
			"doing",
			"done",
		},
	}

	var err = commandExecutor.Execute(command)

	assert.Nil(t, err)
	assert.Equal(t, 1, eventPublisher.NumberOfEventsPublished())
	assert.IsType(t, domain.BoardCreated{}, eventPublisher.GetEvent(0))
}
