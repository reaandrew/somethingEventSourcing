package test_test

import (
	"testing"

	"github.com/reaandrew/eventsourcing-in-go/domain/models"
	"github.com/reaandrew/eventsourcing-in-go/test"
	"github.com/stretchr/testify/assert"
)

func TestCreateSampleBoard(t *testing.T) {
	var sut = test.NewSystemUnderTest()

	var boardID = sut.CreateSampleBoard("fubar")

	assert.NotNil(t, boardID)
	assert.Equal(t, 1, sut.NumberOfEventsPublished())
	assert.IsType(t, models.BoardCreated{}, sut.GetEvent(0))
}
