package services_test

import (
	"testing"

	"github.com/reaandrew/eventsourcing-in-go/domain"
	"github.com/reaandrew/eventsourcing-in-go/test"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

func TestDomainRepositoryGetAggregateReplaysEvents(t *testing.T) {
	var sut = test.NewSystemUnderTest()

	var board = domain.NewBoard(domain.BoardInfo{
		Columns: []string{"A", "B", "C"},
	})
	sut.DomainRepository.Save(board)

	savedBoard, err := sut.DomainRepository.GetBoard(board.GetID())

	assert.Nil(t, err)

	assert.Equal(t, savedBoard.GetVersion(), 1)
}

func TestDomainRepositoryGetAggregateReturnsErrorWhenIDNotFound(t *testing.T) {
	var sut = test.NewSystemUnderTest()

	_, err := sut.DomainRepository.GetBoard(uuid.NewV4())

	assert.NotNil(t, err)
}
