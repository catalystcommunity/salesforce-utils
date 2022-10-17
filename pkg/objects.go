package pkg

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/joomcode/errorx"
	"github.com/valyala/fasthttp"
)

type ObjectResponse struct {
	Id      string
	Errors  []string
	Success bool
}

func (s *SalesforceUtils) CreateObject(typeName string, jsonBytes []byte) (response ObjectResponse, err error) {
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)
	uri := s.getTypeUrl(typeName)
	req.SetRequestURI(uri)
	req.Header.SetMethod(http.MethodPost)
	req.Header.Set("Content-Type", "application/json")
	req.SetBody(jsonBytes)
	body, statusCode, deferredFunc, requestErr := s.sendRequest(req)
	defer deferredFunc()
	if requestErr != nil {
		err = requestErr
		return
	}
	if statusCode != http.StatusCreated {
		err = errorx.IllegalState.New("unexpected status code: %d with body: %s", statusCode, body)
		return
	}
	err = json.Unmarshal(body, &response)
	return
}

func (s *SalesforceUtils) UpdateObject(typeName, id string, jsonBytes []byte) error {
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)
	uri := s.getObjectIdUrl(typeName, id)
	req.SetRequestURI(uri)
	req.Header.SetMethod(http.MethodPatch)
	req.Header.Set("Content-Type", "application/json")
	req.SetBody(jsonBytes)
	body, statusCode, deferredFunc, requestErr := s.sendRequest(req)
	defer deferredFunc()
	if requestErr != nil {
		return requestErr
	}
	if statusCode != http.StatusNoContent {
		return errorx.IllegalState.New("unexpected status code: %d with body: %s", statusCode, body)
	}
	return nil
}

func (s *SalesforceUtils) DeleteObject(typeName, id string) error {
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)
	uri := s.getObjectIdUrl(typeName, id)
	req.SetRequestURI(uri)
	req.Header.SetMethod(http.MethodDelete)
	body, statusCode, deferredFunc, err := s.sendRequest(req)
	defer deferredFunc()
	if err != nil {
		return err
	}
	if statusCode != http.StatusNoContent {
		return errorx.IllegalState.New("unexpected status code: %d with body: %s", statusCode, body)
	}
	return nil
}

// DescribeObjectResponse is a simplified struct representation of the json
// response from the "sObject Describe" API call. Currently only contains the
// "name" and "fields" fields.
type DescribeObjectResponse struct {
	Name   string                         `json:"name"`
	Fields []DescribeObjectResponseFields `json:"fields"`
}

// DescribeObjectResponseFields is a nested struct for the Fields field in the
// "sObject Describe" API response
type DescribeObjectResponseFields struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

// DescribeObject describes the object type, returning all of the field and types
func (s *SalesforceUtils) DescribeObject(typeName string) (response DescribeObjectResponse, err error) {
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)
	uri := s.getDescribeUrl(typeName)
	req.SetRequestURI(uri)
	req.Header.SetMethod(http.MethodGet)
	req.Header.Set("Content-Type", "application/json")
	body, statusCode, deferredFunc, requestErr := s.sendRequest(req)
	defer deferredFunc()
	if requestErr != nil {
		err = requestErr
		return
	}
	if statusCode != http.StatusOK {
		err = errorx.IllegalState.New("unexpected status code: %d with body: %s", statusCode, body)
		return
	}
	err = json.Unmarshal(body, &response)
	return
}

// getDataUrl gets a formatted url to the data endpoint
func (s *SalesforceUtils) getDataUrl() string {
	return fmt.Sprintf("%s/services/data/v%s/sobjects", s.Config.BaseUrl, s.Config.ApiVersion)
}

func (s *SalesforceUtils) getTypeUrl(typeName string) string {
	return fmt.Sprintf("%s/%s", s.getDataUrl(), typeName)
}

// getObjectIdUrl gets a formatted url to the endoint for a specific object by id
func (s *SalesforceUtils) getObjectIdUrl(typeName, id string) string {
	return fmt.Sprintf("%s/%s", s.getTypeUrl(typeName), id)
}

func (s *SalesforceUtils) getDescribeUrl(typeName string) string {
	return fmt.Sprintf("%s/describe", s.getTypeUrl(typeName))
}
