package main

import (
	"github.com/reaandrew/forora/api/http/rest"
	"github.com/reaandrew/forora/commands"
	"github.com/reaandrew/forora/domain/services"
	"github.com/reaandrew/forora/infrastructure/inmemory"
)

func main() {
	var queryStore = map[string]interface{}{}
	var queryExecutor = inmemory.NewInMemoryQueryExecutor(queryStore)

	var eventStore = inmemory.NewInMemoryEventStore()
	var eventPublisher = inmemory.NewInMemoryEventPublisher(queryStore)
	var domainRepository = services.NewDomainRepository(eventStore, eventPublisher)
	var commandExecutor = commands.NewCommandExecutor(domainRepository)

	var r = rest.SetupRouter(commandExecutor, queryExecutor)
	r.Run() // listen and serve on 0.0.0.0:8080
}
