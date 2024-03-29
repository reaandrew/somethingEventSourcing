package services_test

import (
	"testing"

	"github.com/reaandrew/forora/domain/models"
	"github.com/reaandrew/forora/domain/services"
	"github.com/reaandrew/forora/infrastructure/inmemory"
	"github.com/reaandrew/forora/test"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

func TestDomainRepositoryGetAggregateReplaysEvents(t *testing.T) {
	var sut = test.NewSystemUnderTest()

	var board = models.NewBoard(models.BoardInfo{
		Columns: []models.BoardColumnInfo{
			models.BoardColumnInfo{Name: "A", ID: uuid.NewV4()},
			models.BoardColumnInfo{Name: "B", ID: uuid.NewV4()},
			models.BoardColumnInfo{Name: "C", ID: uuid.NewV4()},
		},
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

func TestCommitClearsCommittedEventsFromAllAggregatesSaved(t *testing.T) {
	var domainRepository = services.NewDomainRepository(
		inmemory.NewInMemoryEventStore(),
		inmemory.NewInMemoryEventPublisher(map[string]interface{}{}))

	var sample1 = test.NewSampleAggregate()
	var sample2 = test.NewSampleAggregate()

	domainRepository.Save(sample1)
	domainRepository.Save(sample2)

	domainRepository.Commit()

	assert.Equal(t, 0, len(sample1.GetUncommittedEvents()))
	assert.Equal(t, 0, len(sample2.GetUncommittedEvents()))
}
