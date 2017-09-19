package rest_test

import (
	"fmt"
	"testing"

	"github.com/reaandrew/forora/api/http/rest"
	"github.com/reaandrew/forora/commands"
	"github.com/reaandrew/forora/test"
	"github.com/stretchr/testify/assert"
)

func TestGetAllBoards(t *testing.T) {
	var sut = test.NewSystemUnderTest()

	var numberOfBoards = 10
	sut.CreateSampleBoards(numberOfBoards)

	var resp = sut.Get("/v1/boards")

	assert.Equal(t, 200, resp.Code)

	var apiResponse = WrapApiResponseWithAssertions(rest.LoadApiResponse(resp.Body.Bytes()), t)

	apiResponse.
		AssertData("boards").
		AssertDataLength("boards", numberOfBoards)

}

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

	var url = "/v1/boards/" + boardID.String() + "/tickets/"
	var resp = sut.Post(command, url)

	assert.Equal(t, resp.Code, 202)

	var apiResponse = WrapApiResponseWithAssertions(rest.LoadApiResponse(resp.Body.Bytes()), t)
	apiResponse.
		AssertLinkCount(1).
		AssertLink("self", "/v1/boards/:id/tickets/:id")
}

func TestGettingABoard(t *testing.T) {
	var sut = test.NewSystemUnderTest()
	var boardID = sut.CreateSampleBoard("some board")

	var url = fmt.Sprintf("/v1/boards/%s", boardID)
	var resp = sut.Get(url)

	var apiResponse = WrapApiResponseWithAssertions(rest.LoadApiResponse(resp.Body.Bytes()), t)

	apiResponse.
		AssertLink("self", "/v1/boards/:id").
		AssertLink("tickets", "/v1/boards/:id/tickets").
		AssertData("board")
}
