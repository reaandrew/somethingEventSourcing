package main

import (
	"github.com/coreos/go-systemd/daemon"
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

	daemon.SdNotify(false, "READY=1")
	r.Run(":9000")
}
