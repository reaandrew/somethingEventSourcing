package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/reaandrew/forora/commands"
	"github.com/reaandrew/forora/queries"
	uuid "github.com/satori/go.uuid"
)

func SetupRouter(commandExecutor commands.CommandExecutor,
	queryExecutor queries.QueryExecutor) *gin.Engine {
	r := gin.Default()
	var v1 = r.Group("/v1")
	{
		var boards = v1.Group("/boards")
		{
			boards.POST("", func(c *gin.Context) {
				var createBoardCommand commands.CreateBoardCommand
				c.BindJSON(&createBoardCommand)
				createBoardCommand.BoardID = uuid.NewV4().String()

				commandExecutor.Execute(createBoardCommand)

				c.JSON(202, NewApiResponse().AddLink(Link{
					Rel:  "self",
					Href: "/v1/boards/" + createBoardCommand.BoardID,
				}).AddLink(Link{
					Rel:  "tickets",
					Href: "/v1/boards/" + createBoardCommand.BoardID + "/tickets",
				}))
			})

			boards.GET("", func(c *gin.Context) {
				var query = queries.GetAllBoardsRequest{}

				var response, err = queryExecutor.Execute(query)

				if err != nil {
					//These need to change to return server errors
					panic(err)
				}

				var getAllBoardsResponse = response.(queries.GetAllBoardsResponse)

				c.JSON(200, NewApiResponse().
					SetData("boards", getAllBoardsResponse.Boards))
			})

			boards.GET("/:id", func(c *gin.Context) {
				var query = queries.GetBoardByIDRequest{
					BoardID: c.Param("id"),
				}

				var response, err = queryExecutor.Execute(query)

				if err != nil {
					panic(err)
				}

				var getBoardResponse = response.(queries.GetBoardByIDResponse)

				c.JSON(200, NewApiResponse().
					AddLink(Link{
						Rel:  "self",
						Href: "/v1/boards/" + getBoardResponse.Board.ID,
					}).
					AddLink(Link{
						Rel:  "tickets",
						Href: "/v1/boards/" + getBoardResponse.Board.ID + "/tickets",
					}).
					SetData("board", getBoardResponse.Board))
			})

			boards.POST("/:id/tickets/", func(c *gin.Context) {
				var createTicketCommand commands.CreateTicketCommand

				c.BindJSON(&createTicketCommand)
				createTicketCommand.BoardID = c.Param("id")
				createTicketCommand.TicketID = uuid.NewV4().String()

				commandExecutor.Execute(createTicketCommand)

				var boardLink = "/v1/boards/" + createTicketCommand.BoardID
				var ticketLink = boardLink + "/tickets/" + createTicketCommand.TicketID

				c.JSON(202, NewApiResponse().AddLink(
					Link{
						Rel:  "self",
						Href: ticketLink,
					}))
			})
		}
	}

	return r
}
