package model

import (
	"encoding/json"
	"io/ioutil"
	"perftest/http/logger"
)

var log = logger.LOGGER

/**
* struct HttpRequest 
*/
type HttpRequest struct {
	TransId string			 `json:"transId"`
	Url    string            `json:"url"`
	Method string            `json:"method"`
	Params map[string]string `json:"params,omitempty"`
	Body   string            `json:"body,omitempty"`

	Headers map[string]string `json:"headers,omitempty"`

	HttpResponse string `json:"httpResponse,omitempty"`
}

func (request *HttpRequest) ToJSON() string {
	data, err := json.Marshal(*request)
	if err != nil {
		log.Errorf("json marshal happens error, error msg: %s", err.Error())
		return ""
	}
	return string(data)
}

func FromJsonFile(filename string) (httpRequest *HttpRequest, err error) {
	log.Infof("Get HttpRequest object from json file: %s", filename)

	byteContent, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Errorln("Read file happen error, ", err.Error())
		return nil, err
	}
	if len(byteContent) == 0 {
		log.Errorln("File content is empty.")
	}
	_httpReq := &HttpRequest{}
	json.Unmarshal(byteContent, _httpReq)
	return _httpReq, nil
}

/**
* struct HttpResponse 
*/
type HttpResponse struct {
	HttpRequest HttpRequest `json:"httpRequest"`
	StatusCode  int         `json:"statusCode"`
	Body        string      `json:"body,omitempty"`

	Headers map[string][]string `json:"headers,omitempty"`
}

func (response *HttpResponse) ToJSON() string {
	data, err := json.Marshal(*response)
	if err != nil {
		log.Errorf("json marshal happens error, error msg: %s", err.Error())
		return ""
	}
	return string(data)
}