package rest

import (
	"encoding/json"
)

const ApiResponseLinkKey = "links"

type ApiResponse map[string]interface{}

func NewApiResponse() (newResponse ApiResponse) {
	newResponse = map[string]interface{}{}
	newResponse[ApiResponseLinkKey] = []Link{}
	return
}

func LoadApiResponse(data []byte) (response ApiResponse) {
	json.Unmarshal(data, &response)
	var links = response[ApiResponseLinkKey].([]interface{})
	var returnArray = []Link{}
	for index, _ := range links {
		var obj = links[index].(map[string]interface{})
		var link = Link{
			Href: obj["href"].(string),
			Rel:  obj["rel"].(string),
		}
		returnArray = append(returnArray, link)
	}

	response[ApiResponseLinkKey] = returnArray

	return
}

func (response ApiResponse) AddLink(link Link) (next ApiResponse) {
	next = response
	var linkArray = response.Links()
	linkArray = append(linkArray, link)
	next[ApiResponseLinkKey] = linkArray
	return
}

func (response ApiResponse) SetData(key string, data interface{}) (next ApiResponse) {
	next = response
	next[key] = data
	return
}

func (response ApiResponse) Links() (links []Link) {
	return response[ApiResponseLinkKey].([]Link)
}
