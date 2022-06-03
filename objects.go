package salesforce_utils

import (
	"encoding/json"
	"fmt"
	"github.com/joomcode/errorx"
	"github.com/valyala/fasthttp"
	"net/http"
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
	body, statusCode, requestErr := s.sendRequest(req)
	if requestErr != nil {
		err = requestErr
		return
	}
	if statusCode != http.StatusCreated {
		err = errorx.Decorate(err, "unexpected status code: %d with body: %s", statusCode, body)
		return
	}
	err = json.Unmarshal(body, &response)
	return
}

func (s *SalesforceUtils) UpdateObject(typeName, id string, jsonBytes []byte) (response ObjectResponse, err error) {
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)
	uri := s.getObjectIdUrl(typeName, id)
	req.SetRequestURI(uri)
	req.Header.SetMethod(http.MethodPatch)
	req.Header.Set("Content-Type", "application/json")
	req.SetBody(jsonBytes)
	body, statusCode, requestErr := s.sendRequest(req)
	if requestErr != nil {
		err = requestErr
		return
	}
	if statusCode != http.StatusCreated {
		err = errorx.Decorate(err, "unexpected status code: %d with body: %s", statusCode, body)
		return
	}
	err = json.Unmarshal(body, &response)
	return
}

func (s *SalesforceUtils) DeleteObject(typeName, id string) error {
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)
	uri := s.getObjectIdUrl(typeName, id)
	req.SetRequestURI(uri)
	req.Header.SetMethod(http.MethodDelete)
	body, statusCode, err := s.sendRequest(req)
	if err != nil {
		return err
	}
	if statusCode != http.StatusNoContent {
		return errorx.Decorate(err, "unexpected status code: %d with body: %s", statusCode, body)
	}
	return nil
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
