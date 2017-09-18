package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/reaandrew/eventsourcing-in-go/commands"
	"github.com/reaandrew/eventsourcing-in-go/domain/services"
	"github.com/reaandrew/eventsourcing-in-go/infrastructure/inmemory"
	uuid "github.com/satori/go.uuid"
)

type HttpLink struct {
	Rel  string `json:"rel"`
	Href string `json:"href"`
}

type ApiCommandResponse struct {
	Links []HttpLink        `json:"links"`
	Meta  map[string]string `json:"meta"`
}

func SetupRouter(commandExecutor commands.CommandExecutor) *gin.Engine {
	r := gin.Default()
	var v1 = r.Group("/v1")
	{
		var boards = v1.Group("/boards")
		{
			boards.POST("/", func(c *gin.Context) {
				var createBoardCommand commands.CreateBoardCommand
				c.BindJSON(&createBoardCommand)
				createBoardCommand.BoardID = uuid.NewV4().String()

				commandExecutor.Execute(createBoardCommand)

				c.JSON(202, ApiCommandResponse{
					Links: []HttpLink{
						HttpLink{
							Rel:  "self",
							Href: "/v1/boards/" + createBoardCommand.BoardID,
						},
					},
				})
			})
		}
	}

	return r
}

func main() {
	var eventStore = inmemory.NewInMemoryEventStore()
	var eventPublisher = inmemory.NewInMemoryEventPublisher()
	var domainRepository = services.NewDomainRepository(eventStore, eventPublisher)
	var commandExecutor = commands.NewCommandExecutor(domainRepository)
	var r = SetupRouter(commandExecutor)
	r.Run() // listen and serve on 0.0.0.0:8080
}
