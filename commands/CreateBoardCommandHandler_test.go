package commands_test

import (
	"testing"

	"github.com/reaandrew/forora/commands"
	"github.com/reaandrew/forora/domain/models"
	"github.com/reaandrew/forora/test"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

func TestCreateBoardCommandPublishesBoardCreated(t *testing.T) {
	var sut = test.NewSystemUnderTest()
	var command = commands.CreateBoardCommand{
		BoardID: uuid.NewV4().String(),
		Name:    "some board",
		Columns: []string{
			"todo",
			"doing",
			"done",
		},
	}

	var err = sut.CommandExecutor.Execute(command)

	assert.Nil(t, err)
	assert.Equal(t, 1, sut.NumberOfEventsPublished())
	assert.IsType(t, models.BoardCreated{}, sut.GetEvent(0))
}

func TestReturnErrInvalidBoardID(t *testing.T) {
	var sut = test.NewSystemUnderTest()

	var command = commands.CreateBoardCommand{
		BoardID: "",
	}

	var err = sut.CommandExecutor.Execute(command)

	assert.Equal(t, err, models.ErrInvalidBoardID)
}
