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
	var request, _ = http.NewRequest("POST", "/v1/boards/", reader)
	var resp = httptest.NewRecorder()
	var router *gin.Engine
	router = rest.SetupRouter(sut.CommandExecutor)

	router.ServeHTTP(resp, request)

	assert.Equal(t, resp.Code, 202)

	var apiCommandResponse rest.ApiCommandResponse

	json.Unmarshal(resp.Body.Bytes(), &apiCommandResponse)

	assert.Equal(t, len(apiCommandResponse.Links), 1)
}
