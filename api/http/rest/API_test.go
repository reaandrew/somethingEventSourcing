package rest_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/reaandrew/eventsourcing-in-go/api/http/rest"
	"github.com/reaandrew/eventsourcing-in-go/commands"
	"github.com/reaandrew/eventsourcing-in-go/test"
	"github.com/stretchr/testify/assert"
)

func TestCreatingABoard(t *testing.T) {
	var sut = test.NewSystemUnderTest()

	var command = commands.CreateBoardCommand{
		Name: "Some Board",
		Columns: []string{
			"todo",
			"doing",
			"done",
		},
	}

	var resp = sut.Post(command, "/v1/boards")

	assert.Equal(t, resp.Code, 202)

	var apiResponse = WrapApiResponseWithAssertions(rest.LoadApiResponse(resp.Body.Bytes()), t)
	apiResponse.
		AssertLinkCount(2).
		AssertLink("self", "/v1/boards/:id").
		AssertLink("tickets", "/v1/boards/:id/tickets")
}

func TestAddingATicketToABoard(t *testing.T) {
	var sut = test.NewSystemUnderTest()
	var boardID = sut.CreateSampleBoard("some name")

	var command = commands.CreateTicketCommand{
		Column: "todo",
		Title:  "some ticket",
	}

	var jsonBytes, _ = json.Marshal(command)
	var reader = bytes.NewReader(jsonBytes)
	var url = "/v1/boards/" + boardID.String() + "/tickets/"
	var request, _ = http.NewRequest("POST", url, reader)
	var resp = httptest.NewRecorder()
	var router *gin.Engine
	router = rest.SetupRouter(sut.CommandExecutor, sut.QueryExecutor)

	router.ServeHTTP(resp, request)

	assert.Equal(t, resp.Code, 202)

	var apiResponse = WrapApiResponseWithAssertions(rest.LoadApiResponse(resp.Body.Bytes()), t)
	apiResponse.
		AssertLinkCount(1).
		AssertLink("self", "/v1/boards/:id/tickets/:id")
}

func TestGettingABoard(t *testing.T) {
	var sut = test.NewSystemUnderTest()
	var boardID = sut.CreateSampleBoard("some board")

	var req, _ = http.NewRequest("GET", fmt.Sprintf("/v1/boards/%s", boardID), nil)
	var resp = httptest.NewRecorder()
	var router *gin.Engine
	router = rest.SetupRouter(sut.CommandExecutor, sut.QueryExecutor)
	router.ServeHTTP(resp, req)

	var apiResponse = WrapApiResponseWithAssertions(rest.LoadApiResponse(resp.Body.Bytes()), t)

	apiResponse.
		AssertLink("self", "/v1/boards/:id").
		AssertLink("tickets", "/v1/boards/:id/tickets").
		AssertData("board")
}
