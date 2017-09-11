package services_test

import (
	"testing"

	"github.com/reaandrew/eventsourcing-in-go/domain"
	"github.com/reaandrew/eventsourcing-in-go/domain/services"
	"github.com/reaandrew/eventsourcing-in-go/infrastructure/inmemory"
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
	var eventStore = inmemory.NewInMemoryEventStore()
	var eventPublisher = inmemory.NewInMemoryEventPublisher()
	var domainRepository = services.NewDomainRepository(eventStore, eventPublisher)
	_, err := domainRepository.GetBoard(uuid.NewV4())

	assert.NotNil(t, err)
}
