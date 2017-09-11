package commands_test

import (
	"testing"

	"github.com/reaandrew/eventsourcing-in-go/commands"
	"github.com/reaandrew/eventsourcing-in-go/domain"
	"github.com/reaandrew/eventsourcing-in-go/test"
	"github.com/stretchr/testify/assert"
)

func TestCreateBoardCommandPublishesBoardCreated(t *testing.T) {
	var sut = test.NewSystemUnderTest()
	var command = commands.CreateBoardCommand{
		Name: "some board",
		Columns: []string{
			"todo",
			"doing",
			"done",
		},
	}

	var err = sut.CommandExecutor.Execute(command)

	assert.Nil(t, err)
	assert.Equal(t, 1, sut.NumberOfEventsPublished())
	assert.IsType(t, domain.BoardCreated{}, sut.GetEvent(0))
}
