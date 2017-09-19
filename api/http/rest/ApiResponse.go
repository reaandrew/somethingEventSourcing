package rest

import "encoding/json"

type HttpLink struct {
	Rel  string `json:"rel"`
	Href string `json:"href"`
}

const ApiResponseLinkKey = "links"

type ApiResponse map[string]interface{}

func NewApiResponse() (newResponse ApiResponse) {
	newResponse = map[string]interface{}{}
	newResponse[ApiResponseLinkKey] = []HttpLink{}
	return
}

func LoadApiResponse(data []byte) (response ApiResponse) {
	json.Unmarshal(data, &response)
	var links = response[ApiResponseLinkKey].([]interface{})
	var returnArray = []HttpLink{}
	for index, _ := range links {
		var obj = links[index].(map[string]interface{})
		var link = HttpLink{
			Href: obj["href"].(string),
			Rel:  obj["rel"].(string),
		}
		returnArray = append(returnArray, link)
	}

	response[ApiResponseLinkKey] = returnArray

	return
}

func (response ApiResponse) AddLink(link HttpLink) (next ApiResponse) {
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

func (response ApiResponse) LinkForRel(rel string) (returnLink HttpLink) {
	for _, link := range response.Links() {
		if link.Rel == rel {
			returnLink = link
			break
		}
	}
	return
}

func (response ApiResponse) Links() (links []HttpLink) {
	return response[ApiResponseLinkKey].([]HttpLink)
}
