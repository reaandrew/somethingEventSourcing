package rest

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/reaandrew/forora/commands"
	"github.com/reaandrew/forora/queries"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

func SetupRouter(commandExecutor commands.CommandExecutor,
	queryExecutor queries.QueryExecutor) *gin.Engine {
	r := gin.Default()
	r.LoadHTMLGlob("./data/templates/*")
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
	// Credentials which stores google ids.
	type Credentials struct {
		Cid     string `json:"cid"`
		Csecret string `json:"csecret"`
	}

	var c = Credentials{
		Cid:     os.Getenv("GOOGLE_CID"),
		Csecret: os.Getenv("GOOGLE_CSECRET"),
	}
	conf := &oauth2.Config{
		ClientID:     c.Cid,
		ClientSecret: c.Csecret,
		RedirectURL:  "http://localhost:8000/auth",
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email", // You have to select your own scope from here -> https://developers.google.com/identity/protocols/googlescopes#google_sign-in
		},
		Endpoint: google.Endpoint,
	}
	var auth = r.Group("")
	{
		auth.GET("/", func(c *gin.Context) {
			c.HTML(http.StatusOK, "index.tmpl", gin.H{
				"LoginURL": conf.AuthCodeURL(uuid.NewV4().String()),
			})
		})
	}

	return r
}
