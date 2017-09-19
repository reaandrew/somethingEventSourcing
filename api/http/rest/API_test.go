package rest_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/reaandrew/eventsourcing-in-go/api/http/rest"
	"github.com/reaandrew/eventsourcing-in-go/commands"
	"github.com/reaandrew/eventsourcing-in-go/queries"
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

	var jsonBytes, _ = json.Marshal(command)
	var reader = bytes.NewReader(jsonBytes)
	var request, _ = http.NewRequest("POST", "/v1/boards", reader)
	var resp = httptest.NewRecorder()
	var router *gin.Engine
	router = rest.SetupRouter(sut.CommandExecutor, sut.QueryExecutor)

	router.ServeHTTP(resp, request)

	assert.Equal(t, resp.Code, 202)

	var apiResponse = rest.LoadApiResponse(resp.Body.Bytes())

	assert.Equal(t, len(apiResponse.Links()), 1)
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
}

func TestRetrievingABoardAfterCreation(t *testing.T) {
	var sut = test.NewSystemUnderTest()
	var command = commands.CreateBoardCommand{
		Name: "Some Board",
		Columns: []string{
			"todo",
			"doing",
			"done",
		},
	}

	var jsonBytes, _ = json.Marshal(command)
	var reader = bytes.NewReader(jsonBytes)
	var request, _ = http.NewRequest("POST", "/v1/boards", reader)
	var resp = httptest.NewRecorder()
	var router *gin.Engine
	router = rest.SetupRouter(sut.CommandExecutor, sut.QueryExecutor)

	router.ServeHTTP(resp, request)

	var apiResponse = rest.LoadApiResponse(resp.Body.Bytes())

	var link = apiResponse.LinkForRel("self")

	var getBoardRequest, _ = http.NewRequest("GET", link.Href, nil)
	var getBoardResp = httptest.NewRecorder()
	router.ServeHTTP(getBoardResp, getBoardRequest)

	var response queries.GetBoardByIDResponse

	json.Unmarshal(getBoardResp.Body.Bytes(), &response)

	assert.Equal(t, len(command.Columns), len(response.Board.Columns))
}
