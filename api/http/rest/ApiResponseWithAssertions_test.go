package rest_test

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/reaandrew/eventsourcing-in-go/api/http/rest"
	"github.com/stretchr/testify/assert"
)

var (
	ErrNoLinkFound = errors.New("ErrNoLinkFound")
)

type ApiResponseWithAssertions struct {
	obj rest.ApiResponse
	t   *testing.T
}

func (response ApiResponseWithAssertions) LinkForRel(rel string) (returnLink rest.Link, err error) {
	for _, link := range response.obj.Links() {
		if link.Rel == rel {
			returnLink = link
			break
		}
	}
	return
}

func (response ApiResponseWithAssertions) AssertLinkCount(count int) (next ApiResponseWithAssertions) {
	next = response
	assert.Equal(response.t, len(response.obj.Links()), count)
	return
}

func (response ApiResponseWithAssertions) AssertData(key string) (next ApiResponseWithAssertions) {
	next = response
	if _, ok := response.obj[key]; !ok {
		assert.Fail(response.t, fmt.Sprintf("data not found key=%s", key))
	}
	return
}

func (response ApiResponseWithAssertions) AssertLink(rel string, pattern string) (next ApiResponseWithAssertions) {
	next = response

	var _, err = response.LinkForRel(rel)
	assert.Nil(response.t, err)

	var runPattern = strings.Replace(pattern, ":id", `[\w\d-^/]{36}`, -1)
	var urlPattern = regexp.MustCompile(runPattern)
	for _, link := range response.obj.Links() {
		if link.Rel == rel &&
			urlPattern.MatchString(link.Href) {
			return
		}
	}
	assert.Fail(response.t, fmt.Sprintf("Not Found rel=%s pattern=%s", rel, pattern))
	return
}

func WrapApiResponseWithAssertions(response rest.ApiResponse, t *testing.T) (responseWithAssertions ApiResponseWithAssertions) {
	responseWithAssertions = ApiResponseWithAssertions{
		obj: response,
		t:   t,
	}

	return
}
