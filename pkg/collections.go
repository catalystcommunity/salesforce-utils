package pkg

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/joomcode/errorx"
	"github.com/valyala/fasthttp"
)

// CollectionsRequest is used by the CollectionsCreateObjects and
// CollectionsUpdateObjects methods when interacting with the composite
// collections api
type CollectionsRequest struct {
	AllOrNone bool `json:"allOrNone"`
	// the use of json.RawMessage here is to avoid the need to unmarshal the
	// provided records json and remarshal it.
	Records []json.RawMessage `json:"records"`
}

// CollectionsResponseItem is the response item for a single object in the
// response from the collections api
type CollectionsResponseItem struct {
	Id      string                     `json:"id"`
	Success bool                       `json:"success"`
	Errors  []CollectionsResponseError `json:"errors"`
}

// CollectionsResponseError is the error response item for a single object in
// the response from the collections api
type CollectionsResponseError struct {
	StatusCode string   `json:"statusCode"`
	Message    string   `json:"message"`
	Fields     []string `json:"fields"`
}

// CollectionsCreateObjects creates objects in salesforce using the composite
// "collections" api. this implementation requires that you marshal your
// objects to json yourself with an attributes field to your object to define
// the type.
//
// ex:
//
//	{
//	  "attributes" : {"type" : "Account"},
//	  "Name" : "Example"
//	  ...
//	}
//
// ref: https://developer.salesforce.com/docs/atlas.en-us.api_rest.meta/api_rest/resources_composite_sobjects_collections_create.htm
func (s *SalesforceUtils) CollectionsCreateObjects(recordsJsonBytes [][]byte) (response []CollectionsResponseItem, err error) {
	err = validateCollectionsRequestLength(len(recordsJsonBytes))
	if err != nil {
		return nil, err
	}

	body, err := jsonRecordsToCollectionsRequestJson(recordsJsonBytes)
	if err != nil {
		return nil, err
	}

	return s.doCollectionsRequest(s.getCollectionsUrl(), fasthttp.MethodPost, body)
}

// CollectionsUpdateObjects updates objects in salesforce using the composite
// "collections" api. this implementation requires that you marshal your
// objects to json yourself with an attributes field to your object to define
// the type. the objects must also have an id field.
//
// ex:
//
//	{
//	  "attributes" : {"type" : "Account"},
//	  "id" : "001RM0000068xVCYAY",
//	  "Name" : "Example"
//	  ...
//	}
//
// ref: https://developer.salesforce.com/docs/atlas.en-us.api_rest.meta/api_rest/resources_composite_sobjects_collections_update.htm
func (s *SalesforceUtils) CollectionsUpdateObjects(recordsJsonBytes [][]byte) (response []CollectionsResponseItem, err error) {
	err = validateCollectionsRequestLength(len(recordsJsonBytes))
	if err != nil {
		return nil, err
	}

	body, err := jsonRecordsToCollectionsRequestJson(recordsJsonBytes)
	if err != nil {
		return nil, err
	}

	return s.doCollectionsRequest(s.getCollectionsUrl(), fasthttp.MethodPatch, body)
}

// CollectionsDeleteRequest is used by the CollectionsDeleteObjects method when
// interacting with the composite collections api
type CollectionsDeleteRequest struct {
	AllOrNone bool     `json:"allOrNone"`
	Ids       []string `json:"ids"`
}

// CollectionsDeleteObjects deletes objects in salesforce using the composite
// "collections" api. all that is required is the IDs of the objects to delete.
func (s *SalesforceUtils) CollectionsDeleteObjects(ids []string) (response []CollectionsResponseItem, err error) {
	err = validateCollectionsRequestLength(len(ids))
	if err != nil {
		return nil, err
	}
	deleteUrl := s.getCollectionsDeleteUrl(ids)
	return s.doCollectionsRequest(deleteUrl, fasthttp.MethodDelete, nil)
}

// doCollectionsRequest is a helper method for making requests to the the
// composite collections api. the only difference between creates and updates
// is the method, and with deletes the url requires the ids to be passed as
// query parameters. this parameterizes the url, method, and body so each
// method can make use of it. everything else is the same.
func (s *SalesforceUtils) doCollectionsRequest(url string, method string, body []byte) (response []CollectionsResponseItem, err error) {
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)
	req.SetRequestURI(url)
	req.Header.SetMethod(method)
	req.Header.Set("Content-Type", "application/json")
	if body != nil {
		req.SetBody(body)
	}

	body, statusCode, deferredFunc, err := s.sendRequest(req)
	defer deferredFunc()
	if err != nil {
		return response, err
	}
	if statusCode != fasthttp.StatusOK {
		return response, errorx.IllegalState.New("unexpected status code: %d with body: %s", statusCode, body)
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		return response, errorx.IllegalState.New("failed to unmarshal response: %s", err)
	}

	// check for errors in the response and return an error if any are found so
	// that the caller doesn't have to check the response for errors
	for _, respItem := range response {
		if !respItem.Success || len(respItem.Errors) > 0 {
			// return the first error, since allOrNone is always true
			return response, errorx.IllegalState.New("failed to update object: %s", respItem.Errors)
		}
	}

	return response, nil
}

func jsonRecordsToCollectionsRequestJson(recordsJsonBytes [][]byte) ([]byte, error) {
	collectionsReq := CollectionsRequest{
		AllOrNone: true,
	}
	for _, recordJsonBytes := range recordsJsonBytes {
		collectionsReq.Records = append(collectionsReq.Records, recordJsonBytes)
	}

	return json.Marshal(collectionsReq)
}

// validateCollectionsRequestLength is a simple helper method to validate the
// length of input for any collections request. pass the length of the input.
// salesforce collections api has a limit of 200 records per request.
func validateCollectionsRequestLength(length int) error {
	if length == 0 {
		return errorx.IllegalArgument.New("input must not be empty")
	}
	if length > 200 {
		return errorx.IllegalArgument.New("input must not be larger than 200")
	}
	return nil
}

// getDeleteUrl gets a formatted full url to the collections api for deleting.
// builds query parameters to include each id and adds the allOrNone=true
// parameter.
func (s *SalesforceUtils) getCollectionsDeleteUrl(ids []string) string {
	idsAsCommaSeparatedString := strings.Join(ids, ",")
	queryParams := fmt.Sprintf("?allOrNone=true&ids=%s", idsAsCommaSeparatedString)
	return fmt.Sprintf("%s%s", s.getCollectionsUrl(), queryParams)
}

// getCollectionsUrl gets a formatted full url to the collections api.
func (s *SalesforceUtils) getCollectionsUrl() string {
	return fmt.Sprintf("%s/services/data/v%s/composite/sobjects", s.Config.BaseUrl, s.Config.ApiVersion)
}
